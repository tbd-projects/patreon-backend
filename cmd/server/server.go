package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"os"
	server "patreon/internal/app/server"
	"patreon/internal/app/store"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
}
func main() {
	flag.Parse()

	config := server.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	logger := log.New()
	logger.SetLevel(level)

	handler := server.NewMainHandler()
	handler.SetLogger(logger)

	router := server.NewRouter()
	router.Configure()
	handler.SetRouter(router)

	st := store.New(config.Store)
	err = st.Open()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	handler.SetStore(st)
	s := server.New(config, handler)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
