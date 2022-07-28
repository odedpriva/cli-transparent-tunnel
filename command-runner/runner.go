package command_runner

import (
	"io"
	"os"
	"os/exec"

	"github.com/odedpriva/cli-transparent-tunnel/logging"
)

type CommandRunner interface {
	RunCommand(command string, args []string) error
}
type CommandRunnerImpl struct {
	log *logging.Logging
}

func NewCommandRunnerImpl() *CommandRunnerImpl {
	log := logging.GetLogger()
	return &CommandRunnerImpl{log: log}
}

func (c *CommandRunnerImpl) RunCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	c.log.Debug(cmd.String())

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.log.WithError(err).Fatal("failed to get the Stdout of the subprocess")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.log.WithError(err).Fatal("failed to get the Stderr of the subprocess")
	}

	err = cmd.Start()
	if err != nil {
		c.log.WithError(err).Fatal("failed to start the subprocess")
	}

	go func() {
		_, err = io.Copy(os.Stdout, stdout)
		if err != nil {
			c.log.WithError(err).Errorf("error copying to stdout")
		}
	}()
	go func() {
		_, err = io.Copy(os.Stderr, stderr)
		if err != nil {
			c.log.WithError(err).Errorf("error copying to stderr")
		}
	}()

	return cmd.Wait()
}
