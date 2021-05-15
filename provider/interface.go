package provider

import (
	"github.com/exepirit/cf-ddns/domain"
	"github.com/exepirit/cf-ddns/plan"
)

type Provider interface {
	ApplyPlan(plan plan.Plan) error
	CurrentEndpoints() ([]*domain.Endpoint, error)
}
