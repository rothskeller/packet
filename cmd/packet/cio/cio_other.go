//go:build !windows

package cio

import (
	"os"

	"golang.org/x/term"
)

func (cio *CIO) detect() {
	var istate, ostate *term.State
	var err error

	istate, err = term.GetState(int(os.Stdin.Fd()))
	cio.InputIsTerm = err == nil
	if cio.InputIsTerm && cio.initialStateIn == nil {
		cio.initialStateIn = istate
	}
	ostate, err = term.GetState(int(os.Stdout.Fd()))
	cio.OutputIsTerm = err == nil
	if cio.OutputIsTerm && cio.initialStateOut == nil {
		cio.initialStateOut = ostate
	}
	if cio.OutputIsTerm {
		cio.Width, _, _ = term.GetSize(int(os.Stdout.Fd()))
	}
}

func (cio *CIO) rawMode() {
	term.MakeRaw(int(os.Stdin.Fd()))
}

func (cio *CIO) restoreTerminal() {
	term.Restore(int(os.Stdin.Fd()), cio.initialStateIn)
}
