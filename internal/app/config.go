package app

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	LogAddr     string `toml:"log_path"`
	DataBaseUrl string `toml:"database_url"`
	RedisUrl    string `toml:"redis_url"`
}

func NewConfig() *Config {
	return &Config{}
}
