package file

import (
	"encoding/json"
	"os"

	"github.com/exepirit/cf-ddns/domain"
	"github.com/pkg/errors"
)

//Source is a factory that makes endpoints from provided configuartion file.
type Source struct {
	filepath string
}

//NewSource create new file.Source.
func NewSource(filepath string) *Source {
	return &Source{
		filepath: filepath,
	}
}

func (s Source) GetEndpoints() ([]*domain.Endpoint, error) {
	data, err := s.readFile()
	if err != nil {
		return nil, errors.WithMessage(err, "read config file")
	}

	endpoints := make([]*domain.Endpoint, len(data.Domains))
	for i := 0; i < len(endpoints); i++ {
		endpoints[i] = &domain.Endpoint{
			DNSName:    data.Domains[i],
			Target:     data.Targets,
			RecordType: domain.RecordTypeA,
		}
	}
	return endpoints, nil
}

func (s Source) readFile() (*fileContent, error) {
	file, err := os.Open(s.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content := &fileContent{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(content)
	return content, errors.WithMessage(err, "decode file")
}
