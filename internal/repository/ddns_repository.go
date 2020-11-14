package repository

type ObjectRepository interface {
	Init() error
	Reset() error
}

type BindingRepository interface {
	ObjectRepository
	BindingGetter
	BindingUpdater
	Add(b DnsBinding) error
}

type BindingGetter interface {
	GetAll() []DnsBinding
}

type BindingUpdater interface {
	Update(b DnsBinding) error
}
