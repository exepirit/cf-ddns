package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/exepirit/cf-ddns/config"
	"github.com/exepirit/cf-ddns/control"
	"github.com/exepirit/cf-ddns/provider"
	"github.com/exepirit/cf-ddns/source"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	provider, err := provider.NewFromConfig(&provider.Config{
		ProviderType:     cfg.Provider,
		CloudflareZoneID: cfg.CfZoneID,
		CloudflareApiKey: cfg.CfApiKey,
		CloudflareEmail:  cfg.CfEmail,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	src, err := source.NewFromConfig(&source.Config{
		SourceType: cfg.Source,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	controller := control.Controller{
		Source:     src,
		Provider:   provider,
		TimePeriod: time.Minute,
	}

	err = controller.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
