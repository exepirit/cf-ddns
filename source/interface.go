package source

import "github.com/exepirit/cf-ddns/domain"

// Source provide domains that must bound with current IP.
type Source interface {
	GetEndpoints() ([]*domain.Endpoint, error)
}
