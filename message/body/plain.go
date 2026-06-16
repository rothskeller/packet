package body

import (
	"iter"
	"strings"

	"github.com/rothskeller/packet/v4/message/cachetrack"
	"github.com/rothskeller/packet/v4/message/field"
	"github.com/rothskeller/packet/v4/message/msgifc"
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
	return &PlainBody{body: strings.ReplaceAll(body, "\r", "")}
}

// EncodedBody returns the message body.
func (b *PlainBody) EncodedBody() string { return b.body }

// SetBody sets the plain text message body.
func (b *PlainBody) SetBody(body string) {
	body = strings.ReplaceAll(body, "\r", "")
	if b.body != body {
		b.body = body
		b.MarkDirty("body.PlainBody.Body")
	}
}

// Clone returns a copy of the body.
func (b *PlainBody) Clone() Body {
	return NewPlainBody(b.EncodedBody())
}

// Fields returns an interator on the one and only body field.
func (b *PlainBody) Fields() iter.Seq[field.Field] {
	return func(yield func(field.Field) bool) {
		yield(plainBodyField)
	}
}

var plainBodyField = field.NewField("", "Body").
	Common(field.CDefaultBody).
	ValueFunc(func(m msgifc.Message) string { return m.Body().(*PlainBody).body }).
	SetValueFunc(func(m msgifc.Message, s string) { m.Body().(*PlainBody).body = s }).
	Required().
	EditHelp("This is the body of the message.").
	Multiline().
	CompareFunc(field.CompareText).
	MakeField()
