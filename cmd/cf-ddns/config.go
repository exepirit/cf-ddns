package main

import "github.com/spf13/viper"

type config struct {
	CloudflareToken string
	CloudflareZone  string
	BindAddress     string
	DnsServer       string
}

func loadConfig() (*config, error) {
	loader := viper.New()
	loader.SetConfigType("json")
	loader.AddConfigPath("config/")
	loader.AddConfigPath("$HOME/.config/cf-ddns")

	tryLoadCfg := func(name string) (*config, error) {
		loader.SetConfigName(name)
		if err := loader.ReadInConfig(); err != nil {
			return nil, err
		}

		cfg := &config{}
		if err := loader.Unmarshal(cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	if cfg, err := tryLoadCfg("config.local"); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			return tryLoadCfg("config")
		default:
			return nil, err
		}
	} else {
		return cfg, nil
	}
}
