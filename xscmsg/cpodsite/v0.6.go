// Package cpodsite defines the CPOD Site Information Form message type.
package cpodsite

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for a CPOD site information form.
var Type = message.Type{
	Tag:     "CPODSite",
	HTML:    "form-cpod-site-info.html",
	Version: "0.6",
	Name:    "CPOD site information form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.",
		"8c.", "7d.", "8d.", "20.", "21d.", "21t.", "22.", "23.", "24.",
		"25a.", "25c.", "25z.", "26.", "27.", "28.", "30.", "31.",
		"40.", "41.", "42a.", "42b.", "42c.", "42d.", "43.", "44.",
		"45.", "46.", "47.", "50a.", "50b.", "50c.", "50d.", "50e.",
		"50f.", "50g.", "51.", "52a.", "52b.", "52c.", "60d.", "60t.",
		"61d.", "61t.", "70a.", "70b.", "70c.", "70d.", "71a.", "71b.",
		"71c.", "71d.", "72a.", "72b.", "72c.", "72d.", "73a.", "73b.",
		"73c.", "73d.", "74a.", "74b.", "74c.", "74d.", "75a.", "75b.",
		"75c.", "75d.", "76a.", "76b.", "76c.", "76d.", "77a.", "77b.",
		"77c.", "77d.", "OpRelayRcvd", "OpRelaySent", "OpName",
		"OpCall", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type, decode, create)
}

var basePDFRenderers = baseform.BaseFormPDF{
	OriginMsgID: &message.PDFMultiRenderer{
		&message.PDFTextRenderer{X: 219.96, Y: 44.04, W: 100.80, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
		&message.PDFTextRenderer{Page: 2, X: 447.72, Y: 18.00, W: 111.72, H: 19.44},
	},
	DestinationMsgID: &message.PDFTextRenderer{X: 423.60, Y: 44.04, W: 137.16, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 76.56, Y: 96.24, W: 61.20, H: 19.56, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageTime:      &message.PDFTextRenderer{X: 207.84, Y: 96.24, W: 33.48, H: 19.56, Style: message.PDFTextStyle{VAlign: "baseline"}},
	Handling: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
		"IMMEDIATE": {307.8, 105.8},
		"PRIORITY":  {404.9, 105.8},
		"ROUTINE":   {486.0, 105.8},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 128.28, Y: 122.40, W: 171.48, H: 15.96, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 128.28, Y: 144.84, W: 171.48, H: 15.96, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 128.28, Y: 167.40, W: 171.48, H: 15.96, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 128.28, Y: 189.84, W: 171.48, H: 15.60, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 393.84, Y: 122.40, W: 165.60, H: 15.96, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 393.84, Y: 144.84, W: 165.60, H: 15.96, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 393.84, Y: 167.40, W: 165.60, H: 15.96, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 393.84, Y: 189.84, W: 165.60, H: 15.60, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{Page: 2, X: 111.00, Y: 516.60, W: 205.56, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{Page: 2, X: 357.12, Y: 516.60, W: 203.64, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{Page: 2, X: 77.04, Y: 535.44, W: 123.72, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{Page: 2, X: 259.08, Y: 535.44, W: 67.68, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{Page: 2, X: 368.52, Y: 535.44, W: 57.24, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{Page: 2, X: 495.84, Y: 535.44, W: 64.92, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

// CPODSite holds a CPOD site information form.
type CPODSite struct {
	message.BaseMessage
	baseform.BaseForm
	Jurisdiction          string
	PreparedDate          string
	PreparedTime          string
	SiteName              string
	CPODType              string
	Status                string
	Address               string
	City                  string
	ZIPCode               string
	PoliceJurisdiction    string
	FireJurisdiction      string
	AdditionalInfo        string
	PartnerContact        string
	SiteContact           string
	Dimensions            string
	Size                  string
	AccessConcrete        string
	AccessGravel          string
	AccessPaved           string
	AccessOther           string
	Accessible            string
	AccessGate            string
	GateClosedContact     string
	Driveways             string
	SpikeStrips           string
	SafetyFencing         string
	SafetyOutsideLighting string
	SafetyInsideLighting  string
	SafetyCCTV            string
	SafetyPA              string
	SafetyCovered         string
	SafetyNoMove          string
	Fencing               string
	AccessWheelchair      string
	AccessUneven          string
	AccessRamp            string
	DateOpened            string
	TimeOpened            string
	DateClosed            string
	TimeClosed            string
	Commodities           [8]Commodity
}

// A Commodity is the description of a single commodity in a CPOD site
// information form.
type Commodity struct {
	Type           string
	StartingQty    string
	QtyDistributed string
	QtyAvailable   string
}

var zipCodeRE = regexp.MustCompile(`^\d{5}(?:-\d{4})?$`)

func create() message.Message {
	f := makeF()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.PreparedDate = f.MessageDate
	f.Handling = "ROUTINE"
	f.ToLocation = "County EOC"
	return f
}

func makeF() *CPODSite {
	const fieldCount = 100
	f := CPODSite{BaseMessage: message.BaseMessage{Type: &Type}}
	f.FSubject = &f.SiteName
	f.FBody = &f.AdditionalInfo
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:       "Jurisdiction",
			Value:       &f.Jurisdiction,
			Choices:     message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara (City)", "Saratoga", "Sunnyvale", "Santa Clara County", "County unincorporated"},
			Presence:    message.Required,
			PIFOTag:     "20.",
			PDFRenderer: &message.PDFTextRenderer{X: 121.20, Y: 230.04, W: 438.24, H: 11.64},
			EditHelp:    `This is the name of the jurisdiction responsible for the CPOD.  It is required.`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Prepared Date",
			Value:       &f.PreparedDate,
			Presence:    message.Required,
			PIFOTag:     "21d.",
			PDFRenderer: &message.PDFTextRenderer{X: 121.20, Y: 248.52, W: 178.08, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the form was prepared.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Prepared Time",
			Value:       &f.PreparedTime,
			Presence:    message.Required,
			PIFOTag:     "21t.",
			PDFRenderer: &message.PDFTextRenderer{X: 387.60, Y: 248.52, W: 171.84, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
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
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 121.20, Y: 267.00, W: 438.24, H: 11.52},
			EditWidth:   80,
			EditHelp:    `This is the name of the CPOD site.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "CPOD Type",
			Value:    &f.CPODType,
			Presence: message.Required,
			PIFOTag:  "23.",
			Choices:  message.Choices{"Type I", "Type II", "Type III", "LSA", "Non-CPOD Distribution Point"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Type I":                      {124.1, 290.1},
				"Type II":                     {124.1, 305.6},
				"Type III":                    {232.1, 290.1},
				"LSA":                         {232.1, 305.6},
				"Non-CPOD Distribution Point": {412.1, 290.1},
			}},
			EditHelp: `This is the type of CPOD.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Status",
			Value:    &f.Status,
			Presence: message.Required,
			PIFOTag:  "24.",
			Choices:  message.Choices{"Activated", "Pending Activation", "Pending Demobilization", "Demobilization", "Not Activated"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Activated":              {124.1, 322.2},
				"Pending Activation":     {124.1, 337.6},
				"Pending Demobilization": {268.1, 322.2},
				"Demobilization":         {268.1, 337.6},
				"Not Activated":          {412.1, 322.2},
			}},
			EditHelp: `This is the status of the CPOD.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Address",
			Value:       &f.Address,
			Presence:    message.Required,
			PIFOTag:     "25a.",
			PDFRenderer: &message.PDFTextRenderer{X: 91.08, Y: 349.68, W: 167.76, H: 11.28, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   35,
			EditHelp:    `This is the street address of the CPOD site.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "City",
			Value:       &f.City,
			Presence:    message.Required,
			PIFOTag:     "25c.",
			PDFRenderer: &message.PDFTextRenderer{X: 296.88, Y: 349.68, W: 146.40, H: 11.28, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   31,
			EditHelp:    `This is the city portion of the address of the CPOD site.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:   "ZIP Code",
			Value:   &f.ZIPCode,
			PIFOTag: "25z.",
			PIFOValid: func(field *message.Field) string {
				if p := field.PresenceValid(); p != "" {
					return p
				}
				if *field.Value != "" && !zipCodeRE.MatchString(*field.Value) {
					return fmt.Sprintf("The %q field does not contain a valid ZIP code.", field.Label)
				}
				return ""
			},
			PDFRenderer: &message.PDFTextRenderer{X: 502.32, Y: 349.68, W: 57.12, H: 11.28, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   10,
			EditHelp:    `This is the ZIP code of the address of the CPOD site.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Police Jurisdiction",
			Value:       &f.PoliceJurisdiction,
			PIFOTag:     "26.",
			PDFRenderer: &message.PDFTextRenderer{X: 279.12, Y: 367.92, W: 280.32, H: 11.76},
			EditWidth:   59,
			EditHelp:    `This is the police department, and division or region if applicable, with jurisdiction over the CPOD site.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Fire Jurisdiction",
			Value:       &f.FireJurisdiction,
			PIFOTag:     "27.",
			PDFRenderer: &message.PDFTextRenderer{X: 279.12, Y: 386.52, W: 280.32, H: 11.64},
			EditWidth:   59,
			EditHelp:    `This is the fire department, and division or region if applicable, with jurisdiction over the CPOD site.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Additional Information",
			Value:       &f.AdditionalInfo,
			PIFOTag:     "28.",
			PDFRenderer: &message.PDFTextRenderer{X: 161.54, Y: 405.00, W: 397.90, H: 23.64, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   80,
			EditHelp:    `This is any additional pertinent information regarding the CPOD and its site.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Partner Point of Contact",
			Value:       &f.PartnerContact,
			PIFOTag:     "30.",
			PDFRenderer: &message.PDFTextRenderer{X: 165.00, Y: 453.12, W: 134.28, H: 11.40, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   28,
			EditHelp:    `This is the name and contact information of the contact person with the partner operating the CPOD.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Site Point of Contact",
			Value:       &f.SiteContact,
			PIFOTag:     "31.",
			PDFRenderer: &message.PDFTextRenderer{X: 413.52, Y: 453.12, W: 145.92, H: 11.40, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   31,
			EditHelp:    `This is the name and contact information of the contact person responsible for the site hosting the CPOD.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Dimensions of Site in Feet",
			Value:       &f.Dimensions,
			PIFOTag:     "40.",
			PDFRenderer: &message.PDFTextRenderer{X: 172.80, Y: 489.00, W: 126.48, H: 11.16, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   27,
			EditHelp:    `This is the dimensions of the CPOD site, in feet (e.g. "200x150").`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Size of Site in Acres",
			Value:       &f.Size,
			PIFOTag:     "41.",
			PDFRenderer: &message.PDFTextRenderer{X: 408.00, Y: 489.00, W: 151.44, H: 11.16, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   32,
			EditHelp:    `This is the size of the CPOD site in acres.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Access: Concrete",
			Value:       &f.AccessConcrete,
			PIFOTag:     "42a.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {155.0, 507.1}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that some access roads or walkways are concrete.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Access: Gravel Hard-Stand",
			Value:       &f.AccessGravel,
			PIFOTag:     "42b.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {227.0, 507.1}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that some access roads or walkways are gravel hard-stand.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Access: Paved",
			Value:       &f.AccessPaved,
			PIFOTag:     "42c.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {371.0, 507.1}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that some access roads or walkways are paved.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Access: Other",
			Value:       &f.AccessOther,
			PIFOTag:     "42d.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {479.0, 507.1}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that some access roads or walkways have other surfaces.  Give the details in the Additional Information field.`,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Access",
			TableValue: func(*message.Field) string {
				var access []string
				if f.AccessConcrete != "" {
					access = append(access, "Concrete")
				}
				if f.AccessGravel != "" {
					access = append(access, "Gravel Hard-Stand")
				}
				if f.AccessPaved != "" {
					access = append(access, "Paved")
				}
				if f.AccessOther != "" {
					access = append(access, "Other")
				}
				return strings.Join(access, ", ")
			},
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Accessible at All Times",
			Value:   &f.Accessible,
			PIFOTag: "43.",
			Choices: message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Yes": {196.1, 530.1},
				"No":  {242.6, 530.1},
			}},
			EditHelp: `This indicates whether the site is accessible at all times.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Access Controlled by Gate",
			Value:   &f.AccessGate,
			PIFOTag: "44.",
			Choices: message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Yes": {196.1, 548.1},
				"No":  {242.6, 548.1},
			}},
			EditHelp: `This indicates whether site access is controlled by a gate.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Site Contact if Gate is Closed",
			Value:       &f.GateClosedContact,
			PIFOTag:     "45.",
			PDFRenderer: &message.PDFTextRenderer{X: 450.72, Y: 542.64, W: 108.72, H: 11.64},
			EditWidth:   23,
			EditHelp:    `This gives the name and contact information of the person to contact for access if the gate is closed.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Location of Driveway(s)",
			Value:       &f.Driveways,
			PIFOTag:     "46.",
			PDFRenderer: &message.PDFTextRenderer{X: 165.35, Y: 560.64, W: 394.09, H: 11.28, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   82,
			EditHelp:    `This gives the locations of all driveways into the site.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Spike Strips at any of the driveways",
			Value:   &f.SpikeStrips,
			PIFOTag: "47.",
			Choices: message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Yes": {268.1, 583.9},
				"No":  {314.6, 583.9},
			}},
			EditHelp: `This indicates whether any of the driveways have spike strips.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Safety: Has perimeter fencing",
			Value:       &f.SafetyFencing,
			PIFOTag:     "50a.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {83.03, 626.93}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether the site has perimeter fencing.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Safety: Has outside lighting",
			Value:       &f.SafetyOutsideLighting,
			PIFOTag:     "50b.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {83.03, 642.41}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether the site has fixed lighting throughout all outside areas.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Safety: Has inside lighting",
			Value:       &f.SafetyInsideLighting,
			PIFOTag:     "50c.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {83.03, 658.01}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether the site has fixed lighting throughout all inside areas.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Safety: Monitored with CCTV",
			Value:       &f.SafetyCCTV,
			PIFOTag:     "50d.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {83.03, 673.49}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether the site is monitored using closed-circuit TV cameras.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Safety: Has PA system",
			Value:       &f.SafetyPA,
			PIFOTag:     "50e.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {335.03, 626.93}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether the site has a public address system installed.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Safety: Has covered areas",
			Value:       &f.SafetyCovered,
			PIFOTag:     "50f.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {335.03, 642.41}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether the site has covered areas.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Safety: Has fixed equipment",
			Value:       &f.SafetyNoMove,
			PIFOTag:     "50g.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {335.03, 658.01}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether there is fixed or non-fixed equipment located on the site that may be difficult to move.`,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Safety",
			TableValue: func(*message.Field) string {
				var safety []string
				if f.SafetyFencing != "" {
					safety = append(safety, "has perimeter fencing")
				}
				if f.SafetyOutsideLighting != "" {
					safety = append(safety, "has outside lighting")
				}
				if f.SafetyInsideLighting != "" {
					safety = append(safety, "has inside lighting")
				}
				if f.SafetyCCTV != "" {
					safety = append(safety, "monitored with CCTV")
				}
				if f.SafetyPA != "" {
					safety = append(safety, "has PA system")
				}
				if f.SafetyCovered != "" {
					safety = append(safety, "has covered areas")
				}
				if f.SafetyNoMove != "" {
					safety = append(safety, "has fixed equipment")
				}
				if len(safety) != 0 {
					safety[0] = strings.ToUpper(safety[0][:1]) + safety[0][1:]
				}
				return strings.Join(safety, ", ")
			},
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Perimeter Fencing Details",
			Value:       &f.Fencing,
			PIFOTag:     "51.",
			PDFRenderer: &message.PDFTextRenderer{X: 174.90, Y: 689.88, W: 384.54, H: 21.12, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   80,
			EditHelp:    `This gives details of the perimeter fencing.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Accessibility: Wheelchair access",
			Value:       &f.AccessWheelchair,
			PIFOTag:     "52a.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {155.03, 716.93}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether there are sidewalks leading to the site that have wheelchair access.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Accessibility: Uneven surfaces",
			Value:       &f.AccessUneven,
			PIFOTag:     "52b.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {155.03, 732.53}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether there are uneven surfaces leading up to the site.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Accessibility: Ramp from staff parking",
			Value:       &f.AccessRamp,
			PIFOTag:     "52c.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {155.03, 745.78}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether there is a ramp from the staff parking location leading up to the POD location.`,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Accessibility",
			TableValue: func(*message.Field) string {
				var access []string
				if f.AccessWheelchair != "" {
					access = append(access, "wheelchair access")
				}
				if f.AccessUneven != "" {
					access = append(access, "uneven surfaces")
				}
				if f.AccessRamp != "" {
					access = append(access, "ramp from staff parking")
				}
				if len(access) != 0 {
					access[0] = strings.ToUpper(access[0][:1]) + access[0][1:]
				}
				return strings.Join(access, ", ")
			},
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Date Opened",
			Value:       &f.DateOpened,
			PIFOTag:     "60d.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 117.12, Y: 62.88, W: 137.64, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date on which the CPOD site opened or will open.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time Opened",
			Value:       &f.TimeOpened,
			PIFOTag:     "60t.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 336.12, Y: 62.88, W: 223.32, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time of day when the CPOD site opened or will open.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time Opened",
			EditHelp: `This is the date and time when the CPOD site opened or will open, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.DateOpened, &f.TimeOpened),
		message.NewDateField(true, &message.Field{
			Label:       "Date Closed",
			Value:       &f.DateClosed,
			PIFOTag:     "61d.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 117.12, Y: 80.88, W: 137.64, H: 11.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date on which the CPOD site closed or will close.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time Closed",
			Value:       &f.TimeClosed,
			PIFOTag:     "61t.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 336.12, Y: 80.88, W: 223.32, H: 11.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time of day when the CPOD site closed or will close.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time Closed",
			EditHelp: `This is the date and time when the CPOD site closed or will close, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.DateClosed, &f.TimeClosed),
	)
	for i := range f.Commodities {
		f.Fields = append(f.Fields, f.Commodities[i].Fields(&f, i+1)...)
	}
	f.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update CPODSite fieldCount")
	}
	return &f
}

func decode(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *CPODSite

	if form == nil || form.HTMLIdent != Type.HTML || form.FormVersion != Type.Version {
		return nil
	}
	df = makeF()
	message.DecodeForm(form, df)
	return df
}

func (c *Commodity) Fields(m *CPODSite, index int) []*message.Field {
	var typePresence, qtyPresence func() (message.Presence, string)
	if index == 1 {
		typePresence = message.Required
		qtyPresence = message.Required
	} else {
		typePresence = message.Optional
		qtyPresence = c.requiredIfTypeElseNotAllowed
	}
	var offset = 46.354 * float64(index-1)
	return []*message.Field{
		message.NewMultilineField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Type of Commodity", index),
			Value:       &c.Type,
			Presence:    typePresence,
			PIFOTag:     fmt.Sprintf("%da.", 69+index),
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 189.60, Y: 121.80 + offset, W: 369.84, H: 21.00, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   78,
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
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 178.80, Y: 149.28 + offset, W: 63.72, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   4,
			EditHelp:    `This is the quantity of the commodity that the CPOD site had when it opened.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Qty Distributed", index),
			Value:       &c.QtyDistributed,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("%dc.", 69+index),
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 333.96, Y: 149.28 + offset, W: 72.12, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   4,
			EditHelp:    `This is the quantity of the commodity that the CPOD site has distributed to visitors.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Item %d: Qty Available", index),
			Value:       &c.QtyAvailable,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("%dd.", 69+index),
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 488.16, Y: 149.28 + offset, W: 71.28, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
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
