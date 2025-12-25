// Package payload defines the Payload interface, for a wrapper around a
// message Body that handles base64 encoding, quoted-printable encoding, and
// Outpost body flags.
package payload

import (
	"net/mail"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/cachetrack"
)

// Payload is the interface satisfied by all payload implementations.
type Payload interface {
	cachetrack.CacheTracker
	// Body returns the payload body.
	Body() body.Body
	// Encode encodes the payload.
	Encode() string
	// Clone returns a copy of the payload.
	Clone() Payload
}

type EditablePayload interface {
	Payload
}

// A Decoder is a function that can be registered with RegisterDecoder to
// recognize and decode a body wrapper, i.e., content transfer encoding.  The
// function is given the message's headers and encoded body.  If the function
// recognizes the transfer encoding, and is able to decode it and the contained
// body successfully, it
// should return a non-nil Payload and a possibly non-nil error (indicating
// non-fatal issues with the decoding of the Payload or the contained body).  If the function recognizes
// the transfer encoding but is not able to decode it or the contained body successfully, it should
// return a nil Payload and a non-nil error.  If the function does not recognize
// the transfer encoding, it should return nil, nil.
type Decoder func(mail.Header, string) (Payload, error)

// decoders is the list of registered Decoders.
var decoders []Decoder

// RegisterDecoder registers a Decoder that can be used to recognize and decode
// a content transfer encoding.
func RegisterDecoder(fn Decoder) {
	decoders = append(decoders, fn)
}

// Decode decodes a wrapped (transfer-encoded) body and returns a Payload
// containing it.  If any registered Decoder recognizes the transfer encoding,
// and is able to decode it successfull, Decode returns the resulting decoded
// Payload and a possibly non-nil error (indicating non-fatal issues with the
// decoding of the Payload or its contained body).  If no Decoder recognizes
// and successfully decodes the body, Decode returns a nil Payload and an error.
func Decode(headers mail.Header, body string) (_ Payload, err error) {
	for _, fn := range decoders {
		p, e := fn(headers, body)
		err = errors.Join(err, e)
		if p != nil {
			return p, e
		}
	}
	p, e := decodeOutpostPayload(headers, body)
	if e != nil {
		err = errors.Join(err, e)
	}
	if p != nil {
		return p, e
	}
	return nil, err
}
