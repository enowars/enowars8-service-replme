package server

import (
	"net/http"
	"os"

	"replme/controller"
	"replme/database"
	"replme/model"
	"replme/service"
	"replme/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func NewRouter(docker *service.DockerService, pgUrl string, containerLogsPath string, devenvFilesPath string, devenvFilesTmpPath string) *gin.Engine {

	logLevel, exists := os.LookupEnv("REPL_LOG")
	if !exists {
		logLevel = "info"
	}

	util.LoggerInit(logLevel)

	util.SLogger.Info("Connecting to DB ..")
	database.Connect(pgUrl)
	util.SLogger.Info("Migrating DB ..")
	database.Migrate()

	setupCors := false
	if _, exists := os.LookupEnv("REPL_CORS"); exists {
		setupCors = true
	}

	replState := service.ReplState()
	authController := controller.NewAuthController()
	devenvController := controller.NewDevenvController(docker, devenvFilesPath, devenvFilesTmpPath)
	replController := controller.NewReplController(docker, &replState)

	cleanup := service.Cleanup(docker, &replState, containerLogsPath, devenvFilesPath, devenvFilesTmpPath)
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

	engine.POST("/api/auth/register", authController.Register)
	engine.POST("/api/auth/login", authController.Login)
	engine.GET("/api/auth/user", authController.GetUser)
	engine.POST("/api/auth/logout", authController.Logout)

	devenvs := engine.Group("/api/devenv", func(ctx *gin.Context) {
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

		ctx.Set("current_user", user[0])
		ctx.Next()
	})

	devenvs.GET("", devenvController.GetAll)
	devenvs.POST("", devenvController.Create)

	devenv := devenvs.Group("/:uuid", func(ctx *gin.Context) {
		_user, _ := ctx.Get("current_user")
		user := _user.(model.User)

		id := ctx.Query("uuid")
		if id == "" {
			id = ctx.Param("uuid")
		}
		util.SLogger.Debugf("id: %s", id)
		uuid := util.ExtractUuid(id)

		if uuid == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, &gin.H{
				"error": "Invalid uuid",
			})
			return
		}

		var devenvs []model.Devenv
		err := database.DB.Model(&user).Where("id = ?", uuid[:36]).Association("Devenvs").Find(&devenvs)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
				"error": err.Error(),
			})
			return
		}

		if len(devenvs) == 0 {
			ctx.AbortWithStatusJSON(http.StatusNotFound, &gin.H{
				"error": "Devenv not found",
			})
			return
		}

		ctx.Set("uuid", uuid)
		ctx.Set("current_devenv", devenvs[0])
		ctx.Next()
	})

	devenv.GET("", devenvController.GetOne)
	devenv.PATCH("", devenvController.Patch)
	devenv.GET("/files", devenvController.GetFiles)
	devenv.POST("/files", devenvController.CreateFile)
	devenv.GET("/files/:name", devenvController.GetFileContent)
	devenv.POST("/files/:name", devenvController.SetFileContent)
	devenv.DELETE("/files/:name", devenvController.DeleteFile)
	devenv.GET("/exec", devenvController.Exec)

	engine.POST("/api/repl", replController.Create)
	engine.GET("/api/repl/sessions", replController.Sessions)

	engine.GET("/api/repl/:name", replController.Websocket)

	return engine
}
