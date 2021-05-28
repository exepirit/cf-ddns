package config

import (
	"github.com/kelseyhightower/envconfig"
)

func Load() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("ddns", cfg)
	return cfg, err
}
