//go:build integration
// +build integration

package kubectl_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path"
	"testing"
)

var binaryDirectory = os.Getenv("CTT_TEST_BIN_DIR")
var kubeConfig = os.Getenv("CTT_TEST_KUBECONFIG")
var cttConfig = os.Getenv("CTT_TEST_CTT_CONFIG")

func setup(t *testing.T) (nativeCommand, cttCommand *exec.Cmd, cleanup func()) {
	if binaryDirectory == "" || kubeConfig == "" || cttConfig == "" {
		t.Logf("missing one of %s, %s,%s did you set them up ?",
			"CTT_TEST_BIN_DIR", "CTT_TEST_KUBECONFIG", "CTT_TEST_CTT_CONFIG")
		t.FailNow()
	}
	kubectl := path.Join(binaryDirectory, "kubectl")

	nativeCommand = exec.Command(kubectl)
	nativeCommand.Env = append(nativeCommand.Env, fmt.Sprintf("%s=%s", "CTT_CONFIG", cttConfig), fmt.Sprintf("%s=%s", "KUBECONFIG", kubeConfig))

	ctt := path.Join(binaryDirectory, "ctt")
	cttCommand = exec.Command(ctt, "kubectl")
	cttCommand.Env = os.Environ()
	cttCommand.Env = append(cttCommand.Env, fmt.Sprintf("%s=%s", "CTT_CONFIG", cttConfig),
		fmt.Sprintf("%s=%s", "KUBECONFIG", kubeConfig),
		fmt.Sprintf("%s=%s", "CTT_LOG_LEVEL", "DEBUG"),
	)

	cleanup = func() {}
	return nativeCommand, cttCommand, cleanup
}

func TestBasicKubeCtl(t *testing.T) {

	tests := []struct {
		command               string
		expectedExistCode     int
		expectedOutputContain string
	}{
		{
			command:               "cluster-info",
			expectedExistCode:     0,
			expectedOutputContain: "Kubernetes control plane",
		},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			_, cttCommand, cleanup := setup(t)
			defer cleanup()
			cttCommand.Args = append(cttCommand.Args, tt.command)
			output, existCode := runCommand(t, cttCommand)
			assert.Equal(t, tt.expectedExistCode, existCode)
			t.Logf("command output %s", output)
			assert.Contains(t, output, tt.expectedOutputContain)
		})
	}

}

func runCommand(t *testing.T, cmd *exec.Cmd) (output string, exitCode int) {
	var stdOutput bytes.Buffer
	cmd.Stdout = &stdOutput
	cmd.Stderr = &stdOutput
	t.Log(cmd.String())
	if err := cmd.Start(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return stdOutput.String(), err.(*exec.ExitError).ExitCode()
		} else {
			t.Logf("cmd.Wait: %v", err)
		}
	}

	return stdOutput.String(), 0

}
