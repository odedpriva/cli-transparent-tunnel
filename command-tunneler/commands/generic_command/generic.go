package generic_command

import (
	"fmt"
	"net"

	"github.com/odedpriva/cli-transparent-tunnel/logging"
	"github.com/odedpriva/cli-transparent-tunnel/mytypes"
	args_utils "github.com/odedpriva/cli-transparent-tunnel/utils/args-utils"
)

type GenericCommand struct {
	command       string
	args          []string
	log           *logging.Logging
	genericConfig *mytypes.Command
}

func NewGenericCommand(command string, args []string, genericConfig *mytypes.Command) *GenericCommand {
	log := logging.GetLogger()
	return &GenericCommand{command: command, args: args, log: log, genericConfig: genericConfig}
}

func (g *GenericCommand) GetCommandWithTunnel(scheme, tunnelAddress string, originalServer string) (string, []string, error) {

	tunnelHost, tunnelPort, err := net.SplitHostPort(tunnelAddress)
	if err != nil {
		return "", nil, err
	}

	originalHost, _, err := net.SplitHostPort(originalServer)
	if err != nil {
		return "", nil, err
	}

	if scheme != "" {
		tunnelAddress = fmt.Sprintf("%s://%s", scheme, tunnelAddress)
	}


	builder := args_utils.NewArgsBuilder(g.args)
	args := builder.
		ReplaceOrAddArgForFlags("host-flags", g.genericConfig.HostFlags, tunnelHost).
		ReplaceOrAddArgForFlags("port-flags", g.genericConfig.PortFlags, tunnelPort).
		ReplaceOrAddArgForFlags("sni-flags", g.genericConfig.SNIFlags, originalHost).
		ReplaceOrAddArgForFlags("address-flags", g.genericConfig.AddressFlags, tunnelAddress).
		Build()
	return g.command, args.GetArgs(), nil

}
