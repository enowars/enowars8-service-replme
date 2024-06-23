package controller

import (
	"fmt"
	"net/http"
	"replme/service"
	"replme/types"
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
	auth_type := session.Get("auth_type")
	if session.ID() == "" || auth_type == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}
	util.SLogger.Debugf("[%-25s] Get sessions", fmt.Sprintf("ID:%s..", session.ID()[:5]))
	names := user.ReplState.GetContainerNames(session.ID())
	ctx.JSON(http.StatusOK, names)
}
