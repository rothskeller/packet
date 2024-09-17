// Package racesmar defines the RACES Mutual Aid Request Form message type.
package racesmar

import (
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for a RACES mutual aid request form.
var Type = message.Type{
	Tag:     "RACES-MAR",
	Name:    "RACES mutual aid request form",
	Article: "a",
}

func init() {
	Type.Create = New
	Type.Decode = decode
}

// versions is the list of supported versions.  The first one is used when
// creating new forms.
var versions = []*message.FormVersion{
	{HTML: "form-oa-mutual-aid-request-v2.html", Version: "2.4", Tag: "RACES-MAR", FieldOrder: fieldOrder},
	{HTML: "form-oa-mutual-aid-request-v2.html", Version: "2.3", Tag: "RACES-MAR", FieldOrder: fieldOrder},
	{HTML: "form-oa-mutual-aid-request-v2.html", Version: "2.1", Tag: "RACES-MAR", FieldOrder: fieldOrder},
	{HTML: "form-oa-mutual-aid-request.html", Version: "1.6", Tag: "RACES-MAR", FieldOrder: fieldOrder},
}
var fieldOrder = []string{
	"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.",
	"7d.", "8d.", "15.", "16a.", "16b.", "17.", "18a.", "18b.", "18c.",
	"18d.", "18.1a.", "18.1e.", "18.1f.", "18.1b.", "18.1c.", "18.1d.",
	"18.2a.", "18.2e.", "18.2f.", "18.2b.", "18.2c.", "18.2d.", "18.3a.",
	"18.3e.", "18.3f.", "18.3b.", "18.3c.", "18.3d.", "18.4a.", "18.4e.",
	"18.4f.", "18.4b.", "18.4c.", "18.4d.", "18.5a.", "18.5e.", "18.5f.",
	"18.5b.", "18.5c.", "18.5d.", "19a.", "19b.", "20a.", "20b.", "21.",
	"22.", "23.", "24a.", "24b.", "24c.", "25a.", "25b.", "25c.", "25s.",
	"26a.", "26b.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall",
	"OpDate", "OpTime",
}
var basePDFRenderers = baseform.BaseFormPDF{
	OriginMsgID: &message.PDFMultiRenderer{
		&message.PDFTextRenderer{X: 227, Y: 53, R: 346, B: 64, Style: message.PDFTextStyle{VAlign: "baseline"}},
		&message.PDFTextRenderer{Page: 2, X: 394, Y: 33, R: 503, B: 45},
	},
	DestinationMsgID: &message.PDFTextRenderer{X: 456, Y: 53, R: 570, B: 64, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 78, Y: 91, R: 125, B: 101, Style: message.PDFTextStyle{VAlign: "baseline", FontSize: 9}},
	MessageTime:      &message.PDFTextRenderer{X: 178, Y: 91, R: 214, B: 101, Style: message.PDFTextStyle{VAlign: "baseline", FontSize: 9}},
	Handling: &message.PDFRadioRenderer{Radius: 3.5, Points: map[string][]float64{
		"IMMEDIATE": {308, 96},
		"PRIORITY":  {407, 96},
		"ROUTINE":   {493, 96},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 139, Y: 109, R: 290, B: 122, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 139, Y: 126, R: 288, B: 138, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 139, Y: 142, R: 290, B: 155, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 139, Y: 161, R: 290, B: 173, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 400, Y: 109, R: 551, B: 122, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 400, Y: 126, R: 549, B: 138, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 400, Y: 142, R: 551, B: 155, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 400, Y: 161, R: 551, B: 173, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{X: 110, Y: 716, R: 272, B: 728, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{X: 356, Y: 716, R: 520, B: 728, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{X: 75, Y: 734, R: 240, B: 748, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{X: 301, Y: 734, R: 358, B: 748, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{X: 430, Y: 734, R: 479, B: 748, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{X: 519, Y: 734, R: 572, B: 748, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

// RACESMAR holds a RACES mutual aid request form.
type RACESMAR struct {
	message.BaseMessage
	baseform.BaseForm
	AgencyName            string
	EventName             string
	EventNumber           string
	Assignment            string
	Resources             [5]Resource // made multiple in v2.1
	RequestedArrivalDates string
	RequestedArrivalTimes string
	NeededUntilDates      string
	NeededUntilTimes      string
	ReportingLocation     string
	ContactOnArrival      string
	TravelInfo            string
	RequestedByName       string
	RequestedByTitle      string
	RequestedByContact    string
	ApprovedByName        string
	ApprovedByTitle       string
	ApprovedByContact     string
	WithSignature         string // added in v2.4
	ApprovedByDate        string
	ApprovedByTime        string
}

func New() (f *RACESMAR) {
	f = create(versions[0]).(*RACESMAR)
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Handling = "ROUTINE"
	f.ToLocation = "County EOC"
	return f
}

func create(version *message.FormVersion) message.Message {
	const fieldCount = 74
	var f = RACESMAR{BaseMessage: message.BaseMessage{
		Type: &Type,
		Form: version,
	}}
	f.BaseMessage.FSubject = &f.AgencyName
	f.BaseMessage.FBody = &f.Assignment
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:       "Agency Name",
			Value:       &f.AgencyName,
			Presence:    message.Required,
			PIFOTag:     "15.",
			PDFRenderer: &message.PDFTextRenderer{X: 198, Y: 180, R: 551, B: 192},
			EditWidth:   80,
			EditHelp:    `This is the name of the agency requesting mutual aid.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Event Name",
			Value:       &f.EventName,
			Presence:    message.Required,
			PIFOTag:     "16a.",
			PDFRenderer: &message.PDFTextRenderer{X: 198, Y: 199, R: 429, B: 211},
			TableValue:  message.TableOmit,
			EditWidth:   52,
			EditHelp:    `This is the name of the event for which mutual aid is being requested.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Event Number",
			Value:       &f.EventNumber,
			PIFOTag:     "16b.",
			Compare:     message.CompareExact,
			PDFRenderer: &message.PDFTextRenderer{X: 490, Y: 197, R: 563, B: 208},
			TableValue:  message.TableOmit,
			EditWidth:   17,
			EditHelp:    `This is the requesting agency's activation number for the event for which mutual aid is being requested.`,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Event Name/Number",
			TableValue: func(*message.Field) string {
				return message.SmartJoin(f.EventName, f.EventNumber, " ")
			},
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Assignment",
			Value:       &f.Assignment,
			Presence:    message.Required,
			PIFOTag:     "17.",
			PDFRenderer: &message.PDFTextRenderer{X: 139, Y: 217, R: 562, B: 319, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   94,
			EditHelp:    `This is a description of the assignment.  Describe the type of duties, conditions, any special equipment needed (other than 12-hour Go Kit). If multiple shifts are involved, give details. Provide enough detail for volunteer to decide if they are willing and able to accept the assignment.  This field is required.`,
		}),
	)
	switch f.Form.Version {
	case "1.6":
		f.Fields = append(f.Fields,
			message.NewMultilineField(&message.Field{
				Label:       "Resource Quantity",
				Value:       &f.Resources[0].Qty,
				Presence:    message.Required,
				PIFOTag:     "18a.",
				PDFRenderer: &message.PDFTextRenderer{X: 139, Y: 351, R: 150, B: 362, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
				EditWidth:   2,
				EditHelp:    `This is the number of people requested.  It is required.`,
			}),
			message.NewMultilineField(&message.Field{
				Label:       "Role/Position",
				Value:       &f.Resources[0].RolePos,
				Presence:    message.Required,
				PIFOTag:     "18b.",
				PDFRenderer: &message.PDFTextRenderer{X: 170, Y: 351, R: 263, B: 362, Style: message.PDFTextStyle{VAlign: "baseline"}},
				EditWidth:   31,
				EditHelp:    `This is the role and position for which people are requested.  It is required.`,
			}),
			message.NewMultilineField(&message.Field{
				Label:       "Preferred Type",
				Value:       &f.Resources[0].PreferredType,
				Presence:    message.Required,
				PIFOTag:     "18c.",
				PDFRenderer: &message.PDFTextRenderer{X: 431, Y: 351, R: 489, B: 362, Style: message.PDFTextStyle{VAlign: "baseline"}},
				EditWidth:   7,
				EditHelp:    `This is the preferred resource type (credential) for the people being requested.  It is required.`,
			}),
			message.NewMultilineField(&message.Field{
				Label:       "Minimum Type",
				Value:       &f.Resources[0].MinimumType,
				Presence:    message.Required,
				PIFOTag:     "18d.",
				PDFRenderer: &message.PDFTextRenderer{X: 509, Y: 351, R: 567, B: 362, Style: message.PDFTextStyle{VAlign: "baseline"}},
				EditWidth:   7,
				EditHelp:    `This is the minimum resource type (credential) for the people being requested.  It is required.`,
			}),
		)
	case "2.1":
		for i := range f.Resources {
			f.Fields = append(f.Fields, f.Resources[i].Fields21(i+1)...)
		}
	default:
		for i := range f.Resources {
			f.Fields = append(f.Fields, f.Resources[i].Fields23(i+1)...)
		}
	}
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:       "Requested Arrival Dates",
			Value:       &f.RequestedArrivalDates,
			Presence:    message.Required,
			PIFOTag:     "19a.",
			PDFRenderer: &message.PDFTextRenderer{X: 170, Y: 442, R: 380, B: 456, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   46,
			EditHelp:    `This is the date(s) by when the requested people need to arrive.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Requested Arrival Times",
			Value:       &f.RequestedArrivalTimes,
			Presence:    message.Required,
			PIFOTag:     "19b.",
			PDFRenderer: &message.PDFTextRenderer{X: 432, Y: 442, R: 560, B: 456, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   28,
			EditHelp:    `This is the time(s) by when the requested people need to arrive.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Needed Until Dates",
			Value:       &f.NeededUntilDates,
			Presence:    message.Required,
			PIFOTag:     "20a.",
			PDFRenderer: &message.PDFTextRenderer{X: 170, Y: 461, R: 380, B: 473, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   46,
			EditHelp:    `This is the date(s) until which the requested people will be needed.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Needed Until Times",
			Value:       &f.NeededUntilTimes,
			Presence:    message.Required,
			PIFOTag:     "20b.",
			PDFRenderer: &message.PDFTextRenderer{X: 432, Y: 461, R: 560, B: 473, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   28,
			EditHelp:    `This is the time(s) until which the requested people will be needed.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Reporting Location",
			Value:       &f.ReportingLocation,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 139, Y: 479, R: 560, B: 510, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   94,
			EditHelp:    `This is the location to which the requested people should report.  Include street address, parking info, and entry instructions.  This field is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Contact on Arrival",
			Value:       &f.ContactOnArrival,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 139, Y: 520, R: 562, B: 538, Style: message.PDFTextStyle{VAlign: "top", MinFontSize: 8}},
			EditWidth:   94,
			EditHelp:    `This is the name, position, and contact info (phone, frequency, ...) for the official that the requested people should contact upon arrival. This is typically a net control on a radio frequency or a specific person or function at a telephone number.  This field is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Travel Info",
			Value:       &f.TravelInfo,
			Presence:    message.Required,
			PIFOTag:     "23.",
			PDFRenderer: &message.PDFTextRenderer{X: 139, Y: 547, R: 562, B: 565, Style: message.PDFTextStyle{VAlign: "top", MinFontSize: 8}},
			EditWidth:   94,
			EditHelp:    `This field describes how to travel to the reporting location.  Identify preferred routes, road closures, and hazards to be avoided during travel.  If an overnight stay is included, specify how lodging will be provided.  This field is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Requested By Name",
			Value:       &f.RequestedByName,
			Presence:    message.Required,
			PIFOTag:     "24a.",
			PDFRenderer: &message.PDFTextRenderer{X: 166, Y: 575, R: 370, B: 587, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   45,
			EditHelp:    `This is the name of the official requesting mutual aid (typically the Radio Officer of the requesting agency).  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Requested By Title",
			Value:       &f.RequestedByTitle,
			Presence:    message.Required,
			PIFOTag:     "24b.",
			PDFRenderer: &message.PDFTextRenderer{X: 420, Y: 575, R: 559, B: 587, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   31,
			EditHelp:    `This is the title of the official requesting mutual aid.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Requested By Contact",
			Value:       &f.RequestedByContact,
			PIFOTag:     "24c.",
			Compare:     message.CompareText,
			PDFRenderer: &message.PDFTextRenderer{X: 139, Y: 604, R: 560, B: 613, Style: message.PDFTextStyle{FontSize: 9}},
			EditWidth:   94,
			EditHelp:    `This is the contact information (email, phone, frequency) of the official requesting mutual aid.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Approved By Name",
			Value:       &f.ApprovedByName,
			Presence:    message.Required,
			PIFOTag:     "25a.",
			PDFRenderer: &message.PDFTextRenderer{X: 166, Y: 623, R: 370, B: 637, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   45,
			EditHelp:    `This is the name of the agency official approving the mutual aid request.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Approved By Title",
			Value:       &f.ApprovedByTitle,
			Presence:    message.Required,
			PIFOTag:     "25b.",
			PDFRenderer: &message.PDFTextRenderer{X: 420, Y: 623, R: 559, B: 637, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   31,
			EditHelp:    `This is the title of the agency official approving the mutual aid request.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Approved By Contact",
			Value:       &f.ApprovedByContact,
			Presence:    message.Required,
			PIFOTag:     "25c.",
			PDFRenderer: &message.PDFTextRenderer{X: 139, Y: 658, R: 560, B: 667, Style: message.PDFTextStyle{FontSize: 9}},
			EditWidth:   94,
			EditHelp:    `This is the contact information (email, phone, frequency) for the agency official approving the mutual aid request.  It is required.`,
		}),
	)
	if f.Form.Version >= "2.4" {
		f.Fields = append(f.Fields,
			message.NewRestrictedField(&message.Field{
				Label:       "With Signature",
				Value:       &f.WithSignature,
				PIFOTag:     "25s.",
				Choices:     message.Choices{"checked"},
				PDFRenderer: &message.PDFMappedTextRenderer{X: 139, Y: 682, B: 695, Map: map[string]string{"checked": "[with signature]"}},
				EditHelp:    `This indicates that the original resource request form has been signed.`,
			}),
		)
	} else {
		f.Fields = append(f.Fields,
			message.NewCalculatedField(new(message.Field)),
		)
	}
	f.Fields = append(f.Fields,
		message.NewDateField(true, &message.Field{
			Label:       "Approved By Date",
			Value:       &f.ApprovedByDate,
			Presence:    message.Required,
			PIFOTag:     "26a.",
			PDFRenderer: &message.PDFTextRenderer{X: 403, Y: 682, R: 459, B: 695, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the mutual aid request was approved by the official listed above, in MM/DD/YYYY format.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Approved By Time",
			Value:       &f.ApprovedByTime,
			Presence:    message.Required,
			PIFOTag:     "26b.",
			PDFRenderer: &message.PDFTextRenderer{X: 513, Y: 682, R: 565, B: 695, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time when the mutual aid request was approved by the official listed above, in HH:MM format (24-hour clock).  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Approved Date/Time",
			EditHelp: `This is the date and time when the mutual aid request was approved by the official listed above, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.ApprovedByDate, &f.ApprovedByTime),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update RACESMAR fieldCount")
	}
	return &f
}

func decode(subject, body string) (f *RACESMAR) {
	// Quick check to avoid overhead of creating the form object if it's not
	// our type of form.
	if !strings.Contains(body, "form-oa-mutual-aid-request") {
		return nil
	}
	if df, ok := message.DecodeForm(body, versions, create).(*RACESMAR); ok {
		return df
	} else {
		return nil
	}
}
