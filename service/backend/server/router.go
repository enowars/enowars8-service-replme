package server

import (
	"net/http"
	"path"

	"replme/controller"
	"replme/service"
	"replme/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func NewRouter(docker *service.DockerService, dist string) *gin.Engine {

	replState := service.ReplState()
	userController := controller.NewUserController(&replState)
	replController := controller.NewReplController(docker, &replState)

	cleanup := service.Cleanup(docker, &replState)
	cleanup.DoCleanup()
	cleanup.StartTask()

	engine := gin.Default()

	secret, _ := util.RandomBytes(64)
	store := memstore.NewStore(secret)
	engine.Use(sessions.Sessions("session", store))

	engine.Static("/static", path.Join(dist, "static"))
	engine.LoadHTMLGlob(path.Join(dist, "*.html"))

	engine.Use(func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		temp_auth := session.Get("temp_auth")
		if temp_auth == nil {
			session.Set("temp_auth", true)
			session.Save()
		}
		ctx.Next()
	})

	engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	engine.GET("/term", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "term.html", gin.H{})
	})

	engine.GET("/term/:name", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "term.html", gin.H{})
	})

	/////////////////////// API ///////////////////////

	engine.GET(
		"/api/user/sessions", userController.Sessions,
	)

	// engine.GET(
	// 	"/api/user/sessions/debug", userController.SessionsDebug,
	// )

	engine.POST(
		"/api/repl", replController.Create,
	)

	engine.POST(
		"/api/repl/name/:name", replController.CreateFromName,
	)

	engine.GET(
		"/api/repl/:replUuid", replController.Websocket,
	)

	engine.POST(
		"/api/repl/:replUuid/size", replController.Resize,
	)

	return engine
}
