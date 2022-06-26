package config

import (
	"os"
	"path"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	dir, _ := os.Getwd()
	os.Setenv("CTT_LOG_LEVEL", "DEBUG")
	os.Setenv("CTT_CONFIG", path.Join(dir, "for-tests", "config.yaml"))
	config, err := LoadConfig()
	assert.Nil(t, err)

	assert.Equal(t, config.LogLevel, logrus.DebugLevel)
	assert.Equal(t, "kubectl-path", config.KubeCtl.CliPath)
	assert.Equal(t, "path/id_rsa", config.SshConfig.KeyPath)
	assert.Equal(t, []string{"create"}, config.KubeCtl.SupportedSubCommands)
}

func TestLoadConfigErrorFlowCouldNotOpen(t *testing.T) {
	dir, _ := os.Getwd()
	os.Setenv("CTT_CONFIG", path.Join(dir, "..", "for-tests", "not-existing-file.yaml"))
	_, err := LoadConfig()
	assert.ErrorIs(t, err, ErrReadingConfig)
}
