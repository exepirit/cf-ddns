package repository

import (
	"github.com/pkg/errors"
	"sync"
)

type inMemoryRepo struct {
	records []DnsBinding
	lock    sync.RWMutex
}

func NewMemory() BindingRepository {
	return &inMemoryRepo{records: make([]DnsBinding, 0)}
}

func (i *inMemoryRepo) Init() error {
	return nil
}

func (i *inMemoryRepo) Reset() error {
	i.records = make([]DnsBinding, 0)
	return nil
}

func (i *inMemoryRepo) GetAll() []DnsBinding {
	i.lock.RLock()
	defer i.lock.RUnlock()
	result := make([]DnsBinding, len(i.records))
	copy(result, i.records)
	return result
}

func (i *inMemoryRepo) Get(domain string) (DnsBinding, error) {
	for _, record := range i.records {
		if record.Domain == domain {
			return record, nil
		}
	}
	return DnsBinding{}, errors.New("record not found")
}

func (i *inMemoryRepo) Add(binding DnsBinding) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.records = append(i.records, binding)
	return nil
}

func (i *inMemoryRepo) Update(binding DnsBinding) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	for idx, old := range i.records {
		if old.Domain == binding.Domain {
			i.records[idx] = binding
			return nil
		}
	}
	return errors.New("binding not found")
}
