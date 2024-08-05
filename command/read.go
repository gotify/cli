package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/gotify/cli/v2/utils"
	"github.com/mattn/go-isatty"
)

func readMessage(args []string, r io.Reader, output chan<- string, splitOnNull bool) {
	defer close(output)

	if len(args) > 0 {
		if utils.ProbeStdin(r) {
			utils.Exit1With("message is set via arguments and stdin, use only one of them")
		}

		output <- strings.Join(args, " ")
		return
	}

	if isatty.IsTerminal(os.Stdin.Fd()) {
		eofKey := "Ctrl+D"
		if runtime.GOOS == "windows" {
			eofKey = "Ctrl+Z"
		}
		fmt.Fprintf(os.Stderr, "Enter your message, press Enter and then %s to finish:\n", eofKey)
	}

	if splitOnNull {
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
	} else {
		bytes, err := io.ReadAll(r)
		if err != nil {
			utils.Exit1With("cannot read", err)
		}
		output <- string(bytes)
		return
	}

}
