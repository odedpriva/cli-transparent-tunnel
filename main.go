package main

import (
	"ctt/config"
	"ctt/kubectl"
	"ctt/mytypes"
	"ctt/tunnling"
	"fmt"
	"os"

	"errors"
	log "github.com/sirupsen/logrus"
)

func main() {

	if len(os.Args) > 1 {
		if os.Args[1] == "init" {
			file, err := config.InitConfig()
			if err != nil {
				fmt.Printf("error setting up ctt config %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("check config files: %s and update tunnel configurations\n", file)
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
	var targetMachine string
	if val, ok := config.IsTunnelConfigurationExist(appName, tunnelConfigurations); ok {
		sshTunnelServer := mytypes.ConvertToEndpointWithDefault(val.TunnelServer, "ssh")
		originEndpoint := mytypes.ConvertToEndpointWithDefault(val.OriginServer, "https")
		targetMachine = originEndpoint.Host
		sshTunnel, err := tunnling.NewSSHTunnel(c.SshConfig, mylog)
		if err != nil {
			fmt.Errorf("failed creating local tunnel listener %s", err)
			os.Exit(1)
		}
		go sshTunnel.WithUser(sshTunnelServer.User).Start(sshTunnelServer, originEndpoint)
		localAddress, err = sshTunnel.Wait()
		if err != nil {
			fmt.Errorf("failed creating local tunnel listener %s", err)
			os.Exit(1)
		}
	}
	output, err := k.RunCommand(localAddress, targetMachine)
	fmt.Printf("%s", output)
	if err != nil {
		os.Exit(1)
	}
}
