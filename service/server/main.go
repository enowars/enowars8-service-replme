package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cafedodo/client"
	"cafedodo/orchestrator"
	"cafedodo/renderer"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type LoginRequest struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type CreateTermProcessRequest struct {
	Columns int `form:"cols" json:"cols" xml:"cols"`
	Rows int `form:"rows" json:"rows" xml:"rows"`
}

type UpdateTermSizeRequest struct {
	Columns int `form:"cols" json:"cols" xml:"cols" binding:"required"`
	Rows int `form:"rows" json:"rows" xml:"rows" binding:"required"`
}

type SuccessResponse struct {
	Success string `json:"success"`
	Port uint16 `json:"port"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendUserCreateRequest(ip string, port uint16, request LoginRequest) (*http.Response, error) {
	payload, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	return http.Post(
		fmt.Sprintf("http://%s:%d/user/create", ip, port),
		"application/json",
		bytes.NewBuffer(payload),
	)
}


func SendUserLoginRequest(ip string, port uint16, request LoginRequest) (*http.Response, error) {
	payload, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	return http.Post(
		fmt.Sprintf("http://%s:%d/user/login", ip, port),
		"application/json",
		bytes.NewBuffer(payload),
	)
}

func SendCreateTermProcessRequest(ip string, port uint16, request CreateTermProcessRequest) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d/terminals?rows=%d&cols=%d", ip, port, request.Rows, request.Columns)
	req, _ := http.NewRequest("POST", url, nil)

	return http.DefaultClient.Do(req)
}

func SendUpdateTermSizeRequest(ip string, port uint16, pid string, request UpdateTermSizeRequest) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d/terminals/%s/size?rows=%d&cols=%d", ip, port, pid, request.Rows, request.Columns)
	req, _ := http.NewRequest("POST", url, nil)

	return http.DefaultClient.Do(req)
}

func main() {
	imageTag := "ptwhy"

	o := orchestrator.Instance()
	o.BuildImage("./ptwhy/", imageTag)

	var upgrader = websocket.Upgrader{
		// ReadBufferSize:  1024,
		// WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	engine := gin.Default()

	engine.Static("/static", "./static")

	term := engine.Group("/api/ptwhy", func(ctx *gin.Context) {
		authHeader := ctx.Request.Header["Authorization"]

		if len(authHeader) == 0 {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{ "error": "Authorization header missing" },
			)
			ctx.Abort()
			return
		}

		authValues := strings.Split(authHeader[0], " ")

		if len(authValues) < 2 {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{ "error": "Authorization header invalid" },
			)
			ctx.Abort()
			return
		}

		credsBase64 := authValues[1]
		credsString, err := base64.StdEncoding.DecodeString(credsBase64)

		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{ "error": "Authorization value not decodeable" },
			)
			ctx.Abort()
			return
		}

		credsSlice := strings.Split(string(credsString), ":")

		if len(credsSlice) < 2 {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{ "error": "Authorization value invalid" },
			)
			ctx.Abort()
			return
		}

		username := credsSlice[0]
		password := credsSlice[1]

		hash := crc32.ChecksumIEEE([]byte(username))
		name := fmt.Sprint(hash)

		container, _, running := o.GetContainer(name)

		if container == nil {
			ctx.JSON(
				http.StatusPreconditionFailed,
				gin.H{ "error": "Container not found" },
			)
			ctx.Abort()
			return
		}

		if !running {
			ctx.JSON(
				http.StatusPreconditionFailed,
				gin.H{ "error": "Container not running" },
			)
			ctx.Abort()
			return
		}

		_, port, err := o.GetContainerAddress(container.ID)

		if err != nil {
			ctx.JSON(
				http.StatusPreconditionFailed,
				gin.H{ "error": "Could not get container address" },
			)
			ctx.Abort()
			return
		}

		response, err := SendUserLoginRequest(o.HostIP, *port, LoginRequest{
			Username: username,
			Password: password,
		})

		if err != nil {
			ctx.JSON(400, ErrorResponse {
				Error: err.Error(),
			})
			ctx.Abort()
			return
		}

		if response.StatusCode >= 400 {
			payload, err := io.ReadAll(response.Body)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, ErrorResponse {
					Error: "Container communication failed",
				})
				ctx.Abort()
				return
			}
			ctx.Data(response.StatusCode, "application/json", payload)
			ctx.Abort()
			return
		}

		ctx.Set("port", *port)
		ctx.Next()
	})

	ginHtmlRenderer := engine.HTMLRender
	engine.HTMLRender = &renderer.HTMLTemplRenderer {
		FallbackHtmlRenderer: ginHtmlRenderer,
	}

	engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "", client.Home())
	})

	engine.GET("/term/:username/:port", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "", client.Session())
	})

	engine.POST("/api/term/private", func(ctx *gin.Context) {

		var loginRequest LoginRequest
		if err := ctx.ShouldBind(&loginRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash := crc32.ChecksumIEEE([]byte(loginRequest.Username))
		name := fmt.Sprint(hash)

		_, port, err := o.EnsureContainerStarted(
			name,
			loginRequest.Username,
			loginRequest.Password,
		)

		if err != nil {
			ctx.JSON(400, ErrorResponse {
				Error: err.Error(),
			})
			return
		}

		var response *http.Response

		for i := 0; i < 5; i++ {
			response, err = SendUserCreateRequest(o.HostIP, *port, loginRequest)
			if err == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		if err != nil {
			ctx.JSON(400, ErrorResponse {
				Error: err.Error(),
			})
			return
		}

		if response.StatusCode >= 400 {
			payload, err := io.ReadAll(response.Body)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, ErrorResponse {
					Error: "Container communication failed",
				})
				return
			}
			ctx.Data(response.StatusCode, "application/json", payload)
			return
		}

		ctx.JSON(response.StatusCode, SuccessResponse {
			Success: "Success",
			Port: *port,
		})
		
	})

	term.POST("/terminals", func(ctx *gin.Context) {
		port := ctx.MustGet("port").(uint16)

		createTermProcessRequest := CreateTermProcessRequest {
			Rows: 80,
			Columns: 24,
		}

		if err := ctx.ShouldBind(&createTermProcessRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := SendCreateTermProcessRequest(o.HostIP, port, createTermProcessRequest)

		if err != nil {
			ctx.JSON(400, ErrorResponse {
				Error: err.Error(),
			})
			return
		}

		payload, err := io.ReadAll(response.Body)

		if err != nil {
			ctx.JSON(400, ErrorResponse {
				Error: err.Error(),
			})
			return
		}

		ctx.Data(response.StatusCode, response.Header.Get("Content-Type"), payload)
	})

	engine.GET("/ws/ptwhy/terminals/:port/:pid", func(ctx *gin.Context) {
		port := ctx.Param("port")
		pid := ctx.Param("pid")

		clientConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer clientConn.Close()

		targetURL, _ := url.Parse(fmt.Sprintf("ws://%s:%s/terminals/%s", o.HostIP, port, pid))

		targetConn, _, err := websocket.DefaultDialer.Dial(targetURL.String(), nil)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadGateway)
			return
		}
		defer targetConn.Close()

		errc := make(chan error, 2)

		go func() {
			_, err := io.Copy(clientConn.UnderlyingConn(), targetConn.UnderlyingConn())
			errc <- err
		}()
		go func() {
			_, err := io.Copy(targetConn.UnderlyingConn(), clientConn.UnderlyingConn())
			errc <- err
		}()

		if err := <-errc; err != nil {
			log.Printf("WebSocket proxy error: %v", err)
		}
	})

	term.POST("/terminals/:pid/size", func(ctx *gin.Context) {
		port := ctx.MustGet("port").(uint16)
		pid := ctx.Param("pid")

		var updateTermSizeRequest UpdateTermSizeRequest

		if err := ctx.ShouldBind(&updateTermSizeRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := SendUpdateTermSizeRequest(o.HostIP, port, pid, updateTermSizeRequest)

		if err != nil {
			ctx.JSON(400, ErrorResponse {
				Error: err.Error(),
			})
			return
		}

		payload, err := io.ReadAll(response.Body)

		if err != nil {
			ctx.JSON(400, ErrorResponse {
				Error: err.Error(),
			})
			return
		}

		ctx.Data(response.StatusCode, response.Header.Get("Content-Type"), payload)
	})

	engine.Run(":6969")
}

