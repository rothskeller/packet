package pseudomsg

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
)

// A ConfigMessage is a fake message.Message that implements only Fields, used
// for manipulating incident configurations in the command line as if they were
// messages.
type ConfigMessage struct{ Config *incident.Config }

var _ message.Message = (*ConfigMessage)(nil)

var configFields []field.Field
var configFieldsOnce sync.Once

func (lem *ConfigMessage) Fields() iter.Seq[field.Field] {
	configFieldsOnce.Do(func() {
		configFields = []field.Field{
			field.NewField("", "Incident Name").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.IncidentName
				}).MakeField(),
			field.NewField("", "Activation Number").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.ActivationNum
				}).MakeField(),
			field.NewField("", "Operation Start").
				ToHumanFunc(func(_ message.Message, s string) string {
					if t, err := time.ParseInLocation(time.RFC3339, s, time.Local); err == nil {
						return t.Format("01/02/2006 15:04")
					} else {
						return s
					}
				}).
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.OpStart.Format(time.RFC3339)
				}).
				MakeField(),
			field.NewField("", "Operation End").
				ToHumanFunc(func(_ message.Message, s string) string {
					if t, err := time.ParseInLocation(time.RFC3339, s, time.Local); err == nil {
						return t.Format("01/02/2006 15:04")
					} else {
						return s
					}
				}).
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.OpEnd.Format(time.RFC3339)
				}).
				MakeField(),
			field.NewField("", "Operator Call Sign").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.OpCall
				}).MakeField(),
			field.NewField("", "Operator Name").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.OpName
				}).MakeField(),
			field.NewField("", "Tactical Call Sign").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.TacCall
				}).MakeField(),
			field.NewField("", "Tactical Station Name").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.TacName
				}).MakeField(),
			field.NewField("", "Tx Message ID").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.TxMessageID
				}).MakeField(),
			field.NewField("", "Rx Message ID").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.RxMessageID
				}).MakeField(),
			field.NewField("", "Connection Type").
				ToHumanFunc(func(_ message.Message, s string) string {
					switch s {
					case incident.ConnectNone:
						return "Manual"
					case incident.ConnectSerialTNC:
						return "Serial+TNC"
					case incident.ConnectTelnet:
						return "Telnet"
					default:
						return s
					}
				}).
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.ConnectType
				}).MakeField(),
			field.NewField("", "BBS Call Sign").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.ConnectBBS
				}).MakeField(),
			field.NewField("", "Automatic Receipts").
				ValueFunc(func(m message.Message) string {
					if m.(*ConfigMessage).Config.NoSendReceipts {
						return "checked"
					}
					return ""
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType == incident.ConnectNone
				}).MakeField(),
			field.NewField("", "TNC Type").
				ToHumanFunc(func(_ message.Message, s string) string {
					switch s {
					case incident.TNCKPC3Plus:
						return "Kantronics KPC-3 Plus"
					default:
						return s
					}
				}).
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.TNCType
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType == incident.ConnectSerialTNC
				}).MakeField(),
			field.NewField("", "Serial Port").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.SerialPort
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType == incident.ConnectSerialTNC
				}).MakeField(),
			field.NewField("", "BBS Address").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.ConnectAddress
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectNone
				}).MakeField(),
			field.NewField("", "Bulletin Check 1").
				ValueFunc(func(m message.Message) string {
					bc := m.(*ConfigMessage).Config.BulletinChecks
					areas := slices.Collect(maps.Keys(bc))
					if len(areas) >= 1 {
						slices.Sort(areas)
						return areas[0] + " every " + fmtCheckFrequency(bc[areas[0]])
					}
					return ""
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectNone
				}).MakeField(),
			field.NewField("", "Bulletin Check 2").
				ValueFunc(func(m message.Message) string {
					bc := m.(*ConfigMessage).Config.BulletinChecks
					areas := slices.Collect(maps.Keys(bc))
					if len(areas) >= 2 {
						slices.Sort(areas)
						return areas[1] + " every " + fmtCheckFrequency(bc[areas[1]])
					}
					return ""
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectNone
				}).MakeField(),
			field.NewField("", "Bulletin Check 3").
				ValueFunc(func(m message.Message) string {
					bc := m.(*ConfigMessage).Config.BulletinChecks
					areas := slices.Collect(maps.Keys(bc))
					if len(areas) >= 3 {
						slices.Sort(areas)
						return areas[2] + " every " + fmtCheckFrequency(bc[areas[2]])
					}
					return ""
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectNone
				}).MakeField(),
			field.NewField("", "Bulletin Check 4").
				ValueFunc(func(m message.Message) string {
					bc := m.(*ConfigMessage).Config.BulletinChecks
					areas := slices.Collect(maps.Keys(bc))
					if len(areas) >= 4 {
						slices.Sort(areas)
						return areas[3] + " every " + fmtCheckFrequency(bc[areas[3]])
					}
					return ""
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectNone
				}).MakeField(),
			field.NewField("", "Bulletin Check 5").
				ValueFunc(func(m message.Message) string {
					bc := m.(*ConfigMessage).Config.BulletinChecks
					areas := slices.Collect(maps.Keys(bc))
					if len(areas) >= 5 {
						slices.Sort(areas)
						return areas[4] + " every " + fmtCheckFrequency(bc[areas[4]])
					}
					return ""
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectNone
				}).MakeField(),
			field.NewField("", "Default To Address").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultTo
				}).MakeField(),
			field.NewField("", "Default To ICS Position").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultToPos
				}).MakeField(),
			field.NewField("", "Default To Location").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultToLoc
				}).MakeField(),
			field.NewField("", "Default From ICS Position").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultFromPos
				}).MakeField(),
			field.NewField("", "Default From Location").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultFromLoc
				}).MakeField(),
			field.NewField("", "Default Body").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultBody
				}).MakeField(),
		}
	})
	return slices.Values(configFields)
}

func (lem *ConfigMessage) Dirty() bool                  { panic("not implemented") }
func (lem *ConfigMessage) MarkDirty(_ string)           { panic("not implemented") }
func (lem *ConfigMessage) MarkClean()                   { panic("not implemented") }
func (lem *ConfigMessage) OnDirty(_ func(string))       { panic("not implemented") }
func (lem *ConfigMessage) Recognize(_ message.Message)  { panic("not implemented") }
func (lem *ConfigMessage) Name() string                 { panic("not implemented") }
func (lem *ConfigMessage) Type() message.MType          { panic("not implemented") }
func (lem *ConfigMessage) SetType(_ message.MType)      { panic("not implemented") }
func (lem *ConfigMessage) RFC5322() string              { panic("not implemented") }
func (lem *ConfigMessage) To() string                   { panic("not implemented") }
func (lem *ConfigMessage) Subject() subject.Subject     { panic("not implemented") }
func (lem *ConfigMessage) SetSubject(_ subject.Subject) { panic("not implemented") }
func (lem *ConfigMessage) Bulletin() bool               { panic("not implemented") }
func (lem *ConfigMessage) Payload() payload.Payload     { panic("not implemented") }
func (lem *ConfigMessage) Body() body.Body              { panic("not implemented") }
func (lem *ConfigMessage) RenderPDF(m message.Message, filename string, copyname string) error {
	panic("not implemented")
}

func fmtCheckFrequency(cf incident.CheckFrequency) (s string) {
	dur := cf.Duration

	if dur >= time.Hour {
		hr := dur / time.Hour
		s = fmt.Sprintf("%dh", hr)
		dur -= hr * time.Hour
	}
	if dur >= time.Minute {
		min := dur / time.Minute
		s += fmt.Sprintf("%dm", min)
	}
	return s
}
