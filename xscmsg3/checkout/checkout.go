package checkout

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// Tag identifies check-out messages.
const Tag = "Check-Out"
const name = "check-out message"

var mtype = xscmsg.MessageType{Tag: Tag, Name: name, Article: "a"}

var checkOutRE = regexp.MustCompile(`(?i)^Check-Out\s+([A-Z][A-Z0-9]{2,5})\s*,(.*)(?:\n([AKNW][A-Z0-9]{2,5})\s*,(.*))?`)

func init() {
	xscmsg.Register(Tag, name, create, recognize)
}

type checkout struct {
	xscmsg.Message
	tacCall tacCallField
	tacName tacNameField
	opCall  opCallField
	opName  opNameField
}

func create() xscmsg.Message {
	var m = checkout{Message: xscmsg.NewBaseMessage(mtype, nil)}
	m.KeyedField(xscmsg.FHandling).SetValue("ROUTINE")
	return &m
}

func recognize(pm pktmsg.Message) xscmsg.Message {
	if pktmsg.KeyValue(pm, xscmsg.FFormTag) != "" {
		return nil
	}
	if !strings.HasPrefix(strings.ToLower(pktmsg.KeyValue(pm, pktmsg.FSubject)), "check-out ") {
		return nil
	}
	var m = checkout{Message: xscmsg.NewBaseMessage(mtype, pm)}
	if match := checkOutRE.FindStringSubmatch(pktmsg.KeyValue(pm, pktmsg.FBody)); match != nil {
		if match[3] != "" {
			m.tacCall.SetValue(match[1])
			m.tacName.SetValue(strings.TrimSpace(match[2]))
			m.opCall.SetValue(match[3])
			m.opName.SetValue(strings.TrimSpace(match[4]))
		} else {
			m.opCall.SetValue(match[1])
			m.opName.SetValue(strings.TrimSpace(match[2]))
		}
	}
	return m
}

func (m checkout) Iterate(fn func(pktmsg.Field)) {
	m.Message.Iterate(func(f pktmsg.Field) {
		if f, ok := f.(pktmsg.KeyedField); ok {
			if key := f.Key(); key == pktmsg.FSubject || key == pktmsg.FBody {
				return
			}
		}
		fn(f)
	})
	fn(&m.tacCall)
	fn(&m.tacName)
	fn(&m.opCall)
	fn(&m.opName)
}
func (m checkout) KeyedField(key pktmsg.FieldKey) pktmsg.KeyedField {
	switch key {
	case pktmsg.FSubject, pktmsg.FBody:
		return nil
	case xscmsg.FTacCall:
		return &m.tacCall
	case xscmsg.FTacName:
		return &m.tacName
	case xscmsg.FOpCall:
		return &m.opCall
	case xscmsg.FOpName:
		return &m.opName
	}
	return m.Message.KeyedField(key)
}
func (m checkout) Save() string {
	m.encode()
	return m.Message.Save()
}
func (m checkout) Transmit() (to []string, subject string, body string) {
	m.encode()
	return m.Message.Transmit()
}
func (m checkout) encode() {
	if m.tacCall.Value() != "" {
		m.KeyedField(pktmsg.FSubject).SetValue(fmt.Sprintf("Check-Out %s, %s", m.tacCall.Value(), m.tacName.Value()))
		m.KeyedField(pktmsg.FBody).SetValue(fmt.Sprintf("Check-Out %s, %s\n%s, %s\n", m.tacCall.Value(), m.tacName.Value(), m.opCall.Value(), m.opName.Value()))
	} else {
		m.KeyedField(pktmsg.FSubject).SetValue(fmt.Sprintf("Check-Out %s, %s", m.opCall.Value(), m.opName.Value()))
		m.KeyedField(pktmsg.FBody).SetValue(fmt.Sprintf("Check-Out %s, %s\n", m.opCall.Value(), m.opName.Value()))
	}
}

type tacCallField struct{ xscmsg.TacCallSignField }

func (f tacCallField) Key() pktmsg.FieldKey { return xscmsg.FTacCall }
func (f tacCallField) Label() string        { return "Tactical Call Sign" }
func (f tacCallField) Help(m pktmsg.Message) string {
	return "This is the tactical call sign of the local station that handled the message.  Tactical call signs comprise up to six alphanumeric characters, starting with a letter."
}
func (f tacCallField) Validate(_ pktmsg.Message, pifo bool) string {
	if f.TacCallSignField.PValidate(pifo) != "" {
		return "The tactical call sign is not valid."
	}
	return ""
}

type tacNameField struct{ xscmsg.StringField }

func (f tacNameField) Key() pktmsg.FieldKey { return xscmsg.FTacName }
func (f tacNameField) Label() string        { return "Tactical Station Name" }
func (f tacNameField) Help(m pktmsg.Message) string {
	return "This is the tactical name of the local station that handled the message.  It is required if a tactical call sign is supplied, and not allowed otherwise."
}
func (f tacNameField) Validate(m pktmsg.Message, pifo bool) string {
	tacCall := pktmsg.KeyValue(m, xscmsg.FTacCall)
	if !pifo && f.Value() != "" && tacCall == "" {
		return "The tactical station name may not be set unless the tactical call sign is also set."
	}
	if !pifo && f.Value() == "" && tacCall != "" {
		return "The tactical station name is required when the tactical call sign is set."
	}
	return ""
}

// Note this is different from xscmsg.NewOpCallField, which is not editable.
type opCallField struct{ xscmsg.FCCCallSignField }

func (f opCallField) Key() pktmsg.FieldKey { return xscmsg.FOpCall }
func (f opCallField) Label() string        { return "Operator Call Sign" }
func (f opCallField) Help(m pktmsg.Message) string {
	return "This is the FCC call sign of the operator of the local station that handled the message."
}
func (f opCallField) Validate(_ pktmsg.Message, pifo bool) string {
	if f.Value() == "" {
		return "The operator call sign is required."
	}
	if f.FCCCallSignField.PValidate(pifo) != "" {
		return "The operator call sign is not a valid FCC call sign."
	}
	return ""
}

// Note this is different from xscmsg.NewOpNameField, which is not editable.
type opNameField struct{ xscmsg.StringField }

func (f opNameField) Key() pktmsg.FieldKey { return xscmsg.FOpName }
func (f opNameField) Label() string        { return "Operator Name" }
func (f opNameField) Help(m pktmsg.Message) string {
	return "This is the name of the operator of the local station that handled the message.  It is required."
}
func (f opNameField) Validate(_ pktmsg.Message, _ bool) string {
	if f.Value() == "" {
		return "The operator name is required."
	}
	return ""
}
