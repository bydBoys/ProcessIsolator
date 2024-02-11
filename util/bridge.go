package util

import (
	"io"
	"os"
	"strings"
)

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func SendCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}

func ReadCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer func(pipe *os.File) {
		_ = pipe.Close()
	}(pipe)
	msg, err := io.ReadAll(pipe)
	if err != nil {
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
