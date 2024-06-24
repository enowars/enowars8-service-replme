package controller

import (
	"net/http"
	"regexp"
	"replme/database"
	"replme/model"
	"replme/types"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func NewAuthController() AuthController {
	return AuthController{}
}

func (auth *AuthController) Register(ctx *gin.Context) {
	var registerRequest types.RegisterRequest
	if err := ctx.ShouldBind(&registerRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(registerRequest.Username) < 4 || len(registerRequest.Username) > 64 || !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(registerRequest.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Illegal username"})
		return
	}

	if len(registerRequest.Password) < 4 || len(registerRequest.Password) > 64 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Illegal username"})
		return
	}

	var userFound model.User
	database.DB.Where("username = ?", registerRequest.Username).Find(&userFound)

	if userFound.ID != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User exists"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Username: registerRequest.Username,
		Password: string(passwordHash),
	}

	database.DB.Create(&user)

	ctx.Status(http.StatusOK)
}

func (auth *AuthController) Login(ctx *gin.Context) {
	var loginRequest types.LoginRequest
	if err := ctx.ShouldBind(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(loginRequest.Username) < 4 || len(loginRequest.Username) > 64 || !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(loginRequest.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Illegal username"})
		return
	}

	if len(loginRequest.Password) < 4 || len(loginRequest.Password) > 64 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Illegal username"})
		return
	}

	var userFound model.User
	database.DB.Where("username = ?", loginRequest.Username).Find(&userFound)

	if userFound.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(loginRequest.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	session := sessions.Default(ctx)
	authType := session.Get("auth_type")
	if authType == nil || authType == "bare" {
		session.Set("auth_type", "full")
		session.Set("current_user_id", userFound.ID)
		session.Save()
	}

	ctx.Status(http.StatusOK)
}

func (auth *AuthController) GetUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	authType := session.Get("auth_type")
	currentUserId := session.Get("current_user_id")

	if authType == nil || authType != "full" || currentUserId == nil || currentUserId == uint(0) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var user []model.User
	database.DB.Where("id = ?", currentUserId).Find(&user)

	if len(user) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &gin.H{
			"error": "User not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, user[0])
}

func (auth *AuthController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Set("auth_type", "bare")
	session.Delete("current_user_id")
	session.Save()
	ctx.Status(http.StatusOK)
}
