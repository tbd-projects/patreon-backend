package app

const LoadFileUrl = "media/"

type RepositoryConnections struct {
	DataBaseUrl     string `toml:"database_url"`
	SessionRedisUrl string `toml:"session-redis_url"`
	AccessRedisUrl  string `toml:"access-redis_url"`
}

type CorsConfig struct {
	Urls    []string `toml:"urls"`
	Headers []string `toml:"headers"`
	Methods []string `toml:"methods"`
}

type Config struct {
	MediaDir         string                `toml:"media_dir"`
	BindAddr         string                `toml:"bind_addr"`
	LogLevel         string                `toml:"log_level"`
	LogAddr          string                `toml:"log_path"`
	ServerRepository RepositoryConnections `toml:"server"`
	LocalRepository  RepositoryConnections `toml:"local"`
	Cors             CorsConfig            `toml:"cors"`
}

func NewConfig() *Config {
	return &Config{}
}
