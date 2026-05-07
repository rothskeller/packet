package htmlb

import (
	"bytes"
	"reflect"
	"strings"
	"unsafe"
)

// writeBytes writes a byte slice to the output, with error handling.
func (e *Element) writeBytes(bs []byte) {
	if e.err != nil {
		return
	}
	_, e.err = e.w.Write(bs)
}

// writeString writes a string to the output, with no allocations and with error
// handling.
func (e *Element) writeString(s string) {
	if e.err != nil || s == "" {
		return
	}
	const max = 0x7fff0000
	var bs = (*[max]byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data))[:len(s):len(s)]
	_, e.err = e.w.Write(bs)
}

// writeBytesEnc writes a byte slice to the output, with HTML encoding and error
// handling.
func (e *Element) writeBytesEnc(bs []byte) {
	for len(bs) != 0 {
		idx := bytes.IndexAny(bs, `"'&<>`)
		if idx < 0 {
			e.writeBytes(bs)
			break
		}
		if idx > 0 {
			e.writeBytes(bs[:idx])
			bs = bs[idx:]
		}
		switch bs[0] {
		case '"':
			e.writeBytes(escQUOT)
		case '\'':
			e.writeBytes(escAPOS)
		case '&':
			e.writeBytes(escAMP)
		case '<':
			e.writeBytes(escLT)
		case '>':
			e.writeBytes(escGT)
		}
		bs = bs[1:]
	}
}

// writeStringEnc writes a string to the output, with HTML encoding and error
// handling.
func (e *Element) writeStringEnc(s string) {
	for s != "" {
		idx := strings.IndexAny(s, `"'&<>`)
		if idx < 0 {
			e.writeString(s)
			break
		}
		if idx > 0 {
			e.writeString(s[:idx])
			s = s[idx:]
		}
		switch s[0] {
		case '"':
			e.writeBytes(escQUOT)
		case '\'':
			e.writeBytes(escAPOS)
		case '&':
			e.writeBytes(escAMP)
		case '<':
			e.writeBytes(escLT)
		case '>':
			e.writeBytes(escGT)
		}
		s = s[1:]
	}
}

// writeStringEncQuoted writes a string to the output, with HTML encoding and
// error handling.
func (e *Element) writeStringEncQuoted(s string) {
	if s == "" || strings.IndexAny(s, " \t\f\r\n=`") >= 0 {
		e.writeBytes(bsQUOT)
		e.writeStringEnc(s)
		e.writeBytes(bsQUOT)
	} else {
		e.writeStringEnc(s)
	}
}
