package common

import (
	"fmt"
	"io"
	"strings"
)

// EncodeHeader encodes the standard form header fields to the supplied PIFO
// encoder.
func (s *StdFields) EncodeHeader(enc *PIFOEncoder) {
	enc.Write("MsgNo", s.OriginMsgID)
	enc.Write("DestMsgNo", s.DestinationMsgID)
	enc.Write("1a.", s.MessageDate)
	enc.Write("1b.", s.MessageTime)
	enc.Write("5.", s.Handling)
	enc.Write("7a.", s.ToICSPosition)
	enc.Write("8a.", s.FromICSPosition)
	enc.Write("7b.", s.ToLocation)
	enc.Write("8b.", s.FromLocation)
	enc.Write("7c.", s.ToName)
	enc.Write("8c.", s.FromName)
	enc.Write("7d.", s.ToContact)
	enc.Write("8d.", s.FromContact)
}

// EncodeFooter encodes the standard form footer fields to the supplied PIFO
// encoder.
func (s *StdFields) EncodeFooter(enc *PIFOEncoder) {
	enc.Write("OpRelayRcvd", s.OpRelayRcvd)
	enc.Write("OpRelaySent", s.OpRelaySent)
	enc.Write("OpName", s.OpName)
	enc.Write("OpCall", s.OpCall)
	enc.Write("OpDate", s.OpDate)
	enc.Write("OpTime", s.OpTime)
}

// EncodeSubject encodes an XSC-standard message subject line.
func EncodeSubject(msgid, handling, formtag, subject string) string {
	if handling != "" {
		handling = handling[:1]
	}
	if formtag == "" {
		return fmt.Sprintf("%s_%s_%s", msgid, handling, subject)
	}
	return fmt.Sprintf("%s_%s_%s_%s", msgid, handling, formtag, subject)
}

// NewPIFOEncoder creates a new PackItForms encoder that writes form data to the
// supplied output stream.  html is the form HTML filename, which identifies the
// form type.  version is the form version number.
func NewPIFOEncoder(w io.Writer, html, version string) *PIFOEncoder {
	e := PIFOEncoder{w: w}
	_, e.err = fmt.Fprintf(w, "!SCCoPIFO!\n#T: %s\n#V: 3.9-%s\n", html, version)
	return &e
}

// PIFOEncoder is a PackItForms form encoder.
type PIFOEncoder struct {
	w   io.Writer
	err error
}

var quoteSCCoPIFO = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "]", "`]")

// Write writes a single tag/value pair to the form.
func (e *PIFOEncoder) Write(tag, value string) {
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
func (e *PIFOEncoder) Close() error {
	if e.err == nil {
		_, e.err = io.WriteString(e.w, "!/ADDON!\n")
	}
	return e.err
}
