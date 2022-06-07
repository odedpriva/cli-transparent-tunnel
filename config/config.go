package config

import (
	"bytes"
	"ctt/utils"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type TunnelConfiguration struct {
	ExternalHostName string `mapstructure:"ssh-client-host-name"`
	ExternalUserName string `mapstructure:"ssh-client-user-name"`
	TargetMachine    string `mapstructure:"target-machine"`
	TargetPort       string `mapstructure:"target-port"`
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

	loadConfigFile()

	config := newConfig(viper.GetString("log_level"))

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func newConfig(logLevel string) *Config {
	l, err := logrus.ParseLevel(logLevel)
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

func loadConfigFile() {
	viper.SetConfigType("yaml")

	if val := os.Getenv("CTT_CONFIG"); val != "" {
		readConfigFile(val)
		return
	}

	if utils.IsFileExists(path.Join(utils.GetCWD(), "config.yaml")) {
		readConfigFile(path.Join(utils.GetCWD(), "config.yaml"))
		return
	}

	if utils.IsFileExists(path.Join(utils.GetHomeDir(), ".ctt/config.yaml")) {
		readConfigFile(path.Join(utils.GetHomeDir(), ".ctt/config.yaml"))
		return
	}

}

func readConfigFile(configLocation string) {

	con, err := os.ReadFile(configLocation)
	if err != nil {
		fmt.Errorf("error reading file %s %s", configLocation, err)
		os.Exit(1)
	}
	err = viper.ReadConfig(bytes.NewBuffer(con))
	if err != nil {
		fmt.Errorf("error reading file %s %s", configLocation, err)
		os.Exit(1)
	}
}
