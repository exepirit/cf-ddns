package echoip

import "net"

type Resolver interface {
	GetIP() (net.IP, error)
}
