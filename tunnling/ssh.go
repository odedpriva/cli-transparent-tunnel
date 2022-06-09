package tunnling

import (
	"errors"
	"fmt"
	"github.com/odedpriva/cli-transparent-tunnel/config"
	"github.com/odedpriva/cli-transparent-tunnel/mytypes"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
)

type SSHTunnel struct {
	Server *mytypes.Endpoint
	Remote *mytypes.Endpoint
	Config *ssh.ClientConfig
	Log    *logrus.Logger

	AddressChan chan string
	ErrChan     chan error
}

func NewSSHTunnel(sshConfig *config.SshConfig, logger *logrus.Logger) (*SSHTunnel, error) {
	var auth ssh.AuthMethod
	var err error
	if sshConfig.KeyPath == "" {
		return nil, fmt.Errorf("no auth method for ssh")
	}
	logger.Debugf("using ssh key %s", sshConfig.KeyPath)
	auth, err = privateKeyFile(sshConfig.KeyPath)
	if err != nil {
		return nil, err
	}
	return &SSHTunnel{
		Config: &ssh.ClientConfig{
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		},
		Log:         logger,
		AddressChan: make(chan string),
		ErrChan:     make(chan error),
	}, nil
}

func (t *SSHTunnel) Start(sshServer, originServer *mytypes.Endpoint) {
	var err error
	var listener net.Listener
	var conn net.Conn

	localEndpoint := &mytypes.Endpoint{
		Host: "localhost",
	}

	listener, err = net.Listen("tcp", localEndpoint.String())
	if err != nil {
		t.ErrChan <- err
	}
	defer listener.Close()
	t.AddressChan <- listener.Addr().String()
	for {
		conn, err = listener.Accept()
		if err != nil {
			t.ErrChan <- err
		}
		t.Log.Debug("accepted connection")
		go t.forward(conn, sshServer, originServer)
	}
}

func (t *SSHTunnel) WithUser(u string) *SSHTunnel {
	t.Config.User = u
	return t
}

func (t *SSHTunnel) Wait() (string, error) {
	select {
	case err := <-t.ErrChan:
		return "", err
	case localAddress := <-t.AddressChan:
		return localAddress, nil
	}
}

func (t *SSHTunnel) forward(localConn net.Conn, sshServer, originServer *mytypes.Endpoint) {
	t.Log.Debugf("forwording connection on ssh tunnel %s", sshServer.String())
	t.Log.Debugf("using ssh configuration %+v", t.Config)
	serverConn, err := ssh.Dial("tcp", sshServer.String(), t.Config)
	if err != nil {
		t.Log.Errorf("server dial error: %s", err)
		return
	}
	t.Log.Debugf("connected to %s (1 of 2)\n", sshServer.String())
	remoteConn, err := serverConn.Dial("tcp", originServer.String())
	if err != nil {
		t.Log.Errorf("remote dial error: %s", err)
		return
	}
	t.Log.Debugf("connected to %s (2 of 2)\n", originServer.String())
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				t.Log.Infof("io.Copy error: %s", err)
			}
			t.Log.Errorf("io.Copy error: %s", err)
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
