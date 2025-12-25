// Package body defines the Body interface, which represents a decoded packet
// message Body, and the most basic implementation of it (PlainBody).  Callers
// can define and register additional implementations.
package body

import (
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/cachetrack"
)

// Body is the interface honored by all different types of message bodies.  It
// allows queries of the body encoding, metadata, and fields.
type Body interface {
	cachetrack.CacheTracker
	// EncodedBody returns the encoded message body.
	EncodedBody() string
	// Clone returns a copy of the body.
	Clone() Body
}

// A Decoder is a function that can be registered with RegisterDecoder to
// recognize and decode a body encoding.  The function is given the message's
// encoded body.  If the function recognizes the body encoding, and is able to
// decode it successfully, it should return a non-nil Body and a possibly
// non-nil error (indicating non-fatal issues with the decoding).  If the
// function recognizes the body encoding but is not able to decode it
// successfully, it should return a nil Body and a non-nil error.  If the
// function does not recognize the body encoding, it should return nil, nil.
type Decoder func(string) (Body, error)

// decoders is the list of registered Decoders.
var decoders []Decoder

// RegisterDecoder registers a Decoder that can be used to recognize and decode
// a body encoding.
func RegisterDecoder(fn Decoder) {
	decoders = append(decoders, fn)
}

// Decode decodes an encoded body and returns a Body containing it.  If any
// registered Decoder recognizes the body encoding, and is able to decode it
// successfully, Decode returns the resulting decoded Body and a possibly
// non-nil error (indicating non-fatal issues with the decoding).  If no Decoder
// recognizes and successfully decodes the body, Decode returns a PlainBody and
// the cumulative set of errors reported by Decoders that recognized but failed
// to decode the body.
func Decode(body string) (_ Body, err error) {
	for _, fn := range decoders {
		b, e := fn(body)
		err = errors.Join(err, e)
		if b != nil {
			return b, e
		}
	}
	return NewPlainBody(body), err
}
