package domain

const (
	RecordTypeA = "A"
)

// Target is array of network endpoints that domain points to.
type Target []string

func (target Target) Equal(other Target) bool {
	if len(target) != len(other) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(target))
	for _, _x := range target {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range other {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}

// Endpoint is a IP (Service) - Domain connection.
type Endpoint struct {
	ID         string `json:"id,omitempty"`
	DNSName    string `json:"dnsName"`
	Target     Target `json:"target"`
	RecordType string `json:"recordType"`
	TTL        int    `json:"ttl"`
}

func NewEndpoint(dnsName, recordType string, targets ...string) *Endpoint {
	return &Endpoint{
		DNSName:    dnsName,
		Target:     targets,
		RecordType: recordType,
		TTL:        0,
	}
}

// Equal returns true if other Endpoint is same than target.
func (e *Endpoint) Equal(other *Endpoint) bool {
	return e.DNSName == other.DNSName &&
		e.RecordType == other.RecordType &&
		e.Target.Equal(other.Target)
}
