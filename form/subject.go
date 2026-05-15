package form

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
	"github.com/rothskeller/packet/message/subject"
)

var (
	handlingCodes = map[string]string{"I": "IMMEDIATE", "P": "PRIORITY", "R": "ROUTINE"}
	oldSeverityRE = regexp.MustCompile(`^[A-Z]/[A-Z]$`)
)

// A FormSubject is a subject that is encoded according to the Santa Clara
// County standard for packet forms message subject lines.
type FormSubject struct {
	cachetrack.Tracker
	encoded  string
	msgID    string
	handling string
	formtag  string
	summary  string
}

var _ msgifc.Subject = (*FormSubject)(nil)

// NewFormSubject creates a new SCCo-standard form message subject with the
// specified parameters.  It returns any problems with the parameters.
func NewFormSubject(msgID, handling, formtag, summary string) (s *FormSubject, err error) {
	s = new(FormSubject)
	err = errors.Join(
		s.SetSubjectMessageID(msgID),
		s.SetSubjectHandling(handling),
		s.SetSubjectFormTag(formtag),
		s.SetSubjectSummary(summary),
	)
	return s, err
}

// formSubjectFromPlainSubject converts a PlainSubject to a FormSubject.
func formSubjectFromPlainSubject(ps *subject.PlainSubject) (fs *FormSubject) {
	fs = new(FormSubject)
	summary := ps.SubjectSummary()
	first, rest, space := strings.Cut(summary, " ")
	tag, nontag, found := strings.Cut(first, "_")
	if found {
		if space {
			summary = nontag + " " + rest
		} else {
			summary = nontag
		}
	} else {
		tag = ""
	}
	fs.SetSubjectMessageID(ps.SubjectMessageID())
	fs.SetSubjectHandling(ps.SubjectHandling())
	fs.SetSubjectFormTag(tag)
	fs.SetSubjectSummary(summary)
	fs.EncodedSubject() // make sure result is marked clean
	return fs
}

// EncodedSubject returns the encoded subject line.
func (s *FormSubject) EncodedSubject() string {
	if s.Dirty() {
		s.encoded = strings.TrimSpace(fmt.Sprintf("%s_%s_%s_%s", s.msgID, s.handling, s.formtag, s.summary))
		s.MarkClean()
	}
	return s.encoded
}

// SubjectMessageID returns the message ID encoded in the subject line.
func (s *FormSubject) SubjectMessageID() string { return s.msgID }

// SetSubjectMessageID sets the message ID encoded in the subject line.
func (s *FormSubject) SetSubjectMessageID(msgID string) (err error) {
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
		s.MarkDirty("form.FormSubject.MessageID")
	}
	return err
}

// SubjectHandling returns the handling order encoded in the subject line.
// If the subject line contained a known handling order code, SubjectHandling
// returns the corresponding full word.
func (s *FormSubject) SubjectHandling() string {
	if long := handlingCodes[s.handling]; long != "" {
		return long
	}
	return s.handling
}

// SetSubjectHandling sets the handling order encoded in the subject line.  If
// it is sent to a known handling order word, the corresponding code is
// encoded in the subject line.
func (s *FormSubject) SetSubjectHandling(handling string) (err error) {
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
				err = errors.NewF("The handling order on the subject line (%q) is not one of the standard handling order codes (I, P, or R).", handling)
			}
		}
	}
	if s.handling != handling {
		s.handling = handling
		s.MarkDirty("form.SCCoSubject.Handling")
	}
	return err
}

// SubjectFormTag returns the form tag encoded in the subject line.
func (s *FormSubject) SubjectFormTag() string { return s.formtag }

// SetSubjectFormTag sets the form tag encoded in the subject line.
func (s *FormSubject) SetSubjectFormTag(formtag string) (err error) {
	if formtag == "" {
		err = errors.New("The subject line does not have a form tag.")
	} else if trim := strings.Map(removeNewlineUnderline, formtag); len(trim) < len(formtag) {
		err = errors.New("Invalid characters were removed from the form tag on the subject line.")
		formtag = trim
	}
	if s.msgID != formtag {
		s.formtag = formtag
		s.MarkDirty("form.FormSubject.FormTag")
	}
	return err
}

// SubjectSummary returns the message summary encoded in the subject line.
func (s *FormSubject) SubjectSummary() string { return s.summary }

// SetSubjectSummary sets the message summary encoded in the subject line.
func (s *FormSubject) SetSubjectSummary(summary string) (err error) {
	if summary == "" {
		err = errors.New("The subject line does not have a message summary.")
	} else if s := strings.Map(removeNewlines, summary); len(s) < len(summary) {
		err = errors.New("Invalid characters were removed from the message summary on the subject line.")
		summary = s
	}
	if s.summary != summary {
		s.summary = summary
		s.MarkDirty("form.FormSubject.Summary")
	}
	return err
}

// Clone creates a copy of the subject.
func (s *FormSubject) Clone() msgifc.Subject {
	ns := *s
	ns.Tracker = cachetrack.Tracker{}
	return &ns
}

func (s *FormSubject) Fields() iter.Seq[msgifc.Field] {
	return slices.Values(subjectFields)
}

func DecodeFormSubject(subject string) (s *FormSubject) {
	var (
		fields string
		rest   string
	)
	s = &FormSubject{encoded: subject}
	if idx := strings.IndexByte(subject, ' '); idx >= 0 {
		fields, rest = subject[:idx], subject[idx:]
	} else {
		fields = subject
	}
	parts := strings.SplitN(fields, "_", 4)
	if len(parts) < 4 {
		s.summary = subject
	} else {
		s.msgID = parts[0]
		s.handling = parts[1]
		s.formtag = parts[2]
		s.summary = parts[3] + rest
	}
	return s
}

var subjectFields = []msgifc.Field{
	field.NewMessageID("", "Message ID").
		Common(field.CSubjectMessageID).
		ValueFunc(func(m msgifc.Message) string { return m.Subject().SubjectMessageID() }).
		SetValueFunc(func(m msgifc.Message, s string) {
			m.Subject().SetSubjectMessageID(s)
			m.Subject().MarkDirty("from.SCCoSubject.MessageID")
		}).
		VisibleWhen(field.Invisible).
		ValidateFunc(func(m msgifc.Message, f msgifc.Field, vf msgifc.ValidateFlags) error {
			if m.Bulletin() && f.Value(m) == "" {
				return nil
			}
			return errors.AddPrefix(field.ValidateMessageID(m, f, vf), "On the subject line: ")
		}).
		CompareFunc(field.CompareNone).
		MakeField(),
	field.NewField("", "Handling").
		Common(field.CSubjectHandling).
		AllowedValues("ROUTINE", "PRIORITY", "IMMEDIATE").
		ValueFunc(func(m msgifc.Message) string { return m.Subject().SubjectHandling() }).
		SetValueFunc(func(m msgifc.Message, s string) {
			m.Subject().SetSubjectHandling(s)
			m.Subject().MarkDirty("from.SCCoSubject.Handling")
		}).
		VisibleWhen(field.Invisible).
		ValidateFunc(func(m msgifc.Message, f msgifc.Field, vf msgifc.ValidateFlags) error {
			if vf&msgifc.VPIFOOnly != 0 {
				return nil
			}
			if h := f.Value(m); m.Bulletin() && h == "" {
				return nil
			} else if h == "ROUTINE" || h == "PRIORITY" || h == "IMMEDIATE" {
				return nil
			} else if h == "" {
				return errors.New("The subject line does not have a handling order code.")
			} else if oldSeverityRE.MatchString(h) {
				return errors.New("The subject line has an old-style severity code, which is no longer part of the County standard.")
			} else {
				return errors.NewF("The handling order on the subject line (%q) is not one of the standard handling order codes (I, P, or R).", h)
			}
		}).
		CompareFunc(field.CompareNone).
		MakeField(),
	field.NewField("", "Form Tag").
		Common(field.CSubjectFormTag).
		ValueFunc(func(m msgifc.Message) string { return m.Subject().(*FormSubject).SubjectFormTag() }).
		VisibleWhen(field.Invisible).
		ValidateFunc(func(m msgifc.Message, f msgifc.Field, vf msgifc.ValidateFlags) error {
			if vf&msgifc.VPIFOOnly != 0 {
				return nil
			}
			if f.Value(m) == "" {
				return errors.New("The subject line does not have a form tag.")
			}
			return nil
		}).
		CompareFunc(field.CompareNone).
		MakeField(),
	field.NewField("", "Message Summary").
		Common(field.CSubjectSummary).
		ValueFunc(func(m msgifc.Message) string { return m.Subject().SubjectSummary() }).
		SetValueFunc(func(m msgifc.Message, s string) {
			m.Subject().SetSubjectSummary(s)
			m.Subject().MarkDirty("from.SCCoSubject.Summary")
		}).
		VisibleWhen(field.Invisible).
		ValidateFunc(func(m msgifc.Message, f msgifc.Field, vf msgifc.ValidateFlags) error {
			if vf&msgifc.VPIFOOnly != 0 {
				return nil
			}
			if f.Value(m) == "" {
				return errors.New("The subject line does not have a message summary.")
			}
			return nil
		}).
		CompareFunc(field.CompareNone).
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
