package static

import (
	"github.com/exepirit/cf-ddns/domain"
	"os"
	"strings"
)

type Source struct {
	domains []string
	ip      domain.Target
}

func NewSourceFromEnv() *Source {
	domainsStr := os.Getenv("DOMAINS")
	targetsStr := os.Getenv("TARGETS")
	return &Source{
		domains: strings.Split(domainsStr, ","),
		ip:      strings.Split(targetsStr, ","),
	}
}

func (s *Source) GetEndpoints() ([]*domain.Endpoint, error) {
	endpoints := make([]*domain.Endpoint, len(s.domains))
	for i, d := range s.domains {
		endpoints[i] = &domain.Endpoint{
			DNSName:    d,
			Target:     s.ip,
			RecordType: domain.RecordTypeA,
		}
	}
	return endpoints, nil
}
