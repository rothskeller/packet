package pseudomsg

import (
	"cmp"
	"fmt"
	"io/fs"
	"iter"
	"maps"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
	"go.bug.st/serial"
)

// A ConfigMessage is a fake message.Message that implements only Fields, used
// for manipulating incident configurations in the command line as if they were
// messages.
type ConfigMessage struct {
	Config   *incident.Config
	startstr string
	endstr   string
	defrx    bool
	deftx    bool
	bchecks  string
}

var _ message.Message = (*ConfigMessage)(nil)

// NewConfigMessage constructs a ConfigMessage for a Config.
func NewConfigMessage(c *incident.Config) (cm *ConfigMessage) {
	cm = &ConfigMessage{Config: c}
	if !c.OpStart.IsZero() {
		cm.startstr = c.OpStart.Format("01/02/2006 15:04")
	}
	if !c.OpEnd.IsZero() {
		cm.endstr = c.OpEnd.Format("01/02/2006 15:04")
	}
	cm.setDefaultMessageNumbers()
	areas := slices.Collect(maps.Keys(c.BulletinChecks))
	slices.Sort(areas)
	for i, a := range areas {
		areas[i] = fmt.Sprintf("%s:%s", a, fmtCheckFrequency(c.BulletinChecks[areas[i]]))
	}
	cm.bchecks = strings.Join(areas, " ")
	return cm
}

var configFields []field.Field
var configFieldsOnce sync.Once

var comPortRE = regexp.MustCompile(`^COM\d+$`)
var ax25RE = regexp.MustCompile(`^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})-(?:[0-9]|1[0-5])$`)

func (cm *ConfigMessage) Fields() iter.Seq[field.Field] {
	configFieldsOnce.Do(func() {
		configFields = []field.Field{
			field.NewField("", "Incident Name").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.IncidentName
				}).
				EditHelp(`This is the name of the incident, printed on top of the ICS-309 form.`).
				EditWidth(ics309FieldWidth("IncidentName")).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.IncidentName = s
				}).
				MakeField(),
			field.NewField("", "Activation Number").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.ActivationNum
				}).
				EditHelp(`This is the activation number for the incident, printed on top of the ICS-309 form.`).
				EditWidth(ics309FieldWidth("ActivationNum")).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.ActivationNum = s
				}).
				MakeField(),
			field.NewField("", "Operation Start").
				ValueFunc(func(m message.Message) string { return cm.startstr }).
				EditHelp(`This is the date and time when the operational period starts, in MM/DD/YYYY HH:MM format (24-hour clock), printed on top of the ICS-309 form.`).
				EditWidth(16).
				FromHumanFunc(func(m message.Message, s string) string {
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
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).startstr = s
					if t, err := time.ParseInLocation("01/02/2006 15:04", s, time.Local); err == nil {
						m.(*ConfigMessage).Config.OpStart = t
					}
				}).
				ValidateFunc(func(m message.Message, f field.Field, vf message.ValidateFlags) error {
					startstr := m.(*ConfigMessage).startstr
					if startstr != "" && !dateTimeRE.MatchString(startstr) {
						return NoBypassValidationError{errors.New("The date/time string must be in MM/DD/YYYY HH:MM format (24-hour clock).")}
					}
					return nil
				}).
				MakeField(),
			field.NewField("", "Operation End").
				ValueFunc(func(m message.Message) string { return cm.endstr }).
				EditHelp(`This is the date and time when the operational period ends, in MM/DD/YYYY HH:MM format (24-hour clock), printed on top of the ICS-309 form.`).
				EditWidth(16).
				FromHumanFunc(func(m message.Message, s string) string {
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
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).endstr = s
					if t, err := time.ParseInLocation("01/02/2006 15:04", s, time.Local); err == nil {
						m.(*ConfigMessage).Config.OpEnd = t
					}
				}).
				ValidateFunc(func(m message.Message, f field.Field, vf message.ValidateFlags) error {
					endstr := m.(*ConfigMessage).endstr
					if endstr != "" && !dateTimeRE.MatchString(endstr) {
						return NoBypassValidationError{errors.New("The date/time string must be in MM/DD/YYYY HH:MM format (24-hour clock).")}
					}
					return nil
				}).
				MakeField(),
			field.NewFCCCallSign("", "Operator Call Sign").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.OpCall
				}).
				EditHelp(`This is the FCC call sign of the station operator.  It is required.`).
				Required().
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.OpCall = s
					m.(*ConfigMessage).setDefaultMessageNumbers()
				}).
				MakeField(),
			field.NewField("", "Operator Name").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.OpName
				}).
				EditHelp(`This is the name of the station operator.  It is required.`).
				EditWidth(ics309FieldWidth("OpNameCall") - 7).
				Required().
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.OpName = s
				}).
				MakeField(),
			field.NewTacticalCallSign("", "Tactical Call Sign").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.TacCall
				}).
				EditHelp(`This is the call sign of the tactical station, if any.`).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.TacCall = s
					m.(*ConfigMessage).setDefaultMessageNumbers()
				}).
				MakeField(),
			field.NewField("", "Tactical Station Name").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.TacName
				}).
				EditHelp(`This is the name of the tactical station.  It is required if a tactical call sign is provided.`).
				EditWidth(ics309FieldWidth("TacNameCall")-7).
				Required().
				DisallowedUnless("a tactical call sign is provided", func(m message.Message) bool {
					return m.(*ConfigMessage).Config.TacCall != ""
				}).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.TacName = s
				}).
				MakeField(),
			field.NewMessageID("", "Tx Message ID").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.TxMessageID
				}).
				EditHelp(`This is the message number that will be assigned to the next outgoing message, in the form XXX-###P.  It is required.`).
				Required().
				SetValueFunc(func(m message.Message, s string) {
					if s != m.(*ConfigMessage).Config.TxMessageID {
						m.(*ConfigMessage).deftx = false
					}
					m.(*ConfigMessage).Config.TxMessageID = s
				}).
				MakeField(),
			field.NewField("", "Rx Message ID").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.RxMessageID
				}).
				EditHelp(`This is the message number that will be assigned to the next incoming message, in the form XXX-###P.  It is required.`).
				Required().
				SetValueFunc(func(m message.Message, s string) {
					if s != m.(*ConfigMessage).Config.RxMessageID {
						m.(*ConfigMessage).deftx = false
					}
					m.(*ConfigMessage).Config.RxMessageID = s
				}).
				MakeField(),
			field.NewField("", "Connection Type").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.ConnectType
				}).
				EditHelp(`This indicates how the software connects to the BBS.  "Serial+TNC" means it talks to a TNC through a serial port and through that, to the BBS over the radio.  "Telnet" means it connects to JNOS over the Internet.  "Manual" means it does not connect to the BBS at all.`).
				AllowedValues(
					field.ChoicePair{PIFO: incident.ConnectNone, Human: "Manual"},
					field.ChoicePair{PIFO: incident.ConnectSerialTNC, Human: "Serial+TNC"},
					field.ChoicePair{PIFO: incident.ConnectTelnet, Human: "Telnet"},
				).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.ConnectType = s
				}).
				MakeField(),
			field.NewFCCCallSign("", "BBS Call Sign").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.ConnectBBS
				}).
				EditHelp(`This is the call sign of the BBS in use.  It is required.`).
				EditableWhen(func(m message.Message, b bool) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectSerialTNC
				}).
				Required().
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.ConnectBBS = s
				}).
				MakeField(),
			field.NewField("", "Automatic Receipts").
				ValueFunc(func(m message.Message) string {
					if !m.(*ConfigMessage).Config.NoSendReceipts {
						return "checked"
					}
					return ""
				}).
				AllowedValues("checked").
				EditHelp(`This indicates whether delivery receipts will be generated automatically for manually received messages.`).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType == incident.ConnectNone
				}).
				EditableWhen(func(m message.Message, _ bool) bool {
					return m.(*ConfigMessage).Config.ConnectType == incident.ConnectNone
				}).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.NoSendReceipts = s == ""
				}).
				MakeField(),
			field.NewField("", "TNC Type").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.TNCType
				}).
				AllowedValues(field.ChoicePair{PIFO: incident.TNCKPC3Plus, Human: "Kantronics KPC-3 Plus"}).
				Required().
				DisallowedUnless(`"Connection Type" is "Serial+TNC"`, func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType == incident.ConnectSerialTNC
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType == incident.ConnectSerialTNC
				}).
				EditHelp(`This indicates the type of TNC connected to the host computer.  Currently the only supported value is "Kantronics KPC-3 Plus".  It is required when "Connection Type" is "Serial+TNC".`).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.TNCType = s
				}).
				MakeField(),
			field.NewField("", "Serial Port").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.SerialPort
				}).
				FromHumanFunc(func(_ message.Message, s string) string {
					if runtime.GOOS == "windows" {
						s = strings.ToUpper(s)
					}
					return s
				}).
				SuggestedValuesFunc(func(m message.Message) (pairs []field.ChoicePair) {
					for _, port := range GuessSerialPorts() {
						pairs = append(pairs, field.ChoicePair{PIFO: port, Human: port})
					}
					return pairs
				}).
				Required().
				DisallowedUnless(`"Connection Type" is "Serial+TNC"`, func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType == incident.ConnectSerialTNC
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType == incident.ConnectSerialTNC
				}).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.SerialPort = s
				}).
				ValidateFunc(func(m message.Message, f field.Field, vf message.ValidateFlags) error {
					value := f.Value(m)
					if runtime.GOOS == "windows" {
						if !comPortRE.MatchString(value) {
							return errors.NewF("%q is not a valid COM port.", value)
						}
					} else if stat, err := os.Stat(value); err != nil || stat.Mode().Type() != fs.ModeDevice|fs.ModeCharDevice {
						return errors.NewF("%q is not a character device file.", value)
					}
					return nil
				}).
				EditHelp(`This is the serial port to which the TNC is connected.  On Windows, this will be "COM#", where # is a number.  On Mac and Linux, this will be a pathname to a character device file in /dev.  It is required when "Connection Type" is "Serial+TNC".`).
				MakeField(),
			field.NewField("", "BBS Address").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.ConnectAddress
				}).
				FromHumanFunc(func(m message.Message, s string) string {
					if m.(*ConfigMessage).Config.ConnectType == incident.ConnectSerialTNC {
						s = strings.ToUpper(s)
					}
					return s
				}).
				Required().
				DisallowedUnless(`"Connection Type" is not "Manual"`, func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectNone
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectNone
				}).
				ValidateFunc(func(m message.Message, f field.Field, vf message.ValidateFlags) error {
					value := f.Value(m)
					switch m.(*ConfigMessage).Config.ConnectType {
					case incident.ConnectSerialTNC:
						if !ax25RE.MatchString(value) {
							return errors.NewF("%q is not a valid AX.25 address (i.e., FCC call sign followed by a dash and a number between 0 and 15).", value)
						}
					case incident.ConnectTelnet:
						if _, _, err := net.SplitHostPort(value); err != nil {
							return errors.NewF("%q is not a valid address:port or hostname:port string.", value)
						}
					}
					return nil
				}).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.ConnectAddress = s
					if m.(*ConfigMessage).Config.ConnectType == incident.ConnectSerialTNC {
						m.(*ConfigMessage).Config.ConnectBBS, _, _ = strings.Cut(s, "-")
					}
				}).
				EditHelp(`This is the address of the BBS.  For "Connection Type" "Serial+TNC", it must be an AX.25 address (i.e., a call sign followed by a dash an a number from 0 to 15).  For "Connection Type" "Telnet", it must be a hostname:port or IPaddr:port string.  For "Connection Type" "Manual" it must not be set.`).
				MakeField(),
			field.NewField("", "Bulletin Checks").
				ValueFunc(func(m message.Message) string {
					bc := m.(*ConfigMessage).Config.BulletinChecks
					areas := slices.Collect(maps.Keys(bc))
					slices.Sort(areas)
					for i, a := range areas {
						areas[i] = fmt.Sprintf("%s:%s", a, fmtCheckFrequency(bc[areas[i]]))
					}
					return strings.Join(areas, " ")
				}).
				VisibleWhen(func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectNone
				}).
				EditHelp(`This a whitespace-separated list of bulletin areas that should be periodically checked for new bulletins.  Each one can be a bulletin or notice area name (e.g., "ALLXSC" or "XSCEVENT"), or a topic@distribution string (e.g., "XND@XSC").  Each one can be followed by a colon and the frequency with which the area should be checked for new bulletins.  The frequency is specified as "3h", "1h30m", or similar.  If the frequency is omitted, it defaults to one hour.`).
				FromHumanFunc(func(_ message.Message, s string) string {
					terms := strings.Fields(s)
					for i, t := range terms {
						if area, freq, found := strings.Cut(t, ":"); found {
							area = strings.ToUpper(area)
							if cf := parseCheckFrequency(freq); cf.Duration != 0 {
								freq = fmtCheckFrequency(cf)
							}
							terms[i] = area + ":" + freq
						} else {
							terms[i] += ":1h"
						}
					}
					return strings.Join(terms, " ")
				}).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).bchecks = s
					if bc, err := parseBulletinChecks(s); err == nil {
						m.(*ConfigMessage).Config.BulletinChecks = bc
					}
				}).
				ValidateFunc(func(m message.Message, f field.Field, vf message.ValidateFlags) (err error) {
					_, err = parseBulletinChecks(f.Value(m))
					return err
				}).
				DisallowedUnless(`"Connection Type" is not "Manual"`, func(m message.Message) bool {
					return m.(*ConfigMessage).Config.ConnectType != incident.ConnectNone
				}).
				MakeField(),
			field.NewField("", "Default To Address").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultTo
				}).
				EditHelp(`This is the address to which new outgoing messages will be addressed by default.  It is optional.`).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.DefaultTo = s
				}).
				ValidateFunc(func(m message.Message, f field.Field, vf message.ValidateFlags) error {
					value := strings.TrimSpace(f.Value(m))
					if _, err := address.ParseList(value); value != "" && err != nil {
						return errors.NewF("%q is not a valid packet or email address.", value)
					}
					return nil
				}).
				MakeField(),
			field.NewField("", "Default To ICS Position").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultToPos
				}).
				EditHelp(`This is the value that will be filled into the "To ICS Position" field of new outgoing forms messages that have that field.  If empty, the defaults for each message type are used.`).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.DefaultToPos = s
				}).
				MakeField(),
			field.NewField("", "Default To Location").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultToLoc
				}).
				EditHelp(`This is the value that will be filled into the "To Location" field of new outgoing forms messages that have that field.  If empty, the defaults for each message type are used.`).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.DefaultToLoc = s
				}).
				MakeField(),
			field.NewField("", "Default From ICS Position").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultFromPos
				}).
				EditHelp(`This is the value that will be filled into the "From ICS Position" field of new outgoing forms messages that have that field.`).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.DefaultFromPos = s
				}).
				MakeField(),
			field.NewField("", "Default From Location").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultFromLoc
				}).
				EditHelp(`This is the value that will be filled into the "From Location" field of new outgoing forms messages that have that field.`).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.DefaultFromLoc = s
				}).
				MakeField(),
			field.NewField("", "Default Body").
				ValueFunc(func(m message.Message) string {
					return m.(*ConfigMessage).Config.DefaultBody
				}).
				EditHelp(`This is text that will be filled into the first or most prominent multi-line text field of every new outgoing message.  It is most commonly used for "**** This is drill traffic ****" and similar notes.`).
				SetValueFunc(func(m message.Message, s string) {
					m.(*ConfigMessage).Config.DefaultBody = s
				}).
				Multiline().
				MakeField(),
		}
	})
	return slices.Values(configFields)
}

var endsNumericRE = regexp.MustCompile(`\d\d\d$`)

func (cm *ConfigMessage) setDefaultMessageNumbers() {
	prefix := cm.msgIDPrefix()
	if prefix == "" {
		return
	}
	if cm.defrx || cm.Config.RxMessageID == "" {
		cm.Config.RxMessageID = prefix + "-100P"
		cm.defrx = true
	}
	if cm.deftx || cm.Config.TxMessageID == "" {
		cm.Config.TxMessageID = prefix + "-100P"
		cm.deftx = true
	}
}
func (cm *ConfigMessage) msgIDPrefix() string {
	if len(cm.Config.TacCall) >= 3 {
		if strings.HasSuffix(cm.Config.TacCall, "DOC") || strings.HasSuffix(cm.Config.TacCall, "EOC") {
			return cm.Config.TacCall[:3]
		} else if endsNumericRE.MatchString(cm.Config.TacCall) {
			return cm.Config.TacCall[:1] + cm.Config.TacCall[len(cm.Config.TacCall)-2:]
		} else {
			return cm.Config.TacCall[len(cm.Config.TacCall)-3:]
		}
	} else if len(cm.Config.OpCall) >= 3 {
		return cm.Config.OpCall[len(cm.Config.OpCall):]
	} else {
		return ""
	}
}

func (cm *ConfigMessage) Dirty() bool                  { panic("not implemented") }
func (cm *ConfigMessage) MarkDirty(_ string)           { panic("not implemented") }
func (cm *ConfigMessage) MarkClean()                   { panic("not implemented") }
func (cm *ConfigMessage) OnDirty(_ func(string))       { panic("not implemented") }
func (cm *ConfigMessage) Recognize(_ message.Message)  { panic("not implemented") }
func (cm *ConfigMessage) Name() string                 { panic("not implemented") }
func (cm *ConfigMessage) Type() message.MType          { panic("not implemented") }
func (cm *ConfigMessage) SetType(_ message.MType)      { panic("not implemented") }
func (cm *ConfigMessage) RFC5322() string              { panic("not implemented") }
func (cm *ConfigMessage) To() string                   { panic("not implemented") }
func (cm *ConfigMessage) Subject() subject.Subject     { panic("not implemented") }
func (cm *ConfigMessage) SetSubject(_ subject.Subject) { panic("not implemented") }
func (cm *ConfigMessage) Bulletin() bool               { panic("not implemented") }
func (cm *ConfigMessage) Payload() payload.Payload     { panic("not implemented") }
func (cm *ConfigMessage) Body() body.Body              { panic("not implemented") }
func (cm *ConfigMessage) RenderPDF(m message.Message, filename string, copyname string) error {
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

var checkFreqRE = regexp.MustCompile(`^(?:(\d+)h)?(?:(\d+)m)?$`)

func parseCheckFrequency(s string) (cf incident.CheckFrequency) {
	var err error
	if match := checkFreqRE.FindStringSubmatch(s); match != nil {
		var h, m int
		if h, err = strconv.Atoi(match[1]); err != nil && match[1] != "" {
			return cf
		}
		if m, err = strconv.Atoi(match[2]); err != nil && match[2] != "" {
			return cf
		}
		cf.Duration = time.Duration(h)*time.Hour + time.Duration(m)*time.Minute
	}
	return cf
}

var areaRE = regexp.MustCompile(`^(?:[A-Z]+@)?[A-Z]+$`)

func parseBulletinChecks(s string) (bc map[string]incident.CheckFrequency, err error) {
	bc = make(map[string]incident.CheckFrequency)
	for _, term := range strings.Fields(s) {
		area, freq, _ := strings.Cut(term, ":")
		if _, ok := bc[area]; ok {
			return nil, errors.NewF("Bulletin area %q is specified more than once.", area)
		}
		if !areaRE.MatchString(area) {
			return nil, errors.NewF("%q is not a valid bulletin or notice area name.", area)
		}
		cf := parseCheckFrequency(freq)
		if cf.Duration == 0 {
			return nil, errors.NewF("%q is not a valid frequency.  Frequencies are 3h, 1h30m, etc.", freq)
		}
		bc[area] = cf
	}
	return bc, nil
}

func ics309FieldWidth(tag string) int {
	for f := range incident.ICS309FormDef().AllFields() {
		if f.Tag == tag {
			return f.CharWidth(0)
		}
	}
	panic("no such field " + tag)
}

// GuessSerialPorts makes a swag at the possible device files for serial ports.
func GuessSerialPorts() (ports []string) {
	if names, err := serial.GetPortsList(); err == nil && len(names) > 0 {
		return names
	}
	if runtime.GOOS == "windows" {
		// Just give a list of COM numbers.
		return []string{"COM1", "COM2", "COM3", "COM4", "COM5", "COM6"}
	}
	// On any other OS, look for /dev/tty* files with USB in the name.
	ports, _ = filepath.Glob("/dev/tty*usb*")
	if ports2, _ := filepath.Glob("/dev/tty*USB*"); len(ports2) != 0 {
		ports = append(ports, ports2...)
		slices.SortFunc(ports, func(a, b string) int {
			return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
		})
	}
	return ports
}
