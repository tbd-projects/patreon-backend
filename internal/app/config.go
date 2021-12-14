package app

import "patreon/internal"

const (
	LoadFileUrl  = "media/"
	DefaultImage = ""
)

type Payments struct {
	AccountNumber string `toml:"account_number"`
}

type Microservice struct {
	SessionServerUrl string `toml:"session_url"`
	FilesUrl         string `toml:"files_url"`
}

type RepositoryConnections struct {
	DataBaseUrl     string `toml:"database_url"`
	SessionRedisUrl string `toml:"session-redis_url"`
	AccessRedisUrl  string `toml:"access-redis_url"`
}

type Config struct {
	internal.Config
	MediaDir         string                `toml:"media_dir"`
	Microservices    Microservice          `toml:"microservice"`
	ServerRepository RepositoryConnections `toml:"server"`
	LocalRepository  RepositoryConnections `toml:"local"`
	Cors             internal.CorsConfig   `toml:"cors"`
	PaymentsInfo     Payments              `toml:"payments"`
}

func NewConfig() *Config {
	return &Config{}
}
