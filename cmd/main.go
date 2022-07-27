package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/odedpriva/cli-transparent-tunnel/cli/command/setup"
	"github.com/odedpriva/cli-transparent-tunnel/cli/command/tunnel"
	"github.com/odedpriva/cli-transparent-tunnel/config"
	"github.com/odedpriva/cli-transparent-tunnel/logging"
	"github.com/odedpriva/cli-transparent-tunnel/utils/args-utils"
)

var debug bool

func main() {
	args := os.Args

	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	level := logrus.InfoLevel
	if debug {
		level = logrus.DebugLevel
	}

	_ = logging.NewLogger(level)

	cttAllArgs, commandAllArgs := args_utils.SplitArgs(args, conf.CliKeys())

	commands := []*cli.Command{
		setup.NewSetupCmd().GetCommnad(),
		tunnel.NewTunnelCmd(commandAllArgs, conf).GetCommnad(),
	}

	app := &cli.App{
		Name:  "cli-transparent-tunnel",
		Usage: "make a command run through an ssh tunnel",
		Commands: commands,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Value:       false,
				Usage:       "debug mode",
				Destination: &debug,
			},
		},
	}

	err = app.Run(cttAllArgs)
	if err != nil {
		var exErr *exec.ExitError
		if errors.As(err, &exErr) {
			logrus.WithError(err).Errorf("underlying process return %d", exErr.ExitCode())
			os.Exit(exErr.ExitCode())
		}

		logrus.WithError(err).Fatal("Failed to run")
	}
}
