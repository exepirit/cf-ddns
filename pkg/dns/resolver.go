package dns

import (
	"context"
	"net"
	"sync"
)

type Resolver struct {
	ServerAddr string
	Dial       func(network, address string) (net.Conn, error)
	resolver   *net.Resolver
	lock       sync.Mutex
}

func (r *Resolver) LookupIP(ctx context.Context, domain string) ([]net.IPAddr, error) {
	resolver := r.defaultResolver()
	return resolver.LookupIPAddr(ctx, domain)
}

func (r *Resolver) defaultResolver() *net.Resolver {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.resolver == nil {
		r.resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				// TODO: add context handling
				return r.dial()(network, r.ServerAddr)
			},
		}
	}
	return r.resolver
}

func (r *Resolver) dial() func(network, address string) (net.Conn, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.Dial == nil {
		r.Dial = net.Dial
	}
	return r.Dial
}
