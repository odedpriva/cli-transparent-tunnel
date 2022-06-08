package mytypes

import (
	"fmt"
	"strconv"
	"strings"
)

type Endpoint struct {
	Host string
	User string
	Port int
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

func ConvertToEndpoint(s string) *Endpoint {
	endpoint := &Endpoint{
		Host: s,
	}
	if parts := strings.Split(endpoint.Host, "@"); len(parts) > 1 {
		endpoint.User = strings.Join(parts[0:len(parts)-1], "@")
		endpoint.Host = parts[len(parts)-1]
	}
	if parts := strings.Split(endpoint.Host, ":"); len(parts) > 1 {
		endpoint.Host = parts[0]
		endpoint.Port, _ = strconv.Atoi(parts[1])
	}
	return endpoint
}

func ConvertToEndpointWithDefault(s string, proto string) *Endpoint {
	e := ConvertToEndpoint(s)
	if e.Port == 0 {
		switch proto {
		case "http":
			e.Port = 80
		case "https":
			e.Port = 443
		case "ssh":
			e.Port = 22
		}
	}

	return e
}
