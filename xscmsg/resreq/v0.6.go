// Package resreq defines the Resource Request Form message type.
package resreq

import (
	"fmt"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for a resource request form.
var Type = message.Type{
	Tag:     "ResReq",
	HTML:    "form-resource-request.html",
	Version: "0.6",
	Name:    "resource request form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.",
		"8c.", "7d.", "8d.", "20.", "21.", "22d.", "22t.", "23n.",
		"23q.", "23d.", "23i.", "24n.", "24q.", "24d.", "24i.", "25n.",
		"25q.", "25d.", "25i.", "26n.", "26q.", "26d.", "26i.", "27n.",
		"27q.", "27d.", "27i.", "28n.", "28q.", "28d.", "28i.", "30.",
		"40.", "40s.", "41.", "42.", "43.", "50.", "51.", "52.", "53.",
		"60.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall",
		"OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type, decode, create)
}

var basePDFRenderers = baseform.BaseFormPDF{
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

// ResourceRequest holds a resource request form.
type ResourceRequest struct {
	message.BaseMessage
	baseform.BaseForm
	Title                         string
	Jurisdiction                  string
	RequestDate                   string
	RequestTime                   string
	Resources                     [6]Resource
	Priority                      string
	RequestedByName               string
	WithSignature                 string
	RequestedByPhone              string
	RequestedByEmail              string
	RequestedByPosition           string
	RelatedToJurisdictionIncident string
	JurisdictionIncidentName      string
	RelatedToCountyIncident       string
	CountyIncidentName            string
	Comments                      string
}

// A Resource is the description of a single resource in a resource request
// form.
type Resource struct {
	ItemName       string
	QtyRequested   string
	Description    string
	Demobilization string
}

func create() message.Message {
	f := makeF()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.RequestDate = f.MessageDate
	f.ToLocation = "County EOC"
	return f
}

func makeF() *ResourceRequest {
	const fieldCount = 63
	f := ResourceRequest{BaseMessage: message.BaseMessage{Type: &Type}}
	f.FSubject = &f.Title
	f.FBody = &f.Comments
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:       "Title",
			Value:       &f.Title,
			Presence:    message.Required,
			PIFOTag:     "20.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the title of the resource request.  It should be unique among all resource requests from the jurisdiction.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Jurisdiction",
			Value:       &f.Jurisdiction,
			Choices:     message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara (City)", "Saratoga", "Sunnyvale", "Santa Clara County", "County unincorporated"},
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 190, Y: 184, R: 373, B: 206, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   21,
			EditHelp:    `This is the name of the jurisdiction requesting resources.  It is required.`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Date",
			Value:       &f.RequestDate,
			Presence:    message.Required,
			PIFOTag:     "22d.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the date of the resource request.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time",
			Value:       &f.RequestTime,
			Presence:    message.Required,
			PIFOTag:     "22t.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the time of day of the resource request.  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time",
			EditHelp: `This is the date and time of the resource request, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.RequestDate, &f.RequestTime),
	)
	for i := range f.Resources {
		f.Fields = append(f.Fields, f.Resources[i].Fields(&f, i+1)...)
	}
	f.Fields = append(f.Fields,
		message.NewRestrictedField(&message.Field{
			Label:       "Priority",
			Value:       &f.Priority,
			PIFOTag:     "30.",
			Choices:     message.Choices{"Urgent", "High", "Medium", "Low"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the priority of the resource request.  Urgent means as soon as possible.  High means within 24 hours.  Medium means within 1-2 days.  Low means within 2-3 days.  The priority is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Requested By Name",
			Value:       &f.RequestedByName,
			Presence:    message.Required,
			PIFOTag:     "40.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the name of the person requesting the resources.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "with signature",
			Value:       &f.WithSignature,
			PIFOTag:     "40s.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the original resource request form was signed.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Requested By Phone",
			Value:       &f.RequestedByPhone,
			Presence:    message.Required,
			PIFOTag:     "41.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the phone number of the person requesting the resources.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Requested By Email",
			Value:       &f.RequestedByEmail,
			Presence:    message.Required,
			PIFOTag:     "42.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the email address of the person requesting the resources.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Requested By Position",
			Value:       &f.RequestedByPosition,
			Presence:    message.Required,
			PIFOTag:     "43.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the position held by the person requesting the resources.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Related to Jurisdiction Incident?",
			Value:       &f.RelatedToJurisdictionIncident,
			Presence:    message.Required,
			PIFOTag:     "50.",
			Choices:     message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the resource request is related to a current jurisdiction activation.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label: "Name of Jurisdiction's Incident",
			Value: &f.JurisdictionIncidentName,
			Presence: func() (message.Presence, string) {
				if f.RelatedToJurisdictionIncident == "Yes" {
					return message.PresenceRequired, `"Related to Jurisdiction Incident" is "Yes"`
				} else {
					return message.PresenceNotAllowed, `"Related to Jurisdiction Incident" is not "Yes"`
				}
			},
			PIFOTag:     "51.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the name of the current jurisdiction incident that this resource request relates to.  It is required if "Relates to Jurisdiction Incident" is "Yes", otherwise not allowed.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Related to County Incident?",
			Value:       &f.RelatedToCountyIncident,
			Presence:    message.Required,
			PIFOTag:     "52.",
			Choices:     message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the resource request is related to a current county activation.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label: "Name of County's Incident",
			Value: &f.CountyIncidentName,
			Presence: func() (message.Presence, string) {
				if f.RelatedToCountyIncident == "Yes" {
					return message.PresenceRequired, `"Related to County Incident" is "Yes"`
				} else {
					return message.PresenceNotAllowed, `"Related to County Incident" is not "Yes"`
				}
			},
			PIFOTag:     "53.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the name of the current county incident that this resource request relates to.  It is required if "Relates to County Incident" is "Yes", otherwise not allowed.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Comments",
			Value:       &f.Comments,
			PIFOTag:     "60.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This field contains comments about the resource request.`,
		}),
	)
	f.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update ResourceRequest fieldCount")
	}
	return &f
}

func decode(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *ResourceRequest

	if form == nil || form.HTMLIdent != Type.HTML || form.FormVersion != Type.Version {
		return nil
	}
	df = makeF()
	message.DecodeForm(form, df)
	return df
}

func (r *Resource) Fields(m *ResourceRequest, index int) []*message.Field {
	var namePresence, qtyPresence, descPresence func() (message.Presence, string)
	if index == 1 {
		namePresence = message.Required
		qtyPresence = message.Required
		descPresence = message.Optional
	} else {
		namePresence = message.Optional
		qtyPresence = r.requiredIfNameElseNotAllowed
		descPresence = r.allowedIfName
	}
	return []*message.Field{
		message.NewTextField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Item Name", index),
			Value:       &r.ItemName,
			Presence:    namePresence,
			PIFOTag:     fmt.Sprintf("%dn.", 22+index),
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the name of the item being requested.`,
			EditSkip: func(f *message.Field) bool {
				return index > 1 && m.Resources[index-2].ItemName != ""
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Quantity Requested", index),
			Value:       &r.QtyRequested,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("%dq.", 22+index),
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the quantity (including units) of the item being requested.  It is required when there is an item name, otherwise not allowed.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Description", index),
			Value:       &r.Description,
			Presence:    descPresence,
			PIFOTag:     fmt.Sprintf("%dd.", 22+index),
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is a description of the item being requested.  It is allowed only when there is an item name.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Demobilization", index),
			Value:       &r.Demobilization,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("%di.", 22+index),
			Choices:     message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the requested item needs to be tracked for demobilization.  It is required when there is an item name, otherwise not allowed.`,
		}),
	}
}

func (r *Resource) requiredIfNameElseNotAllowed() (message.Presence, string) {
	if r.ItemName != "" {
		return message.PresenceRequired, "there is an item name"
	} else {
		return message.PresenceNotAllowed, "there is no item name"
	}
}

func (r *Resource) allowedIfName() (message.Presence, string) {
	if r.ItemName != "" {
		return message.PresenceOptional, ""
	} else {
		return message.PresenceNotAllowed, "there is no item name"
	}
}
