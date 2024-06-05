package controller

import (
	"fmt"
	"net/http"
	"replme/service"
	"replme/types"
	"replme/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ReplController struct {
	Docker   *service.DockerService
	Upgrader websocket.Upgrader
	CRC      util.CRCUtil
}

func NewReplController(docker *service.DockerService) ReplController {
	return ReplController{
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
	containers := types.Containers{}
	if c := session.Get("containers"); c != nil {
		containers = c.(types.Containers)
	}

	uid := uuid.New().String()

	if _, exists := containers[name]; !exists {
		containers[name] = types.ContainerData{
			Username: loginRequest.Username,
			Password: loginRequest.Password,
			Sessions: map[string]http.Cookie{},
		}
	}

	containers[name].Sessions[uid] = *response.Cookies()[0]

	session.Set("containers", containers)
	session.Save()

	ctx.JSON(response.StatusCode, types.CreateReplResponse{
		ReplUuid: uid,
	})
}

func (repl *ReplController) Resize(ctx *gin.Context) {
	var resizeTermRequest types.ResizeTermRequest

	if err := ctx.ShouldBind(&resizeTermRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	replUuid := ctx.Param("replUuid")

	session := sessions.Default(ctx)
	containers := types.Containers{}
	if c := session.Get("containers"); c != nil {
		containers = c.(types.Containers)
	}

	for name, data := range containers {
		if cookie, exists := data.Sessions[replUuid]; exists {
			port := repl.Docker.GetContainerPort(name)
			if port == nil {
				ctx.JSON(404, types.ErrorResponse{
					Error: "Container/Port not found",
				})
				return
			}
			p := service.Proxy(repl.Docker.HostIP, *port, repl.Docker.ApiKey)

			response, requestError := p.SendResizeTermRequest(resizeTermRequest, types.RequestOptions{
				Cookies: []*http.Cookie{&cookie},
			})
			if requestError != nil {
				ctx.Data(requestError.Code, requestError.ContentType, requestError.Data)
				return
			}
			ctx.Status(response.StatusCode)
			return
		}
	}

	ctx.JSON(404, types.ErrorResponse{
		Error: "Container not found",
	})
}

func (repl *ReplController) Websocket(ctx *gin.Context) {

	replUuid := ctx.Param("replUuid")

	session := sessions.Default(ctx)
	containers := types.Containers{}
	if c := session.Get("containers"); c != nil {
		containers = c.(types.Containers)
	}

	for name, data := range containers {
		if cookie, exists := data.Sessions[replUuid]; exists {
			port := repl.Docker.GetContainerPort(name)
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
			defer clientConn.Close()
			err = p.CreateWebsocketPipeV2(clientConn, cookie)
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
		}
	}
}
