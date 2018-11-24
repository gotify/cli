package utils

import (
	"fmt"
	"time"

	"github.com/tj/go-spin"
)

type fLoader = func(success chan interface{}, failure chan error)

func SpinLoader(loading string, loader fLoader) (interface{}, error) {
	success := make(chan interface{})
	failure := make(chan error)
	go loader(success, failure)
	s := spin.New()
	for {
		select {
		case data := <-success:
			fmt.Printf("\r%s -> Success!\n", loading)
			return data, nil
		case err := <-failure:
			fmt.Printf("\r%s -> Failed\nError: %s\n", loading, err)
			return nil, err
		case <-time.After(time.Millisecond * 100):
			fmt.Printf("\r%s\033[m %s ", loading, s.Next())
		}
	}
}
