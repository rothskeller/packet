package msgifc

import "iter"

// Body is the interface honored by all different types of message bodies.  It
// allows queries of the body encoding, metadata, and fields.
type Body interface {
	CacheTracker
	// EncodedBody returns the encoded message body.
	EncodedBody() string
	// Clone returns a copy of the body.
	Clone() Body
	// Fields returns an iterator on the field(s) of the body.
	Fields() iter.Seq[Field]
}
