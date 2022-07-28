package convertor

import (
	"os"

	"github.com/pkg/errors"

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
		return nil, errors.Errorf("tunnel configurations %+v does not include command %s", c.TunnelConfigurations, commandName)
	}

	cliConfig, ok := c.CliConfigurations[commandName]
	if !ok {
		return nil, errors.Errorf("cli configurations %+v does not include command %s", c.TunnelConfigurations, commandName)
	}

	tunnelConfig, ok := i.isTunnelConfigurationExist(tunnelConfigName, val)
	if !ok {
		return nil, errors.Errorf("tunnel config for comand %s does not include %s", commandName, tunnelConfigName)
	}

	sshEndpoint, err := mytypes.ConvertToSshEndpoint(tunnelConfig.TunnelServer)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert %s to an endpoint", tunnelConfig.TunnelServer)
	}
	applicationEndpoint := mytypes.ConvertTApplicationEndpoint(tunnelConfig.OriginServer)

	sshConf := &mytypes.SshConfig{
		KeyPath: c.SSHConfigurations.KeyPath,
		User:    sshEndpoint.User,
	}

	return &service.RunInput{
		TunnelConfig:        tunnelConfig,
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

	var errorList []error

	if _, ok := config.CliConfigurations[commandName]; !ok {
		errorList = append(errorList, errors.Errorf("cli configuration does not include %s", commandName))
	}

	command := config.CliConfigurations[commandName]
	_, err := os.Stat(command.CliPath)
	if err != nil {
		errorList = append(errorList, errors.Wrapf(err, "cli path error %s", command.CliPath))
	}

	if _, ok := config.TunnelConfigurations[commandName]; !ok {
		errorList = append(errorList, errors.Errorf("tunnel configuration does not include %s", commandName))
	}

	if config.SSHConfigurations.KeyPath != "" {
		_, err = os.Stat(config.SSHConfigurations.KeyPath)
		if err != nil {
			errorList = append(errorList, errors.Wrapf(err, "ssh key error %s", config.SSHConfigurations.KeyPath))
		}
	}

	if len(errorList) > 0 {
		return errorList[0]
	}

	return nil

}
