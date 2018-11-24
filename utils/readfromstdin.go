package utils

import (
	"os"
	"io/ioutil"
)

func ReadFrom(file *os.File) string {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return ""
	}
	if fi.Mode()&os.ModeNamedPipe == 0 && !fi.Mode().IsRegular() {
		return ""
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return ""
	}
	return string(bytes)
}
