// Package racesmar defines the RACES Mutual Aid Request Form message type.
package racesmar

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type24 is the type definition for a RACES mutual aid request form.
var Type24 = message.Type{
	Tag:     "RACES-MAR",
	HTML:    "form-oa-mutual-aid-request-v2.html",
	Version: "2.4",
	Name:    "RACES mutual aid request form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.", "7d.",
		"8d.", "15.", "16a.", "16b.", "17.", "18a.", "18b.", "18c.", "18d.",
		"18.1a.", "18.1e.", "18.1f.", "18.1b.", "18.1c.", "18.1d.",
		"18.2a.", "18.2e.", "18.2f.", "18.2b.", "18.2c.", "18.2d.",
		"18.3a.", "18.3e.", "18.3f.", "18.3b.", "18.3c.", "18.3d.",
		"18.4a.", "18.4e.", "18.4f.", "18.4b.", "18.4c.", "18.4d.",
		"18.5a.", "18.5e.", "18.5f.", "18.5b.", "18.5c.", "18.5d.", "19a.",
		"19b.", "20a.", "20b.", "21.", "22.", "23.", "24a.", "24b.",
		"24c.", "25a.", "25b.", "25c.", "25s.", "26a.", "26b.", "OpRelayRcvd",
		"OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
	},
}

func init() {
	// moved to v3.3.go in order to enforce registration ordering
	// message.Register(&Type24, decode24, create24)
}

// RACESMAR24 holds a RACES mutual aid request form.
type RACESMAR24 struct {
	message.BaseMessage
	baseform.BaseForm
	AgencyName            string
	EventName             string
	EventNumber           string
	Assignment            string
	Resources             [5]Resource24
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
	WithSignature         string
	ApprovedByDate        string
	ApprovedByTime        string
}

// A Resource24 is the description of a single resource in a RACES mutual aid
// request form.
type Resource24 struct {
	Qty           string
	Role          string
	Position      string
	RolePos       string
	PreferredType string
	MinimumType   string
}

var basePDFRenderers24 = baseform.BaseFormPDF{
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

func create24() message.Message {
	var f = make24()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Handling = "ROUTINE"
	f.ToLocation = "County EOC"
	return f
}

func make24() *RACESMAR24 {
	const fieldCount = 74
	var f = RACESMAR24{BaseMessage: message.BaseMessage{Type: &Type24}}
	f.BaseMessage.FSubject = &f.AgencyName
	f.BaseMessage.FBody = &f.Assignment
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFRenderers24)
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
	for i := range f.Resources {
		f.Fields = append(f.Fields, f.Resources[i].Fields(i+1)...)
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
		message.NewRestrictedField(&message.Field{
			Label:       "With Signature",
			Value:       &f.WithSignature,
			PIFOTag:     "25s.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFMappedTextRenderer{X: 139, Y: 682, B: 695, Map: map[string]string{"checked": "[with signature]"}},
			EditHelp:    `This indicates that the original resource request form has been signed.`,
		}),
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
	f.BaseForm.AddFooterFields(&f.BaseMessage, &basePDFRenderers24)
	if len(f.Fields) > fieldCount {
		panic("update RACESMAR24 fieldCount")
	}
	return &f
}

func decode24(_, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *RACESMAR24

	if form == nil || form.HTMLIdent != Type24.HTML || form.FormVersion != Type24.Version {
		return nil
	}
	df = make24()
	message.DecodeForm(form, df)
	return df
}

func (f *RACESMAR24) Compare(actual message.Message) (int, int, []*message.CompareField) {
	if _, ok := actual.(*RACESMAR33); ok {
		return f.convertTo33().Compare(actual)
	}
	return f.BaseMessage.Compare(actual)
}

func (f *RACESMAR24) convertTo33() (c *RACESMAR33) {
	c = create33().(*RACESMAR33)
	c.CopyHeaderFields(&f.BaseForm)
	c.AgencyName = f.AgencyName
	c.EventName = f.EventName
	c.EventNumber = f.EventNumber
	c.Assignment = f.Assignment
	for i := 0; i < 5; i++ {
		c.Resources[i] = f.Resources[i].convertTo33()
	}
	c.RequestedArrivalDate = f.RequestedArrivalDates
	c.RequestedArrivalTime = f.RequestedArrivalTimes
	c.OpEndDate = f.NeededUntilDates
	c.OpEndTime = f.NeededUntilTimes
	c.ReportingLocation = f.ReportingLocation
	c.ContactOnArrival = f.ContactOnArrival
	c.TravelInfo = f.TravelInfo
	c.RequestedByName = f.RequestedByName
	c.RequestedByTitle = f.RequestedByTitle
	c.RequestedByContact = f.RequestedByContact
	c.ApprovedByName = f.ApprovedByName
	c.ApprovedByTitle = f.ApprovedByTitle
	c.ApprovedByContact = f.ApprovedByContact
	c.WithSignature = f.WithSignature
	c.ApprovedByDate = f.ApprovedByDate
	c.ApprovedByTime = f.ApprovedByTime
	c.OpRelayRcvd = f.OpRelayRcvd
	c.OpRelaySent = f.OpRelaySent
	c.CopyFooterFields(&f.BaseForm)
	return c
}

var typeMap = map[string]message.ChoiceMapper{
	"Field Communicator":   message.Choices{"F1", "F2", "F3", "Type IV", "Type V"},
	"Net Control Operator": message.Choices{"N1", "N2", "N3", "Type IV", "Type V"},
	"Packet Operator":      message.Choices{"P1", "P2", "P3", "Type IV", "Type V"},
	"Shadow Communicator":  message.Choices{"S1", "S2", "S3", "Type IV", "Type V"},
}

func (r *Resource24) Fields(index int) []*message.Field {
	var qtyPresence, rolePresence, posPresence, typePresence func() (message.Presence, string)
	if index == 1 {
		qtyPresence = message.Required
		rolePresence = message.Required
		typePresence = message.Required
	} else {
		rolePresence = r.requiredIfQtyElseNotAllowed
		posPresence = r.notAllowedWithoutQty
		typePresence = r.requiredIfQtyElseNotAllowed
	}
	return []*message.Field{
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Quantity", index),
			Value:       &r.Qty,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("18.%da.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 139, Y: 332.6 + 18.6*float64(index), R: 150, H: 11, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			EditWidth:   2,
			EditHelp:    `This is the number of people needed for the role and position requested on this row.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Role", index),
			Value:       &r.Role,
			Presence:    rolePresence,
			PIFOTag:     fmt.Sprintf("18.%de.", index),
			Choices:     message.Choices{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
			PDFRenderer: &message.PDFTextRenderer{X: 170, Y: 332.6 + 18.6*float64(index), R: 263, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
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
			PDFRenderer: &message.PDFTextRenderer{X: 279, Y: 332.6 + 18.6*float64(index), R: 418, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   31,
			EditHelp:    `This is the position to be held by the people requested on this row.`,
			EditApply: func(_ *message.Field, s string) {
				r.Position = strings.TrimSpace(s)
				if r.Position != "" {
					r.RolePos = message.SmartJoin(r.Role, "/ "+r.Position, " ")
				}
			},
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
			Choices:     r,
			PIFOTag:     fmt.Sprintf("18.%dc.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 431, Y: 332.6 + 18.6*float64(index), R: 489, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the preferred resource type (credential) for the people requested on this row.  It is required when there is a quantity on the row.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Minimum Type", index),
			Value:       &r.MinimumType,
			Presence:    typePresence,
			Choices:     r,
			PIFOTag:     fmt.Sprintf("18.%dd.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 509, Y: 332.6 + 18.6*float64(index), R: 567, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the minimum resource type (credential) for the people requested on this row.  It is required when there is a quantity on the row.`,
		}),
	}
}

func (r *Resource24) requiredIfQtyElseNotAllowed() (message.Presence, string) {
	if r.Qty != "" {
		return message.PresenceRequired, "there is a quantity for the resource"
	} else {
		return message.PresenceNotAllowed, "there is no quantity for the resource"
	}
}
func (r *Resource24) notAllowedWithoutQty() (message.Presence, string) {
	if r.Qty == "" {
		return message.PresenceNotAllowed, "there is no quantity for the resource"
	}
	return message.PresenceOptional, ""
}

// Implement ChoiceMapper for Resource, providing the choices for the Preferred
// Type and Minimum Type fields based on the value of the Role field.

func (r *Resource24) IsHuman(s string) bool {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.IsHuman(s)
	}
	return false
}
func (r *Resource24) IsPIFO(s string) bool {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.IsPIFO(s)
	}
	return false
}
func (r *Resource24) ToHuman(s string) string {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.ToHuman(s)
	}
	return s
}
func (r *Resource24) ToPIFO(s string) string {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.ToPIFO(s)
	}
	return s
}
func (r *Resource24) ListHuman() []string {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.ListHuman()
	}
	return nil
}

func (r Resource24) convertTo33() (c Resource33) {
	c.Qty = r.Qty
	c.Role = r.Role
	c.Position = r.Position
	if tm, ok := typeMap[r.Role]; ok {
		if idx := slices.Index(tm.(message.Choices), r.PreferredType); idx >= 0 {
			c.PreferredType = resourceTypes[idx]
		} else {
			c.PreferredType = r.PreferredType
		}
		if idx := slices.Index(tm.(message.Choices), r.MinimumType); idx >= 0 {
			c.MinimumType = resourceTypes[idx]
		} else {
			c.MinimumType = r.MinimumType
		}
	} else {
		c.PreferredType = r.PreferredType
		c.MinimumType = r.MinimumType
	}
	return c
}
