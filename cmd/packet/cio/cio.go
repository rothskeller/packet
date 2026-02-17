// Package cio contains input/output functions for the packet command.  The
// command has different input and output styles depending on whether standard
// input and output are terminals; that knowledge is encapsulated in this
// package.
package cio

import (
	"os"
	"strconv"

	"golang.org/x/term"
)

var spaces = "                                                                                                                                                                                                                                                                                                            "
var openCIO *CIO

// A CIO is a console I/O handler.  Any packet command that writes to the
// console should instantiate a CIO and perform writes through it.
type CIO struct {
	// InputIsTerm is true if the standard input is a terminal.
	InputIsTerm bool
	// OutputIsTerm is true if the standard output is a terminal.
	OutputIsTerm bool
	// Width is the width of the screen in characters (or our best guess).
	Width int

	initialStateIn     *term.State
	initialStateOut    *term.State
	initialStateInWin  uint32
	initialStateOutWin uint32

	curX, curY int
	lastColor  int
	hideCursor bool
	haveStatus bool
	buf        *screenBuf
}

// Open opens a new console I/O handle.
func Open() (cio *CIO) {
	if cio = openCIO; cio == nil {
		cio = new(CIO)
		openCIO = cio
	}
	cio.Detect()
	return cio
}

// Detect determines whether or not standard input and output are terminals, the
// screen width, and the initial state of the terminal.  Detect should be called
// whenever anything might have changed the desired input/output style (e.g.,
// redirection of standard input or output, change in terminal width, etc.).
func (cio *CIO) Detect() {
	cio.Width = 0
	cio.detect() // OS-specific detection
	if cio.Width == 0 {
		if w, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && w > 0 {
			cio.Width = w
		} else {
			cio.Width = 80
		}
	}
}

// setLength returns the input string, truncated or right-padded to be exactly
// the requested length.
func setLength(s string, l int) string {
	if l < 0 {
		return s
	}
	if len(s) < l {
		return s + spaces[:l-len(s)]
	}
	return s[:l]
}

// setMaxLength returns the input string, truncated if necessary to be no longer
// than the requested length.
func setMaxLength(s string, l int) string {
	if l < 0 {
		return s
	}
	if len(s) > l {
		return s[:l]
	}
	return s
}

func prevword(s string, cur int) int {
	for cur > 0 && (s[cur] == ' ' || s[cur] == '\n') {
		cur--
	}
	for cur > 0 && s[cur-1] != ' ' && s[cur] != '\n' {
		cur--
	}
	return cur
}

func nextword(s string, cur int) int {
	l := len(s)
	for cur < l && (s[cur] == ' ' || s[cur] == '\n') {
		cur++
	}
	for cur < l && s[cur] != ' ' && s[cur] != '\n' {
		cur++
	}
	return cur
}

func splitOnSelect(s string, selstart, selend int) (presel, sel, postsel string) {
	if selstart == selend {
		return s, "", ""
	}
	if selstart >= len(s) {
		presel = s
	} else if selstart >= 0 {
		presel = s[:selstart]
	}
	if selstart < len(s) && selend > 0 {
		if selstart > 0 && selend < len(s) {
			sel = s[selstart:selend]
		} else if selstart > 0 {
			sel = s[selstart:]
		} else if selend < len(s) {
			sel = s[:selend]
		} else {
			sel = s
		}
	}
	if selend <= 0 {
		postsel = s
	} else if selend < len(s) {
		postsel = s[selend:]
	}
	if s != presel+sel+postsel {
		panic("assert")
	}
	return presel, sel, postsel
}
