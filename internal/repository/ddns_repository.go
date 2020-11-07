package repository

type ObjectRepository interface {
	Init() error
	Reset() error
}

type RecordRepository interface {
	ObjectRepository
	RecordGetter
	Add(record DDNSRecord) error
}

type RecordGetter interface {
	GetAll() []DDNSRecord
}
