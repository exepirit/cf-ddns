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

	previousDesired []*domain.Endpoint
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

	// managed = current - reviousDesired - currentDesired
	unmanaged := removeEndpoints(removeEndpoints(currentState, ctrl.previousDesired), desiredState)
	managed := removeEndpoints(currentState, unmanaged)

	currentPlan := plan.Plan{
		Current: managed,
		Desired: desiredState,
	}
	currentPlan.Eval()
	ctrl.previousDesired = desiredState

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

func (ctrl *Controller) getManagedEndpoints(desiredState []*domain.Endpoint) []*domain.Endpoint {
	managed := make([]*domain.Endpoint, len(ctrl.previousDesired)+len(desiredState))
	managedEndpointsCount := 0

	inManaged := func(e *domain.Endpoint) bool {
		for _, m := range managed {
			if m.Equal(e) {
				return true
			}
		}
		return false
	}

	copy(managed[:len(ctrl.previousDesired)], ctrl.previousDesired[:])
	managedEndpointsCount += len(ctrl.previousDesired)
	for _, endpoint := range ctrl.previousDesired {
		if !inManaged(endpoint) {
			managed[managedEndpointsCount] = endpoint
			managedEndpointsCount++
		}
	}

	return managed[:managedEndpointsCount]
}

func removeEndpoints(array []*domain.Endpoint, prune []*domain.Endpoint) []*domain.Endpoint {
	arr := make([]*domain.Endpoint, len(array))
	copy(arr[:], array[:])
	for i := 0; i < len(arr); i++ {
		origin := arr[i]
		for _, rem := range prune {
			if origin.Equal(rem) {
				arr = append(arr[:i], arr[i+1:]...)
				i--
				break
			}
		}
	}
	return arr
}
