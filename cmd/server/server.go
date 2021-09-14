package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	server2 "patreon/internal/app/server"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
}
func main() {
	flag.Parse()

	config := server2.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	s := server2.New(config)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
