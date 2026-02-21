package cio

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rothskeller/packet/errors"
)

func (cio *CIO) Error(err error) {
	if cio.OutputIsTerm {
		cio.clearStatus()
	}
	for _, e := range errors.UnwrapJoined(err) {
		if cio.OutputIsTerm {
			cio.print(colorError, cio.WrapText("ERROR: ⇥"+e.Error()))
		} else {
			io.WriteString(os.Stderr, cio.WrapText("ERROR: ⇥"+e.Error()))
		}
	}
	if cio.OutputIsTerm {
		cio.setColor(0)
	}
}

func (cio *CIO) ErrorF(f string, args ...any) {
	var s = f
	if len(args) != 0 {
		s = fmt.Sprintf(f, args...)
	}
	if !strings.HasPrefix(s, "usage: ") {
		s = "ERROR: ⇥" + s
	}
	if cio.OutputIsTerm {
		cio.clearStatus()
		cio.print(colorError, cio.WrapText(s))
		cio.setColor(0)
	} else {
		io.WriteString(os.Stderr, cio.WrapText(s))
	}
}

func (cio *CIO) Warn(f string, args ...any) {
	var s = f
	if len(args) != 0 {
		s = fmt.Sprintf(f, args...)
	}
	s = "WARNING: ⇥" + s
	if cio.OutputIsTerm {
		cio.clearStatus()
		cio.print(colorError, cio.WrapText(s))
		cio.setColor(0)
	} else {
		io.WriteString(os.Stderr, cio.WrapText(s))
	}
}

func (cio *CIO) Confirm(f string, args ...any) {
	if cio.OutputIsTerm {
		var s = f
		if len(args) != 0 {
			s = fmt.Sprintf(f, args...)
		}
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		cio.clearStatus()
		cio.print(0, cio.WrapText(s))
	} // else don't emit
}

func (cio *CIO) Welcome(f string, args ...any) {
	var s = f
	if len(args) != 0 {
		s = fmt.Sprintf(f, args...)
	}
	s = strings.TrimRight(s, "\n")
	if cio.OutputIsTerm {
		cio.print(colorLabel, s)
		cio.print(0, "\n")
	} else {
		io.WriteString(os.Stderr, s)
		io.WriteString(os.Stderr, "\n")
	}
}
