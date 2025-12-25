package subject

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/cachetrack"
	"github.com/rothskeller/packet/message/messageid"
)

var (
	handlingCodes      = map[string]string{"I": "IMMEDIATE", "P": "PRIORITY", "R": "ROUTINE"}
	oldSeverityRE      = regexp.MustCompile(`^[A-Z]/[A-Z]$`)
	ErrNoHandling      = errors.New("The subject line must have a handling order code.")
	ErrHasSeverity     = errors.New("The subject line must not have a severity code; that is an obsolete format.")
	ErrInvalidHandling = errors.New("The handling order on the subject line must be I, P, or R.")
	ErrNoSummary       = errors.New("The subject line must have a message summary.")
)

// An SCCoSubject is a subject that is encoded according to the Santa Clara
// County standard for plain text packet message subject lines.
type SCCoSubject struct {
	cachetrack.Tracker
	encoded  string
	msgID    string
	handling string
	summary  string
}

var _ Subject = (*SCCoSubject)(nil)

// NewSCCoSubject creates a new SCCo-standard plain text message subject with
// the specified parameters.  It returns any problems with the parameters.
func NewSCCoSubject(msgID, handling, summary string) (s *SCCoSubject, err error) {
	s = new(SCCoSubject)
	err = errors.Join(
		s.SetSubjectMessageID(msgID),
		s.SetSubjectHandling(handling),
		s.SetSubjectSummary(summary),
	)
	return s, err
}

// EncodedSubject returns the encoded subject line.
func (s *SCCoSubject) EncodedSubject() string {
	if s.Dirty() {
		s.encoded = fmt.Sprintf("%s_%s_%s", s.msgID, s.handling, s.summary)
		s.MarkClean()
	}
	return s.encoded
}

// SubjectMessageID returns the message ID encoded in the subject line.
func (s *SCCoSubject) SubjectMessageID() string { return s.msgID }

// SetSubjectMessageID sets the message ID encoded in the subject line.
func (s *SCCoSubject) SetSubjectMessageID(msgID string) (err error) {
	if trim := strings.Map(removeNewlineUnderline, msgID); len(trim) < len(msgID) {
		err = ErrNonPrintableSubject
		msgID = trim
	} else {
		msgID, err = messageid.Cleanup(msgID, false)
	}
	if s.msgID != msgID {
		s.msgID = msgID
		s.MarkDirty("subject.SCCoSubject.MessageID")
	}
	return err
}

// SubjectHandling returns the handling order encoded in the subject line.
// If the subject line contained a known handling order code, SubjectHandling
// returns the corresponding full word.
func (s *SCCoSubject) SubjectHandling() string {
	if long := handlingCodes[s.handling]; long != "" {
		return long
	}
	return s.handling
}

// SetSubjectHandling sets the handling order encoded in the subject line.  If
// it is sent to a known handling order word, the corresponding code is
// encoded in the subject line.
func (s *SCCoSubject) SetSubjectHandling(handling string) (err error) {
	if handling == "" {
		err = ErrNoHandling
	} else {
		var found bool

		handling = strings.ToUpper(handling)
		for code, long := range handlingCodes {
			if long == handling {
				handling = code
				found = true
				break
			} else if code == handling {
				found = true
				break
			}
		}
		if !found {
			if oldSeverityRE.MatchString(handling) {
				err = ErrHasSeverity
			} else {
				handling = strings.Map(removeNewlineUnderline, handling)
				err = ErrInvalidHandling
			}
		}
	}
	if s.handling != handling {
		s.handling = handling
		s.MarkDirty("subject.SCCoSubject.Handling")
	}
	return err
}

// SubjectSummary returns the message summary encoded in the subject line.
func (s *SCCoSubject) SubjectSummary() string { return s.summary }

// SetSubjectSummary sets the message summary encoded in the subject line.
func (s *SCCoSubject) SetSubjectSummary(summary string) (err error) {
	if summary == "" {
		err = ErrNoSummary
	} else if s := strings.Map(removeNewlines, summary); len(s) < len(summary) {
		err = ErrNewlineInSubject
		summary = s
	}
	if s.summary != summary {
		s.summary = summary
		s.MarkDirty("subject.SCCoSubject.Summary")
	}
	return err
}

// Validate validates the standards compliance of the fields of the subject.
// If packet is true, the message ID is checked to be sure it has a suffix.
func (s *SCCoSubject) Validate(packet bool) (err error) {
	_, _, _, err = messageid.Decode(s.msgID, false, packet)
	err = errors.AddPrefix(err, "On the subject line: ")
	if handlingCodes[s.handling] == "" {
		if s.handling == "" {
			err = errors.Join(err, ErrNoHandling)
		} else if oldSeverityRE.MatchString(s.handling) {
			err = errors.Join(err, ErrHasSeverity)
		} else {
			err = errors.Join(err, ErrInvalidHandling)
		}
	}
	if s.summary == "" {
		err = errors.Join(err, ErrNoSummary)
	}
	return err
}

func (s *SCCoSubject) Clone() Subject {
	ns, _ := NewSCCoSubject(s.SubjectMessageID(), s.SubjectHandling(), s.SubjectSummary())
	return ns
}

func init() { RegisterDecoder(decodeSCCoSubject) }
func decodeSCCoSubject(subject string) (_ Subject, err error) {
	var (
		fields string
		rest   string
	)
	if idx := strings.IndexByte(subject, ' '); idx >= 0 {
		fields, rest = subject[:idx], subject[idx:]
	} else {
		fields = subject
	}
	parts := strings.SplitN(fields, "_", 3)
	if len(parts) < 3 {
		return nil, nil
	}
	s := &SCCoSubject{
		encoded:  subject,
		msgID:    parts[0],
		handling: parts[1],
		summary:  parts[2] + rest,
	}
	return s, nil
}
