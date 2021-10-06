package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"patreon/internal/app"
	main_server "patreon/internal/app/server"
	"time"

	//_ "patreon/docs"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var (
	configPath string
)

func newLogger(config *app.Config) (log *logrus.Logger, closeResource func() error) {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	currentTime := time.Now()
	formatted := config.LogAddr + fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second()) + ".out"

	f, err := os.OpenFile(formatted, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}

	logger.SetOutput(f)
	logger.SetLevel(level)

	return logger, f.Close
}

func newPostgresConnection(config *app.Config) (db *sql.DB, closeResource func() error) {
	db, err := sql.Open("postgres", config.DataBaseUrl)
	if err != nil {
		logrus.Fatal(err)
	}

	return db, db.Close
}

func newRedisPool(config *app.Config) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(config.RedisUrl)
		},
	}
}

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

	config := app.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		logrus.Fatal(err)
	}

	logger, closeResource := newLogger(config)

	defer func(closer func() error) {
		err := closer()
		if err != nil {
			logger.Fatal(err)
		}
	}(closeResource)

	db, closeResource := newPostgresConnection(config)

	defer func(closer func() error) {
		err := closer()
		if err != nil {
			logger.Fatal(err)
		}
	}(closeResource)

	server := main_server.New(config,
		app.ExpectedConnections{RedisPool: newRedisPool(config), SqlConnection: db},
		logger)

	if err := server.Start(config); err != nil {
		logger.Fatal(err)
	}
}
