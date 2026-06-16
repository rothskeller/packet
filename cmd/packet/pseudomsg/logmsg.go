package pseudomsg

import (
	"fmt"
	"iter"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/rothskeller/packet/v4/message"
	"github.com/rothskeller/packet/v4/message/body"
	"github.com/rothskeller/packet/v4/message/field"
	"github.com/rothskeller/packet/v4/message/msgifc"
	"github.com/rothskeller/packet/v4/message/payload"
	"github.com/rothskeller/packet/v4/message/subject"
)

// A LogEntryMessage is a fake message.Message that implements only Fields, used
// for manipulating log entries in the command line as if they were messages.
type LogEntryMessage struct {
	LogEntry *incident.LogEntry
	timestr  string
}

var _ message.Message = (*LogEntryMessage)(nil)

// NewLogEntryMessage constructs a LogEntryMessage for a LogEntry.
func NewLogEntryMessage(le *incident.LogEntry) (lem *LogEntryMessage) {
	return &LogEntryMessage{LogEntry: le, timestr: le.Time.Format("01/02/2006 15:04")}
}

var (
	dateLooseRE = regexp.MustCompile(`^(0?[1-9]|1[0-2])[-./](0?[1-9]|[12][0-9]|3[01])[-./](?:20)?([0-9][0-9])$`)
	timeLooseRE = regexp.MustCompile(`^([1-9]|[01][0-9]|2[0-4]):?([0-5][0-9])$`)
	dateTimeRE  = regexp.MustCompile(`^(?:0[1-9]|1[0-2])/(?:0[1-9]|[12][0-9]|3[01])/20[0-9][0-9] (?:[01][0-9]|2[0-4]):[0-5][0-9]$`)
)

var logEntryFields []field.Field
var logEntryFieldsOnce sync.Once

func (lem *LogEntryMessage) Fields() iter.Seq[field.Field] {
	logEntryFieldsOnce.Do(func() {
		logEntryFields = []field.Field{
			field.NewField("", "Time").
				ValueFunc(func(m message.Message) string { return lem.timestr }).
				EditHelp(`This is the date and time of the log entry, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.  (The date is not shown, but affects the sorting of entries.)`).
				EditWidth(16).
				FromHumanFunc(func(m msgifc.Message, s string) string {
					parts := strings.Fields(s)
					if len(parts) > 0 {
						if match := dateLooseRE.FindStringSubmatch(parts[0]); match != nil {
							m, d, y := match[1], match[2], match[3]
							if len(y) == 2 {
								y = "20" + y
							}
							parts[0] = fmt.Sprintf("%02s/%02s/%s", m, d, y)
						}
					}
					if len(parts) > 1 {
						if match := timeLooseRE.FindStringSubmatch(parts[1]); match != nil {
							h, m := match[1], match[2]
							parts[1] = fmt.Sprintf("%02s:%s", h, m)
						}
					}
					return strings.Join(parts, " ")
				}).
				Required().
				SetValueFunc(func(m msgifc.Message, s string) {
					m.(*LogEntryMessage).timestr = s
					if t, err := time.ParseInLocation("01/02/2006 15:04", s, time.Local); err == nil {
						m.(*LogEntryMessage).LogEntry.Time = t
					}
				}).
				ValidateFunc(func(m msgifc.Message, f msgifc.Field, vf msgifc.ValidateFlags) error {
					if !dateTimeRE.MatchString(m.(*LogEntryMessage).timestr) {
						return NoBypassValidationError{errors.New("The date/time string must be in MM/DD/YYYY HH:MM format (24-hour clock).")}
					}
					return nil
				}).
				MakeField(),
			field.NewField("", "From Station").
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.FromCall
				}).
				FromHumanFunc(func(_ msgifc.Message, s string) string { return strings.ToUpper(strings.TrimSpace(s)) }).
				SetValueFunc(func(m msgifc.Message, s string) {
					m.(*LogEntryMessage).LogEntry.FromCall = s
				}).
				EditHelp(`This is the call sign of the station that originated the message.`).
				EditWidth(ics309ColumnWidth(1)).
				MakeField(),
			field.NewField("", "From Msg #").
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.FromMsgID
				}).
				FromHumanFunc(func(_ msgifc.Message, s string) string { return strings.ToUpper(strings.TrimSpace(s)) }).
				SetValueFunc(func(m msgifc.Message, s string) {
					m.(*LogEntryMessage).LogEntry.FromMsgID = s
				}).
				EditHelp(`This is the message number assigned by the originating station.`).
				EditWidth(ics309ColumnWidth(2)).
				MakeField(),
			field.NewField("", "To Station").
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.ToCall
				}).
				FromHumanFunc(func(_ msgifc.Message, s string) string { return strings.ToUpper(strings.TrimSpace(s)) }).
				SetValueFunc(func(m msgifc.Message, s string) {
					m.(*LogEntryMessage).LogEntry.ToCall = s
				}).
				EditHelp(`This is the call sign of the destination station.`).
				EditWidth(ics309ColumnWidth(3)).
				MakeField(),
			field.NewField("", "To Msg #").
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.ToMsgID
				}).
				FromHumanFunc(func(_ msgifc.Message, s string) string { return strings.ToUpper(strings.TrimSpace(s)) }).
				SetValueFunc(func(m msgifc.Message, s string) {
					m.(*LogEntryMessage).LogEntry.ToMsgID = s
				}).
				EditHelp(`This is the message number assigned by the destination station.`).
				EditWidth(ics309ColumnWidth(4)).
				MakeField(),
			field.NewField("", "Message").
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.Subject
				}).
				SetValueFunc(func(m msgifc.Message, s string) {
					m.(*LogEntryMessage).LogEntry.Subject = s
				}).
				EditHelp(`This is the description or subject line of the message.`).
				EditWidth(ics309ColumnWidth(5)).
				MakeField(),
			field.NewField("", "Flags").
				ValueFunc(func(m msgifc.Message) string {
					var s []string
					var le = m.(*LogEntryMessage).LogEntry
					if le.Flags&incident.FFollowup != 0 {
						s = append(s, "Needs Followup")
					}
					if le.Flags&incident.FNeedsReceipt != 0 {
						s = append(s, "Needs Receipt")
					}
					if le.Flags&incident.FVoice != 0 {
						s = append(s, "Voice")
					}
					return strings.Join(s, ", ")
				}).MakeField(),
			field.NewField("", "Needs Followup").
				AllowedValues("checked").
				EditHelp(`This flag indicates that the message needs human follow-up.`).
				SetValueFunc(func(m msgifc.Message, s string) {
					if s != "" {
						m.(*LogEntryMessage).LogEntry.Flags |= incident.FFollowup
					} else {
						m.(*LogEntryMessage).LogEntry.Flags &^= incident.FFollowup
					}
				}).
				ValueFunc(func(m msgifc.Message) string {
					if m.(*LogEntryMessage).LogEntry.Flags&incident.FFollowup != 0 {
						return "checked"
					}
					return ""
				}).
				VisibleWhen(field.Invisible).
				MakeField(),
			field.NewField("", "Needs Receipt").
				AllowedValues("checked").
				EditHelp(`This flag indicates that a delivery receipt is expected and has not been received.`).
				EditableWhen(func(m msgifc.Message, _ bool) bool {
					return m.(*LogEntryMessage).LogEntry.Status == incident.StatusSent
				}).
				SetValueFunc(func(m msgifc.Message, s string) {
					if s != "" {
						m.(*LogEntryMessage).LogEntry.Flags |= incident.FNeedsReceipt
					} else {
						m.(*LogEntryMessage).LogEntry.Flags &^= incident.FNeedsReceipt
					}
				}).
				ValueFunc(func(m msgifc.Message) string {
					if m.(*LogEntryMessage).LogEntry.Flags&incident.FNeedsReceipt != 0 {
						return "checked"
					}
					return ""
				}).
				VisibleWhen(field.Invisible).
				MakeField(),
			field.NewField("", "Voice").
				AllowedValues("checked").
				EditHelp(`This flag indicates that the message should be logged on a separate "voice" ICS-309 form.`).
				SetValueFunc(func(m msgifc.Message, s string) {
					if s != "" {
						m.(*LogEntryMessage).LogEntry.Flags |= incident.FVoice
					} else {
						m.(*LogEntryMessage).LogEntry.Flags &^= incident.FVoice
					}
				}).
				ValueFunc(func(m msgifc.Message) string {
					if m.(*LogEntryMessage).LogEntry.Flags&incident.FVoice != 0 {
						return "checked"
					}
					return ""
				}).
				VisibleWhen(field.Invisible).
				MakeField(),
		}
	})
	return slices.Values(logEntryFields)
}

func (lem *LogEntryMessage) Dirty() bool                  { panic("not implemented") }
func (lem *LogEntryMessage) MarkDirty(_ string)           { panic("not implemented") }
func (lem *LogEntryMessage) MarkClean()                   { panic("not implemented") }
func (lem *LogEntryMessage) OnDirty(_ func(string))       { panic("not implemented") }
func (lem *LogEntryMessage) Recognize(_ message.Message)  { panic("not implemented") }
func (lem *LogEntryMessage) Tag() string                  { panic("not implemented") }
func (lem *LogEntryMessage) Name() string                 { panic("not implemented") }
func (lem *LogEntryMessage) Type() message.MType          { panic("not implemented") }
func (lem *LogEntryMessage) SetType(_ message.MType)      { panic("not implemented") }
func (lem *LogEntryMessage) RFC5322() string              { panic("not implemented") }
func (lem *LogEntryMessage) To() string                   { panic("not implemented") }
func (lem *LogEntryMessage) Subject() subject.Subject     { panic("not implemented") }
func (lem *LogEntryMessage) SetSubject(_ subject.Subject) { panic("not implemented") }
func (lem *LogEntryMessage) Bulletin() bool               { panic("not implemented") }
func (lem *LogEntryMessage) Payload() payload.Payload     { panic("not implemented") }
func (lem *LogEntryMessage) Body() body.Body              { panic("not implemented") }
func (lem *LogEntryMessage) RenderPDF(m message.Message, filename string, copyname string) error {
	panic("not implemented")
}
func (lem *LogEntryMessage) CanCompareAgainst(_ message.MType) bool { return false }

// A NoBypassValidationError is a validation error handled specially by the
// editor: the value cannot be accepted.
type NoBypassValidationError struct{ err error }

func (ne NoBypassValidationError) Error() string { return ne.err.Error() }
func (ne NoBypassValidationError) Unwrap() error { return ne.err }

func ics309ColumnWidth(col int) int {
	for f := range incident.ICS309FormDef().AllFields() {
		if f.Label != "Line1" {
			continue
		}
		return f.CharWidth(col)
	}
	return 0
}
