package english

import (
	"bytes"
	"io"
	"regexp"
)

var (
	nl      = []byte{'\n'}
	nlnl    = []byte{'\n', '\n'}
	nlsp    = []byte{'\n', ' '}
	sp      = []byte{' '}
	dotspsp = []byte{'.', ' ', ' '}
)

// Wrapper is a rudimentary word-wrapper.  Text written to it is word-wrapped at
// 78 columns.  Lines starting with space characters are left untouched.
type Wrapper struct {
	w    io.Writer
	held []byte
}

// NewWrapper creates a new Wrapper, writing to the provided writer.  The
// Wrapper must be closed to flush the last line of output.
func NewWrapper(w io.Writer) *Wrapper {
	return &Wrapper{w: w}
}

func (ww *Wrapper) Write(b []byte) (n int, err error) {
	n = len(b)
	ww.held = append(ww.held, b...)
	for len(ww.held) != 0 {
		// If ww.held starts with a space or a newline, the part up to
		// the first newline should be emitted unchanged.
		if ww.held[0] == '\n' || ww.held[0] == ' ' {
			if idx := bytes.IndexByte(ww.held, '\n'); idx >= 0 {
				if _, err = ww.w.Write(ww.held[:idx+1]); err != nil {
					return n, err
				}
				ww.held = ww.held[idx+1:]
				continue
			} else {
				break
			}
		}
		// If ww.held contains a pair of newlines, or a newline followed
		// by a space, then we should word-wrap and emit up to that
		// point.
		stop := bytes.Index(ww.held, nlnl)
		if idx := bytes.Index(ww.held, nlsp); idx >= 0 && (stop < 0 || stop > idx) {
			stop = idx
		}
		if stop < 0 {
			break
		}
		// Word wrap the line up to that point and flush it.
		if err = ww.wrap(ww.held[:stop]); err != nil {
			return n, err
		}
		ww.held = ww.held[stop+1:]
	}
	return n, nil
}

// WriteString writes a string to the Wrapper.
func (ww *Wrapper) WriteString(s string) (int, error) { return ww.Write([]byte(s)) }

// Close flushes the Wrapper output and closes it.
func (ww *Wrapper) Close() (err error) {
	if len(ww.held) != 0 {
		if err = ww.wrap(ww.held); err != nil {
			return err
		}
		ww.held = nil
	}
	if c, ok := ww.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

var wrapRE1 = regexp.MustCompile(`\.\s*\n\s*`)
var wrapRE2 = regexp.MustCompile(`\s*\n\s*`)

// wrap does the actual wrapping.
func (ww *Wrapper) wrap(b []byte) (err error) {
	b = wrapRE1.ReplaceAllLiteral(b, dotspsp)
	b = wrapRE2.ReplaceAllLiteral(b, sp)
	for len(b) > 78 {
		idx := bytes.LastIndexByte(b[:78], ' ')
		if idx < 0 {
			idx = bytes.IndexByte(b, ' ')
		}
		if idx < 0 {
			idx = len(b)
		}
		if _, err = ww.w.Write(bytes.TrimSpace(b[:idx])); err != nil {
			return err
		}
		if _, err = ww.w.Write(nl); err != nil {
			return err
		}
		b = bytes.TrimSpace(b[idx:])
	}
	if len(b) != 0 {
		if _, err = ww.w.Write(bytes.TrimSpace(b)); err != nil {
			return err
		}
		if _, err = ww.w.Write(nl); err != nil {
			return err
		}
	}
	return nil
}
