package message

import (
	"fmt"
	"io"
	"strings"

	"slices"
)

// EncodeSubject encodes the message subject line.
func (bm *BaseMessage) EncodeSubject() string {
	var msgid, handling, formtag, subject string

	if bm.FOriginMsgID != nil {
		msgid = *bm.FOriginMsgID
	}
	if bm.FHandling != nil {
		handling = *bm.FHandling
	}
	if bm.Form != nil {
		formtag = bm.Form.Tag
	}
	if bm.FSubject != nil {
		subject = *bm.FSubject
	}
	return EncodeSubject(msgid, handling, formtag, subject)
}

// EncodeBody encodes the message body, suitable for transmission or
// storage.
func (bm *BaseMessage) EncodeBody() string {
	var (
		sb     strings.Builder
		enc    *PIFOEncoder
		values []string
	)
	if bm.Form.HTML == "" {
		panic("BaseMessage.EncodeBody can only encode PackItForms; other message types must override")
	}
	enc = NewPIFOEncoder(&sb, bm.Form.HTML, bm.Form.Version)
	values = make([]string, len(bm.Form.FieldOrder))
	for _, f := range bm.Fields {
		if f.PIFOTag == "" || *f.Value == "" {
			continue
		}
		if idx := slices.Index(bm.Form.FieldOrder, f.PIFOTag); idx >= 0 {
			values[idx] = *f.Value
		} else {
			enc.Write(f.PIFOTag, *f.Value)
		}
	}
	for i, tag := range bm.Form.FieldOrder {
		if values[i] != "" {
			enc.Write(tag, values[i])
		}
	}
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return sb.String()
}

// EncodeSubject encodes an XSC-standard message subject line.
func EncodeSubject(msgid, handling, formtag, subject string) string {
	if msgid == "" && handling == "" && formtag == "" {
		return subject
	}
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
	_, e.err = fmt.Fprintf(w, "!SCCoPIFO!\n#T: %s\n#V: 3.11-%s\n", html, version)
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
