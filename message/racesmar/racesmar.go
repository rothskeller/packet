// Package racesmar defines the RACES Mutual Aid Request Form message type.
package racesmar

import (
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/baseform"
	"github.com/rothskeller/packet/message/basemsg"
	"github.com/rothskeller/packet/message/common"
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
var versions = []*basemsg.FormVersion{
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
	OriginMsgID: basemsg.PDFMapFunc(func(f *basemsg.Field) []basemsg.PDFField {
		return []basemsg.PDFField{
			{Name: "OriginMsg", Value: *f.Value},
			{Name: "ApprovedDate_2", Value: *f.Value},
		}
	}),
	DestinationMsgID: basemsg.PDFName("DestinationMsg"),
	MessageDate:      basemsg.PDFName("FormDate"),
	MessageTime:      basemsg.PDFName("FormTime"),
	Handling: basemsg.PDFNameMap{"Immediate",
		"", "Off",
		"IMMEDIATE", "1",
		"PRIORITY", "2",
		"ROUTINE", "3",
	},
	ToICSPosition:   basemsg.PDFName("ToICSPosition"),
	ToLocation:      basemsg.PDFName("ToICSLocation"),
	ToName:          basemsg.PDFName("ToICSPosition_3"),
	ToContact:       basemsg.PDFName("ToICSPosition_4"),
	FromICSPosition: basemsg.PDFName("ToICSPosition_2"),
	FromLocation:    basemsg.PDFName("ToICSLocation_2"),
	FromName:        basemsg.PDFName("ToICSPosition_5"),
	FromContact:     basemsg.PDFName("ToICSPosition_6"),
	OpRelayRcvd:     basemsg.PDFName("OperatorRcvd"),
	OpRelaySent:     basemsg.PDFName("OperatorSent"),
	OpName:          basemsg.PDFName("OperatorName"),
	OpCall:          basemsg.PDFName("OperatorCallSign"),
	OpDate:          basemsg.PDFName("OperatorDate"),
	OpTime:          basemsg.PDFName("OperatorTime"),
}

// RACESMAR holds a RACES mutual aid request form.
type RACESMAR struct {
	basemsg.BaseMessage
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

var pdfBase []byte

func create(version *basemsg.FormVersion) message.Message {
	const fieldCount = 74
	var f = RACESMAR{BaseMessage: basemsg.BaseMessage{
		MessageType: &Type,
		PDFBase:     pdfBase,
		Form:        version,
	}}
	f.BaseMessage.FSubject = &f.AgencyName
	f.BaseMessage.FBody = &f.Assignment
	f.Fields = make([]*basemsg.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFMap)
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:     "Agency Name",
			Value:     &f.AgencyName,
			PIFOTag:   "15.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("ToICSPosition_7"),
			EditWidth: 80,
			EditHelp:  `This is the name of the agency requesting mutual aid.  It is required.`,
		},
		&basemsg.Field{
			Label:      "Event Name",
			Value:      &f.EventName,
			PIFOTag:    "16a.",
			Presence:   basemsg.Required,
			Compare:    common.CompareText,
			PDFMap:     basemsg.PDFName("ToICSPosition_8"),
			TableValue: basemsg.OmitFromTable,
			EditWidth:  52,
			EditHelp:   `This is the name of the event for which mutual aid is being requested.  It is required.`,
		},
		&basemsg.Field{
			Label:      "Event Number",
			Value:      &f.EventNumber,
			PIFOTag:    "16b.",
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("Nbr"),
			TableValue: basemsg.OmitFromTable,
			EditWidth:  17,
			EditHelp:   `This is the requesting agency's activation number for the event for which mutual aid is being requested.`,
		},
		&basemsg.Field{
			Label: "Event Name/Number",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.EventName, f.EventNumber, " ")
			},
		},
		&basemsg.Field{
			Label:     "Assignment",
			Value:     &f.Assignment,
			PIFOTag:   "17.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Assignment"),
			EditWidth: 94,
			Multiline: true,
			EditHelp:  `This is a description of the assignment.  Describe the type of duties, conditions, any special equipment needed (other than 12-hour Go Kit). If multiple shifts are involved, give details. Provide enough detail for volunteer to decide if they are willing and able to accept the assignment.  This field is required.`,
		},
	)
	switch f.Form.Version {
	case "1.6":
		f.Fields = append(f.Fields,
			&basemsg.Field{
				Label:     "Resource Quantity",
				Value:     &f.Resources[0].Qty,
				PIFOTag:   "18a.",
				Presence:  basemsg.Required,
				Compare:   common.CompareText,
				PDFMap:    basemsg.PDFName("Qty1"),
				EditWidth: 2,
				Multiline: true,
				EditHelp:  `This is the number of people requested.  It is required.`,
			},
			&basemsg.Field{
				Label:     "Role/Position",
				Value:     &f.Resources[0].RolePos,
				PIFOTag:   "18b.",
				Presence:  basemsg.Required,
				Compare:   common.CompareText,
				PDFMap:    basemsg.PDFName("Position1"),
				EditWidth: 31,
				Multiline: true,
				EditHelp:  `This is the role and position for which people are requested.  It is required.`,
			},
			&basemsg.Field{
				Label:     "Preferred Type",
				Value:     &f.Resources[0].PreferredType,
				PIFOTag:   "18c.",
				Presence:  basemsg.Required,
				Compare:   common.CompareText,
				PDFMap:    basemsg.PDFName("Pref1"),
				EditWidth: 7,
				Multiline: true,
				EditHelp:  `This is the preferred resource type (credential) for the people being requested.  It is required.`,
			},
			&basemsg.Field{
				Label:     "Minimum Type",
				Value:     &f.Resources[0].MinimumType,
				PIFOTag:   "18d.",
				Presence:  basemsg.Required,
				Compare:   common.CompareText,
				PDFMap:    basemsg.PDFName("Min1"),
				EditWidth: 7,
				Multiline: true,
				EditHelp:  `This is the minimum resource type (credential) for the people being requested.  It is required.`,
			},
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
		&basemsg.Field{
			Label:     "Requested Arrival Dates",
			Value:     &f.RequestedArrivalDates,
			PIFOTag:   "19a.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("ReqArriveDates"),
			EditWidth: 46,
			EditHelp:  `This is the date(s) by when the requested people need to arrive.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Requested Arrival Times",
			Value:     &f.RequestedArrivalTimes,
			PIFOTag:   "19b.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("ReqArriveTimes"),
			EditWidth: 28,
			EditHelp:  `This is the time(s) by when the requested people need to arrive.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Needed Until Dates",
			Value:     &f.NeededUntilDates,
			PIFOTag:   "20a.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("NeedUntilDates"),
			EditWidth: 46,
			EditHelp:  `This is the date(s) until which the requested people will be needed.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Needed Until Times",
			Value:     &f.NeededUntilTimes,
			PIFOTag:   "20b.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("NeedUntilTimes"),
			EditWidth: 28,
			EditHelp:  `This is the time(s) until which the requested people will be needed.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Reporting Location",
			Value:     &f.ReportingLocation,
			PIFOTag:   "21.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("ReportingLocation"),
			EditWidth: 94,
			Multiline: true,
			EditHelp:  `This is the location to which the requested people should report.  Include street address, parking info, and entry instructions.  This field is required.`,
		},
		&basemsg.Field{
			Label:     "Contact on Arrival",
			Value:     &f.ContactOnArrival,
			PIFOTag:   "22.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("ContactOnArrival"),
			EditWidth: 94,
			EditHelp:  `This is the name, position, and contact info (phone, frequency, ...) for the official that the requested people should contact upon arrival. This is typically a net control on a radio frequency or a specific person or function at a telephone number.  This field is required.`,
			Multiline: true,
		},
		&basemsg.Field{
			Label:     "Travel Info",
			Value:     &f.TravelInfo,
			PIFOTag:   "23.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("TravelInfo"),
			EditWidth: 94,
			EditHelp:  `This field describes how to travel to the reporting location.  Identify preferred routes, road closures, and hazards to be avoided during travel.  If an overnight stay is included, specify how lodging will be provided.  This field is required.`,
			Multiline: true,
		},
		&basemsg.Field{
			Label:     "Requested By Name",
			Value:     &f.RequestedByName,
			PIFOTag:   "24a.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("RequestedName"),
			EditWidth: 45,
			EditHelp:  `This is the name of the official requesting mutual aid (typically the Radio Officer of the requesting agency).  It is required.`,
		},
		&basemsg.Field{
			Label:     "Requested By Title",
			Value:     &f.RequestedByTitle,
			PIFOTag:   "24b.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("RequestedTitle"),
			EditWidth: 31,
			EditHelp:  `This is the title of the official requesting mutual aid.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Requested By Contact",
			Value:     &f.RequestedByContact,
			PIFOTag:   "24c.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("RequestedContact"),
			EditWidth: 94,
			EditHelp:  `This is the contact information (email, phone, frequency) of the official requesting mutual aid.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Approved By Name",
			Value:     &f.ApprovedByName,
			PIFOTag:   "25a.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("ApprovedName"),
			EditWidth: 45,
			EditHelp:  `This is the name of the agency official approving the mutual aid request.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Approved By Title",
			Value:     &f.ApprovedByTitle,
			PIFOTag:   "25b.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("ApprovedTitle"),
			EditWidth: 31,
			EditHelp:  `This is the title of the agency official approving the mutual aid request.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Approved By Contact",
			Value:     &f.ApprovedByContact,
			PIFOTag:   "25c.",
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("ApprovedContact"),
			EditWidth: 94,
			EditHelp:  `This is the contact information (email, phone, frequency) for the agency official approving the mutual aid request.  It is required.`,
		},
	)
	if f.Form.Version >= "2.4" {
		f.Fields = append(f.Fields,
			&basemsg.Field{
				Label:     "With Signature",
				Value:     &f.WithSignature,
				PIFOTag:   "25s.",
				Choices:   basemsg.Choices{"checked"},
				PIFOValid: basemsg.ValidRestricted,
				Compare:   common.CompareExact,
				PDFMap:    basemsg.PDFNameMap{"ApprovedSignature", "checked", "[with signature]"},
				EditWidth: 7,
				EditHelp:  `This indicates that the original resource request form has been signed.`,
			},
		)
	}
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:      "Approved By Date",
			Value:      &f.ApprovedByDate,
			Presence:   basemsg.Required,
			PIFOTag:    "26a.",
			PIFOValid:  basemsg.ValidDate,
			Compare:    common.CompareDate,
			PDFMap:     basemsg.PDFName("ApprovedDate"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:      "Approved By Time",
			Value:      &f.ApprovedByTime,
			Presence:   basemsg.Required,
			PIFOTag:    "26b.",
			PIFOValid:  basemsg.ValidTime,
			Compare:    common.CompareTime,
			PDFMap:     basemsg.PDFName("ApprovedTime"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Approved Date/Time",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.ApprovedByDate, f.ApprovedByTime, " ")
			},
			EditWidth: 16,
			EditHelp:  `This is the date and time when the mutual aid request was approved by the official listed above, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
			EditHint:  "MM/DD/YYYY HH:MM",
			EditValue: func(_ *basemsg.Field) string {
				return basemsg.ValueDateTime(f.ApprovedByDate, f.ApprovedByTime)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				basemsg.ApplyDateTime(&f.ApprovedByDate, &f.ApprovedByTime, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return basemsg.ValidDateTime(field, f.ApprovedByDate, f.ApprovedByTime)
			},
		},
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
	return basemsg.Decode(body, versions, create).(*RACESMAR)
}
