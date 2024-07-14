package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"replme/database"
	"replme/model"
	"replme/service"
	"replme/types"
	"replme/util"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type DevenvController struct {
	Docker             *service.DockerService
	Upgrader           websocket.Upgrader
	DevenvFilesPath    string
	DevenvFilesPathTmp string
}

func NewDevenvController(docker *service.DockerService, devenvFilesPath string, devenvFilesTmpPath string) DevenvController {

	err := util.MakeDirIfNotExists(devenvFilesPath)

	if err != nil {
		log.Fatal(err)
	}

	err = util.MakeDirIfNotExists(devenvFilesTmpPath)
	if err != nil {
		log.Fatal(err)
	}

	return DevenvController{
		Docker: docker,
		Upgrader: websocket.Upgrader{
			// ReadBufferSize:  1024,
			// WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		DevenvFilesPath:    devenvFilesPath,
		DevenvFilesPathTmp: devenvFilesTmpPath,
	}
}

func (devenv *DevenvController) GetAll(ctx *gin.Context) {
	_user, _ := ctx.Get("current_user")
	user := _user.(model.User)

	var devenvs []model.Devenv
	err := database.DB.Model(user).Association("Devenvs").Find(&devenvs)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, devenvs)
}

func (devenv *DevenvController) GetOne(ctx *gin.Context) {
	_user, _ := ctx.Get("current_user")
	user := _user.(model.User)

	id := ctx.Param("uuid")

	var devenvs []model.Devenv
	err := database.DB.Model(&user).Where("id = ?", id).Association("Devenvs").Find(&devenvs)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(devenvs) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, &gin.H{
			"error": "Not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, devenvs[0])
}

func (devenv *DevenvController) Create(ctx *gin.Context) {
	_user, _ := ctx.Get("current_user")
	user := _user.(model.User)

	var createDevenvRequest types.CreateDevenvRequest
	if err := ctx.ShouldBind(&createDevenvRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentDevenv := model.Devenv{
		Public:   false,
		Name:     createDevenvRequest.Name,
		BuildCmd: createDevenvRequest.BuildCmd,
		RunCmd:   createDevenvRequest.RunCmd,
	}

	err := database.DB.Model(&user).Association("Devenvs").Append(&currentDevenv)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
		return
	}

	dir := filepath.Join(devenv.DevenvFilesPath, currentDevenv.ID)
	err = util.SetFileContent(dir, "main.c", "#include<stdio.h>\n\nint main() {\n  printf(\"Hello, REPL!\\n\");\n  return 0;\n}\n")

	ctx.JSON(http.StatusOK, types.CreateDevenvResponse{
		DevenvUuid: currentDevenv.ID,
	})
}

func (devenv *DevenvController) Patch(ctx *gin.Context) {
	_user, _ := ctx.Get("current_user")
	user := _user.(model.User)

	var patchDevenvRequest types.PatchDevenvRequest
	if err := ctx.ShouldBind(&patchDevenvRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("uuid")

	var devenvs []model.Devenv
	err := database.DB.Model(&user).Where("id = ?", id).Association("Devenvs").Find(&devenvs)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(devenvs) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, &gin.H{
			"error": "Not found",
		})
		return
	}

	if patchDevenvRequest.Name != "" {
		devenvs[0].Name = patchDevenvRequest.Name
	}

	if patchDevenvRequest.BuildCmd != "" {
		devenvs[0].BuildCmd = patchDevenvRequest.BuildCmd
	}

	if patchDevenvRequest.RunCmd != "" {
		devenvs[0].RunCmd = patchDevenvRequest.RunCmd
	}

	database.DB.Save(&devenvs[0])

	ctx.Status(http.StatusOK)
}

func (devenv *DevenvController) GetFiles(ctx *gin.Context) {
	_devenv, _ := ctx.Get("current_devenv")
	currentDevenv := _devenv.(model.Devenv)

	dir := filepath.Join(devenv.DevenvFilesPath, currentDevenv.ID)
	err := util.MakeDirIfNotExists(dir)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
		return
	}

	files := []string{}
	for _, e := range entries {
		files = append(files, e.Name())
	}

	ctx.JSON(http.StatusOK, files)
}

func (devenv *DevenvController) CreateFile(ctx *gin.Context) {
	_devenv, _ := ctx.Get("current_devenv")
	currentDevenv := _devenv.(model.Devenv)

	var createFileRequest types.CreateFileRequest
	if err := ctx.ShouldBind(&createFileRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !util.IsValidFilename(createFileRequest.Name) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Illegal filename"})
		return
	}

	dir := filepath.Join(devenv.DevenvFilesPath, currentDevenv.ID)
	err := util.TouchIfNotExists(dir, createFileRequest.Name)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (devenv *DevenvController) GetFileContent(ctx *gin.Context) {
	_uuid, _ := ctx.Get("uuid")
	uuid := _uuid.(string)
	name := ctx.Param("name")
	path := filepath.Join(devenv.DevenvFilesPath, uuid, name)

	if !strings.HasPrefix(path, devenv.DevenvFilesPath) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &gin.H{
			"error": "Invalid uuid",
		})
		return
	}

	content, err := util.GetFileContent(path)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.String(http.StatusOK, content)
}

func (devenv *DevenvController) SetFileContent(ctx *gin.Context) {
	_devenv, _ := ctx.Get("current_devenv")
	currentDevenv := _devenv.(model.Devenv)

	if ctx.Request.ContentLength > 1024 {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Max filesize is 1KB"})
		return
	}

	name := ctx.Param("name")

	if !util.IsValidFilename(name) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Illegal filename"})
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dir := filepath.Join(devenv.DevenvFilesPath, currentDevenv.ID)
	err = util.SetFileContent(dir, name, string(body))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (devenv *DevenvController) DeleteFile(ctx *gin.Context) {
	_devenv, _ := ctx.Get("current_devenv")
	currentDevenv := _devenv.(model.Devenv)

	name := ctx.Param("name")

	path := filepath.Join(devenv.DevenvFilesPath, currentDevenv.ID, name)
	err := util.DeleteFile(path)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (devenv *DevenvController) Exec(ctx *gin.Context) {
	_devenv, _ := ctx.Get("current_devenv")
	currentDevenv := _devenv.(model.Devenv)

	src := filepath.Join(devenv.DevenvFilesPath, currentDevenv.ID)
	tmpUuid := guuid.New().String()

	target := filepath.Join(devenv.DevenvFilesPathTmp, tmpUuid)
	mount := filepath.Join("/tmp", tmpUuid)
	util.SLogger.Debugf("Copying %s -> %s", src, target)

	err := util.CopyRecurse(src, target, 0777)
	if err != nil {
		util.SLogger.Warnf("Copying devenv container failed, %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, types.ErrorResponse{
			Error: "Could not copy src dir",
		})
		return
	}

	defer func() {
		util.SLogger.Debugf("Deleting dir %s", target)
		err := util.DeleteDir(target)
		if err != nil {
			util.SLogger.Warnf("Failed to delete dir %s, %s", target, err.Error())
		}
	}()

	start := time.Now()
	id, _, port, err := devenv.Docker.EnsureDevenvContainerStarted(target, mount)
	util.SLogger.Debugf("[%-25s] Starting exec container took %v", fmt.Sprintf("UID:%s..", currentDevenv.ID[:5]), time.Since(start))

	if err != nil {
		util.SLogger.Warnf("Creating devenv container failed, %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, types.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	defer func() {
		util.SLogger.Debugf("Removing container with id %s", (*id)[:5])
		err := devenv.Docker.RemoveContainerById(*id)
		if err != nil {
			util.SLogger.Warnf("Removing container with id %s failed, %s", (*id)[:5], err.Error())
		}
	}()

	p := service.Proxy(devenv.Docker.HostIP, *port, devenv.Docker.ApiKey)

	username, _ := util.RandomString(50)
	password, _ := util.RandomString(50)

	response, requestError := p.SendRegisterRequest(
		types.RegisterRequest{
			Username: username,
			Password: password,
		},
		&types.RequestOptions{
			Retries: 10,
		},
	)

	if requestError != nil {
		ctx.Data(requestError.Code, requestError.ContentType, requestError.Data)
		ctx.Abort()
		return
	}

	command := fmt.Sprintf("%s && %s", currentDevenv.BuildCmd, currentDevenv.RunCmd)

	cookie := response.Cookies()[0]

	clientConn, err := devenv.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer clientConn.Close()
	err = p.CreateExecWebsocketPipe(clientConn, *cookie, mount, command)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
}
