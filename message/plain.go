package message

import (
	"fmt"

	"github.com/rothskeller/packet/errors"
)

var (
	ErrHasFormTag         = errors.New("The subject line of a plain text message must not contain a form name.")
	ErrNonStandardSubject = errors.New("The subject line is not in SCCo-standard format.")
	ErrNoBody             = errors.New("The message body is empty.")
)

type ErrNotPlainTextBody string

func (e ErrNotPlainTextBody) Error() string {
	return fmt.Sprintf("A plain text message must have a plain text body, not %s.", string(e))
}

// A PlainMessage is a plain text message with no special structure or
// interpretation.
type plainMessage struct{ *BaseMType }

var PlainMessage plainMessage

func init() {
	PlainMessage.BaseMType = NewBaseMType("a plain text message", "plain")
	PlainMessage.AddField() // TODO
}

// Recognize does nothing.  Messages are assigned the PlainMessage type by
// SetType when nothing else matches.
func (mt plainMessage) Recognize(m Message) {}

/*
// Validate validates a plain message.
func (m *PlainMessage) Validate(pifo bool) (err error) {
	// Unless the message is a bulletin, it should have an SCCo-standard
	// subject line.
	if m.envelope.Bulletin() {
		if m.envelope.Subject().EncodedSubject() == "" && !pifo {
			err = errors.Join(err, subject.ErrNoSubject)
		}
	} else if !pifo {
		switch s := m.envelope.Subject().(type) {
		case *subject.SCCoSubject:
			err = errors.Join(err, s.Validate(!m.envelope.Received()))
		case *subject.SCCoFormSubject:
			err = errors.Join(err,
				ErrHasFormTag,
				s.Validate(!m.envelope.Received()),
			)
		default:
			err = errors.Join(err, ErrNonStandardSubject)
		}
	}
	switch b := m.envelope.Body().(type) {
	case *body.PlainBody:
		if b.EncodedBody() == "" && !pifo {
			err = errors.Join(err, ErrNoBody)
		}
	default:
		err = errors.Join(err, ErrNotPlainTextBody(fmt.Sprintf("%T", b)))
	}
	return err
}
*/
