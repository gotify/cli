package utils

import (
	"bytes"
	"io"
)

func Evaluate(s string) string {
	res := bytes.NewBuffer([]byte{})
	sourceReader := bytes.NewBufferString(s)
	for {
		r, _, err := sourceReader.ReadRune()
		if err == io.EOF {
			break
		}
		if r == '\\' {
			nextRune, _, err := sourceReader.ReadRune()
			if err == nil {
				switch nextRune {
				case '\\':
					// ignore
				case 't':
					r = '\t'
				case 'n':
					r = '\n'
				default:
					sourceReader.UnreadRune()
				}
			}
		}
		res.WriteRune(r)
	}
	return res.String()
}
