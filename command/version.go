package command

import (
	"gopkg.in/urfave/cli.v1"
)

func Version() cli.Command {
	return cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Shows the version",
		Action:  cli.ShowVersion,
	}
}
