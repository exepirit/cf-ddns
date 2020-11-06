package main

import (
	"errors"
	"github.com/cloudflare/cloudflare-go"
	"github.com/exepirit/cf-ddns/internal/app/renewal"
	"github.com/exepirit/cf-ddns/pkg/ddns"
	"github.com/exepirit/cf-ddns/pkg/echoip"
	"github.com/exepirit/cf-ddns/pkg/lookup"
	"os"
)

func makeWorker() (*renewal.Worker, error) {
	ipResolver := &echoip.IfconfigResolver{}
	dnsResolver := &lookup.Resolver{
		ServerAddr: "8.8.8.8",
	}

	cfToken, ok := os.LookupEnv("CF_TOKEN")
	if !ok {
		return nil, errors.New("required environment variable CF_TOKEN")
	}

	cfAPI, err := cloudflare.NewWithAPIToken(cfToken)
	if err != nil {
		return nil, err
	}

	zone := os.Getenv("DNS_ZONE")
	ddnsUpdater, err := ddns.NewUpdater(cfAPI, zone)
	if err != nil {
		return nil, err
	}

	worker := renewal.NewWorker(ipResolver, dnsResolver, ddnsUpdater)
	return worker, nil
}
