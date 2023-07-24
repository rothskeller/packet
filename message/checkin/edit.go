package checkin

import (
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// EditFields returns the set of editable fields of the message.
func (m *CheckIn) EditFields() []*message.EditField {
	if m.edit == nil {
		m.edit = new(checkInEdit)
		m.edit.OriginMsgID = common.OriginMsgIDEditField
		m.edit.Handling = common.HandlingEditField
		m.edit.TacticalCallSign = message.EditField{
			Label: "Tactical Call Sign",
			Width: 6,
			Help:  "This is the tactical call sign assigned to the station being operated, if any.  It is expected to be six letters or digits, starting with a letter.",
		}
		m.edit.TacticalStationName = message.EditField{
			Label: "Tactical Station Name",
			Width: 80,
			Help:  "This is the name of the station being operated.  It must be set when the Tactical Call Sign is set.",
		}
		m.edit.OperatorCallSign = message.EditField{
			Label: "Operator Call Sign",
			Width: 6,
			Help:  "This is the FCC call sign assigned to the operator of the station.  It is required.",
		}
		m.edit.OperatorName = message.EditField{
			Label: "Operator Name",
			Width: 80,
			Help:  "This is the name of the operator of the station.  It is required.",
		}
		m.edit.fields = []*message.EditField{
			&m.edit.OriginMsgID,
			&m.edit.Handling,
			&m.edit.TacticalCallSign,
			&m.edit.TacticalStationName,
			&m.edit.OperatorCallSign,
			&m.edit.OperatorName,
		}
		m.toEdit()
		m.validate()
	}
	return m.edit.fields
}

// ApplyEdits applies the revised Values in the EditFields to the message.
func (m *CheckIn) ApplyEdits() {
	m.fromEdit()
	m.validate()
	m.toEdit()
	m.OriginMsgID = strings.ToUpper(strings.TrimSpace(m.edit.OriginMsgID.Value))
	m.Handling = common.ExpandRestricted(&m.edit.Handling)
	m.TacticalCallSign = strings.ToUpper(strings.TrimSpace(m.edit.TacticalCallSign.Value))
	m.TacticalStationName = strings.TrimSpace(m.edit.TacticalStationName.Value)
	m.OperatorCallSign = strings.ToUpper(strings.TrimSpace(m.edit.OperatorCallSign.Value))
	m.OperatorName = strings.TrimSpace(m.edit.OperatorName.Value)
}

func (m *CheckIn) toEdit() {
	m.edit.OriginMsgID.Value = m.OriginMsgID
	m.edit.Handling.Value = m.Handling
	m.edit.TacticalCallSign.Value = m.TacticalCallSign
	m.edit.TacticalStationName.Value = m.TacticalStationName
	m.edit.OperatorCallSign.Value = m.OperatorCallSign
	m.edit.OperatorName.Value = m.OperatorName
}

func (m *CheckIn) fromEdit() {
	m.OriginMsgID = common.CleanMessageNumber(m.edit.OriginMsgID.Value)
	m.Handling = common.ExpandRestricted(&m.edit.Handling)
	m.TacticalCallSign = common.ExpandRestricted(&m.edit.TacticalCallSign)
	m.TacticalStationName = common.ExpandRestricted(&m.edit.TacticalStationName)
	m.OperatorCallSign = common.ExpandRestricted(&m.edit.OperatorCallSign)
	m.OperatorName = common.ExpandRestricted(&m.edit.OperatorName)
}

func (m *CheckIn) validate() {
	if common.ValidateRequired(&m.edit.OriginMsgID) {
		common.ValidateMessageNumber(&m.edit.OriginMsgID)
	}
	if common.ValidateRequired(&m.edit.Handling) {
		common.ValidateRestricted(&m.edit.Handling)
	}
	if m.edit.TacticalCallSign.Value != "" {
		common.ValidateTacticalCallSign(&m.edit.TacticalCallSign)
	} else {
		m.edit.TacticalCallSign.Problem = ""
	}
	if m.edit.TacticalCallSign.Value == "" && m.edit.TacticalStationName.Value != "" {
		m.edit.TacticalStationName.Problem = `The "Tactical Station Name" field cannot have a value unless the "Tactical Call Sign" field has a value.`
	} else if m.edit.TacticalCallSign.Value != "" && m.edit.TacticalStationName.Value == "" {
		m.edit.TacticalStationName.Problem = `The "Tactical Station Name" field is required when "Tactical Call Sign" is set.`
	} else {
		m.edit.TacticalStationName.Problem = ""
	}
	if common.ValidateRequired(&m.edit.OperatorCallSign) {
		common.ValidateFCCCallSign(&m.edit.OperatorCallSign)
	}
	common.ValidateRequired(&m.edit.OperatorName)
}
