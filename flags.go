package main

import (
	"ProcZygote/constants"
	"ProcZygote/impl/app"
	"ProcZygote/util"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"strings"
)

var (
	childCommand = cli.Command{
		Hidden: true,
		Name:   "init",
		Usage:  "Do not call it outside",
		Action: func(context *cli.Context) error {
			return app.RunChildInitProcess()
		},
	}
	daemonCommand = cli.Command{
		Hidden: true,
		Name:   "daemon",
		Usage:  "Do not call it outside",
		Action: func(context *cli.Context) error {
			return app.RunProcZygote(false, "")
		},
	}

	startCommand = cli.Command{
		Name:  "start",
		Usage: "Start " + constants.Name,
		Action: func(context *cli.Context) error {
			daemon := context.Bool("daemon")
			outFile := context.String("out")

			return app.RunProcZygote(daemon, outFile)
		},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "daemon",
				Usage: "Enable daemon",
			},
			cli.StringFlag{
				Name:  "out",
				Usage: "Specify the output file",
			},
		},
		// 每次运行前，先检查是否有守护进程正在运行
		Before: func(context *cli.Context) error {
			return func(keyword string) error {
				var (
					currentPID = os.Getpid()
					pid        int // 内存优化
					command    string
					err        error
					files      []os.DirEntry
				)

				files, err = os.ReadDir("/proc")
				if err != nil {
					return err
				}

				for _, file := range files {
					if file.IsDir() {
						if pid, err = strconv.Atoi(file.Name()); err == nil {
							if pid == currentPID {
								continue
							}
							command, err = util.ReadCmdline(pid)
							if err != nil {
								continue
							}

							if strings.Contains(command, keyword) {
								color.Cyan("Find exist process %d, try kill it.", pid)
								if err = util.KillProcess(pid); err != nil {
									return err
								}
							}
						}
					}
				}
				return nil
			}(constants.Name)
		},
	}
	versionCommand = cli.Command{
		Name:  "version",
		Usage: "Show version",
		Action: func(context *cli.Context) error {
			color.Cyan(constants.Version)
			return nil
		},
	}
)
