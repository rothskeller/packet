// Package notrep defines the Notable Report Form message type.
package notrep

import (
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for a notable report form.
var Type = message.Type{
	Tag:     "NotRep",
	HTML:    "form-notable-report.html",
	Version: "0.6",
	Name:    "notable report form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.",
		"8c.", "7d.", "8d.", "20.", "21.", "22.", "23.", "24.", "25.",
		"30.", "31d.", "31t.", "32.", "33.", "40.", "41.", "42.", "50.",
		"51.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall",
		"OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type, decode, create)
}

var basePDFRenderers = baseform.BaseFormPDF{
	OriginMsgID: &message.PDFMultiRenderer{
		&message.PDFTextRenderer{X: 223, Y: 50, R: 348, B: 67, Style: message.PDFTextStyle{VAlign: "baseline"}},
		&message.PDFTextRenderer{Page: 2, X: 420, Y: 36, R: 574, B: 48},
	},
	DestinationMsgID: &message.PDFTextRenderer{X: 452, Y: 50, R: 574, B: 67, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 69, Y: 87, R: 128, B: 104, Style: message.PDFTextStyle{VAlign: "baseline", FontSize: 9}},
	MessageTime:      &message.PDFTextRenderer{X: 162, Y: 87, R: 203, B: 104, Style: message.PDFTextStyle{VAlign: "baseline", FontSize: 9}},
	Handling: &message.PDFRadioRenderer{Radius: 5, Points: map[string][]float64{
		"IMMEDIATE": {277, 96},
		"PRIORITY":  {388, 96},
		"ROUTINE":   {493, 96},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 132, Y: 106, R: 292, B: 123, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 132, Y: 125, R: 292, B: 142, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 132, Y: 144, R: 292, B: 161, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 132, Y: 163, R: 292, B: 181, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 382, Y: 106, R: 572, B: 123, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 382, Y: 125, R: 572, B: 142, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 382, Y: 144, R: 572, B: 161, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 382, Y: 163, R: 572, B: 181, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{X: 109, Y: 711, R: 320, B: 728, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{X: 356, Y: 711, R: 574, B: 728, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{X: 76, Y: 730, R: 249, B: 747, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{X: 301, Y: 730, R: 366, B: 747, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{X: 401, Y: 730, R: 479, B: 747, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{X: 542, Y: 730, R: 574, B: 747, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

// NotableReport holds a notable report form.
type NotableReport struct {
	message.BaseMessage
	baseform.BaseForm
	Jurisdiction     string
	Originator       string
	Severity         string
	EscalateToCounty string
	IncidentSpecific string
	IncidentName     string
	Title            string
	EventDate        string
	EventTime        string
	EventType        string
	Details          string
	ContactName      string
	ContactPhone     string
	ContactEmail     string
	Expires          string
	DateToExpire     string
}

func create() message.Message {
	f := makeF()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.EventDate = f.MessageDate
	f.ToLocation = "County EOC"
	return f
}

func makeF() *NotableReport {
	const fieldCount = 39
	f := NotableReport{BaseMessage: message.BaseMessage{Type: &Type}}
	f.FSubject = &f.Title
	f.FBody = &f.Details
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:       "Jurisdiction",
			Value:       &f.Jurisdiction,
			Choices:     message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara (City)", "Saratoga", "Sunnyvale", "Santa Clara County", "County unincorporated"},
			Presence:    message.Required,
			PIFOTag:     "20.",
			PDFRenderer: &message.PDFTextRenderer{X: 190, Y: 184, R: 373, B: 206, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   21,
			EditHelp:    `This is the name of the jurisdiction originating the notable report.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Originator",
			Value:       &f.Originator,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the originator of the notable report.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Severity",
			Value:       &f.Severity,
			Presence:    message.Required,
			PIFOTag:     "22.",
			Choices:     message.Choices{"High", "Medium", "Low"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the severity of the notable report.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Escalate to County Op Area?",
			Value:       &f.EscalateToCounty,
			Presence:    message.Required,
			PIFOTag:     "23.",
			Choices:     message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the notable report should be escalated to the county operational area.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label: "Incident Specific?",
			Value: &f.IncidentSpecific,
			Presence: func() (message.Presence, string) {
				if f.EscalateToCounty == "Yes" {
					return message.PresenceRequired, `"Escalate to County Op Area" is "Yes"`
				} else {
					return message.PresenceNotAllowed, `"Escalate to County Op Area" is not "Yes"`
				}
			},
			PIFOTag:     "24.",
			Choices:     message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the notable report is specific to a particular incident.  It is required if "Escalate to County Op Area" is "Yes", and otherwise not allowed.`,
		}),
		message.NewTextField(&message.Field{
			Label: "Name of Op Area Incident",
			Value: &f.IncidentName,
			Presence: func() (message.Presence, string) {
				if f.IncidentSpecific == "Yes" {
					return message.PresenceRequired, `"Incident Specific" is "Yes"`
				} else {
					return message.PresenceNotAllowed, `"Incident Specific" is not "Yes"`
				}
			},
			PIFOTag:     "25.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the county operational area's name for the incident leading to this notable report.  It is required when "Incident Specific" is "Yes", and otherwise not allowed.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Title",
			Value:       &f.Title,
			Presence:    message.Required,
			PIFOTag:     "30.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the title of the notable report.  It should be unique among all notable reports from the same jurisdiction.  It is required.`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Date",
			Value:       &f.EventDate,
			Presence:    message.Required,
			PIFOTag:     "31d.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the date of the event being reported.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time",
			Value:       &f.EventTime,
			Presence:    message.Required,
			PIFOTag:     "31t.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the time of day of the event being reported.  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time",
			EditHelp: `This is the date and time of the event being reported, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.EventDate, &f.EventTime),
		message.NewRestrictedField(&message.Field{
			Label:    "Event Type",
			Value:    &f.EventType,
			Presence: message.Required,
			PIFOTag:  "32.",
			Choices: message.ChoicePairs{
				"1", "1: Transportation",
				"2", "2: Communications",
				"3", "3: Construction and Engineering",
				"4", "4: Fire and Rescue",
				"5", "5: Management",
				"6", "6: Care and Shelter",
				"7", "7: Resources",
				"8", "8: Public Health and Medical",
				"9", "9: Search and Rescue",
				"10", "10: Hazardous Materials Response",
				"11", "11: Food and Agriculture",
				"12", "12: Utilities",
				"13", "13: Law Enforcement",
				"14", "14: Recovery",
				"15", "15: Public Information",
				"16", "16: Animal Services",
				"17", "17: Volunteer Management",
				"18", "18: Cybersecurity",
				"19", "19: Donations Management",
				"20", "20: Continuity of Operations/Government",
			},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the type of event being reported.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Details",
			Value:       &f.Details,
			Presence:    message.Required,
			PIFOTag:     "33.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This gives the details of the report.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Point of Contact Name",
			Value:       &f.ContactName,
			PIFOTag:     "40.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This gives the name of the contact person for the report.  It is required.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Contact Phone Number",
			Value:       &f.ContactPhone,
			PIFOTag:     "41.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This gives the phone number of the contact person for the report.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Contact Email",
			Value:       &f.ContactEmail,
			PIFOTag:     "42.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This gives the email address of the contact person for the report.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Event Expires?",
			Value:       &f.Expires,
			Presence:    message.Required,
			PIFOTag:     "50.",
			Choices:     message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the event being reported will expire.  It is required.`,
		}),
		message.NewDateField(false, &message.Field{
			Label: "Date to Expire",
			Value: &f.DateToExpire,
			Presence: func() (message.Presence, string) {
				if f.Expires == "Yes" {
					return message.PresenceOptional, ""
				}
				return message.PresenceNotAllowed, `"Event Expires" is not "Yes"`
			},
			PIFOTag:     "51.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the date when the event being reported will expire.  It is only allowed when "Event Expires" is "Yes".`,
		}),
	)
	f.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update NotableReport fieldCount")
	}
	return &f
}

func decode(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *NotableReport

	if form == nil || form.HTMLIdent != Type.HTML || form.FormVersion != Type.Version {
		return nil
	}
	df = makeF()
	message.DecodeForm(form, df)
	return df
}
