package command_tunneler

type TunnelerCommand interface {
	GetTunnelConfiguration() string
	RunCommand(tunnelAddress string, originalServer string) ([]byte, error)
}
