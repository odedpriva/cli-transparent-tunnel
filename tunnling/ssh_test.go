package tunnling

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"os"
	"path"
	"testing"
)

func Test_privateKeyFile(t *testing.T) {
	dir, _ := os.Getwd()
	testAssets := path.Join(dir, "for-tests")
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    ssh.AuthMethod
		wantErr bool
	}{
		{
			name: "happy-flow",
			args: args{
				file: path.Join(testAssets, "id_rsa"),
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "error-flow",
			args: args{
				file: "non-exist",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := privateKeyFile(tt.args.file)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
