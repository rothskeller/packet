// Package eoc213rr defines the Santa Clara County EOC-213RR Resource Request
// Form message type.
package eoc213rr

import (
	"fmt"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for an EOC-213RR resource request form.
var Type = message.Type{
	Tag:         "EOC213RR",
	Name:        "EOC-213RR resource request form",
	Article:     "an",
	PDFBase:     pdfBase,
	PDFFontSize: 12,
}

func init() {
	Type.Create = New
	Type.Decode = decode
}

// versions is the list of supported versions.  The first one is used when
// creating new forms.
var versions = []*message.FormVersion{
	{HTML: "form-scco-eoc-213rr.html", Version: "2.4", Tag: "EOC213RR", FieldOrder: fieldOrder},
	{HTML: "form-scco-eoc-213rr.html", Version: "2.3", Tag: "EOC213RR", FieldOrder: fieldOrder},
	{HTML: "form-scco-eoc-213rr.html", Version: "2.1", Tag: "EOC213RR", FieldOrder: fieldOrder},
	{HTML: "form-scco-eoc-213rr.html", Version: "2.0", Tag: "EOC213RR", FieldOrder: fieldOrder},
}
var fieldOrder = []string{
	"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.", "7d.", "8d.", "21.", "22.", "23.", "24.", "25.",
	"26.", "27.", "27s.", "28.", "29.", "30.", "31.", "32.", "33.", "34.", "35.", "36a.", "36b.", "36c.", "36d.", "36e.",
	"36f.", "36g.", "36h.", "36i.", "37.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
}

// EOC213RR holds an EOC-213RR resource request form.
type EOC213RR struct {
	message.BaseMessage
	baseform.BaseForm
	IncidentName        string
	DateInitiated       string
	TimeInitiated       string
	TrackingNumber      string
	RequestedBy         string
	PreparedBy          string
	ApprovedBy          string
	WithSignature       string // added in v2.4
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

func New() (f *EOC213RR) {
	f = create(versions[0]).(*EOC213RR)
	f.MessageDate = time.Now().Format("01/02/2006")
	f.ToICSPosition = "Planning Section"
	f.ToLocation = "County EOC"
	f.DateInitiated = f.MessageDate
	return f
}

var pdfBase []byte

func create(version *message.FormVersion) message.Message {
	const fieldCount = 51
	var f = EOC213RR{BaseMessage: message.BaseMessage{
		Type: &Type,
		Form: version,
	}}
	f.BaseMessage.FSubject = &f.IncidentName
	f.BaseMessage.FBody = &f.Instructions
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &baseform.DefaultPDFMaps)
	f.Fields = append(f.Fields,
		message.NewStaticPDFContentField(&message.Field{
			PDFMap: message.PDFMapFunc(func(*message.Field) []message.PDFField {
				return []message.PDFField{{Name: "Form Type", Value: "EOC-213RR Resource Request"}}
			}),
		}),
		message.NewTextField(&message.Field{
			Label:    "Incident Name",
			Value:    &f.IncidentName,
			Presence: message.Required,
			PIFOTag:  "21.",
			PDFMap: message.PDFMapFunc(func(f *message.Field) []message.PDFField {
				return []message.PDFField{
					{Name: "Form Topic", Value: *f.Value},
					{Name: "1 Incident Name", Value: *f.Value},
				}
			}),
			EditWidth: 39,
			EditHelp:  `This is the name of the incident for which resources are being requested.  It is required.`,
		}),
		message.NewDateWithTimeField(&message.Field{
			Label:    "Date Initiated",
			Value:    &f.DateInitiated,
			Presence: message.Required,
			PIFOTag:  "22.",
			PDFMap:   message.PDFName("2 Date Initiated"),
		}),
		message.NewTimeWithDateField(&message.Field{
			Label:    "Time Initiated",
			Value:    &f.TimeInitiated,
			Presence: message.Required,
			PIFOTag:  "23.",
			PDFMap:   message.PDFName("3 Time Initiated"),
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
			Label:     "Requested By",
			Value:     &f.RequestedBy,
			Presence:  message.Required,
			PIFOTag:   "25.",
			PDFMap:    message.PDFName("5 Requested By"),
			EditWidth: 37,
			EditHelp:  `This is the name, agency, position, email, and phone number of the person requesting resources.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Prepared By",
			Value:     &f.PreparedBy,
			PIFOTag:   "26.",
			PDFMap:    message.PDFName("6 Prepared by"),
			EditWidth: 37,
			EditHelp:  `This is the name, position, email, and phone number of the person who prepared the resource request form.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Approved By",
			Value:   &f.ApprovedBy,
			PIFOTag: "27.",
			PDFMap: message.PDFMapFunc(func(*message.Field) []message.PDFField {
				if f.WithSignature != "" {
					return []message.PDFField{{
						Name:  "7 Approved By",
						Value: message.SmartJoin(f.ApprovedBy, "[with signature]", "\n"),
					}}
				}
				return []message.PDFField{{Name: "7 Approved By", Value: f.ApprovedBy}}
			}),
			TableValue: func(*message.Field) string {
				if f.WithSignature != "" {
					return message.SmartJoin(f.ApprovedBy, "[with signature]", "\n")
				}
				return f.ApprovedBy
			},
			EditWidth: 37,
			EditHelp:  `This is the name, position, email, and phone number of the person who approved the resource request.`,
		}),
	)
	if f.Form.Version >= "2.4" {
		f.Fields = append(f.Fields,
			message.NewRestrictedField(&message.Field{
				Label:      "With Signature",
				Value:      &f.WithSignature,
				Choices:    message.Choices{"checked"},
				PIFOTag:    "27s.",
				TableValue: message.TableOmit,
				EditHelp:   `This indicates whether the original paper resource request form was signed.`,
			}),
		)
	}
	f.Fields = append(f.Fields,
		message.NewMultilineField(&message.Field{
			Label:     "Qty/Unit",
			Value:     &f.QtyUnit,
			Presence:  message.Required,
			PIFOTag:   "28.",
			Compare:   message.CompareExact,
			PDFMap:    message.PDFName("8 QtyUnit"),
			EditWidth: 9,
			EditHelp:  `This is the quantity (with units where applicable) of the resource requested.  If multiple resources are being requested, enter the quantity of each on a separate line.  This field is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Resource Description",
			Value:     &f.ResourceDescription,
			Presence:  message.Required,
			PIFOTag:   "29.",
			PDFMap:    message.PDFName("9 Resource Description"),
			EditWidth: 34,
			EditHelp:  `This is the description of the resource requested.  If multiple resources are being requested, enter the description of each on a separate line, parallel to the corresponding quantity.  This field is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Resource Arrival",
			Value:     &f.ResourceArrival,
			Presence:  message.Required,
			PIFOTag:   "30.",
			PDFMap:    message.PDFName("10 Arrival"),
			EditWidth: 18,
			EditHelp:  `This is the date and time by which the resource needs to have arrived.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Priority",
			Value:    &f.Priority,
			Choices:  message.Choices{"Now", "High", "Medium", "Low"},
			Presence: message.Required,
			PIFOTag:  "31.",
			PDFMap:   message.PDFNameMap{"11 Priority", "", "Off"},
			EditHelp: `This is the priority of the resource request.  It must have the value "Now", "High" (meaning within the next 4 hours), "Medium" (meaning between 5 and 12 hours), or "Low" (meaning more than 12 hours).  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Estimated Cost",
			Value:     &f.EstdCost,
			PIFOTag:   "32.",
			PDFMap:    message.PDFName("12 Estd Cost"),
			EditWidth: 11,
			EditHelp:  `This is the estimated cost of the resources requested.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Deliver To",
			Value:     &f.DeliverTo,
			Presence:  message.Required,
			PIFOTag:   "33.",
			PDFMap:    message.PDFName("13 Deliver To"),
			EditWidth: 44,
			EditHelp:  `This is the name, agency, position, email, and phone number of the person to whom the requested resources should be delivered.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Deliver To Location",
			Value:     &f.DeliverToLocation,
			Presence:  message.Required,
			PIFOTag:   "34.",
			PDFMap:    message.PDFName("14 Location"),
			EditWidth: 42,
			EditHelp:  `This is the address and/or GPS coordinates of the location to which the requested resources should be delivered.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Substitutes/Sources",
			Value:     &f.Substitutes,
			PIFOTag:   "35.",
			PDFMap:    message.PDFName("15 Sub Sugg Sources"),
			EditWidth: 86,
			EditHelp:  `This is the names, phone numbers, and/or websites of suggested or substitute sources for the requested resources.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Equipment Operator",
			Value:      &f.EquipmentOperator,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36a.",
			PDFMap:     message.PDFNameMap{"Equip Oper", "", "Off", "false", "Off", "checked", "Yes"},
			TableValue: message.TableOmit,
			EditHelp:   `This indicates a supplemental requirement for an equipment operator.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Lodging",
			Value:      &f.Lodging,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36b.",
			PDFMap:     message.PDFNameMap{"Lodging", "", "Off", "false", "Off", "checked", "Yes"},
			TableValue: message.TableOmit,
			EditHelp:   `This indicates a supplemental requirement for lodging.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Fuel",
			Value:      &f.Fuel,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36c.",
			PDFMap:     message.PDFNameMap{"Fuel", "", "Off", "false", "Off", "checked", "Yes"},
			TableValue: message.TableOmit,
			EditHelp:   `This indicates a supplemental requirement for fuel.  The fuel type must be specified in the "Fuel Type" field.`,
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
			PDFMap:     message.PDFName("Fuel Type"),
			TableValue: message.TableOmit,
			EditWidth:  11,
			EditHelp:   `This is the type of fuel required.  It must be and can be set only when "Supplemental: Fuel" is checked.`,
			EditSkip:   func(*message.Field) bool { return f.Fuel == "" },
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Power",
			Value:      &f.Power,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36e.",
			PDFMap:     message.PDFNameMap{"Power", "", "Off", "false", "Off", "checked", "Yes"},
			TableValue: message.TableOmit,
			EditHelp:   `This indicates a supplemental requirement for power.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Meals",
			Value:      &f.Meals,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36f.",
			PDFMap:     message.PDFNameMap{"Meals", "", "Off", "false", "Off", "checked", "Yes"},
			TableValue: message.TableOmit,
			EditHelp:   `This indicates a supplemental requirement for meals.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Maintenance",
			Value:      &f.Maintenance,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36g.",
			PDFMap:     message.PDFNameMap{"Maintenance", "", "Off", "false", "Off", "checked", "Yes"},
			TableValue: message.TableOmit,
			EditHelp:   `"This indicates a supplemental requirement for maintenance.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Water",
			Value:      &f.Water,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36h.",
			PDFMap:     message.PDFNameMap{"Water", "", "Off", "false", "Off", "checked", "Yes"},
			TableValue: message.TableOmit,
			EditHelp:   `This indicates a supplemental requirement for water.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Supplemental: Other",
			Value:      &f.Other,
			Choices:    message.Choices{"checked"},
			PIFOTag:    "36i.",
			PDFMap:     message.PDFNameMap{"Other", "", "Off", "false", "Off", "checked", "Yes"},
			TableValue: message.TableOmit,
			EditHelp:   `This indicates an additional supplemental requirement (described in the "Special Instructions" field).`,
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
			Label:     "Special Instructions",
			Value:     &f.Instructions,
			PIFOTag:   "37.",
			PDFMap:    message.PDFName("17 Special Instructions"),
			EditWidth: 42,
			EditHelp:  `This is any special requirements or instructions for the resource request.`,
		}),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &baseform.DefaultPDFMaps)
	if len(f.Fields) > fieldCount {
		panic("update EOC213RR fieldCount")
	}
	return &f
}

func decode(subject, body string) (f *EOC213RR) {
	// Quick check to avoid overhead of creating the form object if it's not
	// our type of form.
	if !strings.Contains(body, "form-scco-eoc-213rr.html") {
		return nil
	}
	return message.DecodeForm(body, versions, create).(*EOC213RR)
}
