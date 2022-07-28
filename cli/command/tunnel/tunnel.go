package tunnel

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	command_runner "github.com/odedpriva/cli-transparent-tunnel/command-runner"
	"github.com/odedpriva/cli-transparent-tunnel/command-tunneler/commands/generic_command"
	"github.com/odedpriva/cli-transparent-tunnel/config"
	"github.com/odedpriva/cli-transparent-tunnel/mytypes"
	"github.com/odedpriva/cli-transparent-tunnel/network_tunnler"
	"github.com/odedpriva/cli-transparent-tunnel/service"
	"github.com/odedpriva/cli-transparent-tunnel/service/convertor"
	args_utils "github.com/odedpriva/cli-transparent-tunnel/utils/args-utils"
)

type TunnelCmd struct {
	commandArgs []string
	conf        *config.Config
}

func NewTunnelCmd(commandArgs []string, conf *config.Config) *TunnelCmd {
	return &TunnelCmd{commandArgs: commandArgs, conf: conf}
}

var tunnelConfigName string

func (c *TunnelCmd) GetCommnad() *cli.Command {
	return &cli.Command{
		Name:   "tunnel",
		Usage:  "",
		Action: c.runCommand,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tunnel-config",
				Value:       "",
				Usage:       "tunnel config name",
				Destination: &tunnelConfigName,
			},
		},
	}
}

func (c *TunnelCmd) runCommand(*cli.Context) error {

	commandName, commandArgs := args_utils.SplitToCommand(c.commandArgs)

	conv := convertor.NewConvertImpl()
	serviceInput, err := conv.GetServiceInput(commandName, tunnelConfigName, commandArgs, c.conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	command := newCommand(c.conf.CliConfigurations[commandName])
	commandTunnel := generic_command.NewGenericCommand(commandName, commandArgs, command)

	commandRunner := command_runner.NewCommandRunnerImpl()
	tunneler, err := network_tunnler.NewSSHTunnel(serviceInput.SshConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create new SSH tunnel")
	}
	myService := service.NewService(commandTunnel, commandRunner, tunneler)
	return myService.Run(serviceInput)
}

func newCommand(config config.CliConfig) *mytypes.Command {
	return &mytypes.Command{
		CliPath:      config.CliPath,
		HostFlags:    config.FlagsConfig.Host,
		PortFlags:    config.FlagsConfig.Port,
		AddressFlags: config.FlagsConfig.Address,
		SNIFlags:     config.FlagsConfig.SNI,
	}
}
