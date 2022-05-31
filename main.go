package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"sac-cli/config"
	"sac-cli/sac"
	"sac-cli/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/yaml"

	"golang.org/x/crypto/ssh"
)

var errChan chan error
var readyChan chan string
var wg = sync.WaitGroup{}
var localAddress = ""

func main() {

	sshKey := "/Users/odedpriva/dotfiles/secrets/ssh-keys/symchatbotdemo.pem"

	config, _ := config.GetConfig()
	sacSettings := &sac.SecureAccessCloudSettings{
		ClientID:     config.SacSettings.ClientID,
		ClientSecret: config.SacSettings.ClientSecret,
		TenantDomain: config.SacSettings.TenantDomain,
	}
	sacClient := sac.NewSecureAccessCloudClientImpl(sacSettings)
	app, err := sacClient.FindApplicationByName("kubernetes-docker-internal")
	utils.HandleError(err)

	fmt.Printf("%+v\n", app)
	sacExternalDomain := app.ConnectionSettings.ExternalAddress
	sacUser := fmt.Sprintf("tcptunnel@%s", app.ConnectionSettings.Subdomain)

	sacConnectionDomain := sacExternalDomain + ":22"

	//key, err := ioutil.ReadFile(sshKey)
	//utils.HandleError(err)

	tunnel := NewSSHTunnel(fmt.Sprintf("%s@%s", sacUser, sacConnectionDomain),
		PrivateKeyFile(sshKey),
		app.TcpTunnelSettings[0].Target+":6443",
	)
	go tunnel.Start()
	time.Sleep(100 * time.Millisecond)
	// Create the Signer for this private key.
	//signer, err := ssh.ParsePrivateKey(key)
	//utils.HandleError(err)

	//clientConfig := &ssh.ClientConfig{
	//	User: fmt.Sprintf("%s", sacUser),
	//	Auth: []ssh.AuthMethod{
	//		ssh.PublicKeys(signer),
	//	},
	//	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	//}
	//
	//fmt.Printf("establishing ssh client to %s@%s\n", sacUser, sacConnectionDomain)
	//clientConn, err := ssh.Dial("tcp", sacConnectionDomain, clientConfig)
	//utils.HandleError(err)
	//
	//wg.Add(1)
	//go sacTunnel(clientConn, "kubernetes.docker.internal:6443")
	//wg.Wait()
	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	fmt.Println(tunnel.Local.String())
	localAddress = tunnel.Local.String()
	runCommandOnTunnel(os.Args)

}

func sacTunnel(conn *ssh.Client, remote string) {
	pipe := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			log.Printf("failed to copy: %s", err)
		}
	}
	listener, err := net.Listen("tcp4", "127.0.0.1:")
	if err != nil {
		fmt.Errorf("error setting up tcp channel %s", err)
		os.Exit(1)
	}
	fmt.Printf("setup tcp listener on %s\n", listener.Addr().String())
	localAddress = listener.Addr().String()
	wg.Done()
	for {
		here, err := listener.Accept()
		if err != nil {
			log.Printf("failed to copy: %s", err)
		}
		go func(here net.Conn) {
			there, err := conn.Dial("tcp", remote)
			if err != nil {
				log.Printf("failed to copy: %s", err)
			}
			go pipe(there, here)
			go pipe(here, there)
		}(here)
	}
}

func extractKubectlConfig(tunnelAddress string, k8sFilePath string) (*os.File, error) {
	kubeFileStat, _ := os.Stat(k8sFilePath)

	kubeFile, _ := os.Open(k8sFilePath)

	filesize := kubeFileStat.Size()
	buffer := make([]byte, filesize)
	_, err := kubeFile.Read(buffer)
	if err != nil {
		return nil, err
	}

	c := api.NewConfig()
	err = yaml.Unmarshal(buffer, c)
	if err != nil {
		return nil, err
	}

	updateContextCluster(tunnelAddress, c.CurrentContext, c)

	updatedKubeFile, err := ioutil.TempFile(os.TempDir(), "tunnelAddress.*.yaml")
	if err != nil {
		return nil, err
	}

	d, err := yaml.Marshal(c)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(updatedKubeFile.Name(), d, kubeFileStat.Mode())
	if err != nil {
		return nil, err
	}
	return updatedKubeFile, nil
}

func runCommandOnTunnel(args []string) {
	cmd := exec.Command("kubectl")
	cmd.Args = append(args, "--server", "https://"+localAddress, "--insecure-skip-tls-verify")
	fmt.Println(cmd.String())
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed command %s", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", stdoutStderr)
}

func updateContextCluster(tunnelAddress string, context string, config *api.Config) {
	contextName := context

	for s := range config.Contexts {
		if s == contextName {
			cluster := config.Contexts[s].Cluster
			for s2 := range config.Clusters {
				if s2 == cluster {
					config.Clusters[cluster].TLSServerName = config.Clusters[cluster].Server
					config.Clusters[cluster].Server = tunnelAddress
				}
			}
		}
	}

}

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
	Log    *log.Logger
}

func (t *SSHTunnel) logf(fmt string, args ...interface{}) {
	if t.Log != nil {
		t.Log.Printf(fmt, args...)
	}
}

func (tunnel *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()
	tunnel.Local.Port = listener.Addr().(*net.TCPAddr).Port
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		tunnel.logf("accepted connection")
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		tunnel.logf("server dial error: %s", err)
		return
	}
	tunnel.logf("connected to %s (1 of 2)\n", tunnel.Server.String())
	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		tunnel.logf("remote dial error: %s", err)
		return
	}
	tunnel.logf("connected to %s (2 of 2)\n", tunnel.Remote.String())
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			tunnel.logf("io.Copy error: %s", err)
		}
	}
	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func PrivateKeyFile(file string) ssh.AuthMethod {
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

func NewSSHTunnel(tunnel string, auth ssh.AuthMethod, destination string) *SSHTunnel {
	// A random port will be chosen for us.
	localEndpoint := NewEndpoint("localhost:0")
	server := NewEndpoint(tunnel)
	if server.Port == 0 {
		server.Port = 22
	}
	sshTunnel := &SSHTunnel{
		Config: &ssh.ClientConfig{
			User: server.User,
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				// Always accept key.
				return nil
			},
		},
		Local:  localEndpoint,
		Server: server,
		Remote: NewEndpoint(destination),
	}
	return sshTunnel
}
