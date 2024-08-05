package command

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/gotify/cli/v2/utils"
)

func readMessage(args []string, r io.Reader, output chan<- string, splitOnNull bool) {
	defer close(output)

	switch {
	case len(args) > 0:
		if utils.ProbeStdin(r) {
			utils.Exit1With("message is set via arguments and stdin, use only one of them")
		}

		output <- strings.Join(args, " ")
	case splitOnNull:
		read := bufio.NewReader(r)
		for {
			s, err := read.ReadString('\x00')
			if err != nil {
				if !errors.Is(err, io.EOF) {
					utils.Exit1With("read error", err)
				}
				if len(s) > 0 {
					output <- s
				}
				return
			} else {
				if len(s) > 1 {
					output <- strings.TrimSuffix(s, "\x00")
				}
			}
		}
	default:
		bytes, err := io.ReadAll(r)
		if err != nil {
			utils.Exit1With("cannot read", err)
		}
		output <- string(bytes)
		return
	}

}
