// Package checkout implements the Santa Clara County standard check-out
// message.
package checkout

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/basemsg"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for a check-out message.
var Type = message.Type{
	Tag:     "Check-Out",
	Name:    "check-out message",
	Article: "a",
}

func init() {
	Type.Create = New
	Type.Decode = decode
}

// CheckOut holds the details of an XSC-standard check-out message.
type CheckOut struct {
	basemsg.BaseMessage
	OriginMsgID         string
	Handling            string
	TacticalCallSign    string
	TacticalStationName string
	OperatorCallSign    string
	OperatorName        string
}

// New creates a new, outgoing message of the type.
func New() (m *CheckOut) {
	m = create()
	m.Handling = "ROUTINE"
	return m
}

func create() (m *CheckOut) {
	m = &CheckOut{BaseMessage: basemsg.BaseMessage{MessageType: &Type}}
	m.BaseMessage.FOriginMsgID = &m.OriginMsgID
	m.BaseMessage.FHandling = &m.Handling
	m.BaseMessage.FOpCall = &m.OperatorCallSign
	m.BaseMessage.FOpName = &m.OperatorName
	m.Fields = []*basemsg.Field{
		basemsg.NewMessageNumberField(&basemsg.Field{
			Label:    "Origin Message Number",
			Value:    &m.OriginMsgID,
			Presence: basemsg.Required,
			EditHelp: `This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "Handling",
			Value:    &m.Handling,
			Choices:  basemsg.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence: basemsg.Required,
			EditHelp: `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
		}),
		basemsg.NewTacticalCallSignField(&basemsg.Field{
			Label:      "Tactical Call Sign",
			Value:      &m.TacticalCallSign,
			TableValue: basemsg.TableOmit,
			EditHelp:   `This is the tactical call sign assigned to the station being operated, if any.  It is expected to be five or six letters or digits, starting with a letter.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label: "Tactical Station Name",
			Value: &m.TacticalStationName,
			Presence: func() (basemsg.Presence, string) {
				if m.TacticalCallSign == "" {
					return basemsg.PresenceNotAllowed, `unless "Tactical Call Sign" has a value`
				} else {
					return basemsg.PresenceRequired, `when "Tactical Call Sign" has a value`
				}
			},
			TableValue: basemsg.TableOmit,
			EditWidth:  80,
			EditHelp:   `This is the name of the station being operated.  It must be set when the Tactical Call Sign is set.`,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
			Label: "Tactical Station",
			TableValue: func(f *basemsg.Field) string {
				return common.SmartJoin(m.TacticalCallSign, m.TacticalStationName, " ")
			},
		}),
		basemsg.NewFCCCallSignField(&basemsg.Field{
			Label:      "Operator Call Sign",
			Value:      &m.OperatorCallSign,
			Presence:   basemsg.Required,
			TableValue: basemsg.TableOmit,
			EditHelp:   `This is the FCC call sign assigned to the operator of the station.  It is required.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:      "Operator Name",
			Value:      &m.OperatorName,
			Presence:   basemsg.Required,
			TableValue: basemsg.TableOmit,
			EditWidth:  80,
			EditHelp:   `This is the name of the operator of the station.  It is required.`,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
			Label: "Operator",
			TableValue: func(f *basemsg.Field) string {
				return common.SmartJoin(m.OperatorCallSign, m.OperatorName, " ")
			},
		}),
	}
	return m
}

var checkInBodyRE = regexp.MustCompile(`(?im)^Check-Out\s+([A-Z][A-Z0-9]{2,5})\s*,(.*)(?:\n(A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})\s*,(.*))?`)

func decode(subject, body string) (f *CheckOut) {
	var msgid, _, handling, formtag, realsubj = common.DecodeSubject(subject)
	if formtag != "" || !strings.HasPrefix(strings.ToLower(realsubj), "check-out ") {
		return nil
	}
	if h := common.DecodeHandlingMap[handling]; h != "" {
		handling = h
	}
	if match := checkInBodyRE.FindStringSubmatch(body); match != nil {
		var co = create()
		co.OriginMsgID, co.Handling = msgid, handling
		if match[3] != "" {
			co.TacticalCallSign = match[1]
			co.TacticalStationName = strings.TrimSpace(match[2])
			co.OperatorCallSign = match[3]
			co.OperatorName = strings.TrimSpace(match[4])
		} else {
			co.OperatorCallSign = match[1]
			co.OperatorName = strings.TrimSpace(match[2])
		}
		return co
	}
	return nil
}

// EncodeSubject encodes the message subject.
func (m *CheckOut) EncodeSubject() string {
	if m.TacticalCallSign != "" {
		return common.EncodeSubject(m.OriginMsgID, m.Handling, "",
			fmt.Sprintf("Check-Out %s, %s", m.TacticalCallSign, m.TacticalStationName))
	}
	return common.EncodeSubject(m.OriginMsgID, m.Handling, "",
		fmt.Sprintf("Check-Out %s, %s", m.OperatorCallSign, m.OperatorName))
}

// EncodeBody encodes the message body.
func (m *CheckOut) EncodeBody() string {
	if m.TacticalCallSign != "" {
		return fmt.Sprintf("Check-Out %s, %s\n%s, %s\n",
			m.TacticalCallSign, m.TacticalStationName, m.OperatorCallSign, m.OperatorName)
	}
	return fmt.Sprintf("Check-Out %s, %s\n", m.OperatorCallSign, m.OperatorName)
}
