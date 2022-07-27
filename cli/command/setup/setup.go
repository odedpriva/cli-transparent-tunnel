package setup

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type SetupCmd struct {
}

func NewSetupCmd() *SetupCmd {
	return &SetupCmd{}
}

func (c *SetupCmd) GetCommnad() *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "",
		Action: func(*cli.Context) error {
			fmt.Println("will setup ")
			return nil
		},
	}
}
