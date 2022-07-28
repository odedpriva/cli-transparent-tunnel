package command_tunneler

type TunnelerCommand interface {
	GetCommandWithTunnel(scheme, tunnelAddress string, originalServer string) (string, []string, error)
}
