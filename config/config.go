package config

import (
	"bufio"
	"bytes"
	"ctt/utils"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var ErrConfigNotExist = errors.New("config not exist")
var ErrReadingConfig = errors.New("error reading config")

type TunnelConfiguration struct {
	TunnelServer string `mapstructure:"ssh-tunnel-server"`
	OriginServer string `mapstructure:"origin-server"`
}

type KubeCtl struct {
	CliPath              string                         `mapstructure:"path"`
	TunnelConfigurations map[string]TunnelConfiguration `mapstructure:"tunnel-configurations"`
	EligibleSubCommands  []string                       `mapstructure:"eligible-subcommands"`
}

type SshConfig struct {
	KeyPath string `mapstructure:"key-path"`
}

type Config struct {
	SshConfig   *SshConfig   `mapstructure:"ssh-config"`
	KubeCtl     *KubeCtl     `mapstructure:"kubectl"`
	LogLevel    logrus.Level `mapstructure:"log-level"`
	CommandName string
}

func LoadConfig() (*Config, error) {

	viper.SetEnvPrefix("CTT")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	if err := loadConfigFile(); err != nil {
		return nil, err
	}

	config := newConfig()

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func InitConfig() error {
	var err error
	configFileLocation := path.Join(utils.GetHomeDir(), ".ctt/config.yaml")
	if os.Getenv("CTT_CONFIG") != "" {
		configFileLocation = os.Getenv("CTT_CONFIG")
	}

	viper.SetConfigFile(configFileLocation)

	//c := newConfig()
	sshKeyPath, err := getUserInput("which ssh-key to use?", path.Join(utils.GetHomeDir(), ".ssh/id_rsa"))
	if err != nil {
		return err
	}
	viper.Set("ssh-config.key-path", sshKeyPath)
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		kubectlPath = ""
	}
	kubectlPath, err = getUserInput("which kubectl to use?", kubectlPath)
	if err != nil {
		return err
	}

	viper.Set("kubectl.path", sshKeyPath)
	viper.Set("kubectl.eligible-subcommands", kubectlSupportedSubCommands())
	viper.Set("kubectl.tunnel-configurations", map[string]TunnelConfiguration{"example": exampleTunnelConfiguration()})

	err = os.MkdirAll(filepath.Dir(configFileLocation), os.ModePerm)
	if err != nil {
		return err
	}

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func newConfig() *Config {

	l, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		l = logrus.InfoLevel
	}
	return &Config{
		SshConfig: &SshConfig{},
		KubeCtl: &KubeCtl{
			TunnelConfigurations: map[string]TunnelConfiguration{},
		},
		LogLevel:    l,
		CommandName: "kubectl",
	}
}

func loadConfigFile() error {
	viper.SetConfigType("yaml")

	if val := os.Getenv("CTT_CONFIG"); val != "" {
		return readConfigFile(val)
	}

	if utils.IsFileExists(path.Join(utils.GetHomeDir(), ".ctt/config.yaml")) {
		return readConfigFile(path.Join(utils.GetHomeDir(), ".ctt/config.yaml"))
	}

	return ErrConfigNotExist

}

func readConfigFile(configLocation string) error {

	con, err := os.ReadFile(configLocation)
	if err != nil {
		return fmt.Errorf("%w %s", ErrReadingConfig, err)
	}
	err = viper.ReadConfig(bytes.NewBuffer(con))
	if err != nil {
		return fmt.Errorf("%w %s", ErrReadingConfig, err)
	}

	return nil
}

func getUserInput(param, defaultValue string) (string, error) {
	fmt.Printf("%s (%s): ", param, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error occurred while reading input. Please try again %s", err)
	}

	if strings.TrimSuffix(input, "\n") == "" {
		return defaultValue, nil
	}
	return input, nil
}

func kubectlSupportedSubCommands() []string {
	return []string{
		"create",
		"expose",
		"run",
		"set",
		"explain",
		"get",
		"delete",
		"rollout",
		"scale",
		"certificate",
		"cluster-info",
		"cordon",
		"uncordon",
		"drain",
		"taint",
		"describe",
		"logs",
		"auth",
		"apply",
		"patch",
		"replace",
	}
}

func exampleTunnelConfiguration() TunnelConfiguration {
	return TunnelConfiguration{
		TunnelServer: "<my-user>@<my-server.com>",
		OriginServer: "<k8s-api-endpoint>:<port>",
	}
}
