package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"image-go/service"
	"image-go/types"
)

type UserController struct {
	Service service.UserService
}

func NewUserController() UserController {
	return UserController{
		Service: service.NewUserService(),
	}
}

func (user *UserController) Login(ctx *gin.Context) {
	var credentials = types.CredentialRequest{}
	if err := ctx.ShouldBind(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := user.Service.Login(credentials.Username, credentials.Password)

	if err != nil {
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	session := sessions.Default(ctx)
	session.Set("username", credentials.Username)
	session.Save()

	ctx.JSON(result.Code, gin.H{"success": result.Message})
}

func (user *UserController) Register(ctx *gin.Context) {
	var credentials = types.CredentialRequest{}
	if err := ctx.ShouldBind(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := user.Service.Register(credentials.Username, credentials.Password)

	if err != nil {
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	session := sessions.Default(ctx)
	session.Set("username", credentials.Username)
	session.Save()

	ctx.JSON(result.Code, gin.H{"success": result.Message})
}
