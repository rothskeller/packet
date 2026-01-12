// Package subject defines the Subject interface, which represents a subject
// line of a packet message.
package subject

import (
	"fmt"
	"iter"
	"regexp"
	"slices"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/cachetrack"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/messageid"
	"github.com/rothskeller/packet/message/msgifc"
)

// Subject is the interface satisfied by all subject types.
type Subject = msgifc.Subject

// A PlainSubject is a subject line for a packet message that does not contain
// a form.  It should comply with the Santa Clara County standard for packet
// message subject lines.
type PlainSubject struct {
	cachetrack.Tracker
	encoded  string
	msgID    string
	handling string
	summary  string
}

var _ Subject = (*PlainSubject)(nil)

var (
	handlingCodes = map[string]string{"I": "IMMEDIATE", "P": "PRIORITY", "R": "ROUTINE"}
	oldSeverityRE = regexp.MustCompile(`^[A-Z]/[A-Z]$`)
)

// NewPlainSubject creates a new SCCo-standard non-form message subject with
// the specified parameters.  It returns any problems with the parameters.
func NewPlainSubject(msgID, handling, summary string) (s *PlainSubject, err error) {
	s = new(PlainSubject)
	if msgID == "" && handling == "" {
		err = s.SetSubjectSummary(summary)
	} else {
		err = errors.Join(
			s.SetSubjectMessageID(msgID),
			s.SetSubjectHandling(handling),
			s.SetSubjectSummary(summary),
		)
	}
	return s, err
}

// EncodedSubject returns the encoded subject line.
func (s *PlainSubject) EncodedSubject() string {
	if s.Dirty() {
		if s.msgID == "" && s.handling == "" {
			s.encoded = s.summary
		} else {
			s.encoded = fmt.Sprintf("%s_%s_%s", s.msgID, s.handling, s.summary)
		}
		s.MarkClean()
	}
	return s.encoded
}

// SubjectMessageID returns the message ID encoded in the subject line.
func (s *PlainSubject) SubjectMessageID() string { return s.msgID }

// SetSubjectMessageID sets the message ID encoded in the subject line.
func (s *PlainSubject) SetSubjectMessageID(msgID string) (err error) {
	if trim := strings.Map(removeNewlineUnderline, msgID); len(trim) < len(msgID) {
		err = errors.New("Invalid characters were removed from the message ID on the subject line.")
		msgID = trim
	} else {
		if msgID, err = messageid.Cleanup(msgID, false); err != nil {
			err = errors.AddPrefix(err, "On the subject line: ")
		}
	}
	if s.msgID != msgID {
		s.msgID = msgID
		s.MarkDirty("subject.PlainSubject.MessageID")
	}
	return err
}

// SubjectHandling returns the handling order encoded in the subject line.
// If the subject line contained a known handling order code, SubjectHandling
// returns the corresponding full word.
func (s *PlainSubject) SubjectHandling() string {
	if long := handlingCodes[s.handling]; long != "" {
		return long
	}
	return s.handling
}

// SetSubjectHandling sets the handling order encoded in the subject line.  If
// it is sent to a known handling order word, the corresponding code is
// encoded in the subject line.
func (s *PlainSubject) SetSubjectHandling(handling string) (err error) {
	if handling == "" {
		err = errors.New("The subject line does not have a handling order code.")
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
				err = errors.New("The subject line has an old-style severity code, which is no longer part of the County standard.")
			} else {
				handling = strings.Map(removeNewlineUnderline, handling)
				err = errors.New("The handling order on the subject line is not one of the standard handling order codes (I, P, or R).")
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
func (s *PlainSubject) SubjectSummary() string { return s.summary }

// SetSubjectSummary sets the message summary encoded in the subject line.
func (s *PlainSubject) SetSubjectSummary(summary string) (err error) {
	if summary == "" {
		err = errors.New("The subject line does not have a message summary.")
	} else if s := strings.Map(removeNewlines, summary); len(s) < len(summary) {
		err = errors.New("Invalid characters were removed from the message summary on the subject line.")
		summary = s
	}
	if s.summary != summary {
		s.summary = summary
		s.MarkDirty("subject.PlainSubject.Summary")
	}
	return err
}

// Clone creates a copy of the subject.
func (s *PlainSubject) Clone() Subject {
	ns, _ := NewPlainSubject(s.SubjectMessageID(), s.SubjectHandling(), s.SubjectSummary())
	return ns
}

func (s *PlainSubject) Fields() iter.Seq[msgifc.Field] {
	return slices.Values(subjectFields)
}

func DecodePlainSubject(subject string) (s *PlainSubject) {
	var (
		fields string
		rest   string
	)
	s = &PlainSubject{encoded: subject}
	if idx := strings.IndexByte(subject, ' '); idx >= 0 {
		fields, rest = subject[:idx], subject[idx:]
	} else {
		fields = subject
	}
	parts := strings.SplitN(fields, "_", 3)
	if len(parts) < 3 {
		s.summary = subject
	} else {
		s.msgID = parts[0]
		s.handling = parts[1]
		s.summary = parts[2] + rest
	}
	return s
}

var subjectFields = []msgifc.Field{
	field.NewMessageID("", "Message ID").
		Common(field.CSubjectMessageID).
		ValueFunc(func(m msgifc.Message) string { return m.Subject().SubjectMessageID() }).
		SetValueFunc(func(m msgifc.Message, s string) { m.Subject().SetSubjectMessageID(s) }).
		EditHelp("This is the message ID assigned to the message by the originating station.  It has the form XXX-###P, where XXX is the three-character prefix associated with the originating station, ### is a unique number, and P is a suffix letter.").
		ValidateFunc(func(m msgifc.Message, f msgifc.Field, vf msgifc.ValidateFlags) error {
			if m.Bulletin() && f.Value(m) == "" {
				return nil
			}
			return field.ValidateMessageID(m, f, vf)
		}).
		MakeField(),
	field.NewField("", "Handling").
		Common(field.CSubjectHandling).
		AllowedValues("ROUTINE", "PRIORITY", "IMMEDIATE").
		EditHelp("This is the handling order for the message, one of ROUTINE, PRIORITY, or IMMEDIATE.").
		ValueFunc(func(m msgifc.Message) string { return m.Subject().SubjectHandling() }).
		SetValueFunc(func(m msgifc.Message, s string) { m.Subject().SetSubjectHandling(s) }).
		ValidateFunc(func(m msgifc.Message, f msgifc.Field, vf msgifc.ValidateFlags) error {
			if vf&msgifc.VPIFOOnly != 0 {
				return nil
			}
			if h := f.Value(m); m.Bulletin() && h == "" {
				return nil
			} else if handlingCodes[h] != "" {
				return nil
			} else if h == "" {
				return errors.New("The subject line does not have a handling order code.")
			} else if oldSeverityRE.MatchString(h) {
				return errors.New("The subject line has an old-style severity code, which is no longer part of the County standard.")
			} else {
				return errors.New("The handling order on the subject line is not one of the standard handling order codes (I, P, or R).")
			}
		}).
		MakeField(),
	field.NewField("", "Message Summary").
		Common(field.CSubjectSummary).
		EditHelp("This is a brief summary of the content of the message: the content of the subject line after the encoded message ID and handling order.").
		ValueFunc(func(m msgifc.Message) string { return m.Subject().SubjectSummary() }).
		SetValueFunc(func(m msgifc.Message, s string) { m.Subject().SetSubjectSummary(s) }).
		ValidateFunc(func(m msgifc.Message, f msgifc.Field, vf msgifc.ValidateFlags) error {
			if vf&msgifc.VPIFOOnly != 0 {
				return nil
			}
			if f.Value(m) == "" {
				return errors.New("The subject line does not have a message summary.")
			}
			return nil
		}).
		MakeField(),
}

func removeNewlines(r rune) rune {
	if r == '\r' || r == '\n' {
		return -1
	}
	return r
}

func removeNewlineUnderline(r rune) rune {
	if r == '\r' || r == '\n' || r == '_' {
		return -1
	}
	return r
}
