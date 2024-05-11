package server

import (
	"net/http"
	"path"

	"cafedodo/controller"
	"cafedodo/service"
	"cafedodo/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(docker *service.DockerService, dist string) *gin.Engine {

	terminalController := controller.NewTerminalController(docker)

	engine := gin.Default()

	engine.Static("/static", path.Join(dist, "static"))
	engine.LoadHTMLGlob(path.Join(dist, "*.html"))

	engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	engine.GET("/term/:username/:port", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "term.html", gin.H{})
	})

	engine.POST("/api/term/private", terminalController.CreateUser)
	engine.GET("/ws/terminal/:port/:pid", terminalController.CreateWebsocketProxy)

	term := engine.Group("/api/ptwhy", middleware.AuthMiddleware(docker))
	term.POST("/terminals", terminalController.CreateTerminal)
	term.POST("/terminals/:pid/size", terminalController.ResizeTerminal)

	return engine
}
