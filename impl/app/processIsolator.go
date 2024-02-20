package app

import (
	"ProcessIsolator/impl/app/log"
	"ProcessIsolator/impl/server"
	"ProcessIsolator/util"
	"os"
	"os/exec"
)

func RunProcessIsolator(daemon bool, outPath string) error {
	if daemon {
		return runSelfDaemon(outPath)
	}

	server.StartServer(log.GetLogChan())
	// todo: syscall listener
	strings := make(chan string)
	for {
		select {
		case <-strings:

		}
	}
	return nil
}

func runSelfDaemon(outPath string) error {
	selfProc, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return err
	}
	cmd := exec.Command(selfProc, "daemon")
	if cmd.Stdout, err = util.GenerateFile(outPath); err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
