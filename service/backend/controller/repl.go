package controller

import (
	"fmt"
	"net/http"
	"replme/service"
	"replme/types"
	"replme/util"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ReplController struct {
	Docker    *service.DockerService
	ReplState *service.ReplStateService
	Upgrader  websocket.Upgrader
	CRC       util.CRCUtil
}

func NewReplController(docker *service.DockerService, replState *service.ReplStateService) ReplController {
	return ReplController{
		Docker:    docker,
		ReplState: replState,
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

func (repl *ReplController) Create(ctx *gin.Context) {
	var loginRequest types.LoginRequest
	if err := ctx.ShouldBind(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash := repl.CRC.Calculate(util.DecodeSpecialChars([]byte(loginRequest.Username)))
	name := fmt.Sprintf("%x", hash)

	_, port, err := repl.Docker.EnsureContainerStarted(name)

	if err != nil {
		ctx.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	p := service.Proxy(repl.Docker.HostIP, *port, repl.Docker.ApiKey)
	var response *http.Response
	var requestError *types.RequestError
	response, requestError = p.SendRegisterRequest(
		types.RegisterRequest{
			Username: loginRequest.Username,
			Password: loginRequest.Password,
		},
		&types.RequestOptions{
			Retries: 10,
		},
	)

	if requestError != nil {
		ctx.Data(requestError.Code, requestError.ContentType, requestError.Data)
		return
	}

	response, requestError = p.SendCreateTermRequest(types.RequestOptions{
		Cookies: response.Cookies(),
	})

	if requestError != nil {
		ctx.Data(requestError.Code, requestError.ContentType, requestError.Data)
		return
	}

	session := sessions.Default(ctx)
	repl.ReplState.AddContainer(
		session.ID(),
		name,
		loginRequest.Username,
		loginRequest.Password,
	)

	uuid := repl.ReplState.AddReplSession(
		session.ID(),
		name,
		*response.Cookies()[0],
	)

	ctx.JSON(response.StatusCode, types.CreateReplResponse{
		ReplUuid: uuid,
	})
}

func (repl *ReplController) CreateFromName(ctx *gin.Context) {
	session := sessions.Default(ctx)
	name := ctx.Param("name")

	data := repl.ReplState.GetContainer(session.ID(), name)

	if data == nil {
		ctx.JSON(404, types.ErrorResponse{
			Error: "Container not found",
		})
		return
	}

	_, port, err := repl.Docker.EnsureContainerStarted(name)

	if err != nil {
		ctx.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	p := service.Proxy(repl.Docker.HostIP, *port, repl.Docker.ApiKey)
	var response *http.Response
	var requestError *types.RequestError
	response, requestError = p.SendRegisterRequest(
		types.RegisterRequest{
			Username: data.Username,
			Password: data.Password,
		},
		&types.RequestOptions{
			Retries: 10,
		},
	)

	if requestError != nil {
		ctx.Data(requestError.Code, requestError.ContentType, requestError.Data)
		return
	}

	response, requestError = p.SendCreateTermRequest(types.RequestOptions{
		Cookies: response.Cookies(),
	})

	if requestError != nil {
		ctx.Data(requestError.Code, requestError.ContentType, requestError.Data)
		return
	}

	uuid := repl.ReplState.AddReplSession(
		session.ID(),
		name,
		*response.Cookies()[0],
	)

	ctx.JSON(response.StatusCode, types.CreateReplResponse{
		ReplUuid: uuid,
	})
}

func (repl *ReplController) Resize(ctx *gin.Context) {
	session := sessions.Default(ctx)
	replUuid := ctx.Param("replUuid")

	var resizeTermRequest types.ResizeTermRequest
	if err := ctx.ShouldBind(&resizeTermRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name, cookie, exists := repl.ReplState.GetReplSession(session.ID(), replUuid)
	if !exists {
		ctx.JSON(404, types.ErrorResponse{
			Error: "Container not found",
		})
		return
	}

	port := repl.Docker.GetContainerPort(*name)
	if port == nil {
		ctx.JSON(404, types.ErrorResponse{
			Error: "Container/Port not found",
		})
		return
	}
	p := service.Proxy(repl.Docker.HostIP, *port, repl.Docker.ApiKey)

	response, requestError := p.SendResizeTermRequest(resizeTermRequest, types.RequestOptions{
		Cookies: []*http.Cookie{cookie},
	})
	if requestError != nil {
		ctx.Data(requestError.Code, requestError.ContentType, requestError.Data)
		return
	}
	ctx.Status(response.StatusCode)
	return
}

func (repl *ReplController) Websocket(ctx *gin.Context) {
	session := sessions.Default(ctx)
	replUuid := ctx.Param("replUuid")
	name, cookie, exists := repl.ReplState.GetReplSession(session.ID(), replUuid)

	if !exists {
		ctx.JSON(404, types.ErrorResponse{
			Error: "Container not found",
		})
		return
	}

	port := repl.Docker.GetContainerPort(*name)
	if port == nil {
		ctx.JSON(404, types.ErrorResponse{
			Error: "Container/Port not found",
		})
		return
	}
	p := service.Proxy(repl.Docker.HostIP, *port, repl.Docker.ApiKey)
	clientConn, err := repl.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer func() {
		repl.ReplState.DeleteReplSession(
			session.ID(),
			*name,
			replUuid,
		)
		go func() {
			time.Sleep(30 * time.Second)
			if !repl.ReplState.HasActiveReplSessions(session.ID(), *name) {
				repl.Docker.KillContainerByName(*name)
			}
		}()
	}()
	defer clientConn.Close()
	err = p.CreateWebsocketPipe(clientConn, *cookie)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
}
