package command_tunneler

type TunnelerCommand interface {
	GetCommandWithTunnel(tunnelAddress string, originalServer string) (string, []string, error)
}
