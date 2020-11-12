package renewal

import (
	"sync"
	"time"

	"github.com/exepirit/cf-ddns/internal/bus"
	"github.com/exepirit/cf-ddns/internal/repository"
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
	done := make(chan struct{})

	waitForDomain := func(domain string) {
		select {
		case <-d.domainsTickers[domain].C:
			result <- domain
		case <-done:
			return
		}
	}

	d.lock.RLock()
	defer d.lock.RLocker()
	for domain := range d.domainsTickers {
		go waitForDomain(domain)
	}
	r := <-result
	close(done)
	return r
}

func (d *domains) addDomain(name string, checkInterval time.Duration) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.domainsTickers[name] = time.NewTicker(checkInterval)
}

func (d *domains) removeDomain(name string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	_, ok := d.domainsTickers[name]
	if ok {
		delete(d.domainsTickers, name)
	}
	return ok
}

func (d *domains) Consume(event interface{}) {
	switch event.(type) {
	case bus.AddDomainBinding:
		binding := event.(repository.DnsBinding)
		d.addDomain(binding.Domain, binding.UpdatePeriod)
	case bus.RemoveDomainBinding:
		binding := event.(repository.DnsBinding)
		d.removeDomain(binding.Domain)
	case bus.UpdateDomainBinding:
		panic("implement me")
	}
}
