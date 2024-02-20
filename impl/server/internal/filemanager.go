package internal

import (
	"ProcessIsolator/constants"
	"fmt"
	"io"
	"os"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4096)
	},
}

func WriteFile(name string, src io.Reader) error {
	outFile, err := os.Create(fmt.Sprintf(constants.FilePath, name))
	if err != nil {
		return err
	}
	defer func(outFile *os.File) {
		_ = outFile.Close()
	}(outFile)

	buf := bufferPool.Get().([]byte)
	defer func(tmpBuf []byte) {
		bufferPool.Put(buf)
	}(buf)

	_, err = io.CopyBuffer(outFile, src, buf)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(name string) error {
	return os.Remove(fmt.Sprintf(constants.FilePath, name))
}

func CheckFile(name string) bool {
	_, err := os.Stat(fmt.Sprintf(constants.FilePath, name))
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return false
		}
	}
	return true
}
