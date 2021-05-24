package test

import (
	"github.com/exepirit/cf-ddns/domain"
	"github.com/stretchr/testify/mock"
)

type SourceMock struct {
	mock.Mock
}

func (m *SourceMock) GetEndpoints() ([]*domain.Endpoint, error) {
	args := m.MethodCalled("GetEndpoints")
	return args.Get(0).([]*domain.Endpoint), args.Error(1)
}
