package service

import (
	"fmt"
	"github.com/odedpriva/cli-transparent-tunnel/logging"

	command_runner "github.com/odedpriva/cli-transparent-tunnel/command-runner"
	command_tunneler "github.com/odedpriva/cli-transparent-tunnel/command-tunneler"
	"github.com/odedpriva/cli-transparent-tunnel/mytypes"
	"github.com/odedpriva/cli-transparent-tunnel/network_tunnler"
)

type ServiceError struct {
	err     error
	message string
}

func (s ServiceError) Error() string {
	return fmt.Sprintf("%s %s", s.message, s.err)
}

type CommandError struct {
	err     error
	message string
}

func (c CommandError) Error() string {
	return fmt.Sprintf("%s %s", c.message, c.err)
}

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
		s.log.Errorf("failed creating local tunnel listener %s", err)
		return ServiceError{err: err, message: fmt.Sprint("failed creating ssh tunnel")}
	}
	command, args, err := s.GetCommandWithTunnel(localAddress, input.ApplicationEndpoint.String())
	if err != nil {
		return err
	}
	s.RunCommand(command, args)
	return nil
}
