package app

import (
	"ProcZygote/impl/server"
	"ProcZygote/util"
	"os"
	"os/exec"
)

func RunProcZygote(daemon bool, outPath string) error {
	if daemon {
		return runSelfDaemon(outPath)
	}

	return server.StartRPCServer()
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
