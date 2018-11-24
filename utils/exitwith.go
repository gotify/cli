package utils

import (
	"fmt"
	"os"
)

func Exit1With(message ...interface{}) {
	fmt.Println(message...)
	os.Exit(1)
}
