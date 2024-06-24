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
	var devenvFiles string
	var devenvFilesTmp string

	flag.StringVar(&imagePath, "i", "", "Image directory (required)")
	flag.StringVar(&apiKeyPath, "k", "", "Apikey file (required)")
	flag.StringVar(&devenvFiles, "f", "", "Devenv files (required)")
	flag.StringVar(&devenvFilesTmp, "t", "", "Devenv files tmp (required)")

	flag.Parse()

	if imagePath == "" || apiKeyPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	apiKey := util.ApiKey(apiKeyPath)
	imageTag := "ptwhy"

	docker := service.Docker(apiKey)
	docker.BuildImage(imagePath, imageTag)

	server.Init(&docker, devenvFiles, devenvFilesTmp)
}
