package bus

import (
	"github.com/exepirit/cf-ddns/internal/repository"
	"log"
)

type RepositoryConsumer struct {
	repository.BindingRepository
}

func (r RepositoryConsumer) Consume(event interface{}) {
	switch event.(type) {
	case AddDomainBinding:
		record := event.(AddDomainBinding)
		if err := r.Add(repository.DnsBinding(record)); err != nil {
			log.Println(err)
		}
	case RemoveDomainBinding:
		panic("implement me")
	case UpdateDomainBinding:
		b := event.(repository.DnsBinding)
		if err := r.Update(b); err != nil {
			log.Println(err)
		}
	}
}
