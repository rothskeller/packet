// Package racesmar defines the RACES Mutual Aid Request Form message type.
package racesmar

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type16 is the type definition for a RACES mutual aid request form.
var Type16 = message.Type{
	Tag:     "RACES-MAR",
	HTML:    "form-oa-mutual-aid-request.html",
	Version: "1.6",
	Name:    "RACES mutual aid request form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.", "7d.",
		"8d.", "15.", "16a.", "16b.", "17.", "18a.", "18b.", "18c.", "18d.",
		"18a.", "18b.", "18c.", "18d.",
		"19a.",
		"19b.", "20a.", "20b.", "21.", "22.", "23.", "24a.", "24b.",
		"24c.", "25a.", "25b.", "25c.", "26a.", "26b.", "OpRelayRcvd",
		"OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type16, decode16, nil)
}

// RACESMAR16 holds a RACES mutual aid request form.
type RACESMAR16 struct {
	message.BaseMessage
	baseform.BaseForm
	AgencyName            string
	EventName             string
	EventNumber           string
	Assignment            string
	ResourceQty           string
	ResourceRolePos       string
	ResourcePreferredType string
	ResourceMinimumType   string
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

func make16() *RACESMAR16 {
	const fieldCount = 47
	var f = RACESMAR16{BaseMessage: message.BaseMessage{Type: &Type16}}
	f.BaseMessage.FSubject = &f.AgencyName
	f.BaseMessage.FBody = &f.Assignment
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
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
		message.NewMultilineField(&message.Field{
			Label:    "Resource Quantity",
			Value:    &f.ResourceQty,
			Presence: message.Required,
			PIFOTag:  "18a.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Role/Position",
			Value:    &f.ResourceRolePos,
			Presence: message.Required,
			PIFOTag:  "18b.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Preferred Type",
			Value:    &f.ResourcePreferredType,
			Presence: message.Required,
			PIFOTag:  "18c.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Minimum Type",
			Value:    &f.ResourceMinimumType,
			Presence: message.Required,
			PIFOTag:  "18d.",
		}),
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
	f.BaseForm.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update RACESMAR16 fieldCount")
	}
	return &f
}

func decode16(_, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *RACESMAR16

	if form == nil || form.HTMLIdent != Type16.HTML || form.FormVersion != Type16.Version {
		return nil
	}
	df = make16()
	message.DecodeForm(form, df)
	return df
}

func (f *RACESMAR16) Compare(actual message.Message) (int, int, []*message.CompareField) {
	return f.convertTo33().Compare(actual)
}

func (f *RACESMAR16) RenderPDF(env *envelope.Envelope, filename string) error {
	return f.convertTo33().RenderPDF(env, filename)
}

func (f *RACESMAR16) convertTo33() (c *RACESMAR33) {
	c = create33().(*RACESMAR33)
	c.CopyHeaderFields(&f.BaseForm)
	c.AgencyName = f.AgencyName
	c.EventName = f.EventName
	c.EventNumber = f.EventNumber
	c.Assignment = f.Assignment
	c.Resources[0].Qty = f.ResourceQty
	c.Resources[0].Position = f.ResourceRolePos
	c.Resources[0].RolePos = f.ResourceRolePos
	c.Resources[0].PreferredType = f.ResourcePreferredType
	c.Resources[0].MinimumType = f.ResourceMinimumType
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
