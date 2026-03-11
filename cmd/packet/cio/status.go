package cio

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var SuppressStatus bool

func (cio *CIO) Status(f string, args ...any) {
	if cio.OutputIsTerm && !SuppressStatus {
		cio.clearStatus()
		if f != "" {
			var s = f
			if len(args) != 0 {
				s = fmt.Sprintf(f, args...)
			}
			io.WriteString(os.Stdout, strings.TrimRight(s, "\n"))
			cio.haveStatus = true
		}
	} // else print nothing
}
