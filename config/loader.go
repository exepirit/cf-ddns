package config

import (
	"strings"

	"github.com/spf13/viper"
)

//Load pulls configuration from environment and put it into Config struct.
func Load() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("ddns")

	cfg := &Config{
		Provider: viper.GetString("provider"),
		CfApiKey: viper.GetString("provider.cf.apiKey"),
		CfEmail:  viper.GetString("provider.cf.email"),
		CfZoneID: viper.GetString("provider.cf.zoneId"),
		Source:   viper.GetString("source"),
		FilePath: viper.GetString("source.file"),
	}
	return cfg, nil
}
