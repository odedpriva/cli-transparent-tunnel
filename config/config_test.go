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
	os.Setenv("CTT_CONFIG", path.Join(dir, "..", "for-tests", "config.yaml"))
	config, err := LoadConfig()
	assert.Nil(t, err)

	assert.Equal(t, config.LogLevel, logrus.DebugLevel)
	assert.Equal(t, config.KubeCtl.CliPath, "path")
	assert.Equal(t, config.KubeCtl.EligibleSubCommands, []string{"create"})
}

func TestLoadConfigErrorFlow(t *testing.T) {
	if os.Getenv("CI") != "true" {
		t.Skip("")
	}
	_, err := LoadConfig()
	assert.ErrorIs(t, err, ErrConfigNotExist)
}

func TestLoadConfigErrorFlowCouldNotOpen(t *testing.T) {
	dir, _ := os.Getwd()
	os.Setenv("CTT_CONFIG", path.Join(dir, "..", "for-tests", "not-existing-file.yaml"))
	_, err := LoadConfig()
	assert.ErrorIs(t, err, ErrReadingConfig)
}
