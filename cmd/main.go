package main

import (
	"fmt"
	"github.com/odedpriva/cli-transparent-tunnel/cli/command/setup"
	"github.com/odedpriva/cli-transparent-tunnel/cli/command/tunnel"
	"github.com/odedpriva/cli-transparent-tunnel/logging"
	"github.com/odedpriva/cli-transparent-tunnel/utils/args-utils"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/maps"
	"log"
	"os"

	"github.com/odedpriva/cli-transparent-tunnel/config"
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

	cttAllArgs, commandAllArgs := args_utils.SplitArgs(args, maps.Keys(conf.CliConfigurations))

	ctt := &cli.App{
		Name:  "ctt",
		Usage: "make a command run through an ssh tunnel",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Value:       false,
				Usage:       "debug mode",
				Destination: &debug,
			},
		},
	}

	setupCmd := setup.NewSetupCmd()
	tunnelCmd := tunnel.NewTunnelCmd(commandAllArgs, conf)
	ctt.Commands = append(ctt.Commands, setupCmd.GetCommnad(), tunnelCmd.GetCommnad())

	if err := ctt.Run(cttAllArgs); err != nil {
		log.Fatal(err)
	}

}
