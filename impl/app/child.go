package app

import (
	"ProcZygote/util"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func RunChildInitProcess() error {
	cmdArray := util.ReadCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("child process get command error")
	}

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		return err
	}
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		return err
	}
	return nil
}
