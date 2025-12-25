package subject

import (
	"strings"

	"github.com/rothskeller/packet/message/cachetrack"
)

// A PlainSubject is a subject that has no special encoding.
type PlainSubject struct {
	cachetrack.Tracker
	subject string
}

var _ Subject = (*PlainSubject)(nil)

// NewPlainSubject creates a new plain subject with the specified content.
func NewPlainSubject(subject string) (s *PlainSubject) {
	s = &PlainSubject{subject: subject}
	return s
}

// EncodedSubject returns the plain subject line.
func (s *PlainSubject) EncodedSubject() string {
	s.MarkClean()
	return s.subject
}

// SetSubject sets the content of the plain subject line.
func (s *PlainSubject) SetSubject(subject string) (err error) {
	if subject == "" {
		err = ErrNoSubject
	} else if trim := strings.Map(removeNewlines, subject); len(trim) < len(subject) {
		err = ErrNewlineInSubject
		subject = trim
	}
	if s.subject != subject {
		s.subject = subject
		s.MarkDirty("subject.PlainSubject.Subject")
	}
	return err
}

func (s *PlainSubject) Clone() Subject {
	return NewPlainSubject(s.EncodedSubject())
}
