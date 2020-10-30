package apiserver

import (
	"vk-go/app/store"
)

// Config ...
type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLeval string `toml:"log_leval"`
	Store    *store.Config
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLeval: "debug",
		Store:    store.NewConfig(),
	}
}
