package cio

import (
	"encoding/csv"
	"os"
)

type FormsList struct {
	cio       *CIO
	recvonly  bool
	pifoident bool
	cw        *csv.Writer
	data      [][]string
}

// NewFormsList starts a new list of forms.
func (cio *CIO) NewFormsList(recvonly, pifoident bool) (fl *FormsList) {
	fl = &FormsList{cio: cio, recvonly: recvonly, pifoident: pifoident}
	if cio.OutputIsTerm {
		var headings []string
		if recvonly {
			headings = append(headings, "*")
		}
		headings = append(headings, "TAG", "CT", "VER", "DESCRIPTION")
		if pifoident {
			headings = append(headings, "ADDON", "HTML")
		}
		fl.data = append(fl.data, headings)
	} else {
		fl.cw = csv.NewWriter(os.Stdout)
		var headings []string
		if recvonly {
			headings = append(headings, "RS")
		}
		headings = append(headings, "TAG", "CT", "VER", "DESCRIPTION")
		if pifoident {
			headings = append(headings, "ADDON", "HTML")
		}
		fl.cw.Write(headings)
	}
	return fl
}

// ShowForm displays one form in the forms list.
func (fl *FormsList) ShowForm(sendable bool, tag, key, version, title, addon, htmlname string) {
	var line []string

	switch {
	case fl.recvonly && sendable && fl.cw == nil:
		line = append(line, "*")
	case fl.recvonly && sendable && fl.cw != nil:
		line = append(line, "RS")
	case fl.recvonly && !sendable && fl.cw == nil:
		line = append(line, "")
	case fl.recvonly && !sendable && fl.cw != nil:
		line = append(line, "R")
	}
	line = append(line, tag, key, version, title)
	if fl.pifoident {
		if fl.cw != nil {
			line = append(line, addon, htmlname)
		} else {
			line = append(line, "!"+addon+"!", "#T: "+htmlname)
		}
	}
	if fl.cw == nil {
		fl.data = append(fl.data, line)
	} else {
		fl.cw.Write(line)
	}
}

// Close closes the forms list.
func (fl *FormsList) Close() {
	if fl.cw != nil {
		fl.cw.Flush()
		return
	}
	fl.cio.clearStatus()
	if len(fl.data) == 1 { // only the headers
		fl.cio.Confirm("No forms installed.")
	}
	var widths = make([]int, len(fl.data[0]))
	for _, line := range fl.data {
		for i, col := range line {
			widths[i] = max(widths[i], len(col))
		}
	}
	var linebreak, indent int
	if len(widths) > 5 {
		linebreak = len(widths) - 2
		for i := range linebreak - 1 {
			indent += widths[i] + 2
		}
	}
	var color = colorLabel
	for _, line := range fl.data {
		for j := range widths {
			if j == linebreak {
				fl.cio.print(color, spaces[:indent])
			}
			if j != linebreak-1 && j != len(widths)-1 {
				fl.cio.print(color, setLength(line[j], widths[j])+"  ")
			} else {
				fl.cio.print(color, line[j]+"\n")
			}
		}
		color = colorNormal
	}
}
