package config

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/odedpriva/cli-transparent-tunnel/utils"
)

func TestLoadConfigV2(t *testing.T) {
	dir, _ := os.Getwd()
	tests := []struct {
		name             string
		want             *Config
		environmentToSet map[string]string
		wantErr          bool
	}{
		{
			name: "happy-flow",
			want: &Config{
				//LogLevel: logrus.InfoLevel,
				CliConfigurations: map[string]CliConfig{
					"redis-cli": {
						CliPath: "/usr/local/bin/redis-cli",
						FlagsConfig: FlagsConfig{
							Host: []string{"-h"},
							Port: []string{"-p"},
							SNI:  []string{"--sni"},
						},
					},
					"psql": {
						CliPath: "/usr/local/bin/psql",
						FlagsConfig: FlagsConfig{
							Host: []string{"--host", "-h"},
							Port: []string{"-p", "--port"},
						},
					},
					"kubectl": {
						CliPath: "/usr/local/bin/kubectl",
						FlagsConfig: FlagsConfig{
							Address: []string{"--server", "-s"},
							SNI:     []string{"--tls-server-name"},
						},
					},
					"oc": {
						CliPath: "/usr/local/bin/oc",
						FlagsConfig: FlagsConfig{
							Address: []string{"--server", "-s"},
							SNI:     []string{"--tls-server-name"},
						},
					},
				},
				TunnelConfigurations: map[string][]TunnelConfiguration{
					"redis-cli": {
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "redis-eu",
						},
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "redis-us",
						},
					},
					"psql": {
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "psql-eu",
						},
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "psql-us",
						},
					},
					"kubectl": {
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "kind-ctt",
						},
					},
					"oc": {
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "oc-ctt",
						},
					},
				},
				SSHConfigurations: &SshConfig{
					KeyPath: "for-tests/scripts/id_rsa",
				},
			},
			environmentToSet: map[string]string{
				"CTT_TUNNEL_CONFIGURATIONS": path.Join(dir, "for-tests", "ctt-tunnel-configurations.yaml"),
				"CTT_CLI_CONFIGURATIONS":    path.Join(dir, "for-tests", "ctt-cli-configurations.yaml"),
			},
			wantErr: false,
		},
		{
			name: "",
			want: &Config{
				//LogLevel: logrus.DebugLevel,
				CliConfigurations: map[string]CliConfig{
					"redis-cli": {
						CliPath: "/usr/local/bin/redis-cli",
						FlagsConfig: FlagsConfig{
							Host: []string{"-h"},
							Port: []string{"-p"},
							SNI:  []string{"--sni"},
						},
					},
					"psql": {
						CliPath: "/usr/local/bin/psql",
						FlagsConfig: FlagsConfig{
							Host: []string{"--host", "-h"},
							Port: []string{"-p", "--port"},
						},
					},
					"kubectl": {
						CliPath: "/usr/local/bin/kubectl",
						FlagsConfig: FlagsConfig{
							Address: []string{"--server", "-s"},
							SNI:     []string{"--tls-server-name"},
						},
					},
					"oc": {
						CliPath: "/usr/local/bin/oc",
						FlagsConfig: FlagsConfig{
							Address: []string{"--server", "-s"},
							SNI:     []string{"--tls-server-name"},
						},
					},
				},
				TunnelConfigurations: map[string][]TunnelConfiguration{
					"redis-cli": {
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "redis-eu",
						},
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "redis-us",
						},
					},
					"psql": {
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "psql-eu",
						},
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "psql-us",
						},
					},
					"kubectl": {
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "kind-ctt",
						},
					},
					"oc": {
						{
							TunnelServer: "linuxserver.io@localhost:2222",
							OriginServer: "192.168.68.109:53077",
							Name:         "oc-ctt",
						},
					},
				},
				SSHConfigurations: &SshConfig{
					KeyPath: "for-tests/scripts/id_rsa",
				},
			},
			environmentToSet: map[string]string{
				"CTT_TUNNEL_CONFIGURATIONS": path.Join(dir, "for-tests", "ctt-tunnel-configurations.yaml"),
				"CTT_CLI_CONFIGURATIONS":    path.Join(dir, "for-tests", "ctt-cli-configurations.yaml"),
				"CTT_DEBUG":                 "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for s, s2 := range tt.environmentToSet {
				os.Setenv(s, s2)
			}
			got, err := LoadConfig()
			require.Truef(t, utils.AssertError(tt.wantErr, err), fmt.Sprintf("wanterr: %t got: %s", tt.wantErr, err))
			assert.Equal(t, tt.want.SSHConfigurations.KeyPath, got.SSHConfigurations.KeyPath)
			assert.Equal(t, tt.want.CliConfigurations, got.CliConfigurations)
			assert.Equal(t, tt.want.TunnelConfigurations, got.TunnelConfigurations)
			//assert.Equal(t, tt.want.LogLevel, got.LogLevel)
		})
	}
}
