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
