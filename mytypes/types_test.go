package mytypes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertToEndpointWithDefault(t *testing.T) {
	type args struct {
		s     string
		proto string
	}
	tests := []struct {
		name string
		args args
		want *Endpoint
	}{
		{
			name: "ssh",
			args: args{
				s:     "tcptunnel@kubernetes-docker-internal@kubernetes-docker-internal.tcp.symchatbotdemo.luminatesite.com",
				proto: "ssh",
			},
			want: &Endpoint{
				Host: "kubernetes-docker-internal.tcp.symchatbotdemo.luminatesite.com",
				User: "tcptunnel@kubernetes-docker-internal",
				Port: 22,
			},
		},
		{
			name: "https",
			args: args{
				s:     "kubernetes.docker.internal",
				proto: "https",
			},
			want: &Endpoint{
				Host: "kubernetes.docker.internal",
				User: "",
				Port: 443,
			},
		},
		{
			name: "given port",
			args: args{
				s:     "kubernetes.docker.internal:6443",
				proto: "https",
			},
			want: &Endpoint{
				Host: "kubernetes.docker.internal",
				User: "",
				Port: 6443,
			},
		},
		{
			name: "custom ssh port",
			args: args{
				s:     "tcptunnel@kubernetes-docker-internal@kubernetes-docker-internal.tcp.symchatbotdemo.luminatesite.com:2222",
				proto: "ssh",
			},
			want: &Endpoint{
				Host: "kubernetes-docker-internal.tcp.symchatbotdemo.luminatesite.com",
				User: "tcptunnel@kubernetes-docker-internal",
				Port: 2222,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToEndpointWithDefault(tt.args.s, tt.args.proto)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestConvertToEndpoint(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *Endpoint
	}{
		{
			name: "no-port",
			args: args{
				s: "user@user@user@host",
			},
			want: &Endpoint{
				Host: "host",
				User: "user@user@user",
				Port: 0,
			},
		},
		{
			name: "with-port",
			args: args{
				s: "user@user@user@host:22",
			},
			want: &Endpoint{
				Host: "host",
				User: "user@user@user",
				Port: 22,
			},
		},
	}
	for _, tt := range tests {
		got := ConvertToEndpoint(tt.args.s)
		assert.Equal(t, got, tt.want)
	}
}
