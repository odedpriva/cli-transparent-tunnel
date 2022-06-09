package kubectl

import (
	"github.com/odedpriva/cli-transparent-tunnel/config"
	"os/exec"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/cobra"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubectl/pkg/cmd"
)

type KubeCtl struct {
	log                 *logrus.Logger
	command             *cobra.Command
	cliConfig           *config.KubeCtl
	eligibleSubcommands []string
	osArgs              []string
}

func NewKubeCtl(ctl *config.KubeCtl, log *logrus.Logger, args []string) *KubeCtl {
	command := cmd.NewDefaultKubectlCommand()
	err := command.ParseFlags(args)
	if err != nil {
		log.Errorf("%s", err)
	}

	return &KubeCtl{
		log:                 log,
		eligibleSubcommands: ctl.SupportedSubCommands,
		command:             command,
		osArgs:              args,
		cliConfig:           ctl,
	}
}

func (k *KubeCtl) GetTunnelConfiguration() string {

	args := k.command.Flags().Args()

	if len(args) <= 1 {
		k.log.Debugf("no args for command line")
		return ""
	}

	if !slices.Contains(k.eligibleSubcommands, args[1]) {
		k.log.Debugf("sub command %s not supported for tunnel %s", args[1], k.eligibleSubcommands)
		return ""
	}

	kubeConfig := k.getKubeConfig()

	return k.findClusterName(kubeConfig)

}

func (k *KubeCtl) RunCommand(tunnelAddress string, originalServer string) ([]byte, error) {
	cmd := exec.Command(k.cliConfig.CliPath)
	cmd.Args = k.osArgs
	if tunnelAddress != "" {
		cmd.Args = append([]string{cmd.Args[0], "--server", "https://" + tunnelAddress, "--tls-server-name", originalServer}, cmd.Args[1:len(cmd.Args)]...)
	}
	k.log.Debugf("%s", cmd.String())
	return cmd.CombinedOutput()
}

func (k *KubeCtl) getKubeConfig() *api.Config {

	c, err := clientcmd.NewDefaultClientConfigLoadingRules().GetStartingConfig()
	if err != nil {
		k.log.Fatal("could not load config file")
	}

	return c

}

func (k *KubeCtl) findClusterName(kubeConfig *api.Config) string {

	if val := k.command.Flags().Lookup("cluster").Value.String(); val != "" {
		k.log.Debugf("using cluster %s", val)
		return val
	}

	contextName := k.command.Flags().Lookup("context").Value.String()
	if contextName == "" {
		contextName = kubeConfig.CurrentContext
	}

	k.log.Debugf("using context %s", contextName)
	if context, ok := kubeConfig.Contexts[contextName]; ok {
		k.log.Debugf("using server %s", context.Cluster)
		return context.Cluster
	}

	return ""
}
