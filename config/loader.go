package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

//Load pulls configuration from environment and put it into Config struct.
func Load() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("ddns")
	viper.BindPFlags(setUpFlag())

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

func setUpFlag() *pflag.FlagSet {
	pflag.String("provider", "", "External DNS provider which should be used")
	pflag.String("provider.cf.apiKey", "", "Cloudflare user API key")
	pflag.String("provider.cf.email", "", "Cloudflare user email")
	pflag.String("provider.cf.zoneId", "", "Cloudflare domain zone ID")

	pflag.String("source", "kubernetes", "Domain records source type")

	pflag.Parse()
	return pflag.CommandLine
}
