package main

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/exepirit/cf-ddns/internal/app/renewal"
	"github.com/exepirit/cf-ddns/pkg/ddns"
	"github.com/exepirit/cf-ddns/pkg/echoip"
	"github.com/exepirit/cf-ddns/pkg/lookup"
)

func makeWorker(cfg *config) (*renewal.Worker, error) {
	ipResolver := &echoip.IfconfigResolver{}
	dnsResolver := &lookup.Resolver{
		ServerAddr: cfg.DnsServer,
	}

	cfAPI, err := cloudflare.NewWithAPIToken(cfg.CloudflareToken)
	if err != nil {
		return nil, err
	}

	ddnsUpdater, err := ddns.NewUpdater(cfAPI, cfg.CloudflareZone)
	if err != nil {
		return nil, err
	}

	worker := renewal.NewWorker(ipResolver, dnsResolver, ddnsUpdater)
	return worker, nil
}
