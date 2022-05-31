package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extractKubectlConfig(t *testing.T) {
	type args struct {
		tunnelAddress string
		k8sFilePath   string
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		{
			name: "",
			args: args{
				tunnelAddress: "https://localhost:6465",
				k8sFilePath:   "for-tests/k8s_config.yaml",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractKubectlConfig(tt.args.tunnelAddress, tt.args.k8sFilePath)
			assert.Nil(t, err)
			t.Logf(got.Name())
		})
	}
}
