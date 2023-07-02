package ics213

import "github.com/rothskeller/packet/message/common"

// Validate checks the contents of the message for compliance with rules
// enforced by standard Santa Clara County packet software (Outpost and
// PackItForms).  It returns a list of strings describing problems that those
// programs would flag or block.
func (f *ICS213) Validate() (problems []string) {
	if f.FormVersion >= "2.2" {
		if f.OriginMsgID == "" {
			problems = append(problems, `The "Origin Message #" field is required.`)
		}
	} else {
		if f.OriginMsgID == "" {
			problems = append(problems, `The "My Msg #" field is required.`)
		}
	}
	if f.Date == "" {
		problems = append(problems, `The "Date" field is required.`)
	} else if !common.PIFODateRE.MatchString(f.Date) {
		problems = append(problems, `The "Date" field does not contain a valid date.`)
	}
	if f.Time == "" {
		problems = append(problems, `The "Time" field is required.`)
	} else if !common.PIFOTimeRE.MatchString(f.Time) {
		problems = append(problems, `The "Time" field does not contain a valid time.`)
	}
	if f.FormVersion >= "2.2" {
		switch f.Handling {
		case "":
			problems = append(problems, `The "Handling" field is required.`)
		case "ROUTINE", "PRIORITY", "IMMEDIATE":
			break
		default:
			problems = append(problems, `The "Handling" field does not contain a valid message handling order.`)
		}
	} else {
		switch f.Severity {
		case "":
			problems = append(problems, `The "Situation Severity" field is required.`)
		case "EMERGENCY", "URGENT", "OTHER":
			break
		default:
			problems = append(problems, `The "Situation Severity" field does not contain a valid severity value.`)
		}
		switch f.Handling {
		case "":
			problems = append(problems, `The "Message Handling Order" field is required.`)
		case "ROUTINE", "PRIORITY", "IMMEDIATE":
			break
		default:
			problems = append(problems, `The "Message Handling Order" field does not contain a valid message handling order.`)
		}
	}
	switch f.TakeAction {
	case "", "Yes", "No":
		break
	default:
		problems = append(problems, `The "Take Action" field does not have a valid value.`)
	}
	switch f.Reply {
	case "", "Yes", "No":
		break
	default:
		problems = append(problems, `The "Reply" field does not have a valid value.`)
	}
	if f.ReplyBy != "" && f.Reply != "Yes" {
		problems = append(problems, `The "Reply By" field is not allowed unless "Reply" is "Yes".`)
	}
	if f.FYI != "" && f.FYI != "checked" {
		problems = append(problems, `The "For your info" field does not contain a valid checkbox value.`)
		// No need to make this conditional for old versions, since the
		// empty value in current versions passes the check.
	}
	if f.ToICSPosition == "" {
		problems = append(problems, `The "To ICS Position" field is required.`)
	}
	if f.FromICSPosition == "" {
		problems = append(problems, `The "From ICS Position" field is required.`)
	}
	if f.ToLocation == "" {
		problems = append(problems, `The "To Location" field is required.`)
	}
	if f.FromLocation == "" {
		problems = append(problems, `The "From Location" field is required.`)
	}
	if f.Subject == "" {
		problems = append(problems, `The "Subject" field is required.`)
	}
	if f.Message == "" {
		problems = append(problems, `The "Message" field is required.`)
	}
	switch f.ReceivedSent {
	case "", "receiver", "sender":
		break
	default:
		problems = append(problems, `The "Received" vs. "Sent" switch is not set to a valid value.`)
	}
	if f.OpCall == "" {
		problems = append(problems, `The "Operator Use Only: Call Sign" field is required.`)
	}
	if f.OpName == "" {
		problems = append(problems, `The "Operator Use Only: Name" field is required.`)
	}
	switch f.TxMethod {
	case "Telephone", "Dispatch Center", "EOC Radio", "FAX", "Courier", "Amateur Radio", "Other":
		break
	default:
		problems = append(problems, `The "How Received or Sent" field does not contain a valid value.`)
	}
	if f.OpDate == "" {
		problems = append(problems, `The "Operator Use Only: Date" field is required.`)
	} else if !common.PIFODateRE.MatchString(f.OpDate) {
		problems = append(problems, `The "Operator Use Only: Date" field does not contain a valid date.`)
	}
	if f.OpTime == "" {
		problems = append(problems, `The "Operator Use Only: Time" field is required.`)
	} else if !common.PIFOTimeRE.MatchString(f.OpTime) {
		problems = append(problems, `The "Operator Use Only: Time" field does not contain a valid time.`)
	}
	return problems
}
