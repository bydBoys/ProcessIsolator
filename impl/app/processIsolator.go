package app

import (
	"ProcessIsolator/impl/app/log"
	"ProcessIsolator/impl/server"
	"ProcessIsolator/util"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func RunProcessIsolator(daemon bool, outPath string) error {
	if daemon {
		return runSelfDaemon(outPath)
	}

	server.StartServer(log.GetLogChan())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	<-quit
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
