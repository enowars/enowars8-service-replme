package server

import (
	"net/http"
	"os"

	"replme/controller"
	"replme/service"
	"replme/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func NewRouter(docker *service.DockerService) *gin.Engine {

	logLevel, exists := os.LookupEnv("REPL_LOG")
	if !exists {
		logLevel = "info"
	}

	setupCors := false
	if _, exists := os.LookupEnv("REPL_CORS"); exists {
		setupCors = true
	}

	util.LoggerInit(logLevel)

	replState := service.ReplState()
	userController := controller.NewUserController(&replState)
	replController := controller.NewReplController(docker, &replState)

	cleanup := service.Cleanup(docker, &replState)
	cleanup.DoCleanup()
	cleanup.StartTask()

	engine := gin.Default()

	if setupCors {
		engine.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	}

	secret, _ := util.RandomBytes(64)
	store := memstore.NewStore(secret)
	if setupCors {
		store.Options(sessions.Options{
			Path:     "/api",
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})
	} else {
		store.Options(sessions.Options{
			Path:   "/api",
			Secure: true,
		})
	}
	engine.Use(sessions.Sessions("session", store))

	/////////////////////// API ///////////////////////

	engine.GET(
		"/api/user/sessions", userController.Sessions,
	)

	engine.POST(
		"/api/repl", replController.Create,
	)

	engine.GET(
		"/api/repl/:name", replController.Websocket,
	)

	return engine
}
