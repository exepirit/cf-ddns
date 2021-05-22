package control

import (
	"testing"

	"github.com/exepirit/cf-ddns/domain"
	"github.com/stretchr/testify/require"
)

func TestController_FilterUnmanagedDomains(t *testing.T) {
	table := []struct {
		name                   string
		currentDomains         []string
		desiredDomains         []string
		filteredCurrentDomains []string
	}{
		{
			name:                   "filter 1 unmanaged domain",
			currentDomains:         []string{"example.com", "test.example.com"},
			desiredDomains:         []string{"test.example.com"},
			filteredCurrentDomains: []string{"test.example.com"},
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := &Controller{}
			current := make([]*domain.Endpoint, len(testCase.currentDomains))
			for i, d := range testCase.currentDomains {
				current[i] = &domain.Endpoint{DNSName: d}
			}

			desired := make([]*domain.Endpoint, len(testCase.desiredDomains))
			for i, d := range testCase.desiredDomains {
				desired[i] = &domain.Endpoint{DNSName: d}
			}

			filtered := make([]*domain.Endpoint, len(testCase.filteredCurrentDomains))
			for i, d := range testCase.filteredCurrentDomains {
				filtered[i] = &domain.Endpoint{DNSName: d}
			}

			result := ctrl.filterUnmanagedDomains(current, desired)

			require.ElementsMatch(t, filtered, result)
		})
	}
}
