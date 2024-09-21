// Package eoc213rr defines the Santa Clara County EOC-213RR Resource Request
// Form message type.
package eoc213rr

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type23 is the type definition for an EOC-213RR resource request form.
var Type23 = message.Type{
	Tag:     "EOC213RR",
	HTML:    "form-scco-eoc-213rr.html",
	Version: "2.3",
	Name:    "EOC-213RR resource request form",
	Article: "an",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.", "7d.", "8d.", "21.", "22.", "23.", "24.", "25.",
		"26.", "27.", "28.", "29.", "30.", "31.", "32.", "33.", "34.", "35.", "36a.", "36b.", "36c.", "36d.", "36e.",
		"36f.", "36g.", "36h.", "36i.", "37.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type23, decode23, nil)
}

// EOC213RR23 holds an EOC-213RR resource request form.
type EOC213RR23 struct {
	message.BaseMessage
	baseform.BaseForm
	IncidentName        string
	DateInitiated       string
	TimeInitiated       string
	TrackingNumber      string
	RequestedBy         string
	PreparedBy          string
	ApprovedBy          string
	QtyUnit             string
	ResourceDescription string
	ResourceArrival     string
	Priority            string
	EstdCost            string
	DeliverTo           string
	DeliverToLocation   string
	Substitutes         string
	EquipmentOperator   string
	Lodging             string
	Fuel                string
	FuelType            string
	Power               string
	Meals               string
	Maintenance         string
	Water               string
	Other               string
	Instructions        string
}

func make23() (f *EOC213RR23) {
	const fieldCount = 49
	f = &EOC213RR23{BaseMessage: message.BaseMessage{Type: &Type23}}
	f.BaseMessage.FSubject = &f.IncidentName
	f.BaseMessage.FBody = &f.Instructions
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, nil)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:    "Incident Name",
			Value:    &f.IncidentName,
			Presence: message.Required,
			PIFOTag:  "21.",
		}),
		message.NewDateField(true, &message.Field{
			Label:    "Date Initiated",
			Value:    &f.DateInitiated,
			Presence: message.Required,
			PIFOTag:  "22.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:    "Time Initiated",
			Value:    &f.TimeInitiated,
			Presence: message.Required,
			PIFOTag:  "23.",
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time Initiated",
			Presence: message.Required,
		}, &f.DateInitiated, &f.TimeInitiated),
		message.NewTextField(&message.Field{
			Label:   "Tracking Number",
			Value:   &f.TrackingNumber,
			PIFOTag: "24.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Requested By",
			Value:    &f.RequestedBy,
			Presence: message.Required,
			PIFOTag:  "25.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Prepared By",
			Value:   &f.PreparedBy,
			PIFOTag: "26.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Approved By",
			Value:   &f.ApprovedBy,
			PIFOTag: "27.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Qty/Unit",
			Value:    &f.QtyUnit,
			Presence: message.Required,
			PIFOTag:  "28.",
			Compare:  message.CompareExact,
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Resource Description",
			Value:    &f.ResourceDescription,
			Presence: message.Required,
			PIFOTag:  "29.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Resource Arrival",
			Value:    &f.ResourceArrival,
			Presence: message.Required,
			PIFOTag:  "30.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Priority",
			Value:    &f.Priority,
			Choices:  message.Choices{"Now", "High", "Medium", "Low"},
			Presence: message.Required,
			PIFOTag:  "31.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Estimated Cost",
			Value:   &f.EstdCost,
			PIFOTag: "32.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Deliver To",
			Value:    &f.DeliverTo,
			Presence: message.Required,
			PIFOTag:  "33.",
		}),
		message.NewMultilineField(&message.Field{
			Label:    "Deliver To Location",
			Value:    &f.DeliverToLocation,
			Presence: message.Required,
			PIFOTag:  "34.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Substitutes/Sources",
			Value:   &f.Substitutes,
			PIFOTag: "35.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Supplemental: Equipment Operator",
			Value:   &f.EquipmentOperator,
			Choices: message.Choices{"checked"},
			PIFOTag: "36a.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Supplemental: Lodging",
			Value:   &f.Lodging,
			Choices: message.Choices{"checked"},
			PIFOTag: "36b.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Fuel",
			Value:      &f.Fuel,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36c.",
			TableValue: message.TableOmit,
		}),
		message.NewTextField(&message.Field{
			Label: "Supplemental: Fuel Type",
			Value: &f.FuelType,
			Presence: func() (message.Presence, string) {
				if f.Fuel == "" {
					return message.PresenceNotAllowed, `when "Fuel" is not checked`
				} else {
					return message.PresenceRequired, `when "Fuel" is checked`
				}
			},
			PIFOTag:    "36d.",
			TableValue: message.TableOmit,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Power",
			Value:      &f.Power,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36e.",
			TableValue: message.TableOmit,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Meals",
			Value:      &f.Meals,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36f.",
			TableValue: message.TableOmit,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Maintenance",
			Value:      &f.Maintenance,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36g.",
			TableValue: message.TableOmit,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Water",
			Value:      &f.Water,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36h.",
			TableValue: message.TableOmit,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Other",
			Value:      &f.Other,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36i.",
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Supplemental Requirements",
			TableValue: func(*message.Field) string {
				var reqs []string
				if f.EquipmentOperator != "" {
					reqs = append(reqs, "Equipment Operator")
				}
				if f.Lodging != "" {
					reqs = append(reqs, "Lodging")
				}
				if f.Fuel != "" {
					if f.FuelType != "" {
						reqs = append(reqs, fmt.Sprintf("Fuel (%s)", f.FuelType))
					} else {
						reqs = append(reqs, "Fuel")
					}
				}
				if f.Power != "" {
					reqs = append(reqs, "Power")
				}
				if f.Meals != "" {
					reqs = append(reqs, "Meals")
				}
				if f.Maintenance != "" {
					reqs = append(reqs, "Maintenance")
				}
				if f.Water != "" {
					reqs = append(reqs, "Water")
				}
				if f.Other != "" {
					reqs = append(reqs, "Other")
				}
				return strings.Join(reqs, ", ")
			},
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Special Instructions",
			Value:   &f.Instructions,
			PIFOTag: "37.",
		}),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, nil)
	if len(f.Fields) > fieldCount {
		panic("update EOC213RR23 fieldCount")
	}
	return f
}

func decode23(_, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type23.HTML || form.FormVersion != Type23.Version {
		return nil
	}
	var df = make23()
	message.DecodeForm(form, df)
	return df
}

func (f *EOC213RR23) Compare(actual message.Message) (int, int, []*message.CompareField) {
	return f.convertTo24().Compare(actual)
}

func (f *EOC213RR23) RenderPDF(env *envelope.Envelope, filename string) error {
	return f.convertTo24().RenderPDF(env, filename)
}

func (f *EOC213RR23) convertTo24() (c *EOC213RR24) {
	c = make24()
	c.CopyHeaderFields(&f.BaseForm)
	c.IncidentName = f.IncidentName
	c.DateInitiated = f.DateInitiated
	c.TimeInitiated = f.TimeInitiated
	c.TrackingNumber = f.TrackingNumber
	c.RequestedBy = f.RequestedBy
	c.PreparedBy = f.PreparedBy
	c.ApprovedBy = f.ApprovedBy
	c.QtyUnit = f.QtyUnit
	c.ResourceDescription = f.ResourceDescription
	c.ResourceArrival = f.ResourceArrival
	c.Priority = f.Priority
	c.EstdCost = f.EstdCost
	c.DeliverTo = f.DeliverTo
	c.DeliverToLocation = f.DeliverToLocation
	c.Substitutes = f.Substitutes
	c.EquipmentOperator = f.EquipmentOperator
	c.Lodging = f.Lodging
	c.Fuel = f.Fuel
	c.FuelType = f.FuelType
	c.Power = f.Power
	c.Meals = f.Meals
	c.Maintenance = f.Maintenance
	c.Water = f.Water
	c.Other = f.Other
	c.Instructions = f.Instructions
	c.CopyFooterFields(&f.BaseForm)
	return c
}
