package util

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func ReadCmdline(pid int) (string, error) {
	content, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(content), "\x00"), nil
}

func KillProcess(pid int) error {
	cmd := exec.Command("kill", strconv.Itoa(pid))
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
