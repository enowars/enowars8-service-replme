package middleware

import (
	"replme/service"
	"replme/types"
	"replme/util"

	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(docker *service.DockerService) gin.HandlerFunc {

	crc := util.CRC()

	return func(ctx *gin.Context) {

		authHeader := ctx.Request.Header["Authorization"]

		if len(authHeader) == 0 {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header missing"},
			)
			ctx.Abort()
			return
		}

		authValues := strings.Split(authHeader[0], " ")

		if len(authValues) < 2 {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Authorization header invalid"},
			)
			ctx.Abort()
			return
		}

		credsBase64 := authValues[1]
		credsString, err := base64.StdEncoding.DecodeString(credsBase64)

		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Authorization value not decodeable"},
			)
			ctx.Abort()
			return
		}

		credsSlice := strings.Split(string(credsString), ":")

		if len(credsSlice) < 2 {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Authorization value invalid"},
			)
			ctx.Abort()
			return
		}

		username := credsSlice[0]
		password := credsSlice[1]

		hash := crc.Calculate(util.DecodeSpecialChars([]byte(username)))
		name := fmt.Sprintf("%x", hash)

		container, _, running := docker.GetContainer(name)

		if container == nil {
			ctx.JSON(
				http.StatusPreconditionFailed,
				gin.H{"error": "Container not found"},
			)
			ctx.Abort()
			return
		}

		if !running {
			ctx.JSON(
				http.StatusPreconditionFailed,
				gin.H{"error": "Container not running"},
			)
			ctx.Abort()
			return
		}

		_, port, err := docker.GetContainerAddress(container.ID)

		if err != nil {
			ctx.JSON(
				http.StatusPreconditionFailed,
				gin.H{"error": "Could not get container address"},
			)
			ctx.Abort()
			return
		}

		p := service.Proxy(docker.HostIP, *port)
		response, err := p.SendUserLoginRequest(types.LoginRequest{
			Username: username,
			Password: password,
		})

		if err != nil {
			ctx.JSON(400, types.ErrorResponse{
				Error: err.Error(),
			})
			ctx.Abort()
			return
		}

		if response.StatusCode >= 400 {
			payload, err := io.ReadAll(response.Body)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
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
	}

}
