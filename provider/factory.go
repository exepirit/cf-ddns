package provider

import (
	"log"

	"github.com/exepirit/cf-ddns/provider/cloudflare"
	"github.com/pkg/errors"
)

func NewFromConfig(cfg *Config) (Provider, error) {
	switch cfg.ProviderType {
	case "cloudflare":
		return cloudflare.NewProvider(cfg.CloudflareZoneID, cfg.CloudflareApiKey, cfg.CloudflareEmail)
	case "stub":
		return &Stub{
			Logger: log.Default(),
		}, nil
	default:
		return nil, errors.Errorf("unknown provider \"%s\"", cfg.ProviderType)
	}
}
