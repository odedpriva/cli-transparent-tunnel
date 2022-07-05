package main

import (
	"github.com/odedpriva/cli-transparent-tunnel/utils/args-utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_splitArgs(t *testing.T) {
	type args struct {
		args              []string
		availableCommands []string
	}
	tests := []struct {
		name            string
		args            args
		wantCttArgs     []string
		wantCommandName string
		wantCommandArgs []string
		wantErr         bool
	}{
		{
			name: "",
			args: args{
				args:              []string{"ctt"},
				availableCommands: []string{"redis-cli"},
			},
			wantCttArgs:     []string{"ctt"},
			wantCommandArgs: nil,
		},
		{
			name: "",
			args: args{
				args:              []string{"ctt", "setup"},
				availableCommands: []string{"redis-cli"},
			},
			wantCttArgs:     []string{"ctt", "setup"},
			wantCommandArgs: nil,
		},
		{
			name: "",
			args: args{
				args:              []string{"ctt", "redis-cli", "get", "keys"},
				availableCommands: []string{"redis-cli"},
			},
			wantCttArgs:     []string{"ctt"},
			wantCommandArgs: []string{"redis-cli", "get", "keys"},
		},
		{
			name: "",
			args: args{
				args:              []string{"ctt", "-t", "test", "redis-cli", "get", "keys"},
				availableCommands: []string{"redis-cli"},
			},
			wantCttArgs:     []string{"ctt", "-t", "test"},
			wantCommandArgs: []string{"redis-cli", "get", "keys"},
		},
		{
			name: "",
			args: args{
				args:              []string{"ctt", "-t", "test", "not-supported", "get", "keys"},
				availableCommands: []string{"redis-cli"},
			},
			wantCttArgs:     []string{"ctt", "-t", "test", "not-supported", "get", "keys"},
			wantCommandArgs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCttArgs, gotCommandArgs := args_utils.splitArgs(tt.args.args, tt.args.availableCommands)
			assert.Equalf(t, tt.wantCttArgs, gotCttArgs, "splitArgs(%v, %v)", tt.args.args, tt.args.availableCommands)
			assert.Equalf(t, tt.wantCommandArgs, gotCommandArgs, "splitArgs(%v, %v)", tt.args.args, tt.args.availableCommands)
		})
	}
}

func Test_splitToCommand(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name            string
		args            args
		wantCmd         string
		wantCommandArgs []string
	}{
		{
			name: "",
			args: args{
				args: []string{},
			},
			wantCmd:         "",
			wantCommandArgs: nil,
		},
		{
			name: "",
			args: args{
				args: []string{"test"},
			},
			wantCmd:         "test",
			wantCommandArgs: nil,
		},
		{
			name: "",
			args: args{
				args: []string{"test", "arg1"},
			},
			wantCmd:         "test",
			wantCommandArgs: []string{"arg1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotCommandArgs := splitToCommand(tt.args.args)
			assert.Equalf(t, tt.wantCmd, gotCmd, "splitToCommand(%v)", tt.args.args)
			assert.Equalf(t, tt.wantCommandArgs, gotCommandArgs, "splitToCommand(%v)", tt.args.args)
		})
	}
}
