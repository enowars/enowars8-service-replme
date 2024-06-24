package main

import (
	"image-go/server"
	"os"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	server.Init(apiKey)
}
