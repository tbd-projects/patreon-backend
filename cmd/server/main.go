package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"patreon/internal/app"
	main_server "patreon/internal/app/server"
	"time"

	"github.com/gomodule/redigo/redis"

	_ "patreon/docs"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var (
	configPath          string
	useServerRepository bool
	runHttps            bool
)

func newLogger(config *app.Config) (log *logrus.Logger, closeResource func() error) {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	currentTime := time.Now().In(time.UTC)
	formatted := config.LogAddr + fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second()) + ".log"

	f, err := os.OpenFile(formatted, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}

	logger.SetOutput(f)
	logger.Writer()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger, f.Close
}

func newPostgresConnection(config *app.RepositoryConnections) (db *sqlx.DB, closeResource func() error) {
	db, err := sqlx.Open("postgres", config.DataBaseUrl)
	if err != nil {
		logrus.Fatal(err)
	}

	return db, db.Close
}

func newRedisPool(redisUrl string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisUrl)
		},
	}
}

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
	flag.BoolVar(&useServerRepository, "server-run", false, "true if it server run, false if it local run")
	flag.BoolVar(&runHttps, "run-https", false, "run https serve with certificates")

}

// @title Patreon
// @version 1.0
// @description Server for Patreon application.

// @host api.pyaterochka-team.site
// @BasePath /api/v1

// @x-extension-openapi {"example": "value on a json format"}

func main() {
	flag.Parse()
	logrus.Info(os.Args[:])

	config := app.NewConfig()
	config.IsHTTPSServer = runHttps

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		logrus.Fatal(err)
	}

	logger, closeResource := newLogger(config)

	defer func(closer func() error, log *logrus.Logger) {
		err := closer()
		if err != nil {
			log.Fatal(err)
		}
	}(closeResource, logger)

	repositoryConfig := &config.LocalRepository
	if useServerRepository {
		repositoryConfig = &config.ServerRepository
	}

	db, closeResource := newPostgresConnection(repositoryConfig)

	defer func(closer func() error, log *logrus.Logger) {
		err := closer()
		if err != nil {
			log.Fatal(err)
		}
	}(closeResource, logger)

	server := main_server.New(config,
		app.ExpectedConnections{
			SessionRedisPool: newRedisPool(repositoryConfig.SessionRedisUrl),
			AccessRedisPool:  newRedisPool(repositoryConfig.AccessRedisUrl),
			SqlConnection:    db,
			PathFiles:        config.MediaDir,
		},
		logger,
	)

	if err = server.Start(config); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Server was stopped")
}
