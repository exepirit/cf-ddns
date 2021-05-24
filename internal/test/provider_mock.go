package test

import (
	"github.com/exepirit/cf-ddns/domain"
	"github.com/exepirit/cf-ddns/plan"
	"github.com/stretchr/testify/mock"
)

type ProviderMock struct {
	mock.Mock
}

func (m *ProviderMock) ApplyPlan(plan plan.Plan) error {
	return m.MethodCalled("ApplyPlan", plan).Error(0)
}

func (m *ProviderMock) CurrentEndpoints() ([]*domain.Endpoint, error) {
	args := m.MethodCalled("CurrentEndpoints")
	return args.Get(0).([]*domain.Endpoint), args.Error(1)
}
