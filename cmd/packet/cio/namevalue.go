package cio

import (
	"encoding/csv"
	"os"
	"strings"
)

var tableCW *csv.Writer

// ShowNameValue shows one entry in a name/value table.  nameWidth is the width
// of the name column.  End the table by calling EndNameValueList.
func (cio *CIO) ShowNameValue(name, value string, nameWidth int) {
	if cio.OutputIsTerm {
		cio.showNameValueTable(name, value, nameWidth)
	} else {
		cio.showNameValueCSV(name, value)
	}
}

func (cio *CIO) showNameValueCSV(name, value string) {
	if tableCW == nil {
		tableCW = csv.NewWriter(os.Stdout)
	}
	tableCW.Write([]string{name, value})
}

func (cio *CIO) showNameValueTable(name, value string, nameWidth int) {
	var (
		linelen int
		lines   []string
		indent  string
	)
	cio.clearStatus()
	// Find the length of the longest line in the value.
	nameWidth = max(nameWidth, len(name))
	value = strings.TrimRight(value, "\n")
	lines = strings.Split(value, "\n")
	for _, line := range strings.Split(value, "\n") {
		linelen = max(linelen, len(line))
	}
	// If the longest line fits to the right of the name, show it that way.
	// Otherwise, show it on the following lines with a 4-space indent.
	if linelen <= cio.Width-nameWidth-3 {
		cio.print(colorLabel, setLength(name, nameWidth)+"  ")
		indent = spaces[:nameWidth+2]
	} else {
		cio.print(colorLabel, name)
		cio.print(0, "\n    ")
		lines, _ = wrap(value, cio.Width-5)
		indent = spaces[:4]
	}
	// Show the lines.
	for i, line := range lines {
		if i != 0 {
			print(0, indent)
		}
		cio.print(0, line)
		cio.print(0, "\n")
	}
}

// EndNameValueList ends a name/value list that consisted of zero or more calls
// to ShowNameValue.
func (cio *CIO) EndNameValueList() {
	if tableCW != nil {
		tableCW.Flush()
		tableCW = nil
	}
}
