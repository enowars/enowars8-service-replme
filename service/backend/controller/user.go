package controller

import (
	"fmt"
	"net/http"
	"replme/service"
	"replme/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	ReplState *service.ReplStateService
}

func NewUserController(replState *service.ReplStateService) UserController {
	return UserController{
		ReplState: replState,
	}
}

func (user *UserController) Sessions(ctx *gin.Context) {
	session := sessions.Default(ctx)
	names := user.ReplState.GetContainerNames(session.ID())
	util.SLogger.Debugf("[%-25s] Get sessions", fmt.Sprintf("ID:%s..", session.ID()[:5]))
	ctx.JSON(http.StatusOK, names)
}
