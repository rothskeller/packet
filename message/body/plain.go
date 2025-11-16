package body

import (
	"github.com/rothskeller/packet/message/cachetrack"
)

// A PlainBody is a body that consists of unencoded plain text.
type PlainBody struct {
	cachetrack.Tracker
	body string
}

var _ Body = (*PlainBody)(nil)

// NewPlainBody creates a new plain text message body with the specified
// contents.
func NewPlainBody(body string) (b *PlainBody) {
	return &PlainBody{body: body}
}

// EncodedBody returns the message body.
func (b *PlainBody) EncodedBody() string { return b.body }

// SetBody sets the plain text message body.
func (b *PlainBody) SetBody(body string) {
	if b.body != body {
		b.body = body
		b.MarkDirty("body.PlainBody.Body")
	}
}
