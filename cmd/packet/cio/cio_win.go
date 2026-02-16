//go:build windows

package cio

import (
	"os"

	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

func (cio *CIO) detect() {
	var istate, ostate uint32

	err := windows.GetConsoleMode(windows.Handle(int(os.Stdin.Fd())), &istate)
	cio.InputIsTerm = err == nil
	if cio.InputIsTerm && cio.initialStateInWin == 0 {
		cio.initialStateInWin = istate
	}
	err = windows.GetConsoleMode(windows.Handle(int(os.Stdout.Fd())), &ostate)
	cio.OutputIsTerm = err == nil
	if cio.OutputIsTerm {
		if ostate&0x0005 != 0x0005 {
			ostate |= 0x0004 // ENABLE_VIRTUAL_TERMINAL_PROCESSING
			ostate |= 0x0001 // ENABLE_PROCESSED_OUTPUT
			windows.SetConsoleMode(windows.Handle(int(os.Stdout.Fd())), ostate)
		}
		if cio.initialStateOutWin == 0 {
			cio.initialStateOutWin = ostate
		}
		cio.Width, _, _ = term.GetSize(int(os.Stdout.Fd()))
	}
}

func (cio *CIO) rawMode() {
	windows.SetConsoleMode(windows.Handle(int(os.Stdin.Fd())), 0x0200) // ENABLE_VIRTUAL_TERMINAL_INPUT
}

func (cio *CIO) restoreTerminal() {
	windows.SetConsoleMode(windows.Handle(int(os.Stdin.Fd())), cio.initialStateInWin)
	windows.SetConsoleMode(windows.Handle(int(os.Stdout.Fd())), cio.initialStateOutWin)
}
