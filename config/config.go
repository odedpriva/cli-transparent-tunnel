package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var ErrConfigNotExist = errors.New("config not exist")
var ErrReadingConfig = errors.New("error reading config")

type TunnelConfiguration struct {
	TunnelServer string `mapstructure:"ssh-tunnel-server"`
	OriginServer string `mapstructure:"origin-server"`
	Name         string `mapstructure:"name"`
}

type SshConfig struct {
	KeyPath string `mapstructure:"key-path"`
}

type FlagsConfig struct {
	Host    []string `mapstructure:"host"`
	Port    []string `mapstructure:"port"`
	Address []string `mapstructure:"address"`
	SNI     []string `mapstructure:"sni"`
}

type Cliconfig struct {
	CliPath     string      `mapstructure:"path"`
	FlagsConfig FlagsConfig `mapstructure:"flags"`
}

type Config struct {
	CliConfigurations    map[string]Cliconfig             `mapstructure:"commands-configuration"`
	TunnelConfigurations map[string][]TunnelConfiguration `mapstructure:"tunnel-configurations"`
	SSHConfigurations    *SshConfig                       `mapstructure:"ssh-config"`
}

func LoadConfig() (*Config, error) {

	viper.SetEnvPrefix("ctt")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	config := newConfig()

	if err := loadConfigFileV2("CTT_TUNNEL_CONFIGURATIONS", "tunnel-configurations"); err != nil {
		return nil, err
	}

	if err := loadConfigFileV2("CTT_CLI_CONFIGURATIONS", "cli-configurations"); err != nil {
		return nil, err
	}

	var Configurations map[string][]TunnelConfiguration
	Configurations = make(map[string][]TunnelConfiguration)

	if err := viper.UnmarshalKey("configurations", &Configurations); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	config.TunnelConfigurations = Configurations
	return config, nil
}

func loadConfigFileV2(environment, configName string) error {

	if val := os.Getenv(environment); val != "" {
		viper.SetConfigFile(val)
		return viper.MergeInConfig()
	}

	if configName != "" {
		viper.SetConfigName(configName)
		return viper.MergeInConfig()
	}

	return fmt.Errorf("failed loading config file")
}

func newConfig() *Config {
	return &Config{
		CliConfigurations:    map[string]Cliconfig{},
		TunnelConfigurations: map[string][]TunnelConfiguration{},
		SSHConfigurations: &SshConfig{
			KeyPath: "",
		},
	}
}


