package server

import (
	"net/http"
	"path"

	"replme/controller"
	"replme/middleware"
	"replme/service"

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

	engine.POST(
		"/api/login/private",
		terminalController.EnsureUser,
	)
	engine.GET(
		"/ws/terminal/:port/:pid",
		terminalController.CreateWebsocketProxy,
	)

	term := engine.Group(
		"/api/terminal",
		middleware.AuthMiddleware(docker),
	)
	term.POST(
		"/",
		terminalController.CreateTerminal,
	)
	term.POST(
		"/:pid/size",
		terminalController.ResizeTerminal,
	)

	return engine
}
