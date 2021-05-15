package plan

import "github.com/exepirit/cf-ddns/domain"

// Plan calculate list of domains changes.
type Plan struct {
	Current []*domain.Endpoint
	Desired []*domain.Endpoint
	Changes *Changes
}

// Eval calculate difference between current and desired state.
func (plan *Plan) Eval() {
	currentMap := makeDomainMap(plan.Current)
	desiredMap := makeDomainMap(plan.Desired)

	changes := &Changes{}
	changes.Create = desiredMap.filterDomains(plan.Current).endpoints()
	changes.Delete = currentMap.filterDomains(plan.Desired).endpoints()

	for domainName, currentEndpoint := range currentMap {
		desiredEndpoint, ok := desiredMap[domainName]
		if !ok {
			continue
		}

		if !currentEndpoint.Target.Equal(desiredEndpoint.Target) {
			changes.Update = append(changes.Update, desiredEndpoint)
		}
	}

	plan.Changes = changes
}

type domainMap map[string]*domain.Endpoint

func makeDomainMap(endpoints []*domain.Endpoint) domainMap {
	result := make(map[string]*domain.Endpoint)
	for _, record := range endpoints {
		result[record.DNSName] = record
	}
	return result
}

func (m domainMap) copy() domainMap {
	result := make(domainMap)
	for key, value := range m {
		result[key] = value
	}
	return result
}

func (m domainMap) filterDomains(endpoints []*domain.Endpoint) domainMap {
	result := m.copy()
	for _, r := range endpoints {
		delete(result, r.DNSName)
	}
	return result
}

func (m domainMap) endpoints() []*domain.Endpoint {
	result := make([]*domain.Endpoint, len(m))
	i := 0
	for _, record := range m {
		result[i] = record
		i++
	}
	return result
}
