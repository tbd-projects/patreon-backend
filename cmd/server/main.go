package main

import (
	"flag"
	"os"
	server "patreon/internal/app/server"

	_ "patreon/docs"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	// gin-swagger middleware
	// swagger embed files
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
}

// @title Patreon
// @version 1.0
// @description Server for Patreon application.

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}

func main() {
	flag.Parse()

	config := server.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if err := server.Start(config); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
