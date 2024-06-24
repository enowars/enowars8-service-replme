package server

import (
	"net/http"

	"image-go/controller"
	"image-go/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func NewRouter(apiKey string) *gin.Engine {

	userController := controller.NewUserController()
	termController := controller.NewTermController()

	engine := gin.Default()

	secret, _ := util.RandomBytes(64)
	store := memstore.NewStore(secret)
	engine.Use(sessions.Sessions("session", store))

	api := engine.Group("/api/:apiKey", func(ctx *gin.Context) {
		key := ctx.Param("apiKey")
		if key != apiKey {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid Apikey"})
			ctx.Abort()
			return
		}
		ctx.Next()
	})

	api.POST("/auth/register", userController.Register)
	api.POST("/auth/login", userController.Login)

	term := api.Group("/term", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		username := session.Get("username")
		if username == nil || username.(string) == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Next()
	})

	term.GET("", termController.Websocket)
	term.GET("/exec", termController.WebsocketExec)

	return engine
}
