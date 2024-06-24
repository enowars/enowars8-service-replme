package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	// xterm "golang.org/x/term"

	"image-go/types"
)

type ResizeMessage struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}

type WebsocketMessage struct {
	Stdin  *string        `json:"stdin"`
	Resize *ResizeMessage `json:"resize"`
}

type WebsocketReadWriter struct {
	Conn *websocket.Conn
	Pty  *os.File
}

func (srw WebsocketReadWriter) Write(p []byte) (n int, err error) {
	err = srw.Conn.WriteMessage(websocket.TextMessage, p)
	fmt.Printf("term -> ws: %s\n", string(p))
	return len(p), err
}

func (srw WebsocketReadWriter) Read(p []byte) (n int, err error) {
	_, b, err := srw.Conn.ReadMessage()

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	var wsMessage WebsocketMessage
	err = json.Unmarshal(b, &wsMessage)

	fmt.Printf("ws -> term: %s\n", string(b))

	if err == nil {
		if wsMessage.Resize != nil {
			pty.Setsize(srw.Pty, &pty.Winsize{
				Rows: uint16(wsMessage.Resize.Rows),
				Cols: uint16(wsMessage.Resize.Cols),
				X:    0,
				Y:    0,
			})
			return 0, nil
		}
		if wsMessage.Stdin != nil {
			b = []byte(*wsMessage.Stdin)
		}
	}

	for i, d := range b {
		p[i] = d
	}

	return len(b), nil
}

type TermService struct{}

func NewTermService() TermService {
	return TermService{}
}

func (term *TermService) Create(ctx *gin.Context, websocket *websocket.Conn, user *types.UserPasswdData) {
	cmd := exec.Command(user.Shell)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: user.Uid,
		Gid: user.Gid,
	}
	cmd.Dir = user.Home
	cmd.Env = []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		"COLORTERM=truecolor",
		"TERM=xterm",
	}

	ptmx, err := pty.Start(cmd)
	defer func() { cmd.Process.Kill() }()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer func() { _ = ptmx.Close() }()

	pty.Setsize(ptmx, &pty.Winsize{})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH
	defer func() { signal.Stop(ch); close(ch) }()

	// oldState, err := xterm.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() { _ = xterm.Restore(int(os.Stdin.Fd()), oldState) }()

	wsrw := WebsocketReadWriter{
		Conn: websocket,
		Pty:  ptmx,
	}

	errc := make(chan error, 2)
	go func() {
		_, err := io.Copy(wsrw, ptmx)
		errc <- err
	}()
	go func() {
		_, err := io.Copy(ptmx, wsrw)
		errc <- err
	}()

	if err := <-errc; err != nil {
		log.Printf("WebSocket proxy error: %v", err)
	}
}

func (term *TermService) Exec(
	ctx *gin.Context,
	websocket *websocket.Conn,
	user *types.UserPasswdData,
	cwd string,
	command string,
) {
	cmd := exec.Command(
		"/bin/sh",
		"-c",
		fmt.Sprintf(
			"%s && echo SUCCEEDED || echo FAILED",
			strings.ReplaceAll(command, "\"", "\\\""),
		),
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: user.Uid,
		Gid: user.Gid,
	}
	cmd.Dir = cwd
	cmd.Env = []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		"COLORTERM=truecolor",
		"TERM=xterm",
	}

	ptmx, err := pty.Start(cmd)
	defer func() { cmd.Process.Kill() }()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer func() { _ = ptmx.Close() }()

	pty.Setsize(ptmx, &pty.Winsize{})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH
	defer func() { signal.Stop(ch); close(ch) }()

	// oldState, err := xterm.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() { _ = xterm.Restore(int(os.Stdin.Fd()), oldState) }()

	wsrw := WebsocketReadWriter{
		Conn: websocket,
		Pty:  ptmx,
	}

	errc := make(chan error, 2)
	go func() {
		_, err := io.Copy(wsrw, ptmx)
		errc <- err
	}()
	go func() {
		_, err := io.Copy(ptmx, wsrw)
		errc <- err
	}()

	if err := <-errc; err != nil {
		log.Printf("WebSocket proxy error: %v", err)
	}
}
