package cloudflare

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/exepirit/cf-ddns/domain"
	"github.com/exepirit/cf-ddns/plan"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Provider struct {
	Zone cloudflare.Zone
	CF   *cloudflare.API
}

func NewProvider(zoneID, apiKey, email string) (*Provider, error) {
	api, err := cloudflare.New(apiKey, email)
	if err != nil {
		return nil, errors.WithMessage(err, "unable connect to cloudflare")
	}

	return &Provider{
		Zone: cloudflare.Zone{
			ID: zoneID,
		},
		CF: api,
	}, nil
}

func (provider *Provider) ApplyPlan(p plan.Plan) error {
	if p.Changes == nil {
		p.Eval()
	}
	provider.resolveRecords(p.Changes.Delete)
	provider.resolveRecords(p.Changes.Update)

	err := provider.createRecords(p.Changes.Create)
	if err != nil {
		return err
	}

	err = provider.deleteRecords(p.Changes.Delete)
	if err != nil {
		return err
	}

	return provider.updateRecords(p.Changes.Update)
}

func (provider *Provider) resolveRecords(endpoints []*domain.Endpoint) {
	records, err := provider.CF.DNSRecords(provider.Zone.ID, cloudflare.DNSRecord{})
	if err != nil {
		log.Error().Err(err).Msg("failed to resolve records")
	}

	for _, endpoint := range endpoints {
		if endpoint.ID != "" {
			continue
		}

		for _, d := range records {
			if d.Name == endpoint.DNSName {
				endpoint.ID = d.ID
			}
		}
	}
}

func (provider *Provider) createRecords(endpoints []*domain.Endpoint) error {
	var err error
	for _, endpoint := range endpoints {
		record := cloudflare.DNSRecord{
			Type:    endpoint.RecordType,
			Name:    endpoint.DNSName,
			Content: endpoint.Target[0],
			TTL:     endpoint.TTL,
		}
		_, err = provider.CF.CreateDNSRecord(provider.Zone.ID, record)
		if err != nil {
			return errors.WithMessagef(err, "create domain %q", endpoint.DNSName)
		}
		log.Info().
			Str("domain", endpoint.DNSName).
			Str("ip", endpoint.Target[0]).
			Msg("new domain created")
	}
	return err
}

func (provider *Provider) deleteRecords(endpoints []*domain.Endpoint) error {
	var err error
	for _, endpoint := range endpoints {
		err = provider.CF.DeleteDNSRecord(provider.Zone.ID, endpoint.ID)
		if err != nil {
			return errors.WithMessagef(err, "delete domain %q", endpoint.DNSName)
		}
		log.Info().Msgf("domain %q removed", endpoint.DNSName)
	}
	return err
}

func (provider *Provider) updateRecords(endpoints []*domain.Endpoint) error {
	var err error
	for _, endpoint := range endpoints {
		record := cloudflare.DNSRecord{
			Type:    endpoint.RecordType,
			Name:    endpoint.DNSName,
			Content: endpoint.Target[0],
			TTL:     endpoint.TTL,
		}
		err = provider.CF.UpdateDNSRecord(provider.Zone.ID, endpoint.ID, record)
		if err != nil {
			return errors.WithMessagef(err, "update domain %q", endpoint.DNSName)
		}
		log.Info().
			Str("domain", endpoint.DNSName).
			Str("ip", endpoint.Target[0]).
			Msg("domain bound updated")
	}
	return err
}

func (provider *Provider) CurrentEndpoints() ([]*domain.Endpoint, error) {
	records, err := provider.CF.DNSRecords(provider.Zone.ID, cloudflare.DNSRecord{})
	if err != nil {
		return nil, errors.WithMessage(err, "get current DNS records")
	}

	endpoints := make([]*domain.Endpoint, len(records))
	for i, record := range records {
		endpoints[i] = &domain.Endpoint{
			ID:         record.ID,
			DNSName:    record.Name,
			Target:     domain.Target{record.Content},
			RecordType: record.Type,
			TTL:        record.TTL,
		}
	}
	return endpoints, nil
}
