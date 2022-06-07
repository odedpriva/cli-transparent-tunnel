package main

import (
	"ctt/config"
	"ctt/kubectl"
	"ctt/tunnling"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {

	mylog := log.New()

	c, err := config.LoadConfig()
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	k := kubectl.NewKubeCtl(c.KubeCtl, mylog, os.Args)
	mylog.SetLevel(c.LogLevel)

	appName := k.GetTunnelConfiguration()
	tunnelConfigurations := c.KubeCtl.TunnelConfigurations

	var localAddress string
	localAddressChan := make(chan string)
	errorChan := make(chan error)
	var targetMachine string
	if val, ok := tunnelConfigurations[appName]; ok {
		targetMachine = val.TargetMachine
		sshTunnel := tunnling.NewSSHTunnel(val.ExternalUserName, c.SshConfig, mylog)
		sshTunnel.SetEndpoints(fmt.Sprintf("%s@%s", val.ExternalUserName, val.ExternalHostName),
			fmt.Sprintf("%s:%s", val.TargetMachine, val.TargetPort),
		)
		go sshTunnel.Start(localAddressChan, errorChan)
		select {
		case err = <-errorChan:
			fmt.Errorf("failed creating local tunnel listener %s", err)
			os.Exit(1)
		case localAddress = <-localAddressChan:
		}
	}
	output, err := k.RunCommand(localAddress, targetMachine)
	fmt.Printf("%s", output)
	if err != nil {
		os.Exit(1)
	}
}
