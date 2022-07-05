package convertor

import (
	"fmt"
	"github.com/odedpriva/cli-transparent-tunnel/config"
	"github.com/odedpriva/cli-transparent-tunnel/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvertImpl_validateConfig(t *testing.T) {
	type args struct {
		commandName string
		config      *config.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				commandName: "test",
				config: &config.Config{
					CliConfigurations: map[string]config.Cliconfig{
						"test": {
							CliPath:     "/usr/bin/env",
							FlagsConfig: config.FlagsConfig{},
						},
					},
					TunnelConfigurations: map[string][]config.TunnelConfiguration{
						"test": {},
					},
					SSHConfigurations: &config.SshConfig{},
					LogLevel:          0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ConvertImpl{}
			err := i.validateConfig(tt.args.commandName, tt.args.config)
			require.True(t, utils.AssertError(tt.wantErr, err), fmt.Sprintf("wantErr: %t, got %+v", tt.wantErr, err))
		})
	}
}
