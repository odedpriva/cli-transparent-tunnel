package command_tunneler

import "github.com/odedpriva/cli-transparent-tunnel/config"

type TunnelerCommand interface {
	GetTunnelConfiguration() (*config.TunnelConfiguration, error)
	GetCommandWithTunnel(tunnelAddress string, originalServer string) (string, []string)
	GetPlainCommand() (string, []string)
}
