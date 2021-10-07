package app

type RepositoryConnections struct {
	DataBaseUrl string `toml:"database_url"`
	RedisUrl    string `toml:"redis_url"`
}

type Config struct {
	BindAddr         string                `toml:"bind_addr"`
	LogLevel         string                `toml:"log_level"`
	LogAddr          string                `toml:"log_path"`
	ServerRepository RepositoryConnections `toml:"server"`
	LocalRepository  RepositoryConnections `toml:"local"`
}

func NewConfig() *Config {
	return &Config{}
}
