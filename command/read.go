package command

import (
	"io"
	"strings"

	"github.com/gotify/cli/v2/utils"
)

func readMessage(args []string, r io.Reader, output chan<- string, split *rune) {
	msgArgs := strings.Join(args, " ")

	if msgArgs != "" {
		if utils.ProbeStdin(r) {
			utils.Exit1With("message is set via arguments and stdin, use only one of them")
		}

		output <- msgArgs
		close(output)
		return
	}

	var buf strings.Builder
	for {
		var tmp [256]byte
		n, err := r.Read(tmp[:])
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			utils.Exit1With(err)
		}
		tmpStr := string(tmp[:n])
		if split != nil {
			// split the message on the null character
			parts := strings.Split(tmpStr, string(*split))
			if len(parts) == 1 {
				buf.WriteString(parts[0])
				continue
			}

			previous := buf.String()
			// fuse previous with parts[0], send parts[1] .. parts[n-2] and set parts[n-1] as new previous
			firstMsg := previous + parts[0]
			output <- firstMsg
			for _, part := range parts[1 : len(parts)-1] {
				output <- part
			}
			buf.Reset()
			buf.WriteString(parts[len(parts)-1])
		} else {
			buf.WriteString(tmpStr)
		}
	}

	if buf.Len() > 0 {
		output <- buf.String()
	}

	close(output)
}
