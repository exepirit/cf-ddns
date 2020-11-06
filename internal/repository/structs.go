package repository

import "time"

type DDNSRecord struct {
	Domain       string
	UpdatePeriod time.Duration
}
