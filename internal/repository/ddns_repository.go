package repository

import (
	"github.com/pkg/errors"
	"sync"
)

type ObjectRepository interface {
	Init() error
	Reset() error
}

type DDNSRepository interface {
	ObjectRepository
	GetAll() []DDNSRecord
	Add(record DDNSRecord) error
}

var repoContainer = struct {
	repo DDNSRepository
	lock sync.Mutex
}{}

func Get() DDNSRepository {
	if repoContainer.repo == nil {
		panic(errors.New("ddns records repository is not set"))
	}
	return repoContainer.repo
}

func Set(r DDNSRepository) {
	repoContainer.repo = r
}
