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

var basePDFMap = baseform.BaseFormPDFMaps{
	OriginMsgID: message.PDFMapFunc(func(f *message.Field) []message.PDFField {
		return []message.PDFField{
			{Name: "OriginMsg", Value: *f.Value},
			{Name: "ApprovedDate_2", Value: *f.Value},
		}
	}),
	DestinationMsgID: message.PDFName("DestinationMsg"),
	MessageDate:      message.PDFName("FormDate"),
	MessageTime:      message.PDFName("FormTime"),
	Handling: message.PDFNameMap{"Immediate",
		"", "Off",
		"IMMEDIATE", "1",
		"PRIORITY", "2",
		"ROUTINE", "3",
	},
	ToICSPosition:   message.PDFName("ToICSPosition"),
	ToLocation:      message.PDFName("ToICSLocation"),
	ToName:          message.PDFName("ToICSPosition_3"),
	ToContact:       message.PDFName("ToICSPosition_4"),
	FromICSPosition: message.PDFName("ToICSPosition_2"),
	FromLocation:    message.PDFName("ToICSLocation_2"),
	FromName:        message.PDFName("ToICSPosition_5"),
	FromContact:     message.PDFName("ToICSPosition_6"),
	OpRelayRcvd:     message.PDFName("OperatorRcvd"),
	OpRelaySent:     message.PDFName("OperatorSent"),
	OpName:          message.PDFName("OperatorName"),
	OpCall:          message.PDFName("OperatorCallSign"),
	OpDate:          message.PDFName("OperatorDate"),
	OpTime:          message.PDFName("OperatorTime"),
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
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFMap)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:     "Agency Name",
			Value:     &f.AgencyName,
			Presence:  message.Required,
			PIFOTag:   "15.",
			PDFMap:    message.PDFName("ToICSPosition_7"),
			EditWidth: 80,
			EditHelp:  `This is the name of the agency requesting mutual aid.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:      "Event Name",
			Value:      &f.EventName,
			Presence:   message.Required,
			PIFOTag:    "16a.",
			PDFMap:     message.PDFName("ToICSPosition_8"),
			TableValue: message.TableOmit,
			EditWidth:  52,
			EditHelp:   `This is the name of the event for which mutual aid is being requested.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:      "Event Number",
			Value:      &f.EventNumber,
			PIFOTag:    "16b.",
			Compare:    message.CompareExact,
			PDFMap:     message.PDFName("Nbr"),
			TableValue: message.TableOmit,
			EditWidth:  17,
			EditHelp:   `This is the requesting agency's activation number for the event for which mutual aid is being requested.`,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Event Name/Number",
			TableValue: func(*message.Field) string {
				return message.SmartJoin(f.EventName, f.EventNumber, " ")
			},
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Assignment",
			Value:     &f.Assignment,
			Presence:  message.Required,
			PIFOTag:   "17.",
			PDFMap:    message.PDFName("Assignment"),
			EditWidth: 94,
			EditHelp:  `This is a description of the assignment.  Describe the type of duties, conditions, any special equipment needed (other than 12-hour Go Kit). If multiple shifts are involved, give details. Provide enough detail for volunteer to decide if they are willing and able to accept the assignment.  This field is required.`,
		}),
	)
	switch f.Form.Version {
	case "1.6":
		f.Fields = append(f.Fields,
			message.NewMultilineField(&message.Field{
				Label:     "Resource Quantity",
				Value:     &f.Resources[0].Qty,
				Presence:  message.Required,
				PIFOTag:   "18a.",
				PDFMap:    message.PDFName("Qty1"),
				EditWidth: 2,
				EditHelp:  `This is the number of people requested.  It is required.`,
			}),
			message.NewMultilineField(&message.Field{
				Label:     "Role/Position",
				Value:     &f.Resources[0].RolePos,
				Presence:  message.Required,
				PIFOTag:   "18b.",
				PDFMap:    message.PDFName("Position1"),
				EditWidth: 31,
				EditHelp:  `This is the role and position for which people are requested.  It is required.`,
			}),
			message.NewMultilineField(&message.Field{
				Label:     "Preferred Type",
				Value:     &f.Resources[0].PreferredType,
				Presence:  message.Required,
				PIFOTag:   "18c.",
				PDFMap:    message.PDFName("Pref1"),
				EditWidth: 7,
				EditHelp:  `This is the preferred resource type (credential) for the people being requested.  It is required.`,
			}),
			message.NewMultilineField(&message.Field{
				Label:     "Minimum Type",
				Value:     &f.Resources[0].MinimumType,
				Presence:  message.Required,
				PIFOTag:   "18d.",
				PDFMap:    message.PDFName("Min1"),
				EditWidth: 7,
				EditHelp:  `This is the minimum resource type (credential) for the people being requested.  It is required.`,
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
			Label:     "Requested Arrival Dates",
			Value:     &f.RequestedArrivalDates,
			Presence:  message.Required,
			PIFOTag:   "19a.",
			PDFMap:    message.PDFName("ReqArriveDates"),
			EditWidth: 46,
			EditHelp:  `This is the date(s) by when the requested people need to arrive.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Requested Arrival Times",
			Value:     &f.RequestedArrivalTimes,
			Presence:  message.Required,
			PIFOTag:   "19b.",
			PDFMap:    message.PDFName("ReqArriveTimes"),
			EditWidth: 28,
			EditHelp:  `This is the time(s) by when the requested people need to arrive.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Needed Until Dates",
			Value:     &f.NeededUntilDates,
			Presence:  message.Required,
			PIFOTag:   "20a.",
			PDFMap:    message.PDFName("NeedUntilDates"),
			EditWidth: 46,
			EditHelp:  `This is the date(s) until which the requested people will be needed.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Needed Until Times",
			Value:     &f.NeededUntilTimes,
			Presence:  message.Required,
			PIFOTag:   "20b.",
			PDFMap:    message.PDFName("NeedUntilTimes"),
			EditWidth: 28,
			EditHelp:  `This is the time(s) until which the requested people will be needed.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Reporting Location",
			Value:     &f.ReportingLocation,
			Presence:  message.Required,
			PIFOTag:   "21.",
			PDFMap:    message.PDFName("ReportingLocation"),
			EditWidth: 94,
			EditHelp:  `This is the location to which the requested people should report.  Include street address, parking info, and entry instructions.  This field is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Contact on Arrival",
			Value:     &f.ContactOnArrival,
			Presence:  message.Required,
			PIFOTag:   "22.",
			PDFMap:    message.PDFName("ContactOnArrival"),
			EditWidth: 94,
			EditHelp:  `This is the name, position, and contact info (phone, frequency, ...) for the official that the requested people should contact upon arrival. This is typically a net control on a radio frequency or a specific person or function at a telephone number.  This field is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Travel Info",
			Value:     &f.TravelInfo,
			Presence:  message.Required,
			PIFOTag:   "23.",
			PDFMap:    message.PDFName("TravelInfo"),
			EditWidth: 94,
			EditHelp:  `This field describes how to travel to the reporting location.  Identify preferred routes, road closures, and hazards to be avoided during travel.  If an overnight stay is included, specify how lodging will be provided.  This field is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Requested By Name",
			Value:     &f.RequestedByName,
			Presence:  message.Required,
			PIFOTag:   "24a.",
			PDFMap:    message.PDFName("RequestedName"),
			EditWidth: 45,
			EditHelp:  `This is the name of the official requesting mutual aid (typically the Radio Officer of the requesting agency).  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Requested By Title",
			Value:     &f.RequestedByTitle,
			Presence:  message.Required,
			PIFOTag:   "24b.",
			PDFMap:    message.PDFName("RequestedTitle"),
			EditWidth: 31,
			EditHelp:  `This is the title of the official requesting mutual aid.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Requested By Contact",
			Value:     &f.RequestedByContact,
			PIFOTag:   "24c.",
			Compare:   message.CompareText,
			PDFMap:    message.PDFName("RequestedContact"),
			EditWidth: 94,
			EditHelp:  `This is the contact information (email, phone, frequency) of the official requesting mutual aid.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Approved By Name",
			Value:     &f.ApprovedByName,
			Presence:  message.Required,
			PIFOTag:   "25a.",
			PDFMap:    message.PDFName("ApprovedName"),
			EditWidth: 45,
			EditHelp:  `This is the name of the agency official approving the mutual aid request.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Approved By Title",
			Value:     &f.ApprovedByTitle,
			Presence:  message.Required,
			PIFOTag:   "25b.",
			PDFMap:    message.PDFName("ApprovedTitle"),
			EditWidth: 31,
			EditHelp:  `This is the title of the agency official approving the mutual aid request.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Approved By Contact",
			Value:     &f.ApprovedByContact,
			Presence:  message.Required,
			PIFOTag:   "25c.",
			PDFMap:    message.PDFName("ApprovedContact"),
			EditWidth: 94,
			EditHelp:  `This is the contact information (email, phone, frequency) for the agency official approving the mutual aid request.  It is required.`,
		}),
	)
	if f.Form.Version >= "2.4" {
		f.Fields = append(f.Fields,
			message.NewRestrictedField(&message.Field{
				Label:    "With Signature",
				Value:    &f.WithSignature,
				PIFOTag:  "25s.",
				Choices:  message.Choices{"checked"},
				PDFMap:   message.PDFNameMap{"ApprovedSignature", "checked", "[with signature]"},
				EditHelp: `This indicates that the original resource request form has been signed.`,
			}),
		)
	}
	f.Fields = append(f.Fields,
		message.NewDateWithTimeField(&message.Field{
			Label:    "Approved By Date",
			Value:    &f.ApprovedByDate,
			Presence: message.Required,
			PIFOTag:  "26a.",
			PDFMap:   message.PDFName("ApprovedDate"),
		}),
		message.NewTimeWithDateField(&message.Field{
			Label:    "Approved By Time",
			Value:    &f.ApprovedByTime,
			Presence: message.Required,
			PIFOTag:  "26b.",
			PDFMap:   message.PDFName("ApprovedTime"),
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Approved Date/Time",
			EditHelp: `This is the date and time when the mutual aid request was approved by the official listed above, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.ApprovedByDate, &f.ApprovedByTime),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &basePDFMap)
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
