package plaintext

import (
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// EditFields returns the set of editable fields of the message.
func (m *PlainText) EditFields() []*message.EditField {
	if m.edit == nil {
		m.edit = new(plainTextEdit)
		m.edit.OriginMsgID = common.OriginMsgIDEditField
		m.edit.Handling = common.HandlingEditField
		m.edit.Subject = message.EditField{
			Label: "Subject",
			Width: 80,
			Help:  "This is the subject of the message.  It is required.",
		}
		m.edit.Body = message.EditField{
			Label:     "Body",
			Width:     80,
			Multiline: true,
			Help:      "This is the body of the message.  It is required.",
		}
		m.editFields = []*message.EditField{
			&m.edit.OriginMsgID,
			&m.edit.Handling,
			&m.edit.Subject,
			&m.edit.Body,
		}
		m.toEdit()
		m.validate()
	}
	return m.editFields
}

// ApplyEdits applies the revised Values in the EditFields to the message.
func (m *PlainText) ApplyEdits() {
	m.fromEdit()
	m.toEdit()
	m.validate()
	m.OriginMsgID = strings.ToUpper(strings.TrimSpace(m.edit.OriginMsgID.Value))
	m.Handling = common.ExpandRestricted(&m.edit.Handling)
	m.Subject = strings.TrimSpace(m.edit.Subject.Value)
	m.Body = strings.TrimSpace(m.edit.Body.Value)
}

func (m *PlainText) toEdit() {
	m.edit.OriginMsgID.Value = m.OriginMsgID
	m.edit.Handling.Value = m.Handling
	m.edit.Subject.Value = m.Subject
	m.edit.Body.Value = m.Body
}

func (m *PlainText) fromEdit() {
	m.OriginMsgID = common.CleanMessageNumber(m.edit.OriginMsgID.Value)
	m.Handling = common.ExpandRestricted(&m.edit.Handling)
	m.Subject = strings.TrimSpace(m.edit.Subject.Value)
	m.Body = strings.TrimSpace(m.edit.Body.Value)
}

func (m *PlainText) validate() {
	if common.ValidateRequired(&m.edit.OriginMsgID) {
		common.ValidateMessageNumber(&m.edit.OriginMsgID)
	}
	if common.ValidateRequired(&m.edit.Handling) {
		common.ValidateRestricted(&m.edit.Handling)
	}
	common.ValidateRequired(&m.edit.Subject)
	common.ValidateRequired(&m.edit.Body)
}
