package utils

import (
	"io"
	"os"
)

func ProbeStdin(file io.Reader) bool {
	if file == nil {
		return false
	}
	if file, ok := file.(*os.File); ok {
		fi, err := file.Stat()
		if err != nil {
			return false
		}
		if fi.Mode()&os.ModeNamedPipe == 0 && !fi.Mode().IsRegular() {
			return false
		}
	}

	return true
}
