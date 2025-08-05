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
	OriginMsgID:      &message.PDFTextRenderer{X: 224.52, Y: 62.04, W: 106.80, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	DestinationMsgID: &message.PDFTextRenderer{X: 437.04, Y: 62.04, W: 123.72, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 73.32, Y: 132.36, W: 42.00, H: 11.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageTime:      &message.PDFTextRenderer{X: 186.00, Y: 132.36, W: 41.76, H: 11.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
	Handling: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
		"IMMEDIATE": {299.30, 137.85},
		"PRIORITY":  {401.42, 137.85},
		"ROUTINE":   {487.46, 137.85},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 128.28, Y: 150.84, W: 144.48, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 128.28, Y: 169.32, W: 144.48, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 128.28, Y: 187.80, W: 144.48, H: 11.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 128.28, Y: 206.40, W: 144.48, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 371.28, Y: 150.84, W: 188.16, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 371.28, Y: 169.32, W: 188.16, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 371.28, Y: 187.80, W: 188.16, H: 11.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 371.28, Y: 206.40, W: 188.16, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{X: 111.00, Y: 665.76, W: 205.56, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{X: 357.12, Y: 665.76, W: 203.64, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{X: 77.04, Y: 684.72, W: 123.72, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{X: 259.08, Y: 684.72, W: 63.24, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{X: 363.96, Y: 684.72, W: 52.80, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{X: 486.84, Y: 684.72, W: 73.92, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
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
			PDFRenderer: &message.PDFTextRenderer{X: 103.48, Y: 224.88, W: 191.54, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the name of the jurisdiction originating the notable report.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Originator",
			Value:       &f.Originator,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 362.28, Y: 224.88, W: 197.16, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   999,
			EditHelp:    `This is the originator of the notable report.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Severity",
			Value:    &f.Severity,
			Presence: message.Required,
			PIFOTag:  "22.",
			Choices:  message.Choices{"High", "Medium", "Low"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"High":   {192.96, 248.85},
				"Medium": {264.86, 248.85},
				"Low":    {372.86, 248.85},
			}},
			EditHelp: `This is the severity of the notable report.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Escalate to County Op Area?",
			Value:    &f.EscalateToCounty,
			Presence: message.Required,
			PIFOTag:  "23.",
			Choices:  message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Yes": {264.86, 265.89},
				"No":  {300.86, 265.89},
			}},
			EditHelp: `This indicates whether the notable report should be escalated to the county operational area.  It is required.`,
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
			PIFOTag: "24.",
			Choices: message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Yes": {264.86, 281.49},
				"No":  {300.86, 281.49},
			}},
			EditHelp: `This indicates whether the notable report is specific to a particular incident.  It is required if "Escalate to County Op Area" is "Yes", and otherwise not allowed.`,
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
			PDFRenderer: &message.PDFTextRenderer{X: 280.80, Y: 289.80, W: 278.64, H: 9.96},
			EditWidth:   59,
			EditHelp:    `This is the county operational area's name for the incident leading to this notable report.  It is required when "Incident Specific" is "Yes", and otherwise not allowed.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Title",
			Value:       &f.Title,
			Presence:    message.Required,
			PIFOTag:     "30.",
			PDFRenderer: &message.PDFTextRenderer{X: 71.28, Y: 324.48, W: 488.16, H: 11.52},
			EditWidth:   102,
			EditHelp:    `This is the title of the notable report.  It should be unique among all notable reports from the same jurisdiction.  It is required.`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Date",
			Value:       &f.EventDate,
			Presence:    message.Required,
			PIFOTag:     "31d.",
			PDFRenderer: &message.PDFTextRenderer{X: 73.32, Y: 342.48, W: 221.28, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date of the event being reported.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time",
			Value:       &f.EventTime,
			Presence:    message.Required,
			PIFOTag:     "31t.",
			PDFRenderer: &message.PDFTextRenderer{X: 338.64, Y: 342.48, W: 220.80, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
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
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"1":  {60.14, 378.45},
				"2":  {60.14, 393.93},
				"3":  {60.14, 409.53},
				"4":  {60.14, 425.01},
				"5":  {60.14, 440.61},
				"6":  {60.14, 456.09},
				"7":  {60.14, 471.69},
				"8":  {228.86, 378.45},
				"9":  {228.86, 393.93},
				"10": {228.86, 409.53},
				"11": {228.86, 425.01},
				"12": {228.86, 440.61},
				"13": {228.86, 456.09},
				"14": {228.86, 471.69},
				"15": {408.86, 378.45},
				"16": {408.86, 393.93},
				"17": {408.86, 409.53},
				"18": {408.86, 425.01},
				"19": {408.86, 440.61},
				"20": {408.86, 456.09},
			}},
			EditHelp: `This is the type of event being reported.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Details",
			Value:       &f.Details,
			Presence:    message.Required,
			PIFOTag:     "33.",
			PDFRenderer: &message.PDFTextRenderer{X: 44.64, Y: 508.92, W: 514.80, H: 22.08, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
			EditHelp:    `This gives the details of the report.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Point of Contact Name",
			Value:       &f.ContactName,
			PIFOTag:     "40.",
			PDFRenderer: &message.PDFTextRenderer{X: 154.56, Y: 564.84, W: 140.76, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   30,
			EditHelp:    `This gives the name of the contact person for the report.  It is required.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Contact Phone Number",
			Value:       &f.ContactPhone,
			PIFOTag:     "41.",
			PDFRenderer: &message.PDFTextRenderer{X: 422.16, Y: 564.84, W: 137.28, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   29,
			EditHelp:    `This gives the phone number of the contact person for the report.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Contact Email",
			Value:       &f.ContactEmail,
			PIFOTag:     "42.",
			PDFRenderer: &message.PDFTextRenderer{X: 154.56, Y: 582.48, W: 404.88, H: 11.40},
			EditWidth:   85,
			EditHelp:    `This gives the email address of the contact person for the report.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Event Expires?",
			Value:    &f.Expires,
			Presence: message.Required,
			PIFOTag:  "50.",
			Choices:  message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Yes": {209.13, 632.01},
				"No":  {255.02, 632.01},
			}},
			EditHelp: `This indicates whether the event being reported will expire.  It is required.`,
		}),
		message.NewDateField(false, &message.Field{
			Label: "Date to Expire",
			Value: &f.DateToExpire,
			Presence: func() (message.Presence, string) {
				if f.Expires == "Yes" {
					return message.PresenceRequired, `"Event Expires" is "Yes"`
				}
				return message.PresenceNotAllowed, `"Event Expires" is not "Yes"`
			},
			PIFOTag:     "51.",
			PDFRenderer: &message.PDFTextRenderer{X: 410.16, Y: 627.48, W: 149.28, H: 9.60},
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
