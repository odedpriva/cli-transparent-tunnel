package tunnling

import (
	"ctt/config"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/ssh"
)

type Endpoint struct {
	Host string
	Port int
	User string
}

func NewEndpoint(s string) *Endpoint {
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

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint
	Config *ssh.ClientConfig
	Log    *logrus.Logger
}

func NewSSHTunnel(serverUser string, sshConfig *config.SshConfig, logger *logrus.Logger) *SSHTunnel {
	var auth ssh.AuthMethod
	if sshConfig.KeyPath != "" {
		auth = privateKeyFile(sshConfig.KeyPath)
	}
	return &SSHTunnel{
		Config: &ssh.ClientConfig{
			User: serverUser,
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		},
		Log: logger,
	}
}

func (t *SSHTunnel) SetEndpoints(server, remote string) {
	t.Local = &Endpoint{
		Host: "localhost",
		Port: 0,
		User: "",
	}
	if !strings.Contains(server, ":") {
		server = server + ":22"
	}
	t.Server = NewEndpoint(server)
	t.Remote = NewEndpoint(remote)
}

func (t *SSHTunnel) Start(addressChan chan string, errChan chan error) (string, error) {
	var err error
	var listener net.Listener
	var conn net.Conn

	listener, err = net.Listen("tcp", t.Local.String())
	if err != nil {
		errChan <- err
	}
	defer listener.Close()
	addressChan <- listener.Addr().String()
	for {
		conn, err = listener.Accept()
		if err != nil {
			errChan <- err
		}
		t.Log.Debug("accepted connection")
		go t.forward(conn)
	}
}

func (t *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", t.Server.String(), t.Config)
	if err != nil {
		t.Log.Errorf("server dial error: %s", err)
		return
	}
	t.Log.Debugf("connected to %s (1 of 2)\n", t.Server.String())
	remoteConn, err := serverConn.Dial("tcp", t.Remote.String())
	if err != nil {
		t.Log.Errorf("remote dial error: %s", err)
		return
	}
	t.Log.Debugf("connected to %s (2 of 2)\n", t.Remote.String())
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			t.Log.Errorf("io.Copy error: %s", err)
		}
	}
	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func privateKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
