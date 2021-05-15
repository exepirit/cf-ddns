package plan

import "github.com/exepirit/cf-ddns/domain"

type Changes struct {
	Create []*domain.Endpoint
	Update []*domain.Endpoint
	Delete []*domain.Endpoint
}
