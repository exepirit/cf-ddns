package repository

import "time"

type DnsBinding struct {
	Domain       string
	UpdatePeriod time.Duration
}
