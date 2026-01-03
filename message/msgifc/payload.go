package msgifc

// Payload is the interface satisfied by all payload implementations.
type Payload interface {
	CacheTracker
	// Body returns the payload body.
	Body() Body
	// Encode encodes the payload.
	Encode() string
	// Clone returns a copy of the payload.
	Clone() Payload
}

type EditablePayload interface {
	Payload
}
