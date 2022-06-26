package command_runner

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"os/exec"
)

type CommandRunner interface {
	RunCommand(command string, args []string)
}
type CommandRunnerImpl struct {
	log *logrus.Logger
}

func NewCommandRunnerImpl(log *logrus.Logger) *CommandRunnerImpl {
	return &CommandRunnerImpl{log: log}
}

func (c *CommandRunnerImpl) RunCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	c.log.Debugf("%s", cmd.String())

	stdout, err := cmd.StdoutPipe()
	checkError(err)
	stderr, err := cmd.StderrPipe()
	checkError(err)

	err = cmd.Start()
	checkError(err)

	defer func() {
		var exerr *exec.ExitError
		err = cmd.Wait()
		if err != nil {
			if errors.As(err, &exerr) {
				os.Exit(exerr.ExitCode())
			}
			c.log.Errorf("error waiting for the command %s", err)
			os.Exit(1)
		}
	}()

	go func() {
		_, err = io.Copy(os.Stdout, stdout)
		if err != nil {
			c.log.Errorf("error copying to stdout %s", err)
		}
	}()
	go func() {
		_, err = io.Copy(os.Stderr, stderr)
		if err != nil {
			c.log.Errorf("error copying to stdout %s", err)
		}
	}()
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
