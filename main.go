package main

import (
	"fmt"
	command_runner "github.com/odedpriva/cli-transparent-tunnel/command-runner"
	command_tunneler "github.com/odedpriva/cli-transparent-tunnel/command-tunneler"
	"github.com/odedpriva/cli-transparent-tunnel/command-tunneler/commands/kubectl"
	"github.com/odedpriva/cli-transparent-tunnel/config"
	"github.com/odedpriva/cli-transparent-tunnel/service"
	"github.com/odedpriva/cli-transparent-tunnel/version"
	"os"

	"errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	args := os.Args

	if len(args) > 1 {
		switch args[1] {
		case "ctt-init":
			file, err := config.InitConfig()
			if err != nil {
				fmt.Printf("error setting up ctt config %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("check config files: %s and update tunnel configurations\n", file)
			os.Exit(0)
		case "ctt-version":
			fmt.Printf("%s \n", version.GetVersion())
			os.Exit(0)
		}
	}

	mylog := log.New()

	c, err := config.LoadConfig()
	if err != nil {
		if errors.Is(err, config.ErrConfigNotExist) {
			fmt.Println("ctt requires config file, have you run ctt init?")
			os.Exit(1)
		}
	}
	mylog.SetLevel(c.LogLevel)
	s, err := serviceFactory(c, mylog, args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.Run()

}

func serviceFactory(conf *config.Config, log *log.Logger, args []string) (*service.Service, error) {

	command, commandLineArgs, err := parseOsArgs(args)
	if err != nil {
		return nil, err
	}
	var t command_tunneler.TunnelerCommand
	c := command_runner.NewCommandRunnerImpl(log)
	switch command {
	case "kubectl":
		t, err = kubectl.NewKubeCtl(conf.KubeCtl, log, commandLineArgs)
		if err != nil {
			return nil, err
		}
	default:
		fmt.Printf("does not support %s as command\n", command)
		os.Exit(1)
	}

	return &service.Service{
		TunnelerCommand: t,
		CommandRunner:   c,
		Log:             log,
		SshConfig:       conf.SshConfig,
	}, nil
}

func parseOsArgs(osArgs []string) (command string, args []string, err error) {

	switch numberOfArgs := len(osArgs); {
	case numberOfArgs == 0 || numberOfArgs == 1:
		return command, args, fmt.Errorf("ctt expects subcommand")
	case numberOfArgs == 2:
		return osArgs[1], []string{}, nil
	default:
		return osArgs[1], osArgs[2:], nil
	}

}
