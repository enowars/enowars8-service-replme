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
	var imageTag string
	var dbPath string
	var apiKeyPath string
	var devenvsPath string
	var devenvsTmpPath string
	var containerLogsPath string

	flag.StringVar(&imagePath, "i", "", "Image dir (required)")
	flag.StringVar(&imagePath, "n", "", "Image tag (required), env: REPL_IMG_TAG")
	flag.StringVar(&dbPath, "d", "", "Database file (required), env: REPL_SQLITE")
	flag.StringVar(&apiKeyPath, "k", "", "Apikey file (required), env: REPL_API_KEY")
	flag.StringVar(&devenvsPath, "f", "", "Devenv files dir (required), env: REPL_DEVENVS")
	flag.StringVar(&devenvsTmpPath, "t", "", "Tmp devenv files dir (required), env: REPL_DEVENVS_TMP")
	flag.StringVar(&containerLogsPath, "l", "", "Container logs dir (required), env: REPL_CONTAINER_LOGS")

	flag.Parse()

	if imagePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if imageTag == "" {
		imageTagEnv := os.Getenv("REPL_IMG_TAG")
		if imageTagEnv == "" {
			flag.Usage()
			os.Exit(1)
		}
		imageTag = imageTagEnv
	}

	if dbPath == "" {
		dbPathEnv := os.Getenv("REPL_SQLITE")
		if dbPathEnv == "" {
			flag.Usage()
			os.Exit(1)
		}
		dbPath = dbPathEnv
	}

	if apiKeyPath == "" {
		apiKeyPathEnv := os.Getenv("REPL_API_KEY")
		if apiKeyPathEnv == "" {
			flag.Usage()
			os.Exit(1)
		}
		apiKeyPath = apiKeyPathEnv
	}

	if devenvsPath == "" {
		devenvsPathEnv := os.Getenv("REPL_DEVENVS")
		if devenvsPathEnv == "" {
			flag.Usage()
			os.Exit(1)
		}
		devenvsPath = devenvsPathEnv
	}

	if devenvsTmpPath == "" {
		devenvsTmpPathEnv := os.Getenv("REPL_DEVENVS_TMP")
		if devenvsTmpPathEnv == "" {
			flag.Usage()
			os.Exit(1)
		}
		devenvsTmpPath = devenvsTmpPathEnv
	}

	// REPL_CONTAINER_LOGS

	if containerLogsPath == "" {
		containerLogsPathTmp := os.Getenv("REPL_CONTAINER_LOGS")
		if containerLogsPathTmp == "" {
			flag.Usage()
			os.Exit(1)
		}
		containerLogsPath = containerLogsPathTmp
	}

	apiKey := util.ApiKey(apiKeyPath)

	docker := service.Docker(apiKey, imagePath, imageTag, containerLogsPath)
	docker.BuildImage()

	server.Init(&docker, dbPath, containerLogsPath, devenvsPath, devenvsTmpPath)
}
