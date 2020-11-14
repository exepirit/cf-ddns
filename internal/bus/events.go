package bus

import "github.com/exepirit/cf-ddns/internal/repository"

type AddDomainBinding repository.DnsBinding
type RemoveDomainBinding repository.DnsBinding
type UpdateDomainBinding repository.DnsBinding

type DnsRecordUpdated string
