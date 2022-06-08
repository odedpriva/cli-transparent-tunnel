package main

import (
	"ctt/config"
	"ctt/kubectl"
	"ctt/tunnling"
	"fmt"
	"os"

	"errors"
	log "github.com/sirupsen/logrus"
)

func main() {

	if len(os.Args) > 1 {
		if os.Args[1] == "init" {
			err := config.InitConfig()
			if err != nil {
				fmt.Printf("error setting up ctt config %s\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	mylog := log.New()

	c, err := config.LoadConfig()
	if err != nil {
		if errors.Is(err, config.ErrConfigNotExist) {
			fmt.Println("ctt requires config file, have you run ctt init?")
			os.Exit(1)
		}
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
		targetMachine = val.OriginServer
		sshTunnel := tunnling.NewSSHTunnel(val.ExternalUserName, c.SshConfig, mylog)
		sshTunnel.SetEndpoints(fmt.Sprintf("%s@%s", val.ExternalUserName, val.TunnelServer),
			fmt.Sprintf("%s:%s", val.OriginServer, val.TargetPort),
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
