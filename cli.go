package main

import (
	"fmt"
	"os"

	"github.com/gotify/cli/command"
	"github.com/gotify/cli/utils"
	"gopkg.in/urfave/cli.v1"
)

var (
	// Version the version of Gotify-CLI.
	Version = "unknown"
	// Commit the git commit hash of this version.
	Commit = "unknown"
	// BuildDate the date on which this binary was build.
	BuildDate = "unknown"
)

func main() {
	cli.VersionPrinter = versionPrinter
	app := cli.NewApp()
	app.Name = "Gotify"
	app.Version = Version
	app.Usage = "The official Gotify-CLI"
	app.Commands = []cli.Command{
		command.Init(),
		command.Version(),
		command.Config(),
		command.Push(),
	}
	err := app.Run(os.Args)
	if err != nil {
		utils.Exit1With(err)
	}
}

func versionPrinter(ctx *cli.Context) {
	fmt.Println("Version:   " + Version)
	fmt.Println("Commit:    " + Commit)
	fmt.Println("BuildDate: " + BuildDate)
}
