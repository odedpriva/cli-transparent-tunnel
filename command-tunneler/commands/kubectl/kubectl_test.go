package kubectl

import (
	"github.com/odedpriva/cli-transparent-tunnel/config"
	"github.com/odedpriva/cli-transparent-tunnel/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path"
	"runtime"
	"testing"
)

func Test_parseCommandLineArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    *kubectlArgs
		wantErr bool
	}{
		{
			name: "",
			args: args{
				args: []string{"--context", "c1", "--kubeconfig", "k1", "--cluster", "cl1"},
			},
			want: &kubectlArgs{
				Cluster:    "cl1",
				Kubeconfig: "k1",
				Context:    "c1",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				args: []string{"--kubeconfig", "k1"},
			},
			want: &kubectlArgs{
				Kubeconfig: "k1",
			},
			wantErr: false,
		},
		//{
		//	name: "not existing flag",
		//	args: args{
		//		args: []string{"--nonexisting", "k1"},
		//	},
		//	want:    &kubectlArgs{},
		//	wantErr: false,
		//},
		//{
		//	name: "existing and non existing flag",
		//	args: args{
		//		args: []string{"--kubeconfig", "k12", "nonexistingarg", "--test", "test"},
		//	},
		//	want: &kubectlArgs{
		//		Kubeconfig: "k12",
		//	},
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCommandLineArgs(tt.args.args)
			require.True(t, utils.AssertError(tt.wantErr, err), err)
			assert.Equal(t, tt.want.Cluster, got.Cluster)
			assert.Equal(t, tt.want.Context, got.Context)
			assert.Equal(t, tt.want.Kubeconfig, got.Kubeconfig)
		})
	}
}

func TestKubeCtl_getTunnelConfigurationName(t *testing.T) {
	testLog := logrus.New()
	_, filename, _, _ := runtime.Caller(0)
	currentDir := path.Dir(filename)
	type fields struct {
		log        *logrus.Logger
		cliConfig  *config.CliConfig
		osArgs     []string
		parsedArgs kubectlArgs
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "no args",
			fields: fields{
				log: testLog,
				cliConfig: &config.CliConfig{
					CliPath:     "",
					FlagsConfig: config.FlagsConfig{},
				},
				osArgs: []string{""},
				parsedArgs: kubectlArgs{
					Cluster:    "",
					Kubeconfig: "",
					Context:    "",
				},
			},
			want: "",
		},
		{
			name: "no args",
			fields: fields{
				log: testLog,
				cliConfig: &config.CliConfig{
					CliPath:     "",
					FlagsConfig: config.FlagsConfig{},
				},
				osArgs: []string{"notsupported-sub-commmand"},
				parsedArgs: kubectlArgs{
					Cluster:    "",
					Kubeconfig: "",
					Context:    "",
				},
			},
			want: "",
		},
		{
			name: "config-not-exist",
			fields: fields{
				log: testLog,
				cliConfig: &config.CliConfig{
					CliPath:     "",
					FlagsConfig: config.FlagsConfig{},
				},
				osArgs: []string{"test", "", ""},
				parsedArgs: kubectlArgs{
					Cluster:    "",
					Kubeconfig: "config-not-exist",
					Context:    "",
				},
			},
			want: "",
		},
		{
			name: "with-config-file",
			fields: fields{
				log: testLog,
				cliConfig: &config.CliConfig{
					CliPath:     "",
					FlagsConfig: config.FlagsConfig{},
				},
				osArgs: []string{"test"},
				parsedArgs: kubectlArgs{
					Cluster:    "",
					Kubeconfig: path.Join(currentDir, "k8s_config.yaml"),
					Context:    "",
				},
			},
			want: "cluster1",
		},
		{
			name: "with-config-file-and-context",
			fields: fields{
				log: testLog,
				cliConfig: &config.CliConfig{
					CliPath:     "",
					FlagsConfig: config.FlagsConfig{},
				},
				osArgs: []string{"test"},
				parsedArgs: kubectlArgs{
					Cluster:    "",
					Kubeconfig: path.Join(currentDir, "k8s_config.yaml"),
					Context:    "context2",
				},
			},
			want: "cluster2",
		},
		{
			name: "with-config-file-and-cluster",
			fields: fields{
				log: testLog,
				cliConfig: &config.CliConfig{
					CliPath:     "",
					FlagsConfig: config.FlagsConfig{},
				},
				osArgs: []string{"test"},
				parsedArgs: kubectlArgs{
					Cluster:    "cluster2",
					Kubeconfig: path.Join(currentDir, "k8s_config.yaml"),
					Context:    "",
				},
			},
			want: "cluster2",
		},
		{
			name: "with-config-file-and-cluster",
			fields: fields{
				log: testLog,
				cliConfig: &config.CliConfig{
					CliPath:     "",
					FlagsConfig: config.FlagsConfig{},
				},
				osArgs: []string{"test"},
				parsedArgs: kubectlArgs{
					Cluster:    "",
					Kubeconfig: path.Join(currentDir, "k8s_config.yaml"),
					Context:    "",
				},
			},
			want: "cluster1",
		},
		{
			name: "with-config-file-and-cluster",
			fields: fields{
				log: testLog,
				cliConfig: &config.CliConfig{
					CliPath:     "",
					FlagsConfig: config.FlagsConfig{},
				},
				osArgs: []string{"test"},
				parsedArgs: kubectlArgs{
					Cluster:    "",
					Kubeconfig: path.Join(currentDir, "k8s_config.yaml"),
					Context:    "non-existing",
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KubeCtl{
				log:        tt.fields.log,
				cliConfig:  tt.fields.cliConfig,
				osArgs:     tt.fields.osArgs,
				parsedArgs: tt.fields.parsedArgs,
			}
			assert.Equalf(t, tt.want, k.getTunnelConfigurationName(), "getTunnelConfigurationName()")
		})
	}
}
