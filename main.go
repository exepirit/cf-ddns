package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/exepirit/cf-ddns/config"
	"github.com/exepirit/cf-ddns/control"
	"github.com/exepirit/cf-ddns/provider"
	"github.com/exepirit/cf-ddns/source"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("setting up cf-ddns...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err)
	}

	provider, err := provider.NewFromConfig(&provider.Config{
		ProviderType:     cfg.Provider,
		CloudflareZoneID: cfg.CfZoneID,
		CloudflareApiKey: cfg.CfApiKey,
		CloudflareEmail:  cfg.CfEmail,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("create provider")
	}

	src, err := source.NewFromConfig(&source.Config{
		SourceType: cfg.Source,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("create domains source")
	}

	controller := control.Controller{
		Source:     src,
		Provider:   provider,
		TimePeriod: time.Minute,
	}

	err = controller.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
