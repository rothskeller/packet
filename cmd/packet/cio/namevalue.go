package cio

import (
	"encoding/csv"
	"os"
	"strings"
)

type NameValueList struct {
	cio  *CIO
	cw   *csv.Writer
	data [][]string
}

// NewNameValueList starts a new name/value list.
func (cio *CIO) NewNameValueList() (nv *NameValueList) {
	nv = &NameValueList{cio: cio}
	if !cio.OutputIsTerm {
		nv.cw = csv.NewWriter(os.Stdout)
	}
	return nv
}

// ShowNVPair displays one pair in the name/value list.
func (nv *NameValueList) ShowNVPair(name, value string) {
	if nv.cw == nil {
		nv.data = append(nv.data, []string{name, value})
	} else {
		nv.cw.Write([]string{name, value})
	}
}

// Close closes the forms list.
func (nv *NameValueList) Close() {
	if nv.cw != nil {
		nv.cw.Flush()
		return
	}
	nv.cio.clearStatus()
	var nameWidth int
	for _, row := range nv.data {
		nameWidth = max(nameWidth, len(row[0]))
	}
	for _, row := range nv.data {
		value := strings.TrimRight(row[1], "\n")
		lines := strings.Split(value, "\n")
		var linelen int
		for _, line := range strings.Split(value, "\n") {
			linelen = max(linelen, len(line))
		}
		// If the longest line fits to the right of the name, show it
		// that way.  Otherwise, show it on the following lines with a
		// 4-space indent.
		var indent string
		if linelen <= nv.cio.Width-nameWidth-3 {
			nv.cio.print(colorLabel, setLength(row[0], nameWidth)+"  ")
			indent = spaces[:nameWidth+2]
		} else {
			nv.cio.print(colorLabel, row[0])
			nv.cio.print(0, "\n    ")
			lines, _ = wrap(value, nv.cio.Width-5)
			indent = spaces[:4]
		}
		// Show the lines.
		for i, line := range lines {
			if i != 0 {
				nv.cio.print(0, indent)
			}
			nv.cio.print(0, line)
			nv.cio.print(0, "\n")
		}
	}
}
