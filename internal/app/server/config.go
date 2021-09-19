package server

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DataBaseUrl string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
