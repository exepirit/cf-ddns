package repository

import "sync"

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

func (i *inMemoryRepo) Add(record DDNSRecord) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.records = append(i.records, record)
	return nil
}
