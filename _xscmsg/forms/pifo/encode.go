package pifo

import (
	"fmt"
	"io"
	"strings"
)

// CurrentVersion is the current PackItForms version number.
const CurrentVersion = "3.9"

// NewEncoder creates a new PackItForms encoder that writes form data to the
// supplied output stream.  html is the form HTML filename, which identifies the
// form type.  version is the form version number.
func NewEncoder(w io.Writer, html, version string) *Encoder {
	e := Encoder{w: w}
	_, e.err = fmt.Fprintf(w, "!SCCoPIFO!\n#T: %s\n#V: 3.9-%s\n", html, version)
	return &e
}

// Encoder is a PackItForms form encoder.
type Encoder struct {
	w   io.Writer
	err error
}

var quoteSCCoPIFO = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "]", "`]")

// Write writes a single tag/value pair to the form.
func (e *Encoder) Write(tag, value string) {
	if e.err != nil || value == "" {
		return
	}
	value = quoteSCCoPIFO.Replace(value)
	if strings.HasSuffix(value, "`") {
		value += "]]"
	}
	enc := fmt.Sprintf("%s: [%s]", tag, value)
	for len(enc) != 0 && e.err == nil {
		var toWrite string
		if len(enc) > 128 {
			toWrite, enc = enc[:128], enc[128:]
		} else {
			toWrite, enc = enc, ""
		}
		if _, e.err = io.WriteString(e.w, toWrite); e.err == nil {
			_, e.err = e.w.Write([]byte{'\n'})
		}
	}
}

// Close closes the form encoding.  It returns any error that occurred at any
// point in the form encoding process.  It does not close the underlying output
// stream.
func (e *Encoder) Close() error {
	if e.err == nil {
		_, e.err = io.WriteString(e.w, "!/ADDON!\n")
	}
	return e.err
}
