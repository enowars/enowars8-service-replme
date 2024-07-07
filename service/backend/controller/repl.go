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
	var createReplRequest types.CreateReplRequest
	if err := ctx.ShouldBind(&createReplRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	util.SLogger.Debugf("[%-25s] Creating new REPL user", fmt.Sprintf("UN:%s..", createReplRequest.Username[:5]))

	hash := repl.CRC.Calculate(util.DecodeSpecialChars([]byte(createReplRequest.Username)))
	name := fmt.Sprintf("%x", hash)

	util.SLogger.Debugf("[%-25s] Created new REPL user", fmt.Sprintf("UN:%s.. | NM:%s..", createReplRequest.Username[:5], name[:5]))

	session := sessions.Default(ctx)
	auth_type := session.Get("auth_type")
	if auth_type == nil {
		session.Set("auth_type", "bare")
		session.Save()
	}

	util.SLogger.Debugf("[%-25s] Saving session %s..", fmt.Sprintf("UN:%s.. | NM:%s..", createReplRequest.Username[:5], name[:5]), session.ID()[:5])

	repl.ReplState.AddUserSession(session.ID(), name, createReplRequest.Username, createReplRequest.Password)

	ctx.JSON(http.StatusOK, types.AddReplUserResponse{
		Id: name,
	})
	return
}

func (repl *ReplController) Sessions(ctx *gin.Context) {
	session := sessions.Default(ctx)
	auth_type := session.Get("auth_type")
	if session.ID() == "" || auth_type == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}
	util.SLogger.Debugf("[%-25s] Get sessions", fmt.Sprintf("ID:%s..", session.ID()[:5]))
	names := repl.ReplState.GetContainerNames(session.ID())
	ctx.JSON(http.StatusOK, names)
}

func (repl *ReplController) Websocket(ctx *gin.Context) {
	session := sessions.Default(ctx)

	auth_type := session.Get("auth_type")
	if session.ID() == "" || auth_type == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	name := ctx.Param("name")

	user := repl.ReplState.GetUserSessionData(session.ID(), name)
	util.SLogger.Debugf("Username: %s, Password: %s", user.Username, user.Password)

	if user == nil {
		util.SLogger.Warnf("[%-25s] Attempted websocket: User data not existing", fmt.Sprintf("ID:%s.. | NM:%s..", session.ID()[:5], name[:5]))
		ctx.JSON(http.StatusNotFound, types.ErrorResponse{
			Error: "User data not found",
		})
		return
	}

	util.SLogger.Debugf("Adding container session")
	repl.ReplState.AddContainerSession(name)

	defer func() {
		time.Sleep(5 * time.Second)
		var lock util.Unlocker
		util.SLogger.Debugf("Deleting container session")
		kill := repl.ReplState.DeleteContainerSession(name, func(name string) {
			lock = repl.Docker.MutexMap.Lock(name)
		})

		if kill {
			util.SLogger.Infof("[%-25s] Killing container", fmt.Sprintf("NM:%s..", name[:5]))
			start := time.Now()
			repl.Docker.KillContainerByName(name)
			util.SLogger.Infof("[%-25s] Killing container took %v", fmt.Sprintf("NM:%s..", name[:5]), time.Since(start))
			lock.Unlock()
		}
	}()

	util.SLogger.Debugf("[%-25s] Creating REPL", fmt.Sprintf("UN:%s.. | NM:%s..", user.Username[:5], name[:5]))
	util.SLogger.Debugf("[%-25s] Starting container", fmt.Sprintf("UN:%s.. | NM:%s..", user.Username[:5], name[:5]))
	start := time.Now()
	lock := repl.Docker.MutexMap.Lock(name)
	util.SLogger.Debugf("[%-25s] Locking Mutex took %v", fmt.Sprintf("UN:%s.. | NM:%s..", user.Username[:5], name[:5]), time.Since(start))
	start = time.Now()
	_, port, err := repl.Docker.EnsureReplContainerStarted(name)
	util.SLogger.Debugf("[%-25s] Starting container took %v", fmt.Sprintf("UN:%s.. | NM:%s..", user.Username[:5], name[:5]), time.Since(start))
	lock.Unlock()

	if err != nil {
		util.SLogger.Warnf("[%-25s] Creating container failed, %s", fmt.Sprintf("UN:%s.. | NM:%s..", user.Username[:5], name[:5]), err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, types.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	p := service.Proxy(repl.Docker.HostIP, *port, repl.Docker.ApiKey)

	response, requestError := p.SendRegisterRequest(
		types.RegisterRequest{
			Username: user.Username,
			Password: user.Password,
		},
		&types.RequestOptions{
			Retries: 10,
		},
	)

	if requestError != nil {
		ctx.Data(requestError.Code, requestError.ContentType, requestError.Data)
		ctx.Abort()
		return
	}

	cookie := response.Cookies()[0]

	clientConn, err := repl.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer clientConn.Close()
	err = p.CreateReplWebsocketPipe(clientConn, *cookie)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
}
