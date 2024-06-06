package controller

import (
	"net/http"
	"replme/service"

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
	containers := user.ReplState.GetContainers(session.ID())

	names := []string{}

	for name := range *containers {
		names = append(names, name)
	}

	ctx.JSON(http.StatusOK, names)
}

func (user *UserController) SessionsDebug(ctx *gin.Context) {
	session := sessions.Default(ctx)
	containers := user.ReplState.GetContainers(session.ID())

	ctx.JSON(http.StatusOK, *containers)
}
