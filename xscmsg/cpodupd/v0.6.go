// Package cpodupd defines the CPOD Commodities Update Form message type.
package cpodupd

import (
	"fmt"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for a CPOD commodities update form.
var Type = message.Type{
	Tag:     "CPODUpd",
	HTML:    "form-cpod-commodities.html",
	Version: "0.6",
	Name:    "CPOD commodities update form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.",
		"8c.", "7d.", "8d.", "20.", "21d.", "21t.", "30.", "31.",
		"70a.", "70b.", "70c.", "70d.", "71a.", "71b.", "71c.", "71d.",
		"72a.", "72b.", "72c.", "72d.", "73a.", "73b.", "73c.", "73d.",
		"74a.", "74b.", "74c.", "74d.", "75a.", "75b.", "75c.", "75d.",
		"76a.", "76b.", "76c.", "76d.", "77a.", "77b.", "77c.", "77d.",
		"OpRelayRcvd", "OpRelaySent", "OpName", "OpCall", "OpDate",
		"OpTime",
	},
}

func init() {
	message.Register(&Type, decode, create)
}

var basePDFRenderers = baseform.BaseFormPDF{
	OriginMsgID:      &message.PDFTextRenderer{X: 228.96, Y: 62.04, W: 106.80, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	DestinationMsgID: &message.PDFTextRenderer{X: 438.60, Y: 62.04, W: 122.16, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 72.12, Y: 114.24, W: 56.64, H: 17.88, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageTime:      &message.PDFTextRenderer{X: 171.48, Y: 114.24, W: 65.28, H: 17.88, Style: message.PDFTextStyle{VAlign: "baseline"}},
	Handling: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
		"IMMEDIATE": {303.38, 122.97},
		"PRIORITY":  {402.98, 122.97},
		"ROUTINE":   {486.50, 122.97},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 132.84, Y: 139.08, W: 162.48, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 132.84, Y: 157.08, W: 162.48, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 132.84, Y: 175.08, W: 162.48, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 132.84, Y: 193.08, W: 162.48, H: 15.12, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 389.28, Y: 139.08, W: 170.16, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 389.28, Y: 157.08, W: 170.16, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 389.28, Y: 175.08, W: 170.16, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 389.28, Y: 193.08, W: 170.16, H: 15.12, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{X: 111.00, Y: 685.20, W: 205.56, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{X: 357.12, Y: 685.20, W: 203.64, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{X: 77.04, Y: 704.16, W: 110.28, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{X: 245.52, Y: 704.16, W: 49.80, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{X: 336.96, Y: 704.16, W: 88.80, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{X: 495.84, Y: 704.16, W: 64.92, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

// CPODUpdate holds a CPOD commodities update form.
type CPODUpdate struct {
	message.BaseMessage
	baseform.BaseForm
	Jurisdiction string
	PreparedDate string
	PreparedTime string
	SiteName     string
	Status       string
	Commodities  [8]Commodity
}

// A Commodity is the description of a single commodity in a CPOD commodities
// update form.
type Commodity struct {
	Type           string
	StartingQty    string
	QtyDistributed string
	QtyAvailable   string
}

func create() message.Message {
	f := makeF()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.PreparedDate = f.MessageDate
	f.Handling = "ROUTINE"
	f.ToLocation = "County EOC"
	return f
}

func makeF() *CPODUpdate {
	const fieldCount = 60
	f := CPODUpdate{BaseMessage: message.BaseMessage{Type: &Type}}
	f.FSubject = &f.SiteName
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:       "Jurisdiction",
			Value:       &f.Jurisdiction,
			Choices:     message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara (City)", "Saratoga", "Sunnyvale", "Santa Clara County", "County unincorporated"},
			Presence:    message.Required,
			PIFOTag:     "20.",
			PDFRenderer: &message.PDFTextRenderer{X: 116.64, Y: 233.16, W: 442.80, H: 11.76},
			EditWidth:   21,
			EditHelp:    `This is the name of the jurisdiction responsible for the CPOD.  It is required.`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Prepared Date",
			Value:       &f.PreparedDate,
			Presence:    message.Required,
			PIFOTag:     "21d.",
			PDFRenderer: &message.PDFTextRenderer{X: 116.64, Y: 251.76, W: 178.68, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the form was prepared.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Prepared Time",
			Value:       &f.PreparedTime,
			Presence:    message.Required,
			PIFOTag:     "21t.",
			PDFRenderer: &message.PDFTextRenderer{X: 383.04, Y: 251.76, W: 176.40, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time when the form was prepared.  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Prepared",
			EditHelp: `This is the date and time when the form was prepared, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.PreparedDate, &f.PreparedTime),
		message.NewTextField(&message.Field{
			Label:       "Site Name",
			Value:       &f.SiteName,
			Presence:    message.Required,
			PIFOTag:     "30.",
			PDFRenderer: &message.PDFTextRenderer{X: 96.84, Y: 287.76, W: 462.60, H: 11.52},
			EditWidth:   80,
			EditHelp:    `This is the name of the CPOD site.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Status",
			Value:    &f.Status,
			Presence: message.Required,
			PIFOTag:  "31.",
			Choices:  message.Choices{"Activated", "Pending Activation", "Pending Demobilization", "Demobilization", "Not Activated"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Activated":              {155.66, 310.89},
				"Pending Activation":     {155.66, 326.37},
				"Pending Demobilization": {299.66, 310.89},
				"Demobilization":         {299.66, 326.37},
				"Not Activated":          {443.66, 310.89},
			}},
			EditHelp: `This is the status of the CPOD.  It is required.`,
		}),
	)
	for i := range f.Commodities {
		f.Fields = append(f.Fields, f.Commodities[i].Fields(&f, i+1)...)
	}
	f.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update CPODUpd fieldCount")
	}
	return &f
}

func decode(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *CPODUpdate

	if form == nil || form.HTMLIdent != Type.HTML || form.FormVersion != Type.Version {
		return nil
	}
	df = makeF()
	message.DecodeForm(form, df)
	return df
}

func (c *Commodity) Fields(m *CPODUpdate, index int) []*message.Field {
	var typePresence, qtyPresence func() (message.Presence, string)
	if index == 1 {
		typePresence = message.Required
		qtyPresence = message.Required
	} else {
		typePresence = message.Optional
		qtyPresence = c.requiredIfTypeElseNotAllowed
	}
	var offset = 36.994 * float64(index-1)
	return []*message.Field{
		message.NewTextField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Type of Commodity", index),
			Value:       &c.Type,
			Presence:    typePresence,
			PIFOTag:     fmt.Sprintf("%da.", 69+index),
			PDFRenderer: &message.PDFTextRenderer{X: 189.60, Y: 365.40 + offset, W: 369.84, H: 11.40, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the type of a commodity distributed at the CPOD site.`,
			EditSkip: func(f *message.Field) bool {
				return index > 1 && m.Commodities[index-2].Type == ""
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Starting Quantity", index),
			Value:       &c.StartingQty,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("%db.", 69+index),
			PDFRenderer: &message.PDFTextRenderer{X: 178.80, Y: 383.40 + offset, W: 63.72, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   4,
			EditHelp:    `This is the quantity of the commodity that the CPOD site had when it opened.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Qty Distributed", index),
			Value:       &c.QtyDistributed,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("%dc.", 69+index),
			PDFRenderer: &message.PDFTextRenderer{X: 333.96, Y: 383.40 + offset, W: 72.12, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   4,
			EditHelp:    `This is the quantity of the commodity that the CPOD site has distributed to visitors.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Qty Available", index),
			Value:       &c.QtyAvailable,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("%dd.", 69+index),
			PDFRenderer: &message.PDFTextRenderer{X: 488.16, Y: 383.40 + offset, W: 71.28, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   4,
			EditHelp:    `This is the quantity of the commodity that the CPOD site has available for distribution.  It is required.`,
		}),
	}
}

func (c *Commodity) requiredIfTypeElseNotAllowed() (message.Presence, string) {
	if c.Type != "" {
		return message.PresenceRequired, "there is a commodity type named"
	} else {
		return message.PresenceNotAllowed, "there is no commodity type named"
	}
}
