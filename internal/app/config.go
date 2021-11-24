package app

const (
	LoadFileUrl  = "media/"
	DefaultImage = ""
)

type Microservice struct {
	SessionServerUrl string `toml:"session_url"`
	FilesUrl         string `toml:"files_url"`
}

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
	Microservices    Microservice          `toml:"microservice"`
	ServerRepository RepositoryConnections `toml:"server"`
	LocalRepository  RepositoryConnections `toml:"local"`
	Cors             CorsConfig            `toml:"cors"`
	IsHTTPSServer    bool
	Domen            string `toml:"domen"`
	BindHttpsAddr    string `toml:"bind_addr_https"`
	BindHttpAddr     string `toml:"bind_addr_http"`
}

func NewConfig() *Config {
	return &Config{}
}
