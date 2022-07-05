package mytypes

import (
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToSshEndpoint(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *NetworkEndpoint
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "user@host:2222",
			},
			want: &NetworkEndpoint{
				Host: "host",
				User: "user",
				Port: 2222,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "user@host",
			},
			want: &NetworkEndpoint{
				Host: "host",
				User: "user",
				Port: 22,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "host",
			},
			want: &NetworkEndpoint{
				Host: "host",
				User: func() string { user, _ := user.Current(); return user.Name }(),
				Port: 22,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "whatever",
			},
			want: &NetworkEndpoint{
				Host: "whatever",
				User: func() string { user, _ := user.Current(); return user.Name }(),
				Port: 22,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToSshEndpoint(tt.args.s)
			assert.Equal(t, tt.want, got, t.Name())
		})
	}
}

func Test_getFlag(t *testing.T) {
	type args struct {
		commandLineArgs []string
		flags           []string
	}
	tests := []struct {
		name      string
		args      args
		wantArg   string
		wantValid bool
	}{
		{
			name: "",
			args: args{
				commandLineArgs: []string{"--name", "test", "-h", "test"},
				flags:           []string{"-h"},
			},
			wantArg:   "-h",
			wantValid: true,
		},
		{
			name: "",
			args: args{
				commandLineArgs: []string{"--name", "test", "--host", "test"},
				flags:           []string{"--address"},
			},
			wantArg:   "",
			wantValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArg, gotValid := getFlag(tt.args.commandLineArgs, tt.args.flags)
			assert.Equalf(t, tt.wantArg, gotArg, "getFlag(%v, %v)", tt.args.commandLineArgs, tt.args.flags)
			assert.Equalf(t, tt.wantValid, gotValid, "getFlag(%v, %v)", tt.args.commandLineArgs, tt.args.flags)
		})
	}
}