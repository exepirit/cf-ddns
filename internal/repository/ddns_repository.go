package repository

type ObjectRepository interface {
	Init() error
	Reset() error
}

type RecordRepository interface {
	ObjectRepository
	RecordGetter
	RecordUpdater
	Add(record DDNSRecord) error
}

type RecordGetter interface {
	GetAll() []DDNSRecord
}

type RecordUpdater interface {
	Update(record DDNSRecord) error
}
