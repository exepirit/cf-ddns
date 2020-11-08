package bus

import (
	"github.com/exepirit/cf-ddns/internal/repository"
	"log"
)

type RepositoryConsumer struct {
	repository.RecordRepository
}

func (r RepositoryConsumer) Consume(event interface{}) {
	switch event.(type) {
	case AddDomainRecord:
		record := repository.DDNSRecord(event.(AddDomainRecord))
		if err := r.Add(record); err != nil {
			log.Println(err)
		}
	}
}
