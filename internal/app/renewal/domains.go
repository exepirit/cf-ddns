package renewal

import (
	"sync"
	"time"

	"github.com/exepirit/cf-ddns/internal/bus"
	"github.com/exepirit/cf-ddns/internal/repository"
)

type domains struct {
	domainsTickers    map[string]*time.Ticker
	mapLock           sync.RWMutex
	newBindings       chan repository.DnsBinding
	nextPendingDomain chan string
}

func newDomains() *domains {
	return &domains{
		domainsTickers:    map[string]*time.Ticker{},
		newBindings:       make(chan repository.DnsBinding),
		nextPendingDomain: make(chan string),
	}
}

func (d *domains) Handle(bus bus.Publisher) {
	bus.Subscribe(d)

	go func() {
		for {
			d.putNextPending()
		}
	}()
}

func (d *domains) putNextPending() {
	done := make(chan struct{})
	result := make(chan string)

	waitTicker := func(domain string, t *time.Ticker) {
		select {
		case <-t.C:
			result <- domain
			break
		case <-done:
			return
		}
	}

	d.mapLock.RLock()
	for domain, ticker := range d.domainsTickers {
		go waitTicker(domain, ticker)
	}
	d.mapLock.RUnlock()

	select {
	case r := <-result:
		d.nextPendingDomain <- r
		break
	case b := <-d.newBindings:
		d.nextPendingDomain <- b.Domain
		break
	}

	close(done)
}

func (d *domains) addDomain(name string, checkInterval time.Duration) {
	d.newBindings <- repository.DnsBinding{
		Domain:       name,
		UpdatePeriod: checkInterval,
	}
	d.mapLock.Lock()
	defer d.mapLock.Unlock()
	d.domainsTickers[name] = time.NewTicker(checkInterval)
}

func (d *domains) removeDomain(name string) bool {
	d.mapLock.Lock()
	defer d.mapLock.Unlock()
	_, ok := d.domainsTickers[name]
	if ok {
		delete(d.domainsTickers, name)
	}
	return ok
}

func (d *domains) Consume(event interface{}) {
	switch event.(type) {
	case bus.AddDomainBinding:
		binding := event.(bus.AddDomainBinding)
		d.addDomain(binding.Domain, binding.UpdatePeriod)
		break
	case bus.RemoveDomainBinding:
		binding := event.(bus.RemoveDomainBinding)
		d.removeDomain(binding.Domain)
		break
	case bus.UpdateDomainBinding:
		panic("implement me")
	}
}
