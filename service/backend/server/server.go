package server

import (
	"replme/service"
	"replme/util"
)

func Init(docker *service.DockerService, dbPath string, devenvFilesPath string, devenvFilesTmpPath string) {
	engine := NewRouter(docker, dbPath, devenvFilesPath, devenvFilesTmpPath)
	util.SLogger.Infof("Server is running on port 6969")
	engine.Run(":6969")
}
