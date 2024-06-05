package server

import (
	"net/http"
	"path"

	"replme/controller"
	"replme/middleware"
	"replme/service"
	"replme/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func NewRouter(docker *service.DockerService, dist string) *gin.Engine {

	terminalController := controller.NewTerminalController(docker)
	replController := controller.NewReplController(docker)

	engine := gin.Default()

	secret, _ := util.RandomBytes(64)
	store := memstore.NewStore(secret)
	engine.Use(sessions.Sessions("session", store))

	engine.Static("/static", path.Join(dist, "static"))
	engine.LoadHTMLGlob(path.Join(dist, "*.html"))

	engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	engine.GET("/term", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "term.html", gin.H{})
	})

	/////////////////////// API ///////////////////////

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

	////////////////////// APIV2 //////////////////////

	engine.POST(
		"/api/repl", replController.Create,
	)

	engine.GET(
		"/api/repl/:replUuid", replController.Websocket,
	)

	engine.POST(
		"/api/repl/:replUuid/size", replController.Resize,
	)

	return engine
}
