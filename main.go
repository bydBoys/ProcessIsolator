package main

import (
	"ProcessIsolator/constants"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = constants.Name
	app.Description = constants.Desc

	app.Commands = []cli.Command{
		childCommand,
		daemonCommand,

		startCommand,
		versionCommand,
	}

	app.Before = func(context *cli.Context) error {
		color.Cyan(constants.Desc)
		if os.Getuid() != 0 {
			color.HiRed("Only root user can run it, because of the linux's restriction on namespace")
			os.Exit(-1)
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
