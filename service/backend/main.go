package main

import (
	"flag"
	"os"
	"strings"

	"replme/server"
	"replme/service"
	"replme/util"
)

func main() {
	var imagePath string
	var imageTag string
	var postgresUrl string
	var postgresUser string
	var postgresSecretPath string
	var apiKeyPath string
	var devenvsPath string
	var devenvsTmpPath string
	var containerLogsPath string

	flag.StringVar(&imagePath, "i", "", "Image dir (required)")
	flag.StringVar(&imagePath, "n", "", "Image tag (required), env: REPL_IMG_TAG")
	flag.StringVar(&postgresUrl, "p", "", "Postgres connection url (required), env: REPL_POSTGRES_URL")
	flag.StringVar(&postgresUser, "u", "", "Postgres user (required), env: REPL_POSTGRES_USER")
	flag.StringVar(&postgresSecretPath, "s", "", "Postgres secret file (required), env: REPL_POSTGRES_SECRET")
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

	if postgresUrl == "" {
		postgresUrlEnv := os.Getenv("REPL_POSTGRES_URL")
		if postgresUrlEnv == "" {
			flag.Usage()
			os.Exit(1)
		}
		postgresUrl = postgresUrlEnv
	}

	if postgresUser == "" {
		postgresUserEnv := os.Getenv("REPL_POSTGRES_USER")
		if postgresUserEnv == "" {
			flag.Usage()
			os.Exit(1)
		}
		postgresUser = postgresUserEnv
	}

	if postgresSecretPath == "" {
		postgresSecretPathEnv := os.Getenv("REPL_POSTGRES_SECRET")
		if postgresSecretPathEnv == "" {
			flag.Usage()
			os.Exit(1)
		}
		postgresSecretPath = postgresSecretPathEnv
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

	pgSecret := util.ReadPostgresSecret(postgresSecretPath)
	pgUrl := strings.ReplaceAll(postgresUrl, "{user}", postgresUser)
	pgUrl = strings.ReplaceAll(pgUrl, "{secret}", pgSecret)

	server.Init(&docker, pgUrl, containerLogsPath, devenvsPath, devenvsTmpPath)
}
