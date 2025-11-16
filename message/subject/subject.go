// Package subject defines the Subject interface, which represents a subject
// line of a packet message.  This package defines four implementations of the
// Subject interface (PlainSubject, SCCoPlainSubject, SCCoFormSubject, and
// ReceiptSubject), and allows callers to register additional implementations.
package subject

import (
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/cachetrack"
)

var (
	ErrNoSubject           = errors.New("The subject line must not be empty.")
	ErrNewlineInSubject    = errors.New("The subject line must not contain newlines.")
	ErrNonPrintableSubject = errors.New("Illegal characters were removed from the message ID on the subject line.")
)

// Subject is the interface satisfied by all subject types.
type Subject interface {
	cachetrack.CacheTracker
	// EncodedSubject returns the encoded subject line.
	EncodedSubject() string
}

// A Decoder is a function that can be registered with RegisterDecoder to
// decode a subject line encoding.  The function is given the message's encoded
// subject line.  If the function recognizes the subject line encoding and is
// able to decode it, it should return a non-nil Subject and a possibly non-nil
// error (indicating non-fatal issues with the decoding).  If the function
// recognizes the subject line encoding but is not able to decode it, it should
// return a nil Subject and a non-nil error.  If the function does not
// recognize the subject line encoding, it should return nil, nil.
type Decoder func(string) (Subject, error)

// decoders is the list of registered Decoders.
var decoders []Decoder

// RegisterDecoder registers a Decoder that can be used to recognize a subject
// line encoding.
func RegisterDecoder(fn Decoder) {
	decoders = append(decoders, fn)
}

// Decode decodes an encoded subject line and returns a Subject containing it.
// If any registered Decoder recognizes the subject line encoding, and is able
// to decode it successfully, Decode returns the resulting decoded Subject and
// a possibly non-nil error (indicating non-fatal issues with the decoding).
// If no Decoder recognizes and successfully decodes the subject line, Decode
// returns a PlainSubject and the cumulative set of errors reported by Decoders
// that recognized but failed to decode the subject line.
func Decode(subject string) (_ Subject, err error) {
	if subject == "" {
		return NewPlainSubject(subject), ErrNoSubject
	}
	if trim := strings.Map(removeNewlines, subject); len(trim) < len(subject) {
		err = ErrNewlineInSubject
		subject = trim
	}
	for _, fn := range decoders {
		s, e := fn(subject)
		err = errors.Join(e)
		if s != nil {
			return s, e
		}
	}
	return NewPlainSubject(subject), err
}

func removeNewlines(r rune) rune {
	if r == '\r' || r == '\n' {
		return -1
	}
	return r
}

func removeNewlineUnderline(r rune) rune {
	if r == '\r' || r == '\n' || r == '_' {
		return -1
	}
	return r
}
