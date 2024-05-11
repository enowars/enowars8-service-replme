package server

import "cafedodo/service"

func Init(docker *service.DockerService, dist string) {
	engine := NewRouter(docker, dist)
	engine.Run(":6969")
}
