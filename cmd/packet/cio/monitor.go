package cio

import (
	"io"
	"os"
)

func (cio *CIO) Monitor(line string) {
	if cio.OutputIsTerm {
		cio.clearStatus()
		cio.print(colorReceipt, cio.WrapText(line))
		cio.setColor(0)
	} else {
		io.WriteString(os.Stdout, cio.WrapText(line))
	}
}
