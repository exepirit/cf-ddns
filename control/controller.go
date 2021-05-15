package control

import (
	"github.com/exepirit/cf-ddns/plan"
	"github.com/exepirit/cf-ddns/provider"
	"github.com/exepirit/cf-ddns/source"
	"github.com/pkg/errors"
	"time"
)

type Controller struct {
	Source     source.Source
	Provider   provider.Provider
	TimePeriod time.Duration
}

func (ctrl *Controller) RunOnce() error {
	currentState, err := ctrl.Provider.CurrentEndpoints()
	if err != nil {
		return errors.WithMessage(err, "failed to get bonded domains")
	}

	desiredState, err := ctrl.Source.GetEndpoints()
	if err != nil {
		return errors.WithMessage(err, "failed to get desired domains")
	}

	currentPlan := plan.Plan{
		Current: currentState,
		Desired: desiredState,
	}
	currentPlan.Eval()

	err = ctrl.Provider.ApplyPlan(currentPlan)
	return err
}

func (ctrl *Controller) Run() error {
	ticker := time.NewTicker(ctrl.TimePeriod)
	defer ticker.Stop()

	var err error
	for {
		err = ctrl.RunOnce()
		if err != nil {
			return err
		}
		<-ticker.C
	}
}
