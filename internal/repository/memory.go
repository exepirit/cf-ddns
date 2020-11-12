package repository

import (
	"github.com/pkg/errors"
	"sync"
)

type inMemoryRepo struct {
	records []DDNSRecord
	lock    sync.RWMutex
}

func NewMemory() RecordRepository {
	return &inMemoryRepo{records: make([]DDNSRecord, 0)}
}

func (i *inMemoryRepo) Init() error {
	return nil
}

func (i *inMemoryRepo) Reset() error {
	i.records = make([]DDNSRecord, 0)
	return nil
}

func (i *inMemoryRepo) GetAll() []DDNSRecord {
	i.lock.RLock()
	defer i.lock.RUnlock()
	result := make([]DDNSRecord, len(i.records))
	copy(result, i.records)
	return result
}

func (i *inMemoryRepo) Get(domain string) (DDNSRecord, error) {
	for _, record := range i.records {
		if record.Domain == domain {
			return record, nil
		}
	}
	return DDNSRecord{}, errors.New("record not found")
}

func (i *inMemoryRepo) Add(record DDNSRecord) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.records = append(i.records, record)
	return nil
}

func (i *inMemoryRepo) Update(record DDNSRecord) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	for idx, old := range i.records {
		if old.Domain == record.Domain {
			i.records[idx] = record
			return nil
		}
	}
	return errors.New("record not found")
}
