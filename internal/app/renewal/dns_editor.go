package renewal

import (
	"context"
	"log"
	"net"

	"github.com/exepirit/cf-ddns/pkg/ddns"
	"github.com/exepirit/cf-ddns/pkg/lookup"
)

type dnsEditor struct {
	dns       *lookup.Resolver
	updater   *ddns.Updater
	currentIp net.IP
}

func (e *dnsEditor) updateDomain(ctx context.Context, domain string) error {
	resolvedIPs, err := e.dns.LookupIP(ctx, domain)
	if err != nil {
		return err
	}

	for _, ip := range resolvedIPs {
		if !ip.IP.Equal(e.currentIp) {
			log.Println("Updating domain", domain)
			if err = e.updater.UpdateARecord(domain, e.currentIp); err != nil {
				return err
			}
			log.Printf("Domain %s now refers to %s\n", domain, e.currentIp.String())
			return nil
		}
	}

	return nil
}
