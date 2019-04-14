package command

import (
	"fmt"

	"github.com/gotify/cli/v2/config"
	"github.com/gotify/cli/v2/utils"
	"gopkg.in/urfave/cli.v1"
)

func Config() cli.Command {
	return cli.Command{
		Name:  "config",
		Usage: "Shows the config",
		Action: func(ctx *cli.Context) {
			locations := config.GetLocations()
			conf, err := config.ReadConfig(locations)
			if err != nil {
				utils.Exit1With("cannot read config:", err)
				return
			}
			fmt.Println("Used Config:", conf.FromLocation)
			fmt.Println("URL:", conf.URL)
			fmt.Println("Token:", conf.Token)
		},
	}
}
