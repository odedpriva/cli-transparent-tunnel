package main

import (
	"github.com/odedpriva/cli-transparent-tunnel/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseOsArgs(t *testing.T) {
	type args struct {
		osArgs []string
	}
	tests := []struct {
		name        string
		args        args
		wantCommand string
		wantArgs    []string
		wantErr     bool
	}{
		{
			name: "no sub command",
			args: args{
				osArgs: []string{"ctt"},
			},
			wantCommand: "",
			wantArgs:    []string(nil),
			wantErr:     true,
		},
		{
			name: "no args",
			args: args{
				osArgs: []string{"ctt", "kubectl"},
			},
			wantCommand: "kubectl",
			wantArgs:    []string{},
			wantErr:     false,
		},
		{
			name: "with args",
			args: args{
				osArgs: []string{"ctt", "kubectl", "get", "pods"},
			},
			wantCommand: "kubectl",
			wantArgs:    []string{"get", "pods"},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCommand, gotArgs, err := parseOsArgs(tt.args.osArgs)
			assert.True(t, utils.AssertError(tt.wantErr, err))
			assert.Equal(t, gotArgs, tt.wantArgs)
			assert.Equal(t, gotCommand, tt.wantCommand)
		})
	}
}
