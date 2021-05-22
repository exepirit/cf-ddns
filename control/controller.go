package control

import (
	"time"

	"github.com/exepirit/cf-ddns/domain"
	"github.com/exepirit/cf-ddns/plan"
	"github.com/exepirit/cf-ddns/provider"
	"github.com/exepirit/cf-ddns/source"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	Source     source.Source
	Provider   provider.Provider
	TimePeriod time.Duration
}

func (ctrl *Controller) RunOnce() error {
	desiredState, err := ctrl.Source.GetEndpoints()
	if err != nil {
		return errors.WithMessage(err, "get desired domains")
	}

	currentState, err := ctrl.Provider.CurrentEndpoints()
	if err != nil {
		return errors.WithMessage(err, "get bonded domains")
	}
	currentState = ctrl.filterUnmanagedDomains(currentState, desiredState)

	currentPlan := plan.Plan{
		Current: currentState,
		Desired: desiredState,
	}
	currentPlan.Eval()

	err = ctrl.Provider.ApplyPlan(currentPlan)
	return err
}

func (ctrl *Controller) Run() {
	log.Info().Msg("DNS zone controller started")
	ticker := time.NewTicker(ctrl.TimePeriod)
	defer ticker.Stop()

	var err error
	for {
		err = ctrl.RunOnce()
		if err != nil {
			log.Error().Err(err).Msg("one or more error occured")
		}
		<-ticker.C
	}
}

func (*Controller) filterUnmanagedDomains(inp, managed []*domain.Endpoint) []*domain.Endpoint {
	isManaged := func(domain string) bool {
		for _, e := range managed {
			if e.DNSName == domain {
				return true
			}
		}
		return false
	}

	var result []*domain.Endpoint
	for _, endpoint := range inp {
		if isManaged(endpoint.DNSName) {
			result = append(result, endpoint)
		}
	}
	return result
}
