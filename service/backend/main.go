package main

import (
	"flag"
	"os"

	"replme/server"
	"replme/service"
)

func main() {

	var distPath string
	var imagePath string

	flag.StringVar(&distPath, "d", "", "Dist directory (required)")
	flag.StringVar(&imagePath, "i", "", "Image directory (required)")

	flag.Parse()

	if distPath == "" || imagePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	imageTag := "ptwhy"

	docker := service.Docker()
	docker.BuildImage(imagePath, imageTag)

	server.Init(&docker, distPath)

}
