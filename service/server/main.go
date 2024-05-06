package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"time"

	"cafedodo/client"
	"cafedodo/orchestrator"
	"cafedodo/renderer"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type SuccessResponse struct {
	Success string `json:"success"`
	Port uint16 `json:"port"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendCreateUserRequest(port uint16, request LoginRequest) (*http.Response, error) {
	payload, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	return http.Post(
		fmt.Sprintf("http://127.0.0.1:%d/user/create", port),
		"application/json",
		bytes.NewBuffer(payload),
	)
}

func main() {
	imageTag := "ptwhy"

	o := orchestrator.Instance()
	o.BuildImage("../ptwhy/", imageTag)

	engine := gin.Default()

	engine.Static("/static", "./static")

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

		port, err := o.EnsureContainerStarted(
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

		for range 5 {
			response, err = SendCreateUserRequest(*port, loginRequest)
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

	engine.Run(":8080")
}

