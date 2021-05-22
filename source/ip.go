package source

import (
	"context"

	"github.com/exepirit/cf-ddns/domain"
	"github.com/exepirit/cf-ddns/pkg/echoip"
	"github.com/pkg/errors"
)

type resolvedSource struct {
	echoip.Resolver
	inner Source
}

// WithExternalIp replaces targets, returned by Source, with external ip or any resolved IP.
func WithResolvedTarget(source Source, resolver echoip.Resolver) Source {
	return &resolvedSource{
		Resolver: resolver,
		inner:    source,
	}
}

func (source *resolvedSource) GetEndpoints() ([]*domain.Endpoint, error) {
	endpoints, err := source.GetEndpoints()
	if len(endpoints) == 0 {
		return endpoints, err
	}

	externalIp, err := source.GetIP(context.Background())
	if err != nil {
		return endpoints, errors.WithMessage(err, "resolve external ip")
	}

	for _, endpoint := range endpoints {
		endpoint.Target = []string{externalIp.String()}
	}
	return endpoints, err
}
