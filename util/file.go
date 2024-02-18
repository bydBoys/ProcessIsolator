package util

import (
	"bufio"
	"os"
)

func GenerateFile(path string) (*os.File, error) {
	stdLogFile, err := os.Create(path)

	if err != nil {
		return nil, err
	}
	return stdLogFile, nil
}

func ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
