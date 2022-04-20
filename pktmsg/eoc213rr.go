package pktmsg

// This file defines TxEOC213RRForm and RxEOC213RRForm.

import (
	"fmt"
	"time"
)

// A TxEOC213RRForm is an outgoing PackItForms-encoded message containing an
// SCCo EOC-213RR form.
type TxEOC213RRForm struct {
	TxSCCoForm
	IncidentName             string
	DateTimeInitiated        time.Time
	RequestedBy              string
	PreparedBy               string
	ApprovedBy               string
	QtyUnit                  string
	ResourceDescription      string
	Arrival                  string
	Priority                 string
	EstimatedCost            string
	DeliverTo                string
	DeliverToLocation        string
	SuggestedSources         string
	RequireEquipmentOperator bool
	RequireLodging           bool
	RequireFuel              bool
	FuelRequirement          string
	RequirePower             bool
	RequireMeals             bool
	RequireMaintenance       bool
	RequireWater             bool
	RequireOther             bool
	SpecialInstructions      string
}

// Encode returns the encoded subject line and body of the message.
func (e *TxEOC213RRForm) Encode() (subject, body string, err error) {
	if err = e.checkHeaderFooterFields(); err != nil {
		return "", "", err
	}
	if e.Subject != "" {
		return "", "", ErrDontSet
	}
	if e.IncidentName == "" ||
		e.DateTimeInitiated.IsZero() ||
		e.RequestedBy == "" ||
		e.QtyUnit == "" ||
		e.ResourceDescription == "" ||
		e.Arrival == "" ||
		e.Priority == "" ||
		e.DeliverTo == "" ||
		e.DeliverToLocation == "" {
		return "", "", ErrIncomplete
	}
	if e.Priority != "Now" && e.Priority != "High" && e.Priority != "Medium" && e.Priority != "Low" {
		return "", "", ErrInvalid
	}
	e.FormName = "EOC213RR"
	e.FormHTML = "form-scco-eoc-213rr.html"
	e.FormVersion = "2.3"
	e.Subject = e.IncidentName
	e.encodeHeaderFields()
	e.SetField("21.", e.IncidentName)
	e.SetField("22.", e.DateTimeInitiated.Format("01/02/2006"))
	e.SetField("23.", e.DateTimeInitiated.Format("15:04"))
	e.SetField("25.", e.RequestedBy)
	e.SetField("26.", e.PreparedBy)
	e.SetField("27.", e.ApprovedBy)
	e.SetField("28.", e.QtyUnit)
	e.SetField("29.", e.ResourceDescription)
	e.SetField("30.", e.Arrival)
	e.SetField("31.", e.Priority)
	e.SetField("32.", e.EstimatedCost)
	e.SetField("33.", e.DeliverTo)
	e.SetField("34.", e.DeliverToLocation)
	e.SetField("35.", e.SuggestedSources)
	e.SetField("36a.", boolToChecked(e.RequireEquipmentOperator))
	e.SetField("36b.", boolToChecked(e.RequireLodging))
	e.SetField("36c.", boolToChecked(e.RequireFuel))
	e.SetField("36d.", e.FuelRequirement)
	e.SetField("36e.", boolToChecked(e.RequirePower))
	e.SetField("36f.", boolToChecked(e.RequireMeals))
	e.SetField("36g.", boolToChecked(e.RequireMaintenance))
	e.SetField("36h.", boolToChecked(e.RequireWater))
	e.SetField("36i.", boolToChecked(e.RequireOther))
	e.SetField("37.", e.SpecialInstructions)
	e.encodeFooterFields()
	return e.TxSCCoForm.Encode()
}

//------------------------------------------------------------------------------

// An RxEOC213RRForm is a received PackItForms-encoded message containing an
// SCCo EOC-213RR form.
type RxEOC213RRForm struct {
	RxSCCoForm
	IncidentName             string
	DateInitiated            string
	TimeInitiated            string
	DateTimeInitiated        time.Time
	RequestedBy              string
	PreparedBy               string
	ApprovedBy               string
	QtyUnit                  string
	ResourceDescription      string
	Arrival                  string
	Priority                 string
	EstimatedCost            string
	DeliverTo                string
	DeliverToLocation        string
	SuggestedSources         string
	RequireEquipmentOperator bool
	RequireLodging           bool
	RequireFuel              bool
	FuelRequirement          string
	RequirePower             bool
	RequireMeals             bool
	RequireMaintenance       bool
	RequireWater             bool
	RequireOther             bool
	SpecialInstructions      string
}

// parseRxEOC213RRForm examines an RxForm to see if it contains an EOC-213RR
// form, and if so, wraps it in an RxEOC213RRForm and returns it.  If it is not,
// it returns nil.
func parseRxEOC213RRForm(f *RxForm) *RxEOC213RRForm {
	var e RxEOC213RRForm

	if f.FormHTML != "form-scco-eoc-213rr.html" {
		return nil
	}
	e.RxSCCoForm.RxForm = *f
	e.extractHeaderFields()
	e.IncidentName = e.Fields["21."]
	e.DateInitiated = e.Fields["22."]
	e.TimeInitiated = e.Fields["23."]
	e.DateTimeInitiated = dateTimeParse(e.DateInitiated, e.TimeInitiated)
	e.RequestedBy = e.Fields["25."]
	e.PreparedBy = e.Fields["26."]
	e.ApprovedBy = e.Fields["27."]
	e.QtyUnit = e.Fields["28."]
	e.ResourceDescription = e.Fields["29."]
	e.Arrival = e.Fields["30."]
	e.Priority = e.Fields["31."]
	e.EstimatedCost = e.Fields["32."]
	e.DeliverTo = e.Fields["33."]
	e.DeliverToLocation = e.Fields["34."]
	e.SuggestedSources = e.Fields["35."]
	e.RequireEquipmentOperator = e.Fields["36a."] != ""
	e.RequireLodging = e.Fields["36b."] != ""
	e.RequireFuel = e.Fields["36c."] != ""
	e.FuelRequirement = e.Fields["36d."]
	e.RequirePower = e.Fields["36e."] != ""
	e.RequireMeals = e.Fields["36f."] != ""
	e.RequireMaintenance = e.Fields["36g."] != ""
	e.RequireWater = e.Fields["36h."] != ""
	e.RequireOther = e.Fields["36i."] != ""
	e.SpecialInstructions = e.Fields["37."]
	e.extractFooterFields()
	return &e
}

// Valid returns whether all of the fields of the form have valid values, and
// all required fields are filled in.
func (e *RxEOC213RRForm) Valid() bool {
	return e.RxSCCoForm.Valid() &&
		e.IncidentName != "" &&
		!e.DateTimeInitiated.IsZero() &&
		e.RequestedBy != "" &&
		e.QtyUnit != "" &&
		e.ResourceDescription != "" &&
		e.Arrival != "" &&
		(e.Priority == "Now" || e.Priority == "High" || e.Priority == "Medium" || e.Priority == "Low") &&
		e.Priority != "" &&
		e.DeliverTo != "" &&
		e.DeliverToLocation != ""
}

// EncodeSubjectLine returns what the subject line should be based on the
// received form contents.
func (e *RxEOC213RRForm) EncodeSubjectLine() string {
	return fmt.Sprintf("%s_%s_EOC213RR_%s", e.OriginMessageNumber, e.HandlingOrder.Code(), e.IncidentName)
}

// TypeCode returns the machine-readable code for the message type.
func (*RxEOC213RRForm) TypeCode() string { return "EOC213RR" }

// TypeName returns the human-reading name of the message type.
func (*RxEOC213RRForm) TypeName() string { return "EOC-213RR form" }

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (*RxEOC213RRForm) TypeArticle() string { return "an" }
