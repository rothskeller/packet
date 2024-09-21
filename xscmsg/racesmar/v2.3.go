// Package racesmar defines the RACES Mutual Aid Request Form message type.
package racesmar

import (
	"fmt"
	"slices"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type23 is the type definition for a RACES mutual aid request form.
var Type23 = message.Type{
	Tag:     "RACES-MAR",
	HTML:    "form-oa-mutual-aid-request-v2.html",
	Version: "2.3",
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
		"24c.", "25a.", "25b.", "25c.", "26a.", "26b.", "OpRelayRcvd",
		"OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type23, decode23, nil)
}

// RACESMAR23 holds a RACES mutual aid request form.
type RACESMAR23 struct {
	message.BaseMessage
	baseform.BaseForm
	AgencyName            string
	EventName             string
	EventNumber           string
	Assignment            string
	Resources             [5]Resource23
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
	ApprovedByDate        string
	ApprovedByTime        string
}

// A Resource23 is the description of a single resource in a RACES mutual aid
// request form.
type Resource23 struct {
	Qty           string
	Role          string
	Position      string
	RolePos       string
	PreferredType string
	MinimumType   string
}

func make23() *RACESMAR23 {
	const fieldCount = 73
	var f = RACESMAR23{BaseMessage: message.BaseMessage{Type: &Type23}}
	f.BaseMessage.FSubject = &f.AgencyName
	f.BaseMessage.FBody = &f.Assignment
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, nil)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:    "Agency Name",
			Value:    &f.AgencyName,
			Presence: message.Required,
			PIFOTag:  "15.",
		}),
		message.NewTextField(&message.Field{
			Label:      "Event Name",
			Value:      &f.EventName,
			Presence:   message.Required,
			PIFOTag:    "16a.",
			TableValue: message.TableOmit,
		}),
		message.NewTextField(&message.Field{
			Label:      "Event Number",
			Value:      &f.EventNumber,
			PIFOTag:    "16b.",
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Event Name/Number",
			TableValue: func(*message.Field) string {
				return message.SmartJoin(f.EventName, f.EventNumber, " ")
			},
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Assignment",
			Value:    &f.Assignment,
			Presence: message.Required,
			PIFOTag:  "17.",
		}),
	)
	for i := range f.Resources {
		f.Fields = append(f.Fields, f.Resources[i].Fields(i+1)...)
	}
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:    "Requested Arrival Date",
			Value:    &f.RequestedArrivalDates,
			Presence: message.Required,
			PIFOTag:  "19a.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Requested Arrival Time",
			Value:    &f.RequestedArrivalTimes,
			Presence: message.Required,
			PIFOTag:  "19b.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Needed Until Dates",
			Value:    &f.NeededUntilDates,
			Presence: message.Required,
			PIFOTag:  "20a.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Needed Until Times",
			Value:    &f.NeededUntilTimes,
			Presence: message.Required,
			PIFOTag:  "20b.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Reporting Location",
			Value:    &f.ReportingLocation,
			Presence: message.Required,
			PIFOTag:  "21.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Contact on Arrival",
			Value:    &f.ContactOnArrival,
			Presence: message.Required,
			PIFOTag:  "22.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Travel Info",
			Value:    &f.TravelInfo,
			Presence: message.Required,
			PIFOTag:  "23.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Requested By Name",
			Value:    &f.RequestedByName,
			Presence: message.Required,
			PIFOTag:  "24a.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Requested By Title",
			Value:    &f.RequestedByTitle,
			Presence: message.Required,
			PIFOTag:  "24b.",
		}),
		message.NewTextField(&message.Field{
			Label:   "Requested By Contact",
			Value:   &f.RequestedByContact,
			PIFOTag: "24c.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Approved By Name",
			Value:    &f.ApprovedByName,
			Presence: message.Required,
			PIFOTag:  "25a.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Approved By Title",
			Value:    &f.ApprovedByTitle,
			Presence: message.Required,
			PIFOTag:  "25b.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Approved By Contact",
			Value:    &f.ApprovedByContact,
			Presence: message.Required,
			PIFOTag:  "25c.",
		}),
		message.NewDateField(true, &message.Field{
			Label:    "Approved By Date",
			Value:    &f.ApprovedByDate,
			Presence: message.Required,
			PIFOTag:  "26a.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:    "Approved By Time",
			Value:    &f.ApprovedByTime,
			Presence: message.Required,
			PIFOTag:  "26b.",
		}),
		message.NewDateTimeField(&message.Field{
			Label: "Approved Date/Time",
		}, &f.ApprovedByDate, &f.ApprovedByTime),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, nil)
	if len(f.Fields) > fieldCount {
		panic("update RACESMAR23 fieldCount")
	}
	return &f
}

func decode23(_, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *RACESMAR23

	if form == nil || form.HTMLIdent != Type23.HTML || form.FormVersion != Type23.Version {
		return nil
	}
	df = make23()
	message.DecodeForm(form, df)
	return df
}

func (f *RACESMAR23) Compare(actual message.Message) (int, int, []*message.CompareField) {
	return f.convertTo33().Compare(actual)
}

func (f *RACESMAR23) RenderPDF(env *envelope.Envelope, filename string) error {
	return f.convertTo33().RenderPDF(env, filename)
}

func (f *RACESMAR23) convertTo33() (c *RACESMAR33) {
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
	c.ApprovedByDate = f.ApprovedByDate
	c.ApprovedByTime = f.ApprovedByTime
	c.OpRelayRcvd = f.OpRelayRcvd
	c.OpRelaySent = f.OpRelaySent
	c.CopyFooterFields(&f.BaseForm)
	return c
}

func (r *Resource23) Fields(index int) []*message.Field {
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
			Label:    fmt.Sprintf("Resource %d Quantity", index),
			Value:    &r.Qty,
			Presence: qtyPresence,
			PIFOTag:  fmt.Sprintf("18.%da.", index),
		}),
		message.NewRestrictedField(&message.Field{
			Label:    fmt.Sprintf("Resource %d Role", index),
			Value:    &r.Role,
			Presence: rolePresence,
			PIFOTag:  fmt.Sprintf("18.%de.", index),
			Choices:  message.Choices{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
		}),
		message.NewTextField(&message.Field{
			Label:    fmt.Sprintf("Resource %d Position", index),
			Value:    &r.Position,
			Presence: posPresence,
			PIFOTag:  fmt.Sprintf("18.%df.", index),
		}),
		message.NewCalculatedField(&message.Field{
			Label:   fmt.Sprintf("Resource %d Role/Position", index),
			Value:   &r.RolePos,
			PIFOTag: fmt.Sprintf("18.%db.", index),
		}),
		message.NewRestrictedField(&message.Field{
			Label:    fmt.Sprintf("Resource %d Preferred Type", index),
			Value:    &r.PreferredType,
			Presence: typePresence,
			Choices:  r,
			PIFOTag:  fmt.Sprintf("18.%dc.", index),
		}),
		message.NewRestrictedField(&message.Field{
			Label:    fmt.Sprintf("Resource %d Minimum Type", index),
			Value:    &r.MinimumType,
			Presence: typePresence,
			Choices:  r,
			PIFOTag:  fmt.Sprintf("18.%dd.", index),
		}),
	}
}
func (r *Resource23) requiredIfQtyElseNotAllowed() (message.Presence, string) {
	if r.Qty != "" {
		return message.PresenceRequired, "there is a quantity for the resource"
	} else {
		return message.PresenceNotAllowed, "there is no quantity for the resource"
	}
}
func (r *Resource23) notAllowedWithoutQty() (message.Presence, string) {
	if r.Qty == "" {
		return message.PresenceNotAllowed, "there is no quantity for the resource"
	}
	return message.PresenceOptional, ""
}

// Implement ChoiceMapper for Resource, providing the choices for the Preferred
// Type and Minimum Type fields based on the value of the Role field.

func (r *Resource23) IsHuman(s string) bool {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.IsHuman(s)
	}
	return false
}
func (r *Resource23) IsPIFO(s string) bool {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.IsPIFO(s)
	}
	return false
}
func (r *Resource23) ToHuman(s string) string {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.ToHuman(s)
	}
	return s
}
func (r *Resource23) ToPIFO(s string) string {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.ToPIFO(s)
	}
	return s
}
func (r *Resource23) ListHuman() []string {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.ListHuman()
	}
	return nil
}

func (r Resource23) convertTo33() (c Resource33) {
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
