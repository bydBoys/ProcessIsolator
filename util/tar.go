package util

import (
	"os"
	"os/exec"
)

func UnTar(tarPath, destPath string) error {
	cmd := exec.Command("tar", "-xf", tarPath, "-C", destPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
