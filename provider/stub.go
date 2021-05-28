package provider

import (
	"github.com/exepirit/cf-ddns/domain"
	"github.com/exepirit/cf-ddns/plan"
	"log"
)

type Stub struct {
	Logger *log.Logger
}

func (s *Stub) ApplyPlan(plan plan.Plan) error {
	changes := plan.Changes

	for _, e := range changes.Create {
		s.Logger.Printf("create %q record for %q: %s", e.RecordType, e.DNSName, e.Target)
	}
	for _, e := range changes.Delete {
		s.Logger.Printf("delete %q record for %q", e.RecordType, e.DNSName)
	}
	for _, e := range changes.Update {
		s.Logger.Printf("update %q record for %q: %s", e.RecordType, e.DNSName, e.Target)
	}

	return nil
}

func (*Stub) CurrentEndpoints() ([]*domain.Endpoint, error) {
	return []*domain.Endpoint{}, nil
}
