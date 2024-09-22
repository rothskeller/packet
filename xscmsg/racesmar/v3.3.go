// Package racesmar defines the RACES Mutual Aid Request Form message type.
package racesmar

import (
	"fmt"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type33 is the type definition for a RACES mutual aid request form, version
// 3.3.
var Type33 = message.Type{
	Tag:     "RACES-MAR",
	HTML:    "form-oa-mutual-aid-request-v3.html",
	Version: "3.3",
	Name:    "RACES mutual aid request form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.",
		"8c.", "7d.", "8d.", "15.", "15b.", "16a.", "16b.", "17.",
		"18a.", "18b.", "18c.", "18d.", "18.1a.", "18.1e.", "18.1f.",
		"18.1g.", "18.1b.", "18.1c.", "18.1d.", "18.2a.", "18.2e.",
		"18.2f.", "18.2g.", "18.2b.", "18.2c.", "18.2d.", "18.3a.",
		"18.3e.", "18.3f.", "18.3g.", "18.3b.", "18.3c.", "18.3d.",
		"18.4a.", "18.4e.", "18.4f.", "18.4g.", "18.4b.", "18.4c.",
		"18.4d.", "18.5a.", "18.5e.", "18.5f.", "18.5g.", "18.5b.",
		"18.5c.", "18.5d.", "19a.", "19b.", "20a.", "20b.", "20c.",
		"20d.", "21.", "22.", "23.", "24a.", "24b.", "24c.", "25a.",
		"25b.", "25c.", "25s.", "26a.", "26b.", "OpRelayRcvd",
		"OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
	},
}

func init() {
	// Register 2.4 first so that it's the default for creation if no
	// version number is given.
	message.Register(&Type24, decode24, create24)
	message.Register(&Type33, decode33, create33)
}

var basePDFRenderers33 = baseform.BaseFormPDF{
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

// RACESMAR33 holds a RACES mutual aid request form.
type RACESMAR33 struct {
	message.BaseMessage
	baseform.BaseForm
	AgencyName           string
	OriginalMsgNumber    string
	EventName            string
	EventNumber          string
	Assignment           string
	Resources            [5]Resource33
	RequestedArrivalDate string
	RequestedArrivalTime string
	OpStartDate          string
	OpStartTime          string
	OpEndDate            string
	OpEndTime            string
	ReportingLocation    string
	ContactOnArrival     string
	TravelInfo           string
	RequestedByName      string
	RequestedByTitle     string
	RequestedByContact   string
	ApprovedByName       string
	ApprovedByTitle      string
	ApprovedByContact    string
	WithSignature        string
	ApprovedByDate       string
	ApprovedByTime       string
}

// A Resource33 is the description of a single resource in a RACES mutual aid
// request form, version 3.3.
type Resource33 struct {
	Qty           string
	Role          string
	Position      string
	RolePos       string
	Together      string
	PreferredType string
	MinimumType   string
}

func create33() message.Message {
	var f = make33()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Handling = "ROUTINE"
	f.ToLocation = "County EOC"
	return f
}

func make33() *RACESMAR33 {
	const fieldCount = 85
	var f = RACESMAR33{BaseMessage: message.BaseMessage{Type: &Type33}}
	f.BaseMessage.FSubject = &f.AgencyName
	f.BaseMessage.FBody = &f.Assignment
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFRenderers33)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:       "Agency Name",
			Value:       &f.AgencyName,
			Presence:    message.Required,
			PIFOTag:     "15.",
			PDFRenderer: &message.PDFTextRenderer{X: 190, Y: 184, R: 373, B: 206, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   80,
			EditHelp:    `This is the name of the agency requesting mutual aid.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Original Message Number",
			Value:       &f.OriginalMsgNumber,
			PIFOTag:     "15b.",
			PDFRenderer: &message.PDFTextRenderer{X: 472, Y: 184, R: 572, B: 206, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   80,
			EditHelp:    `If this is a reauthorization, put the message number of the original authorization here.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Event Name",
			Value:       &f.EventName,
			Presence:    message.Required,
			PIFOTag:     "16a.",
			PDFRenderer: &message.PDFTextRenderer{X: 190, Y: 209, R: 420, B: 225, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditWidth:   52,
			EditHelp:    `This is the name of the event for which mutual aid is being requested.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Event Number",
			Value:       &f.EventNumber,
			PIFOTag:     "16b.",
			Compare:     message.CompareExact,
			PDFRenderer: &message.PDFTextRenderer{X: 489, Y: 209, R: 572, B: 225, Style: message.PDFTextStyle{VAlign: "baseline"}},
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
			PDFRenderer: &message.PDFTextRenderer{X: 132, Y: 229, R: 572, B: 330, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   94,
			EditHelp:    `This is a description of the assignment.  Describe the type of duties, conditions, any special equipment needed (other than 12-hour Go Kit). If multiple shifts are involved, give details. Provide enough detail for volunteer to decide if they are willing and able to accept the assignment.  This field is required.`,
		}),
	)
	for i := range f.Resources {
		f.Fields = append(f.Fields, f.Resources[i].Fields(i+1)...)
	}
	f.Fields = append(f.Fields,
		message.NewDateField(true, &message.Field{
			Label:       "Requested Arrival Date",
			Value:       &f.RequestedArrivalDate,
			Presence:    message.Required,
			PIFOTag:     "19a.",
			PDFRenderer: &message.PDFTextRenderer{X: 160, Y: 467, R: 388, B: 483, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date by when the requested people need to arrive.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Requested Arrival Time",
			Value:       &f.RequestedArrivalTime,
			Presence:    message.Required,
			PIFOTag:     "19b.",
			PDFRenderer: &message.PDFTextRenderer{X: 449, Y: 467, R: 572, B: 483, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time by when the requested people need to arrive.  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Requested Arrival",
			EditHelp: `This is the date and time by when the requested people need to arrive, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.RequestedArrivalDate, &f.RequestedArrivalTime),
		message.NewDateField(true, &message.Field{
			Label:       "Operational Period Start Date",
			Value:       &f.OpStartDate,
			Presence:    message.Required,
			PIFOTag:     "20a.",
			PDFRenderer: &message.PDFTextRenderer{X: 162, Y: 486, R: 228, B: 501, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the operational period will start for the requested resources.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Operational Period Start Time",
			Value:       &f.OpStartTime,
			Presence:    message.Required,
			PIFOTag:     "20b.",
			PDFRenderer: &message.PDFTextRenderer{X: 311, Y: 486, R: 348, B: 501, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time when the operational period will start for the requested resources.  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Operational Period Start",
			EditHelp: `This is the date and time when the operational period will start for the requested resources, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.OpStartDate, &f.OpStartTime),
		message.NewDateField(true, &message.Field{
			Label:       "Operational Period End Date",
			Value:       &f.OpEndDate,
			Presence:    message.Required,
			PIFOTag:     "20c.",
			PDFRenderer: &message.PDFTextRenderer{X: 389, Y: 486, R: 451, B: 501, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date by when the operational period will end for the requested resources.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Operational Period End Time",
			Value:       &f.OpEndTime,
			Presence:    message.Required,
			PIFOTag:     "20d.",
			PDFRenderer: &message.PDFTextRenderer{X: 525, Y: 486, R: 572, B: 501, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time by when the operational period will end for the requested resources.  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Operational Period End",
			EditHelp: `This is the date and time by when the operational period will end for the requested resources, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.OpEndDate, &f.OpEndTime),
		message.NewMultilineField(&message.Field{
			Label:       "Reporting Location",
			Value:       &f.ReportingLocation,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 132, Y: 506, R: 572, B: 528, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   94,
			EditHelp:    `This is the location to which the requested people should report.  Include street address, parking info, and entry instructions.  This field is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Contact on Arrival",
			Value:       &f.ContactOnArrival,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 132, Y: 532, R: 572, B: 555, Style: message.PDFTextStyle{VAlign: "top", MinFontSize: 8}},
			EditWidth:   94,
			EditHelp:    `This is the name, position, and contact info (phone, frequency, ...) for the official that the requested people should contact upon arrival. This is typically a net control on a radio frequency or a specific person or function at a telephone number.  This field is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Travel Info",
			Value:       &f.TravelInfo,
			Presence:    message.Required,
			PIFOTag:     "23.",
			PDFRenderer: &message.PDFTextRenderer{X: 132, Y: 560, R: 572, B: 582, Style: message.PDFTextStyle{VAlign: "top", MinFontSize: 8}},
			EditWidth:   94,
			EditHelp:    `This field describes how to travel to the reporting location.  Identify preferred routes, road closures, and hazards to be avoided during travel.  If an overnight stay is included, specify how lodging will be provided.  This field is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Requested By Name",
			Value:       &f.RequestedByName,
			Presence:    message.Required,
			PIFOTag:     "24a.",
			PDFRenderer: &message.PDFTextRenderer{X: 168, Y: 586, R: 388, B: 602, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   45,
			EditHelp:    `This is the name of the official requesting mutual aid (typically the Radio Officer of the requesting agency).  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Requested By Title",
			Value:       &f.RequestedByTitle,
			Presence:    message.Required,
			PIFOTag:     "24b.",
			PDFRenderer: &message.PDFTextRenderer{X: 420, Y: 586, R: 572, B: 602, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   31,
			EditHelp:    `This is the title of the official requesting mutual aid.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Requested By Contact",
			Value:       &f.RequestedByContact,
			PIFOTag:     "24c.",
			Compare:     message.CompareText,
			PDFRenderer: &message.PDFTextRenderer{X: 282, Y: 607, R: 572, B: 628, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   94,
			EditHelp:    `This is the contact information (email, phone, frequency) of the official requesting mutual aid.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Approved By Name",
			Value:       &f.ApprovedByName,
			Presence:    message.Required,
			PIFOTag:     "25a.",
			PDFRenderer: &message.PDFTextRenderer{X: 168, Y: 630, R: 388, B: 647, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   45,
			EditHelp:    `This is the name of the agency official approving the mutual aid request.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Approved By Title",
			Value:       &f.ApprovedByTitle,
			Presence:    message.Required,
			PIFOTag:     "25b.",
			PDFRenderer: &message.PDFTextRenderer{X: 420, Y: 630, R: 572, B: 647, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   31,
			EditHelp:    `This is the title of the agency official approving the mutual aid request.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Approved By Contact",
			Value:       &f.ApprovedByContact,
			Presence:    message.Required,
			PIFOTag:     "25c.",
			PDFRenderer: &message.PDFTextRenderer{X: 238, Y: 651, R: 572, B: 669},
			EditWidth:   94,
			EditHelp:    `This is the contact information (email, phone, frequency) for the agency official approving the mutual aid request.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "With Signature",
			Value:       &f.WithSignature,
			PIFOTag:     "25s.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFMappedTextRenderer{X: 180, Y: 671, B: 693, Map: map[string]string{"checked": "[with signature]"}},
			EditHelp:    `This indicates that the original resource request form has been signed.`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Approved By Date",
			Value:       &f.ApprovedByDate,
			Presence:    message.Required,
			PIFOTag:     "26a.",
			PDFRenderer: &message.PDFTextRenderer{X: 409, Y: 671, R: 479, B: 693, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the mutual aid request was approved by the official listed above, in MM/DD/YYYY format.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Approved By Time",
			Value:       &f.ApprovedByTime,
			Presence:    message.Required,
			PIFOTag:     "26b.",
			PDFRenderer: &message.PDFTextRenderer{X: 537, Y: 671, R: 572, B: 693, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time when the mutual aid request was approved by the official listed above, in HH:MM format (24-hour clock).  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Approved Date/Time",
			EditHelp: `This is the date and time when the mutual aid request was approved by the official listed above, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.ApprovedByDate, &f.ApprovedByTime),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &basePDFRenderers33)
	if len(f.Fields) > fieldCount {
		panic("update RACESMAR fieldCount")
	}
	return &f
}

func decode33(_, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *RACESMAR33

	if form == nil || form.HTMLIdent != Type33.HTML || form.FormVersion != Type33.Version {
		return nil
	}
	df = make33()
	message.DecodeForm(form, df)
	return df
}

func (f *RACESMAR33) Compare(actual message.Message) (int, int, []*message.CompareField) {
	switch act := actual.(type) {
	case *RACESMAR24:
		actual = act.convertTo33()
	case *RACESMAR23:
		actual = act.convertTo33()
	case *RACESMAR21:
		actual = act.convertTo33()
	case *RACESMAR16:
		actual = act.convertTo33()
	}
	return f.BaseMessage.Compare(actual)
}

var resourceTypes = message.Choices{"Type I", "Type II", "Type III", "Type IV", "Type V"}

func (r *Resource33) Fields(index int) []*message.Field {
	var qtyPresence, rolePresence, posPresence, togetherPresence, typePresence func() (message.Presence, string)
	if index == 1 {
		qtyPresence = message.Required
		rolePresence = message.Required
		togetherPresence = message.Optional
		typePresence = message.Required
	} else {
		rolePresence = r.requiredIfQtyElseNotAllowed
		posPresence = r.notAllowedWithoutQty
		togetherPresence = r.notAllowedWithoutQty
		typePresence = r.requiredIfQtyElseNotAllowed
	}
	return []*message.Field{
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Quantity", index),
			Value:       &r.Qty,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("18.%da.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 132, Y: 326.5 + 19.3*float64(index), R: 159, H: 17, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			EditWidth:   2,
			EditHelp:    `This is the number of people needed for the role and position requested on this row.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Role", index),
			Value:       &r.Role,
			Presence:    rolePresence,
			PIFOTag:     fmt.Sprintf("18.%de.", index),
			Choices:     message.Choices{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
			PDFRenderer: &message.PDFTextRenderer{X: 163, Y: 326.5 + 19.3*float64(index), R: 268, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the role of the people requested on this row.  It is required when there is a quantity on the row.`,
			EditApply: func(f *message.Field, s string) {
				*f.Value = f.Choices.ToPIFO(strings.TrimSpace(s))
				if r.Position != "" {
					r.RolePos = message.SmartJoin(r.Role, "/ "+r.Position, " ")
				}
			},
		}),
		message.NewTextField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Position", index),
			Value:       &r.Position,
			Presence:    posPresence,
			PIFOTag:     fmt.Sprintf("18.%df.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 272, Y: 326.5 + 19.3*float64(index), R: 408, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   31,
			EditHelp:    `This is the position to be held by the people requested on this row.`,
			EditApply: func(_ *message.Field, s string) {
				r.Position = strings.TrimSpace(s)
				if r.Position != "" {
					r.RolePos = message.SmartJoin(r.Role, "/ "+r.Position, " ")
				}
			},
		}),
		message.NewRestrictedField(&message.Field{
			Label:    fmt.Sprintf("Resource %d Together", index),
			Value:    &r.Together,
			Presence: togetherPresence,
			PIFOTag:  fmt.Sprintf("18.%dg.", index),
			Choices:  message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10, H: 10, Points: map[string][]float64{
				"checked": {410, 330.5 + 19.3*float64(index)},
			}},
			EditHelp: `If this box is checked, this position should be filled only if the other positions with this box checked can also be filled.`,
		}),
		message.NewCalculatedField(&message.Field{
			Label:   fmt.Sprintf("Resource %d Role/Position", index),
			Value:   &r.RolePos,
			PIFOTag: fmt.Sprintf("18.%db.", index),
		}),
		message.NewRestrictedField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Preferred Type", index),
			Value:       &r.PreferredType,
			Presence:    typePresence,
			Choices:     resourceTypes,
			PIFOTag:     fmt.Sprintf("18.%dc.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 424, Y: 326.5 + 19.3*float64(index), R: 497, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the preferred resource type (credential) for the people requested on this row.  It is required when there is a quantity on the row.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Minimum Type", index),
			Value:       &r.MinimumType,
			Presence:    typePresence,
			Choices:     resourceTypes,
			PIFOTag:     fmt.Sprintf("18.%dd.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 501, Y: 326.5 + 19.3*float64(index), R: 572, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the minimum resource type (credential) for the people requested on this row.  It is required when there is a quantity on the row.`,
		}),
	}
}
func (r *Resource33) requiredIfQtyElseNotAllowed() (message.Presence, string) {
	if r.Qty != "" {
		return message.PresenceRequired, "there is a quantity for the resource"
	} else {
		return message.PresenceNotAllowed, "there is no quantity for the resource"
	}
}
func (r *Resource33) notAllowedWithoutQty() (message.Presence, string) {
	if r.Qty == "" {
		return message.PresenceNotAllowed, "there is no quantity for the resource"
	}
	return message.PresenceOptional, ""
}
