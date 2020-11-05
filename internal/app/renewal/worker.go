package renewal

import (
	"context"
	"log"
	"time"

	"github.com/exepirit/cf-ddns/pkg/ddns"
	"github.com/exepirit/cf-ddns/pkg/echoip"
	"github.com/exepirit/cf-ddns/pkg/lookup"
)

type Worker struct {
	ipResolver echoip.Resolver
	editor     *dnsEditor
	domains    *domains
}

func NewWorker(ip echoip.Resolver, dnsResolver *lookup.Resolver, dnsUpdater *ddns.Updater) *Worker {
	editor := &dnsEditor{
		dns:     dnsResolver,
		updater: dnsUpdater,
	}
	domains := newDomains()
	return &Worker{
		ipResolver: ip,
		editor:     editor,
		domains:    domains,
	}
}

func (w *Worker) Run() {
	for {
		if err := w.updateAllDomains(); err != nil {
			log.Println(err)
		}
	}
}

func (w *Worker) updateAllDomains() error {
	currentIp, err := w.ipResolver.GetIP(context.Background())
	if err != nil {
		return err
	}
	w.editor.currentIp = currentIp

	domain := w.domains.next()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return w.editor.updateDomain(ctx, domain)
}

func (w *Worker) AddDomain(name string, checkInterval time.Duration) {
	w.domains.addDomain(name, checkInterval)
}
