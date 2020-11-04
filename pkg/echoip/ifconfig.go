package echoip

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type IfconfigResolver struct{}

func (IfconfigResolver) GetIP() (net.IP, error) {
	resp, err := http.Get("https://ifconfig.co/")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ifconfig.co returns status %d", resp.StatusCode)
	}

	body := make([]byte, 15)
	n, err := resp.Body.Read(body)
	if err != nil {
		return nil, err
	}
	ipString := strings.TrimSuffix(string(body[:n]), "\n")
	return parseIp(ipString)
}

func parseIp(body string) (net.IP, error) {
	parts := strings.SplitN(body, ".", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("cannot parse ip from %q", body)
	}
	ip := net.IP{0, 0, 0, 0}
	for i, part := range parts {
		octet, err := strconv.Atoi(part)
		if err != nil || octet < 0 || octet > 255 {
			return nil, fmt.Errorf("cannot parse octet from %q", part)
		}
		ip[i] = byte(octet)
	}
	return ip, nil
}
