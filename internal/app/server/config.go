package server

type Config struct {
	BindAddrHTTPS string `toml:"bind_addr_https"`
	BindAddrHTTP  string `toml:"bind_addr_http"`
	LogLevel      string `toml:"log_level"`
	DataBaseUrl   string `toml:"database_url"`
	RedisUrl      string `toml:"redis_url"`
	Domen         string `toml:"domen"`
}

func NewConfig() *Config {
	return &Config{}
}
