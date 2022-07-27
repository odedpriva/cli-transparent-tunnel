package mytypes

import (
	"net"
	"net/url"
	"os/user"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ApplicationEndpoint struct {
	Host string
	Port string
}

func ConvertTApplicationEndpoint(s string) *ApplicationEndpoint {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		logrus.WithError(err).Warnf("failed to split to host:port %s", s)
		return &ApplicationEndpoint{
			Host: s,
		}
	}
	return &ApplicationEndpoint{
		Host: host,
		Port: port,
	}
}

func (e *ApplicationEndpoint) String() string {
	return net.JoinHostPort(e.Host, e.Port)
}

type NetworkEndpoint struct {
	Host string
	User string
	Port string
}

func (s *NetworkEndpoint) String() string {
	return net.JoinHostPort(s.Host, s.Port)
}

func ConvertToSshEndpoint(s string) (*NetworkEndpoint, error) {
	// make sure there is no scheme
	if strings.Contains(s, "://") {
		return nil, errors.Errorf("unexpected scheme in %s", s)
	}

	// parse address
	sshEndpointURL, err := url.Parse("ssh://" + s)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", s)
	}

	// make sure we have a user
	parsedUser := sshEndpointURL.User.Username()
	if parsedUser == "" {
		u, err2 := user.Current()
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to get the current user to use as ssh user")
		}
		if u.Username == "" {
			return nil, errors.Errorf("the current user is empty")
		}

		parsedUser = u.Username
	}

	// make sure we have a port
	parsedPort := sshEndpointURL.Port()
	if parsedPort == "" {
		parsedPort = "22"
	}

	return &NetworkEndpoint{
		Host: sshEndpointURL.Hostname(),
		User: parsedUser,
		Port: parsedPort,
	}, nil
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