package convertor

import (
	"fmt"
	"os"
	"strings"

	"github.com/odedpriva/cli-transparent-tunnel/config"
	"github.com/odedpriva/cli-transparent-tunnel/mytypes"
	"github.com/odedpriva/cli-transparent-tunnel/service"
)

type Converter interface {
	GetServiceInput(commandName, tunnelConfigName string, commandArgs []string, c *config.Config) (*service.RunInput, error)
}

type ConvertImpl struct {
}

func NewConvertImpl() *ConvertImpl {
	return &ConvertImpl{}
}

func (i *ConvertImpl) GetServiceInput(commandName, tunnelConfigName string, commandArgs []string, c *config.Config) (*service.RunInput, error) {

	err := i.validateConfig(commandName, c)
	if err != nil {
		return nil, err
	}

	val, ok := c.TunnelConfigurations[commandName]
	if !ok {
		return nil, fmt.Errorf("tunnel configurations %+v does not include command %s", c.TunnelConfigurations, commandName)
	}

	cliConfig, ok := c.CliConfigurations[commandName]
	if !ok {
		return nil, fmt.Errorf("cli configurations %+v does not include command %s", c.TunnelConfigurations, commandName)
	}

	tunnelConfig, ok := i.isTunnelConfigurationExist(tunnelConfigName, val)
	if !ok {
		return nil, fmt.Errorf("tunnel config for comand %s does not include %s", commandName, tunnelConfigName)
	}

	sshEndpoint := mytypes.ConvertToSshEndpoint(tunnelConfig.TunnelServer)
	applicationEndpoint := mytypes.ConvertTApplicationEndpoint(tunnelConfig.OriginServer)

	sshConf := &mytypes.SshConfig{
		KeyPath: c.SSHConfigurations.KeyPath,
		User:    sshEndpoint.User,
	}

	return &service.RunInput{
		SshConfig:           sshConf,
		Command:             cliConfig.CliPath,
		CommandArgs:         commandArgs,
		SshEndpoint:         sshEndpoint,
		ApplicationEndpoint: applicationEndpoint,
	}, nil

}

func (i *ConvertImpl) isTunnelConfigurationExist(name string, t []config.TunnelConfiguration) (config.TunnelConfiguration, bool) {
	for i := range t {
		if t[i].Name == name {
			return t[i], true
		}
	}
	return config.TunnelConfiguration{}, false
}

func (i *ConvertImpl) validateConfig(commandName string, config *config.Config) error {

	var errors []string

	if _, ok := config.CliConfigurations[commandName]; !ok {
		errors = append(errors, fmt.Sprintf("cli configuration does not include %s", commandName))
	}

	command := config.CliConfigurations[commandName]
	_, err := os.Stat(command.CliPath)
	if err != nil {
		errors = append(errors, fmt.Sprintf("cli path error %s", err))
	}

	if _, ok := config.TunnelConfigurations[commandName]; !ok {
		errors = append(errors, fmt.Sprintf("tunnel configuration does not include %s", commandName))
	}

	if config.SSHConfigurations.KeyPath != "" {
		_, err = os.Stat(config.SSHConfigurations.KeyPath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("ssh key error %s", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}

	return nil

}
