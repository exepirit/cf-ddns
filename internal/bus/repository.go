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
	case AddDomainRecord:
		record := event.(repository.DnsBinding)
		if err := r.Add(record); err != nil {
			log.Println(err)
		}
	}
}
