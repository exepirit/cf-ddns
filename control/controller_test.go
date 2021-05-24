package control

import (
	"testing"

	"github.com/exepirit/cf-ddns/domain"
	"github.com/exepirit/cf-ddns/internal/test"
	"github.com/exepirit/cf-ddns/plan"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func endpointsFromDomains(domains []string) []*domain.Endpoint {
	endpoints := make([]*domain.Endpoint, len(domains))
	for i, d := range domains {
		endpoints[i] = &domain.Endpoint{DNSName: d}
	}
	return endpoints
}

func assertDNSChanges(t *testing.T, expected []string, actual []*domain.Endpoint) {
	require.Len(t, actual, len(expected))
	actualDomains := make([]string, len(actual))
	for i, endpoint := range actual {
		actualDomains[i] = endpoint.DNSName
	}
	require.ElementsMatch(t, expected, actualDomains)
}

func TestController_RunOnce(t *testing.T) {
	table := []struct {
		name           string
		currentDomains []string
		desiredDomains []string

		addedDomains   []string
		removedDomains []string
		updatedDomains []string
	}{
		{
			name:           "create managed domain",
			currentDomains: []string{"example.com"},
			desiredDomains: []string{"example.com", "test.example.com"},
			addedDomains:   []string{"test.example.com"},
		},
		{
			name:           "remove managed domain",
			currentDomains: []string{"example.com", "test.example.com"},
			desiredDomains: []string{"example.com"},
			removedDomains: []string{"test.example.com"},
		},
		{
			name:           "ignore unmanaged domain",
			currentDomains: []string{"unmanaged.com", "test.example.com"},
			desiredDomains: []string{"test.example.com"},
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			desired := endpointsFromDomains(testCase.desiredDomains)
			fakeSource := new(test.SourceMock)
			fakeSource.On("GetEndpoints").Return(desired, nil)

			current := endpointsFromDomains(testCase.currentDomains)
			fakeProvider := new(test.ProviderMock)
			fakeProvider.On("CurrentEndpoints").Return(current, nil)
			fakeProvider.On("ApplyPlan", mock.Anything).Return(nil).Once()

			ctrl := &Controller{
				Source:   fakeSource,
				Provider: fakeProvider,
			}
			err := ctrl.RunOnce()
			require.NoError(t, err)

			fakeSource.AssertExpectations(t)
			fakeProvider.AssertExpectations(t)
			appliedPlan := fakeProvider.Calls[1].Arguments.Get(0).(plan.Plan)
			changes := appliedPlan.Changes
			assertDNSChanges(t, testCase.addedDomains, changes.Create)
			assertDNSChanges(t, testCase.removedDomains, changes.Delete)
			assertDNSChanges(t, testCase.updatedDomains, changes.Update)
		})
	}
}
