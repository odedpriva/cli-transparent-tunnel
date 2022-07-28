package mytypes

import (
	"os/user"
	"strings"
	"testing"

	"github.com/pkg/errors"
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
		err     error
	}{
		{
			name: "",
			args: args{
				s: "user@host:2222",
			},
			want: &NetworkEndpoint{
				Host: "host",
				User: "user",
				Port: "2222",
			},
			err: nil,
		},
		{
			name: "",
			args: args{
				s: "user@host",
			},
			want: &NetworkEndpoint{
				Host: "host",
				User: "user",
				Port: "22",
			},
			err: nil,
		},
		{
			name: "",
			args: args{
				s: "host",
			},
			want: &NetworkEndpoint{
				Host: "host",
				User: func() string { user, _ := user.Current(); return strings.ToLower(user.Username) }(),
				Port: "22",
			},
			err: nil,
		},
		{
			name: "",
			args: args{
				s: "whatever",
			},
			want: &NetworkEndpoint{
				Host: "whatever",
				User: func() string { user, _ := user.Current(); return strings.ToLower(user.Username) }(),
				Port: "22",
			},
			err: nil,
		},
		{
			name: "failed to convert to ssh endpoint",
			args: args{s:
				"bad://user@host:123",
			},
			want: nil,
			err: errors.New("unexpected scheme in bad://user@host:123"),
		},
		{
			name: "bad port",
			args: args{s:
				"user@host:port",
			},
			want: nil,
			err: errors.New("failed to parse user@host:port: parse \"ssh://user@host:port\": invalid port \":port\" after host"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToSshEndpoint(tt.args.s)
			assert.Equal(t, tt.want, got, t.Name())
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error(), t.Name())
			} else {
				assert.NoError(t, err, t.Name())
			}
		})
	}
}

//func Test_getFlag(t *testing.T) {
//	type args struct {
//		commandLineArgs []string
//		flags           []string
//	}
//	tests := []struct {
//		name      string
//		args      args
//		wantArg   string
//		wantValid bool
//	}{
//		{
//			name: "",
//			args: args{
//				commandLineArgs: []string{"--name", "test", "-h", "test"},
//				flags:           []string{"-h"},
//			},
//			wantArg:   "-h",
//			wantValid: true,
//		},
//		{
//			name: "",
//			args: args{
//				commandLineArgs: []string{"--name", "test", "--host", "test"},
//				flags:           []string{"--address"},
//			},
//			wantArg:   "",
//			wantValid: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotArg, gotValid := getFlag(tt.args.commandLineArgs, tt.args.flags)
//			assert.Equalf(t, tt.wantArg, gotArg, "getFlag(%v, %v)", tt.args.commandLineArgs, tt.args.flags)
//			assert.Equalf(t, tt.wantValid, gotValid, "getFlag(%v, %v)", tt.args.commandLineArgs, tt.args.flags)
//		})
//	}
//}