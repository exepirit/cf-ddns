package renewal

import (
	"context"
	"net"
	"time"

	"github.com/exepirit/cf-ddns/internal/bus"
	"github.com/exepirit/cf-ddns/pkg/ddns"
	"github.com/exepirit/cf-ddns/pkg/echoip"
	"github.com/exepirit/cf-ddns/pkg/lookup"
)

// Renewer update domains A records, if they IP addresses, defined in DNS A records are not refers to current host IP.
type Renewer interface {
	RenewAllDomains() error
	AddDomain(domain string, checkPeriod time.Duration) error
}

// NewRenewer produce new Renewer.
func NewRenewer(ip echoip.Resolver, dnsResolver *lookup.Resolver, dnsUpdater *ddns.Updater) Renewer {
	return &renewer{
		ipResolver:  ip,
		dnsResolver: dnsResolver,
		dnsUpdater:  dnsUpdater,
		domains:     newDomains(),
	}
}

type renewer struct {
	ipResolver  echoip.Resolver
	dnsResolver *lookup.Resolver
	dnsUpdater  *ddns.Updater
	domains     *domains
}

func (w *renewer) RenewAllDomains() error {
	domain := <-w.domains.nextPendingDomain
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	currentIp, err := w.ipResolver.GetIP(context.Background())
	if err != nil {
		return err
	}

	err = w.setDomainARecord(ctx, domain, currentIp)
	if err == nil {
		bus.Get().Publish(bus.DnsRecordUpdated(domain))
	}
	return err
}

func (r *renewer) setDomainARecord(ctx context.Context, domain string, ip net.IP) error {
	resolvedIPs, err := r.dnsResolver.LookupIP(ctx, domain)
	if err != nil {
		return err
	}

	for _, currentIp := range resolvedIPs {
		if !currentIp.IP.Equal(ip) {
			if err = r.dnsUpdater.UpdateARecord(domain, ip); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (r *renewer) AddDomain(name string, checkPeriod time.Duration) error {
	r.domains.addDomain(name, checkPeriod)
	return nil
}
