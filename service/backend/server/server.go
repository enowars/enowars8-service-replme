package server

import "replme/service"

func Init(docker *service.DockerService) {
	engine := NewRouter(docker)
	engine.Run(":6969")
}
