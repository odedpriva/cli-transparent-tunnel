package network_tunnler

import (
	"errors"
	"fmt"
	"github.com/odedpriva/cli-transparent-tunnel/logging"
	"io"
	"io/ioutil"
	"net"

	"github.com/odedpriva/cli-transparent-tunnel/mytypes"
	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	Server *mytypes.ApplicationEndpoint
	Remote *mytypes.ApplicationEndpoint
	Config *ssh.ClientConfig

	log         *logging.Logging
	AddressChan chan string
	ErrChan     chan error
}

type NetworkTunneler interface {
	Start(server *mytypes.NetworkEndpoint, originServer string)
	Wait() (string, error)
}

func NewSSHTunnel(c *mytypes.SshConfig) (*SSHTunnel, error) {
	var auth ssh.AuthMethod
	var err error
	log := logging.GetLogger()
	if c.KeyPath == "" {
		return nil, fmt.Errorf("no auth method for ssh")
	}
	log.Debugf("using ssh key %s", c.KeyPath)
	auth, err = privateKeyFile(c.KeyPath)
	if err != nil {
		return nil, err
	}
	return &SSHTunnel{
		Config: &ssh.ClientConfig{
			User: c.User,
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		},
		log:         log,
		AddressChan: make(chan string),
		ErrChan:     make(chan error),
	}, nil
}

func (t *SSHTunnel) Start(server *mytypes.NetworkEndpoint, originServer string) {
	var err error
	var listener net.Listener
	var conn net.Conn

	localEndpoint := &mytypes.ApplicationEndpoint{
		Host: "localhost",
	}

	t.log.Debugf("starting tcp lisener on %s", localEndpoint.String())

	listener, err = net.Listen("tcp", localEndpoint.String())
	if err != nil {
		t.ErrChan <- err
	}
	defer func() {
		err = listener.Close()
		t.log.Errorf("%s\n", err)
	}()
	t.AddressChan <- listener.Addr().String()
	for {
		conn, err = listener.Accept()
		if err != nil {
			t.ErrChan <- err
		}
		t.log.Debug("accepted connection")
		go t.forward(conn, server, originServer)
	}
}

func (t *SSHTunnel) Wait() (string, error) {
	select {
	case err := <-t.ErrChan:
		return "", err
	case localAddress := <-t.AddressChan:
		return localAddress, nil
	}
}

func (t *SSHTunnel) forward(localConn net.Conn, sshServer *mytypes.NetworkEndpoint, originServer string) {
	t.log.Debugf("forwording connection on ssh tunnel %s", sshServer.String())
	t.log.Debugf("using ssh configuration %+v", t.Config)
	serverConn, err := ssh.Dial("tcp", sshServer.String(), t.Config)
	if err != nil {
		t.log.Errorf("server dial error: %s", err)
		return
	}
	t.log.Debugf("connected to %s (1 of 2)\n", sshServer.String())
	remoteConn, err := serverConn.Dial("tcp", originServer)
	if err != nil {
		t.log.Errorf("remote dial error: %s", err)
		return
	}
	t.log.Debugf("connected to %s (2 of 2)\n", originServer)
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				t.log.Infof("io.Copy error: %s", err)
			}
			t.log.Errorf("io.Copy error: %s", err)
		}
	}
	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func privateKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
