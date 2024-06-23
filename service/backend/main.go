package main

import (
	"flag"
	"os"

	"replme/server"
	"replme/service"
	"replme/util"
)

func main() {
	var imagePath string
	var apiKeyPath string

	flag.StringVar(&imagePath, "i", "", "Image directory (required)")
	flag.StringVar(&apiKeyPath, "k", "", "Apikey file (required)")

	flag.Parse()

	if imagePath == "" || apiKeyPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	apiKey := util.ApiKey(apiKeyPath)
	imageTag := "ptwhy"

	docker := service.Docker(apiKey)
	docker.BuildImage(imagePath, imageTag)

	server.Init(&docker)
}
