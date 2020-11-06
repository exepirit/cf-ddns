package repository

type ObjectRepository interface {
	Init() error
	Reset() error
}

type DDNSRepository interface {
	ObjectRepository
	GetAll() []DDNSRecord
	Add(record DDNSRecord) error
}
