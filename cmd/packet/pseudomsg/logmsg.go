package pseudomsg

import (
	"iter"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/msgifc"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
)

// A LogEntryMessage is a fake message.Message that implements only Fields, used
// for manipulating log entries in the command line as if they were messages.
type LogEntryMessage struct{ LogEntry *incident.LogEntry }

var _ message.Message = (*LogEntryMessage)(nil)

var (
	dateLooseRE = regexp.MustCompile(`^(0?[1-9]|1[0-2])[-./](0?[1-9]|[12][0-9]|3[01])[-./](?:20)?([0-9][0-9])$`)
	timeLooseRE = regexp.MustCompile(`^([1-9]:|[01][0-9]:?|2[0-4]:?)([0-5][0-9])$`)
)

var logEntryFields []field.Field
var logEntryFieldsOnce sync.Once

func (lem *LogEntryMessage) Fields() iter.Seq[field.Field] {
	logEntryFieldsOnce.Do(func() {
		logEntryFields = []field.Field{
			field.NewField("", "Time").
				ToHumanFunc(func(s string) string {
					if t, err := time.ParseInLocation(time.RFC3339, s, time.Local); err == nil {
						return t.Format("01/02/2006 15:04")
					} else {
						return s
					}
				}).
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.Time.Format(time.RFC3339)
				}).
				MakeField(),
			field.NewField("", "From Station").
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.FromCall
				}).MakeField(),
			field.NewField("", "From Msg #").
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.FromMsgID
				}).MakeField(),
			field.NewField("", "To Station").
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.ToCall
				}).MakeField(),
			field.NewField("", "To Msg #").
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.ToMsgID
				}).MakeField(),
			field.NewField("", "Message").
				ValueFunc(func(m message.Message) string {
					return m.(*LogEntryMessage).LogEntry.Subject
				}).MakeField(),
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
		}
	})
	return slices.Values(logEntryFields)
}

func (lem *LogEntryMessage) Dirty() bool                  { panic("not implemented") }
func (lem *LogEntryMessage) MarkDirty(_ string)           { panic("not implemented") }
func (lem *LogEntryMessage) MarkClean()                   { panic("not implemented") }
func (lem *LogEntryMessage) OnDirty(_ func(string))       { panic("not implemented") }
func (lem *LogEntryMessage) Recognize(_ message.Message)  { panic("not implemented") }
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
