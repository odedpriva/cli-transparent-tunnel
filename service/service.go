package service

import (
	"github.com/odedpriva/cli-transparent-tunnel/config"
	"github.com/pkg/errors"

	"github.com/odedpriva/cli-transparent-tunnel/logging"

	command_runner "github.com/odedpriva/cli-transparent-tunnel/command-runner"
	command_tunneler "github.com/odedpriva/cli-transparent-tunnel/command-tunneler"
	"github.com/odedpriva/cli-transparent-tunnel/mytypes"
	"github.com/odedpriva/cli-transparent-tunnel/network_tunnler"
)

type Service struct {
	command_tunneler.TunnelerCommand
	command_runner.CommandRunner
	network_tunnler.NetworkTunneler

	log *logging.Logging
}

func NewService(tunnelerCommand command_tunneler.TunnelerCommand, commandRunner command_runner.CommandRunner, network_tunneler network_tunnler.NetworkTunneler) *Service {
	log := logging.GetLogger()
	return &Service{
		TunnelerCommand: tunnelerCommand,
		CommandRunner:   commandRunner,
		NetworkTunneler: network_tunneler,
		log:             log,
	}
}

type RunInput struct {
	TunnelConfig        config.TunnelConfiguration
	SshConfig           *mytypes.SshConfig
	Command             string
	CommandArgs         []string
	SshEndpoint         *mytypes.NetworkEndpoint
	ApplicationEndpoint *mytypes.ApplicationEndpoint
}

func (s *Service) Run(input *RunInput) error {

	go s.Start(input.SshEndpoint, input.ApplicationEndpoint.String())
	localAddress, err := s.Wait()
	if err != nil {
		s.log.WithError(err).Errorf("failed creating local tunnel listener")
		return errors.Wrap(err, "failed creating ssh tunnel")
	}
	command, args, err := s.GetCommandWithTunnel(input.TunnelConfig.Scheme, localAddress, input.ApplicationEndpoint.String())
	if err != nil {
		return err
	}
	return s.RunCommand(command, args)
}
