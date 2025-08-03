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
			PDFRenderer: &message.PDFTextRenderer{X: 190, Y: 184, R: 373, B: 206, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   21,
			EditHelp:    `This is the name of the jurisdiction responsible for the CPOD.  It is required.`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Prepared Date",
			Value:       &f.PreparedDate,
			Presence:    message.Required,
			PIFOTag:     "21d.",
			PDFRenderer: &message.PDFTextRenderer{X: 160, Y: 467, R: 388, B: 483, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the form was prepared.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Prepared Time",
			Value:       &f.PreparedTime,
			Presence:    message.Required,
			PIFOTag:     "21t.",
			PDFRenderer: &message.PDFTextRenderer{X: 449, Y: 467, R: 572, B: 483, Style: message.PDFTextStyle{VAlign: "baseline"}},
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
			PDFRenderer: &message.PDFTextRenderer{X: 190, Y: 184, R: 373, B: 206, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   80,
			EditHelp:    `This is the name of the CPOD site.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Status",
			Value:       &f.Status,
			Presence:    message.Required,
			PIFOTag:     "31.",
			Choices:     message.Choices{"Activated", "Pending Activation", "Pending Demobilization", "Demobilization", "Not Activated"},
			PDFRenderer: &message.PDFMappedTextRenderer{X: 180, Y: 671, B: 693, Map: map[string]string{"checked": "[with signature]"}},
			EditHelp:    `This is the status of the CPOD.  It is required.`,
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
	return []*message.Field{
		message.NewTextField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Type of Commodity", index),
			Value:       &c.Type,
			Presence:    typePresence,
			PIFOTag:     fmt.Sprintf("%da.", 69+index),
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
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
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   4,
			EditHelp:    `This is the quantity of the commodity that the CPOD site had when it opened.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Qty Distributed", index),
			Value:       &c.QtyDistributed,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("%dc.", 69+index),
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   4,
			EditHelp:    `This is the quantity of the commodity that the CPOD site has distributed to visitors.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Qty Available", index),
			Value:       &c.QtyAvailable,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("%dd.", 69+index),
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
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
