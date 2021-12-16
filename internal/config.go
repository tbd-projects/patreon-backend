package internal

type Config struct {
	LogLevel      string `toml:"log_level"`
	LogAddr       string `toml:"log_path"`
	Domen         string `toml:"domen"`
	IsHTTPSServer bool
	BindHttpsAddr string `toml:"bind_addr_https"`
	BindHttpAddr  string `toml:"bind_addr_http"`
}

type CorsConfig struct {
	Urls    []string `toml:"urls"`
	Headers []string `toml:"headers"`
	Methods []string `toml:"methods"`
}
