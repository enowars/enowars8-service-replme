package controller

import (
	"replme/service"
	"replme/types"
	"replme/util"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type TerminalController struct {
	Docker   *service.DockerService
	Upgrader websocket.Upgrader
	CRC      util.CRCUtil
}

func NewTerminalController(docker *service.DockerService) TerminalController {
	return TerminalController{
		Docker: docker,
		Upgrader: websocket.Upgrader{
			// ReadBufferSize:  1024,
			// WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		CRC: util.CRC(),
	}
}

func (term *TerminalController) EnsureUser(ctx *gin.Context) {
	var loginRequest types.LoginRequest
	if err := ctx.ShouldBind(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash := term.CRC.Calculate(util.DecodeSpecialChars([]byte(loginRequest.Username)))
	name := fmt.Sprintf("%x", hash)

	_, port, err := term.Docker.EnsureContainerStarted(
		name,
		loginRequest.Username,
		loginRequest.Password,
	)

	if err != nil {
		ctx.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var response *http.Response
	p := service.Proxy(term.Docker.HostIP, *port)

	for i := 0; i < 5; i++ {
		response, err = p.SendUserCreateRequest(types.LoginRequest{
			Username: loginRequest.Username,
			Password: loginRequest.Password,
		})
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err != nil {
		ctx.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if response.StatusCode >= 400 {
		payload, err := io.ReadAll(response.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
				Error: "Container communication failed",
			})
			return
		}
		ctx.Data(response.StatusCode, "application/json", payload)
		return
	}

	ctx.JSON(response.StatusCode, types.SuccessResponse{
		Success: "Success",
		Port:    *port,
	})

}

func (term *TerminalController) CreateTerminal(ctx *gin.Context) {
	port := ctx.MustGet("port").(uint16)

	createTermProcessRequest := types.CreateTermProcessRequest{
		Rows:    80,
		Columns: 24,
	}

	if err := ctx.ShouldBind(&createTermProcessRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := service.Proxy(term.Docker.HostIP, port)
	response, err := p.SendCreateTermProcessRequest(
		createTermProcessRequest,
	)

	if err != nil {
		ctx.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	payload, err := io.ReadAll(response.Body)

	if err != nil {
		ctx.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.Data(response.StatusCode, response.Header.Get("Content-Type"), payload)
}

func (term *TerminalController) ResizeTerminal(ctx *gin.Context) {
	port := ctx.MustGet("port").(uint16)
	pid := ctx.Param("pid")

	var updateTermSizeRequest types.UpdateTermSizeRequest

	if err := ctx.ShouldBind(&updateTermSizeRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := service.Proxy(term.Docker.HostIP, port)
	response, err := p.SendUpdateTermSizeRequest(
		pid,
		updateTermSizeRequest,
	)

	if err != nil {
		ctx.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	payload, err := io.ReadAll(response.Body)

	if err != nil {
		ctx.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.Data(response.StatusCode, response.Header.Get("Content-Type"), payload)
}

func (term *TerminalController) CreateWebsocketProxy(ctx *gin.Context) {
	portStr := ctx.Param("port")
	portUint, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	port := uint16(portUint)
	pid := ctx.Param("pid")

	clientConn, err := term.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer clientConn.Close()

	p := service.Proxy(term.Docker.HostIP, port)
	err = p.CreateWebsocketPipe(clientConn, pid)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
}
