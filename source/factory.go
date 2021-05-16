package source

import (
	"github.com/exepirit/cf-ddns/source/kubernetes"
	"github.com/exepirit/cf-ddns/source/static"
	"github.com/pkg/errors"
)

func NewFromConfig(cfg *Config) (Source, error) {
	switch cfg.SourceType {
	case "static":
		return static.NewSourceFromEnv(), nil
	case "kubernetes":
		return kubernetes.New()
	default:
		return nil, errors.Errorf("unknown domains source \"%s\"", cfg.SourceType)
	}
}
