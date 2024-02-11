package main

import (
	"ProcZygote/constants"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = constants.Name
	app.Usage = constants.Usage

	app.Commands = []cli.Command{
		childCommand,
		daemonCommand,

		startCommand,
		versionCommand,
	}

	app.Before = func(context *cli.Context) error {
		color.Cyan(constants.Usage)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
