package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"image-go/service"
)

type TermController struct {
	TermService service.TermService
	UserService service.UserService
	Upgrader    websocket.Upgrader
}

func NewTermController() TermController {
	return TermController{
		TermService: service.NewTermService(),
		UserService: service.NewUserService(),
		Upgrader: websocket.Upgrader{
			// ReadBufferSize:  1024,
			// WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (term *TermController) Websocket(ctx *gin.Context) {
	session := sessions.Default(ctx)
	username := session.Get("username").(string)

	user, errResp := term.UserService.GetUserData(username)

	if errResp != nil {
		ctx.JSON(errResp.Code, gin.H{"error": errResp.Message})
		ctx.Abort()
		return
	}

	conn, err := term.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer conn.Close()

	term.TermService.Create(ctx, conn, user)
}

func (term *TermController) WebsocketExec(ctx *gin.Context) {
	cwd := ctx.Query("cwd")
	command := ctx.Query("command")

	if cwd == "" || command == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad query params"})
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)
	username := session.Get("username").(string)

	user, errResp := term.UserService.GetUserData(username)

	if errResp != nil {
		ctx.JSON(errResp.Code, gin.H{"error": errResp.Message})
		ctx.Abort()
		return
	}

	conn, err := term.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer conn.Close()

	term.TermService.Exec(ctx, conn, user, cwd, command)
}
