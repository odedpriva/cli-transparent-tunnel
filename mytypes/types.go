package mytypes

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"
)

type ApplicationEndpoint struct {
	Host string
	Port int
}

func ConvertTApplicationEndpoint(s string) *ApplicationEndpoint {
	host, port := "", 0

	parts := strings.Split(s, ":")
	switch len(parts) {
	case 1:
		host = parts[0]
		port = 0
	case 2:
		host = parts[0]
		port, _ = strconv.Atoi(parts[1])
	}

	return &ApplicationEndpoint{
		Host: host,
		Port: port,
	}
}

func (e *ApplicationEndpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

type NetworkEndpoint struct {
	Host string
	User string
	Port int
}

func (s *NetworkEndpoint) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func ConvertToSshEndpoint(s string) *NetworkEndpoint {
	host, userToUse, port := "", "", 0

	parts := strings.Split(s, "@")
	switch len(parts) {
	case 1:
		userToUse = func() string { username, _ := user.Current(); return username.Name }()
		host = parts[0]
	case 2:
		userToUse = parts[0]
		host = parts[1]
	default:
		userToUse = strings.Join(parts[1:len(parts)-1], "@")
		host = parts[1]
	}

	parts = strings.Split(host, ":")
	switch len(parts) {
	case 1:
		host = parts[0]
		port = 22
	case 2:
		host = parts[0]
		port, _ = strconv.Atoi(parts[1])
	}

	return &NetworkEndpoint{
		Host: host,
		User: userToUse,
		Port: port,
	}
}

type Command struct {
	CliPath      string
	HostFlags    []string
	PortFlags    []string
	AddressFlags []string
	SNIFlags     []string
}

type SshConfig struct {
	KeyPath string
	User    string
}