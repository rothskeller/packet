package msgifc

import "iter"

// Subject is the interface satisfied by all subject types.
type Subject interface {
	CacheTracker
	// EncodedSubject returns the encoded subject line.
	EncodedSubject() string
	// Clone returns a copy of the subject.
	Clone() Subject
	// Fields returns an iterator on the fields encoded in the subject.
	Fields() iter.Seq[Field]
	// SubjectMessageID returns the message ID encoded in the subject, if
	// any.
	SubjectMessageID() string
	// SetSubjectMessageID sets the message ID encoded in the subject.
	SetSubjectMessageID(string) error
	// SubjectHandling returns the handling order encoded in the subject,
	// if any.  If the handling order code in the subject is one of the
	// county standard ones, SubjectHandling returns the corresponding full
	// word handling order.
	SubjectHandling() string
	// SetSubjectHandling sets the handling order code encoded in the
	// subject line.  If the supplied string is one of the county standard
	// handling order words, the corrseponding code is put in the subject.
	SetSubjectHandling(string) error
	// SubjectSummary returns the summary portion of the subject line or,
	// for subject lines not encoded in the county standard format, the
	// entire subject line.
	SubjectSummary() string
	// SetSubjectSummary sets the summary portion of the subject line.
	SetSubjectSummary(string) error
}
