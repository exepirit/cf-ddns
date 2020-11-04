package ddns

import (
	"errors"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"net"
)

type Updater struct {
	Zone cloudflare.Zone
	CF   *cloudflare.API
}

func NewUpdater(api *cloudflare.API, zoneName string) (*Updater, error) {
	zone, err := api.ListZones(zoneName)
	if err != nil {
		return nil, err
	}
	if len(zone) == 0 {
		return nil, errors.New("could not find zone")
	}
	return &Updater{
		Zone: zone[0],
		CF:   api,
	}, nil
}

func (u Updater) UpdateARecord(domain string, ip net.IP) error {
	records, err := u.CF.DNSRecords(u.Zone.ID, cloudflare.DNSRecord{Name: domain})
	if err != nil {
		return err
	}

	updateRecord := func(r cloudflare.DNSRecord, addr net.IP) error {
		r.Content = addr.String()
		return u.CF.UpdateDNSRecord(u.Zone.ID, r.ID, r)
	}

	if len(records) > 0 {
		return updateRecord(records[0], ip)
	} else {
		return fmt.Errorf("requested subdomain %q not found in zone %q", domain, u.Zone)
	}
}
