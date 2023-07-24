package ics213

import (
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

type ics213Edit struct {
	OriginMsgID     message.EditField
	Date            message.EditField
	Time            message.EditField
	Handling        message.EditField
	TakeAction      message.EditField
	Reply           message.EditField
	ReplyBy         message.EditField
	ToICSPosition   message.EditField
	ToLocation      message.EditField
	ToName          message.EditField
	ToTelephone     message.EditField
	FromICSPosition message.EditField
	FromLocation    message.EditField
	FromName        message.EditField
	FromTelephone   message.EditField
	Subject         message.EditField
	Reference       message.EditField
	Message         message.EditField
	OpRelayRcvd     message.EditField
	OpRelaySent     message.EditField
	fields          []*message.EditField
}

// EditFields returns the set of editable fields of the message.
func (f *ICS213) EditFields() []*message.EditField {
	if f.edit == nil {
		f.edit = &ics213Edit{
			OriginMsgID: common.OriginMsgIDEditField,
			Date:        common.MessageDateEditField,
			Time:        common.MessageTimeEditField,
			Handling:    common.HandlingEditField,
			TakeAction: message.EditField{
				Label:   "Take Action",
				Width:   3,
				Help:    "This field indicates whether the sender expects the recipient to take action based on this message.  Allowed values are Yes and No.",
				Choices: []string{"Yes", "No"},
			},
			Reply: message.EditField{
				Label:   "Reply",
				Width:   3,
				Help:    "This field indicates whether the sender expects the recipient to reply to this message.  Allowed values are Yes and No.",
				Choices: []string{"Yes", "No"},
			},
			ReplyBy: message.EditField{
				Label: "Reply By",
				Width: 5,
				Help:  `This is the time by which the sender expects a reply from the recipient.  It can be set only if "Reply" is "Yes".`,
			},
			ToICSPosition: common.ToICSPositionEditField,
			ToLocation:    common.ToLocationEditField,
			ToName:        common.ToNameEditField,
			ToTelephone: message.EditField{
				Label: "To Telephone #",
				Width: 40,
				Help:  "This is the telephone number for the receipient.  It is optional and rarely provided.",
			},
			FromICSPosition: common.FromICSPositionEditField,
			FromLocation:    common.FromLocationEditField,
			FromName:        common.FromNameEditField,
			FromTelephone: message.EditField{
				Label: "From Telephone #",
				Width: 40,
				Help:  "This is the telephone number for the message author.  It is optional and rarely provided.",
			},
			Subject: message.EditField{
				Label: "Subject",
				Width: 80,
				Help:  "This is the subject of the message.  It is required.",
			},
			Reference: message.EditField{
				Label: "Reference",
				Width: 55,
				Help:  "This is the origin message number of the previous message to which this message is responding.",
			},
			Message: message.EditField{
				Label:     "Message",
				Width:     80,
				Multiline: true,
				Help:      "This is the text of the message.  It is required.",
			},
			OpRelayRcvd: common.OpRelayRcvdEditField,
			OpRelaySent: common.OpRelaySentEditField,
		}
		// Apply differences from common fields.
		f.edit.ToICSPosition.Width = 40
		f.edit.ToLocation.Width = 40
		f.edit.ToName.Width = 40
		f.edit.FromICSPosition.Width = 40
		f.edit.FromLocation.Width = 40
		f.edit.FromName.Width = 40
		// Set the field list slice.
		f.edit.fields = []*message.EditField{
			&f.edit.OriginMsgID,
			&f.edit.Date,
			&f.edit.Time,
			&f.edit.Handling,
			&f.edit.TakeAction,
			&f.edit.Reply,
			&f.edit.ReplyBy,
			&f.edit.ToICSPosition,
			&f.edit.ToLocation,
			&f.edit.ToName,
			&f.edit.ToTelephone,
			&f.edit.FromICSPosition,
			&f.edit.FromLocation,
			&f.edit.FromName,
			&f.edit.FromTelephone,
			&f.edit.Subject,
			&f.edit.Reference,
			&f.edit.Message,
			&f.edit.OpRelayRcvd,
			&f.edit.OpRelaySent,
		}
		f.toEdit()
		f.validate()
	}
	return f.edit.fields
}

// ApplyEdits applies the revised Values in the EditFields to the
// message.
func (f *ICS213) ApplyEdits() {
	f.fromEdit()
	f.validate()
	f.toEdit()
}

func (f *ICS213) fromEdit() {
	f.OriginMsgID = common.CleanMessageNumber(f.edit.OriginMsgID.Value)
	f.Date = common.CleanDate(f.edit.Date.Value)
	f.Time = common.CleanTime(f.edit.Time.Value)
	f.Handling = common.ExpandRestricted(&f.edit.Handling)
	f.TakeAction = common.ExpandRestricted(&f.edit.TakeAction)
	f.Reply = common.ExpandRestricted(&f.edit.Reply)
	f.ReplyBy = strings.TrimSpace(f.edit.ReplyBy.Value)
	f.ToICSPosition = strings.TrimSpace(f.edit.ToICSPosition.Value)
	f.ToLocation = strings.TrimSpace(f.edit.ToLocation.Value)
	f.ToName = strings.TrimSpace(f.edit.ToName.Value)
	f.ToTelephone = strings.TrimSpace(f.edit.ToTelephone.Value)
	f.FromICSPosition = strings.TrimSpace(f.edit.FromICSPosition.Value)
	f.FromLocation = strings.TrimSpace(f.edit.FromLocation.Value)
	f.FromName = strings.TrimSpace(f.edit.FromName.Value)
	f.FromTelephone = strings.TrimSpace(f.edit.FromTelephone.Value)
	f.Subject = strings.TrimSpace(f.edit.Subject.Value)
	f.Reference = common.CleanMessageNumber(f.edit.Reference.Value)
	f.Message = strings.TrimSpace(f.edit.Message.Value)
	f.OpRelayRcvd = strings.TrimSpace(f.edit.OpRelayRcvd.Value)
	f.OpRelaySent = strings.TrimSpace(f.edit.OpRelaySent.Value)
}

func (f *ICS213) toEdit() {
	f.edit.OriginMsgID.Value = f.OriginMsgID
	f.edit.Date.Value = f.Date
	f.edit.Time.Value = f.Time
	f.edit.Handling.Value = f.Handling
	f.edit.TakeAction.Value = f.TakeAction
	f.edit.Reply.Value = f.Reply
	f.edit.ReplyBy.Value = f.ReplyBy
	f.edit.ToICSPosition.Value = f.ToICSPosition
	f.edit.ToLocation.Value = f.ToLocation
	f.edit.ToName.Value = f.ToName
	f.edit.ToTelephone.Value = f.ToTelephone
	f.edit.FromICSPosition.Value = f.FromICSPosition
	f.edit.FromLocation.Value = f.FromLocation
	f.edit.FromName.Value = f.FromName
	f.edit.FromTelephone.Value = f.FromTelephone
	f.edit.Subject.Value = f.Subject
	f.edit.Reference.Value = f.Reference
	f.edit.Message.Value = f.Message
	f.edit.OpRelayRcvd.Value = f.OpRelayRcvd
	f.edit.OpRelaySent.Value = f.OpRelaySent
}

func (f *ICS213) validate() {
	if common.ValidateRequired(&f.edit.OriginMsgID) {
		common.ValidateMessageNumber(&f.edit.OriginMsgID)
	}
	if common.ValidateRequired(&f.edit.Date) {
		common.ValidateDate(&f.edit.Date)
	}
	if common.ValidateRequired(&f.edit.Time) {
		common.ValidateTime(&f.edit.Time)
	}
	if common.ValidateRequired(&f.edit.Handling) {
		common.ValidateRestricted(&f.edit.Handling)
	}
	if f.edit.TakeAction.Value != "" {
		common.ValidateRestricted(&f.edit.TakeAction)
	}
	if f.edit.Reply.Value != "" {
		common.ValidateRestricted(&f.edit.Reply)
	}
	if f.edit.Reply.Value != "Yes" && f.edit.ReplyBy.Value != "" {
		f.edit.ReplyBy.Problem = `The "Reply By" field is not allowed unless "Reply" is "Yes".`
	} else {
		f.edit.ReplyBy.Problem = ""
	}
	common.ValidateRequired(&f.edit.ToICSPosition)
	common.ValidateRequired(&f.edit.ToLocation)
	common.ValidateRequired(&f.edit.FromICSPosition)
	common.ValidateRequired(&f.edit.FromLocation)
	common.ValidateRequired(&f.edit.Subject)
	common.ValidateRequired(&f.edit.Message)
}
