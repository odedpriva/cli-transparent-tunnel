package service

import (
	"fmt"
	command_runner "github.com/odedpriva/cli-transparent-tunnel/command-runner"
	command_tunneler "github.com/odedpriva/cli-transparent-tunnel/command-tunneler"
	"github.com/odedpriva/cli-transparent-tunnel/config"
	"github.com/odedpriva/cli-transparent-tunnel/mytypes"
	"github.com/odedpriva/cli-transparent-tunnel/tunnling"
	"github.com/sirupsen/logrus"
	"os"
)

type Service struct {
	command_tunneler.TunnelerCommand
	command_runner.CommandRunner
	Log       *logrus.Logger
	SshConfig *config.SshConfig
}

func (s *Service) Run() {
	var command string
	var args []string
	tunnelConfig, _ := s.GetTunnelConfiguration()
	if tunnelConfig != nil {
		sshTunnelServer := mytypes.ConvertToEndpointWithDefault(tunnelConfig.TunnelServer, "ssh")
		originEndpoint := mytypes.ConvertToEndpointWithDefault(tunnelConfig.OriginServer, "https")
		targetMachine := originEndpoint.Host
		sshTunnel, err := tunnling.NewSSHTunnel(s.SshConfig, s.Log)
		if err != nil {
			fmt.Printf("failed creating local tunnel listener %s", err)
			os.Exit(1)
		}
		go sshTunnel.WithUser(sshTunnelServer.User).Start(sshTunnelServer, originEndpoint)
		localAddress, err := sshTunnel.Wait()
		if err != nil {
			fmt.Printf("failed creating local tunnel listener %s", err)
			os.Exit(1)
		}
		command, args = s.GetCommandWithTunnel(localAddress, targetMachine)
	} else {
		command, args = s.GetPlainCommand()
	}

	s.RunCommand(command, args)

}
