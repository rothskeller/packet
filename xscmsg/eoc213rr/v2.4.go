// Package eoc213rr defines the Santa Clara County EOC-213RR Resource Request
// Form message type.
package eoc213rr

import (
	"fmt"
	"strings"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type24 is the type definition for an EOC-213RR resource request form.
var Type24 = message.Type{
	Tag:     "EOC213RR",
	HTML:    "form-scco-eoc-213rr.html",
	Version: "2.4",
	Name:    "EOC-213RR resource request form",
	Article: "an",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.", "7d.", "8d.", "21.", "22.", "23.", "24.", "25.",
		"26.", "27.", "27s.", "28.", "29.", "30.", "31.", "32.", "33.", "34.", "35.", "36a.", "36b.", "36c.", "36d.", "36e.",
		"36f.", "36g.", "36h.", "36i.", "37.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type24, decode24, create24)
}

// EOC213RR24 holds an EOC-213RR resource request form.
type EOC213RR24 struct {
	message.BaseMessage
	baseform.BaseForm
	IncidentName        string
	DateInitiated       string
	TimeInitiated       string
	TrackingNumber      string
	RequestedBy         string
	PreparedBy          string
	ApprovedBy          string
	WithSignature       string
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

func create24() message.Message {
	f := make24()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.ToICSPosition = "Planning Section"
	f.ToLocation = "County EOC"
	f.DateInitiated = f.MessageDate
	return f
}

func make24() (f *EOC213RR24) {
	const fieldCount = 51
	f = &EOC213RR24{BaseMessage: message.BaseMessage{Type: &Type24}}
	f.FSubject = &f.IncidentName
	f.FBody = &f.Instructions
	f.Fields = make([]*message.Field, 0, fieldCount)
	pdf := baseform.RoutingSlipPDFRenderers
	pdf.OriginMsgID = message.PDFMultiRenderer{
		pdf.OriginMsgID,
		&message.PDFTextRenderer{Page: 2, X: 486, Y: 20, R: 586, H: 12, Style: message.PDFTextStyle{HAlign: "right"}},
		&message.PDFTextRenderer{Page: 3, X: 486, Y: 20, R: 586, H: 12, Style: message.PDFTextStyle{HAlign: "right"}},
	}
	f.AddHeaderFields(&f.BaseMessage, &pdf)
	f.Fields = append(f.Fields,
		message.NewStaticPDFContentField(&message.Field{
			PDFRenderer: &message.PDFStaticTextRenderer{
				Page: 1, X: 119, Y: 225, H: 17,
				Text: "EOC-213RR Resource Request",
			},
		}),
		message.NewTextField(&message.Field{
			Label:    "Incident Name",
			Value:    &f.IncidentName,
			Presence: message.Required,
			PIFOTag:  "21.",
			PDFRenderer: message.PDFMultiRenderer{
				&message.PDFTextRenderer{Page: 1, X: 351, Y: 225, R: 573, B: 242, Style: message.PDFTextStyle{VAlign: "baseline"}},
				&message.PDFTextRenderer{Page: 2, X: 26, Y: 130, R: 247, B: 160, Style: message.PDFTextStyle{VAlign: "baseline"}},
			},
			EditWidth: 39,
			EditHelp:  `This is the name of the incident for which resources are being requested.  It is required.`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Date Initiated",
			Value:       &f.DateInitiated,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 255, Y: 130, R: 355, B: 160, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date on which the resource request was initiated, in MM/DD/YYYY format.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time Initiated",
			Value:       &f.TimeInitiated,
			Presence:    message.Required,
			PIFOTag:     "23.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 364, Y: 130, R: 443, B: 160, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time at which the resource request was initiated, in HH:MM format (24-hour clock).  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time Initiated",
			Presence: message.Required,
			EditHelp: `This is the date and time at which the resource request was initiated, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.DateInitiated, &f.TimeInitiated),
		message.NewTextField(&message.Field{
			Label:   "Tracking Number",
			Value:   &f.TrackingNumber,
			PIFOTag: "24.",
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Requested By",
			Value:       &f.RequestedBy,
			Presence:    message.Required,
			PIFOTag:     "25.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 26, Y: 179, R: 247, B: 300, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   37,
			EditHelp:    `This is the name, agency, position, email, and phone number of the person requesting resources.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Prepared By",
			Value:       &f.PreparedBy,
			PIFOTag:     "26.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 26, Y: 319, R: 247, B: 358, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   37,
			EditHelp:    `This is the name, position, email, and phone number of the person who prepared the resource request form.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Approved By",
			Value:   &f.ApprovedBy,
			PIFOTag: "27.",
			TableValue: func(*message.Field) string {
				if f.WithSignature != "" {
					return message.SmartJoin(f.ApprovedBy, "[with signature]", "\n")
				}
				return f.ApprovedBy
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 26, Y: 378, R: 247, B: 405, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   37,
			EditHelp:    `This is the name, position, email, and phone number of the person who approved the resource request.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "With Signature",
			Value:       &f.WithSignature,
			Choices:     message.Choices{"checked"},
			PIFOTag:     "27s.",
			TableValue:  message.TableOmit,
			PDFRenderer: &message.PDFMappedTextRenderer{Page: 2, X: 84, Y: 408, B: 418, Map: map[string]string{"checked": "[with signature]"}},
			EditHelp:    `This indicates whether the original paper resource request form was signed.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Qty/Unit",
			Value:       &f.QtyUnit,
			Presence:    message.Required,
			PIFOTag:     "28.",
			Compare:     message.CompareExact,
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 73, Y: 460, R: 126, B: 544, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   9,
			EditHelp:    `This is the quantity (with units where applicable) of the resource requested.  If multiple resources are being requested, enter the quantity of each on a separate line.  This field is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Resource Description",
			Value:       &f.ResourceDescription,
			Presence:    message.Required,
			PIFOTag:     "29.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 134, Y: 460, R: 332, B: 544, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   34,
			EditHelp:    `This is the description of the resource requested.  If multiple resources are being requested, enter the description of each on a separate line, parallel to the corresponding quantity.  This field is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Resource Arrival",
			Value:       &f.ResourceArrival,
			Presence:    message.Required,
			PIFOTag:     "30.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 341, Y: 460, R: 446, B: 544, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   18,
			EditHelp:    `This is the date and time by which the resource needs to have arrived.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Priority",
			Value:    &f.Priority,
			Choices:  message.Choices{"Now", "High", "Medium", "Low"},
			Presence: message.Required,
			PIFOTag:  "31.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 2, Points: map[string][]float64{
				"Now":    {508.5, 464},
				"High":   {508.5, 480},
				"Medium": {508.5, 506},
				"Low":    {508.5, 531.5},
			}},
			EditHelp: `This is the priority of the resource request.  It must have the value "Now", "High" (meaning within the next 4 hours), "Medium" (meaning between 5 and 12 hours), or "Low" (meaning more than 12 hours).  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Estimated Cost",
			Value:       &f.EstdCost,
			PIFOTag:     "32.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 525, Y: 460, R: 588, B: 544, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   11,
			EditHelp:    `This is the estimated cost of the resources requested.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Deliver To",
			Value:       &f.DeliverTo,
			Presence:    message.Required,
			PIFOTag:     "33.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 73, Y: 559, R: 332, B: 585, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   44,
			EditHelp:    `This is the name, agency, position, email, and phone number of the person to whom the requested resources should be delivered.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Deliver To Location",
			Value:       &f.DeliverToLocation,
			Presence:    message.Required,
			PIFOTag:     "34.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 341, Y: 559, R: 588, B: 585, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   42,
			EditHelp:    `This is the address and/or GPS coordinates of the location to which the requested resources should be delivered.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Substitutes/Sources",
			Value:       &f.Substitutes,
			PIFOTag:     "35.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 73, Y: 600, R: 588, B: 626, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   86,
			EditHelp:    `This is the names, phone numbers, and/or websites of suggested or substitute sources for the requested resources.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Supplemental: Equipment Operator",
			Value:       &f.EquipmentOperator,
			Choices:     message.Choices{"checked"},
			PIFOTag:     "36a.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 2, Points: map[string][]float64{"checked": {79, 661}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates a supplemental requirement for an equipment operator.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Supplemental: Lodging",
			Value:       &f.Lodging,
			Choices:     message.Choices{"checked"},
			PIFOTag:     "36b.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 2, Points: map[string][]float64{"checked": {222.5, 661}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates a supplemental requirement for lodging.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Supplemental: Fuel",
			Value:       &f.Fuel,
			Choices:     message.Choices{"checked"},
			PIFOTag:     "36c.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 2, Points: map[string][]float64{"checked": {79, 678}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates a supplemental requirement for fuel.  The fuel type must be specified in the "Fuel Type" field.`,
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
			PIFOTag:     "36d.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 132, Y: 684, R: 210, B: 694, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditWidth:   11,
			EditHelp:    `This is the type of fuel required.  It must be and can be set only when "Supplemental: Fuel" is checked.`,
			EditSkip:    func(*message.Field) bool { return f.Fuel == "" },
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Supplemental: Power",
			Value:       &f.Power,
			Choices:     message.Choices{"checked"},
			PIFOTag:     "36e.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 2, Points: map[string][]float64{"checked": {222.5, 678}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates a supplemental requirement for power.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Supplemental: Meals",
			Value:       &f.Meals,
			Choices:     message.Choices{"checked"},
			PIFOTag:     "36f.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 2, Points: map[string][]float64{"checked": {79, 706.5}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates a supplemental requirement for meals.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Supplemental: Maintenance",
			Value:       &f.Maintenance,
			Choices:     message.Choices{"checked"},
			PIFOTag:     "36g.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 2, Points: map[string][]float64{"checked": {222.5, 695}}},
			TableValue:  message.TableOmit,
			EditHelp:    `"This indicates a supplemental requirement for maintenance.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Supplemental: Water",
			Value:       &f.Water,
			Choices:     message.Choices{"checked"},
			PIFOTag:     "36h.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 2, Points: map[string][]float64{"checked": {79, 724}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates a supplemental requirement for water.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Supplemental: Other",
			Value:       &f.Other,
			Choices:     message.Choices{"checked"},
			PIFOTag:     "36i.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 2, Points: map[string][]float64{"checked": {222.5, 713}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates an additional supplemental requirement (described in the "Special Instructions" field).`,
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
			Label:       "Special Instructions",
			Value:       &f.Instructions,
			PIFOTag:     "37.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 341, Y: 644, R: 588, B: 743, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   42,
			EditHelp:    `This is any special requirements or instructions for the resource request.`,
		}),
	)
	f.AddFooterFields(&f.BaseMessage, &pdf)
	if len(f.Fields) > fieldCount {
		panic("update EOC213RR24 fieldCount")
	}
	return f
}

func decode24(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type24.HTML || form.FormVersion != Type24.Version {
		return nil
	}
	df := make24()
	message.DecodeForm(form, df)
	return df
}

func (f *EOC213RR24) Compare(actual message.Message) (int, int, []*message.CompareField) {
	switch act := actual.(type) {
	case *EOC213RR23:
		actual = act.convertTo24()
	}
	return f.BaseMessage.Compare(actual)
}
