package xscmsg

import (
	"strings"

	"github.com/rothskeller/packet/xscmsg/forms/pifo"
	"github.com/rothskeller/packet/xscmsg/forms/xscsubj"
)

// EOC213RR form metadata:
const (
	EOC213RRTag     = "EOC213RR"
	EOC213RRHTML    = "form-scco-eoc-213rr.html"
	EOC213RRVersion = "2.3"
)

// EOC213RR holds an EOC-213RR resource request form.
type EOC213RR struct {
	StdHeader
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
	StdFooter
}

// DecodeEOC213RR decodes the supplied form if it is an EOC213RR form.  It
// returns the decoded form and strings describing any non-fatal decoding
// problems.  It returns nil, nil if the form is not an EOC213RR form or has an
// unknown version.
func DecodeEOC213RR(form *pifo.Form) (f *EOC213RR, problems []string) {
	if form.HTMLIdent != EOC213RRHTML {
		return nil, nil
	}
	switch form.FormVersion {
	case "2.0", "2.1", "2.3":
		break
	default:
		return nil, nil
	}
	f = new(EOC213RR)
	f.FormVersion = form.FormVersion
	f.StdHeader.PullTags(form.TaggedValues)
	f.IncidentName = PullTag(form.TaggedValues, "21.")
	f.DateInitiated = PullTag(form.TaggedValues, "22.")
	f.TimeInitiated = PullTag(form.TaggedValues, "23.")
	f.TrackingNumber = PullTag(form.TaggedValues, "24.")
	f.RequestedBy = PullTag(form.TaggedValues, "25.")
	f.PreparedBy = PullTag(form.TaggedValues, "26.")
	f.ApprovedBy = PullTag(form.TaggedValues, "27.")
	f.QtyUnit = PullTag(form.TaggedValues, "28.")
	f.ResourceDescription = PullTag(form.TaggedValues, "29.")
	f.ResourceArrival = PullTag(form.TaggedValues, "30.")
	f.Priority = PullTag(form.TaggedValues, "31.")
	f.EstdCost = PullTag(form.TaggedValues, "32.")
	f.DeliverTo = PullTag(form.TaggedValues, "33.")
	f.DeliverToLocation = PullTag(form.TaggedValues, "34.")
	f.Substitutes = PullTag(form.TaggedValues, "35.")
	f.EquipmentOperator = PullTag(form.TaggedValues, "36a.")
	f.Lodging = PullTag(form.TaggedValues, "36b.")
	f.Fuel = PullTag(form.TaggedValues, "36c.")
	f.FuelType = PullTag(form.TaggedValues, "36d.")
	f.Power = PullTag(form.TaggedValues, "36e.")
	f.Meals = PullTag(form.TaggedValues, "36f.")
	f.Maintenance = PullTag(form.TaggedValues, "36g.")
	f.Water = PullTag(form.TaggedValues, "36h.")
	f.Other = PullTag(form.TaggedValues, "36i.")
	f.Instructions = PullTag(form.TaggedValues, "37.")
	f.StdFooter.PullTags(form.TaggedValues)
	return f, LeftoverTagProblems(EOC213RRTag, form.FormVersion, form.TaggedValues)
}

// Encode encodes the message content.
func (f *EOC213RR) Encode() (subject, body string) {
	var (
		sb  strings.Builder
		enc *pifo.Encoder
	)
	subject = xscsubj.Encode(f.OriginMsgID, f.Handling, EOC213RRTag, f.IncidentName)
	if f.FormVersion == "" {
		f.FormVersion = "2.3"
	}
	enc = pifo.NewEncoder(&sb, EOC213RRHTML, f.FormVersion)
	f.StdHeader.EncodeBody(enc)
	enc.Write("21.", f.IncidentName)
	enc.Write("22.", f.DateInitiated)
	enc.Write("23.", f.TimeInitiated)
	enc.Write("24.", f.TrackingNumber)
	enc.Write("25.", f.RequestedBy)
	enc.Write("26.", f.PreparedBy)
	enc.Write("27.", f.ApprovedBy)
	enc.Write("28.", f.QtyUnit)
	enc.Write("29.", f.ResourceDescription)
	enc.Write("30.", f.ResourceArrival)
	enc.Write("31.", f.Priority)
	enc.Write("32.", f.EstdCost)
	enc.Write("33.", f.DeliverTo)
	enc.Write("34.", f.DeliverToLocation)
	enc.Write("35.", f.Substitutes)
	enc.Write("36a.", f.EquipmentOperator)
	enc.Write("36b.", f.Lodging)
	enc.Write("36c.", f.Fuel)
	enc.Write("36d.", f.FuelType)
	enc.Write("36e.", f.Power)
	enc.Write("36f.", f.Meals)
	enc.Write("36g.", f.Maintenance)
	enc.Write("36h.", f.Water)
	enc.Write("36i.", f.Other)
	enc.Write("37.", f.Instructions)
	f.StdFooter.EncodeBody(enc)
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return subject, sb.String()
}
