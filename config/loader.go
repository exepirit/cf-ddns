package config

import "github.com/spf13/viper"

func Load() (*Config, error) {
	viper.SetConfigType("env")
	viper.SetEnvPrefix("ddns")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	return cfg, err
}
