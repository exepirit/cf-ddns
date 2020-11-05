package renewal

import (
	"sync"
	"time"
)

type domains struct {
	domainsTickers map[string]*time.Ticker
	newDomain      chan string
	lock           sync.RWMutex
}

func newDomains() *domains {
	return &domains{
		domainsTickers: map[string]*time.Ticker{},
		newDomain:      make(chan string),
		lock:           sync.RWMutex{},
	}
}

func (d *domains) next() string {
	result := make(chan string)

	waitForDomain := func(domain string) {
		select {
		case <-d.domainsTickers[domain].C:
			result <- domain
		}
	}

	d.lock.RLock()
	defer d.lock.RLocker()
	for domain := range d.domainsTickers {
		go waitForDomain(domain)
	}
	return <-result
}

func (d *domains) addDomain(name string, checkInterval time.Duration) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.domainsTickers[name] = time.NewTicker(checkInterval)
}
