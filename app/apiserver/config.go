package apiserver

// Config ...
type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLeval string `toml:"log_leval"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLeval: "debug",
	}
}
