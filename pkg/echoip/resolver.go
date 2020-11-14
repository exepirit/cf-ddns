package echoip

import (
	"context"
	"net"
)

type Resolver interface {
	GetIP(ctx context.Context) (net.IP, error)
}
