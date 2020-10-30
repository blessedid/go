package store

// Config ...
type Config struct {
	dbURL string `toml:"db_url"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}
