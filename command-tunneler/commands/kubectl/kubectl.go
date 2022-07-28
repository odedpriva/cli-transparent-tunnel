package kubectl

import (
	"github.com/alexflint/go-arg"
	"github.com/odedpriva/cli-transparent-tunnel/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type kubectlArgs struct {
	Cluster    string
	Kubeconfig string `arg:"--kubeconfig,env:KUBECONFIG"`
	Context    string
}

type KubeCtl struct {
	log        *logrus.Logger
	cliConfig  *config.CliConfig
	osArgs     []string
	parsedArgs kubectlArgs
}

func NewKubeCtl(ctl *config.CliConfig, log *logrus.Logger, args []string) (*KubeCtl, error) {
	parsedArgs := kubectlArgs{}

	p, err := arg.NewParser(arg.Config{
		//IgnoreUnknownArgs: true,
	}, &parsedArgs)
	if err != nil {
		return nil, err
	}
	err = p.Parse(args)
	if err != nil {
		return nil, err
	}

	return &KubeCtl{
		log:        log,
		osArgs:     args,
		cliConfig:  ctl,
		parsedArgs: parsedArgs,
	}, nil
}

func (k *KubeCtl) GetCommandWithTunnel(tunnelAddress string, originalServer string) (string, []string, error) {
	cmd := k.cliConfig.CliPath
	args := k.osArgs
	args = append([]string{args[0], "--server", "https://" + tunnelAddress, "--tls-server-name", originalServer}, args[1:]...)
	return cmd, args, nil
}

func (k *KubeCtl) GetPlainCommand() (string, []string) {
	return k.cliConfig.CliPath, k.osArgs
}

func (k *KubeCtl) getTunnelConfigurationName() string {

	args := k.osArgs

	if len(args) == 0 {
		k.log.Debugf("no args for command line")
		return ""
	}

	kubeConfig, err := k.getKubeConfig()
	if err != nil {
		k.log.WithError(err).Errorf("ctt could not load kubeconfig file")
		return ""
	}

	return k.findClusterName(kubeConfig)

}

func (k *KubeCtl) getKubeConfig() (*api.Config, error) {

	var c *api.Config
	var err error

	if k.parsedArgs.Kubeconfig != "" {
		c, err = clientcmd.LoadFromFile(k.parsedArgs.Kubeconfig)
		if err != nil {
			return nil, errors.Wrapf(err, "error loading file %s", k.parsedArgs.Kubeconfig)
		}
		return c, nil
	}

	c, err = clientcmd.NewDefaultClientConfigLoadingRules().GetStartingConfig()
	if err != nil {
		return nil, err
	}

	return c, nil

}

func parseCommandLineArgs(args []string) (*kubectlArgs, error) {

	parsedArgs := &kubectlArgs{}

	//p, err := arg.NewParser(arg.Config{IgnoreUnknownArgs: true}, parsedArgs)
	p, err := arg.NewParser(arg.Config{}, parsedArgs)
	if err != nil {
		return nil, err
	}
	err = p.Parse(args)
	if err != nil {
		return nil, err
	}
	return parsedArgs, nil

}

func (k *KubeCtl) findClusterName(kubeConfig *api.Config) string {

	if k.parsedArgs.Cluster != "" {
		return k.parsedArgs.Cluster
	}

	contextName := k.parsedArgs.Context
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
