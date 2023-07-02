package xscmsg

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/xscmsg/forms/pifo"
	"github.com/rothskeller/packet/xscmsg/forms/xscsubj"
)

// RACES-MAR form metadata:
const (
	RACESMARTag     = "RACES-MAR"
	RACESMARHTML    = "form-oa-mutual-aid-request-v2.html"
	racesMARHTML16  = "form-oa-mutual-aid-request.html"
	RACESMARVersion = "2.3"
)

// RACESMAR holds a RACES mutual aid request form.
type RACESMAR struct {
	StdHeader
	AgencyName            string
	EventName             string
	EventNumber           string
	Assignment            string
	Resources             [5]Resource
	ArrivalDates          string
	ArrivalTimes          string
	NeededUntilDates      string
	NeededUntilTimes      string
	ReportingLocation     string
	ContactOnArrival      string
	TravelInfo            string
	RequesterName         string
	RequesterTitle        string
	RequesterContact      string
	AgencyApproverName    string
	AgencyApproverTitle   string
	AgencyApproverContact string
	AgencyApprovedDate    string
	AgencyApprovedTime    string
	StdFooter
}

// A Resource is the description of a single resource in a RACES mutual aid
// request form.
type Resource struct {
	Qty           string
	Role          string
	Position      string // Added in v2.3
	RolePos       string // Added in v2.3
	PreferredType string
	MinimumType   string
}

// DecodeRACESMAR decodes the supplied form if it is a RACES-MAR form.  It
// returns the decoded form and strings describing any non-fatal decoding
// problems.  It returns nil, nil if the form is not a RACES-MAR form or has an
// unknown version.
func DecodeRACESMAR(form *pifo.Form) (f *RACESMAR, problems []string) {
	switch {
	case form.HTMLIdent == racesMARHTML16 && form.FormVersion == "1.6":
		break
	case form.HTMLIdent == RACESMARHTML && (form.FormVersion == "2.1" || form.FormVersion == "2.3"):
		break
	default:
		return nil, nil
	}
	f = new(RACESMAR)
	f.FormVersion = form.FormVersion
	f.StdHeader.PullTags(form.TaggedValues)
	f.AgencyName = PullTag(form.TaggedValues, "15.")
	f.EventName = PullTag(form.TaggedValues, "16a.")
	f.EventNumber = PullTag(form.TaggedValues, "16b.")
	f.Assignment = PullTag(form.TaggedValues, "17.")
	if f.FormVersion == "1.6" {
		f.Resources[0].Qty = PullTag(form.TaggedValues, "18a.")
		f.Resources[0].RolePos = PullTag(form.TaggedValues, "18b.")
		f.Resources[0].PreferredType = PullTag(form.TaggedValues, "18c.")
		f.Resources[0].MinimumType = PullTag(form.TaggedValues, "18d.")
	} else {
		for i := range f.Resources {
			f.Resources[i].Qty = PullTag(form.TaggedValues, fmt.Sprintf("18.%da.", i+1))
			f.Resources[i].RolePos = PullTag(form.TaggedValues, fmt.Sprintf("18.%db.", i+1))
			f.Resources[i].PreferredType = PullTag(form.TaggedValues, fmt.Sprintf("18.%dc.", i+1))
			f.Resources[i].MinimumType = PullTag(form.TaggedValues, fmt.Sprintf("18.%dd.", i+1))
			if f.FormVersion == "2.3" {
				f.Resources[i].Role = PullTag(form.TaggedValues, fmt.Sprintf("18.%de.", i+1))
				f.Resources[i].Position = PullTag(form.TaggedValues, fmt.Sprintf("18.%df.", i+1))
			}
		}
	}
	f.ArrivalDates = PullTag(form.TaggedValues, "19a.")
	f.ArrivalTimes = PullTag(form.TaggedValues, "19b.")
	f.NeededUntilDates = PullTag(form.TaggedValues, "20a.")
	f.NeededUntilTimes = PullTag(form.TaggedValues, "20b.")
	f.ReportingLocation = PullTag(form.TaggedValues, "21.")
	f.ContactOnArrival = PullTag(form.TaggedValues, "22.")
	f.TravelInfo = PullTag(form.TaggedValues, "23.")
	f.RequesterName = PullTag(form.TaggedValues, "24a.")
	f.RequesterTitle = PullTag(form.TaggedValues, "24b.")
	f.RequesterContact = PullTag(form.TaggedValues, "24c.")
	f.AgencyApproverName = PullTag(form.TaggedValues, "25a.")
	f.AgencyApproverTitle = PullTag(form.TaggedValues, "25b.")
	f.AgencyApproverContact = PullTag(form.TaggedValues, "25c.")
	f.AgencyApprovedDate = PullTag(form.TaggedValues, "26a.")
	f.AgencyApprovedTime = PullTag(form.TaggedValues, "26b.")
	f.StdFooter.PullTags(form.TaggedValues)
	return f, LeftoverTagProblems(RACESMARTag, form.FormVersion, form.TaggedValues)
}

// Encode encodes the message contents.
func (f *RACESMAR) Encode() (subject, body string) {
	var (
		sb  strings.Builder
		enc *pifo.Encoder
	)
	subject = xscsubj.Encode(f.OriginMsgID, f.Handling, RACESMARTag, f.AgencyName)
	if f.FormVersion == "" {
		f.FormVersion = "2.3"
	}
	if f.FormVersion == "1.6" {
		enc = pifo.NewEncoder(&sb, racesMARHTML16, f.FormVersion)
	} else {
		enc = pifo.NewEncoder(&sb, RACESMARHTML, f.FormVersion)
	}
	f.StdHeader.EncodeBody(enc)
	enc.Write("15.", f.AgencyName)
	enc.Write("16a.", f.EventName)
	enc.Write("16b.", f.EventNumber)
	enc.Write("17.", f.Assignment)
	if f.FormVersion == "1.6" {
		enc.Write("18a.", f.Resources[0].Qty)
		enc.Write("18b.", f.Resources[0].RolePos)
		enc.Write("18c.", f.Resources[0].PreferredType)
		enc.Write("18d.", f.Resources[0].MinimumType)
	} else {
		for i, r := range f.Resources {
			enc.Write(fmt.Sprintf("18.%da.", i+1), r.Qty)
			enc.Write(fmt.Sprintf("18.%db.", i+1), r.RolePos)
			enc.Write(fmt.Sprintf("18.%dc.", i+1), r.PreferredType)
			enc.Write(fmt.Sprintf("18.%dd.", i+1), r.MinimumType)
			if f.FormVersion == "2.3" {
				enc.Write(fmt.Sprintf("18.%de.", i+1), r.Role)
				enc.Write(fmt.Sprintf("18.%df.", i+1), r.Position)
			}
		}
	}
	enc.Write("19a.", f.ArrivalDates)
	enc.Write("19b.", f.ArrivalTimes)
	enc.Write("20a.", f.NeededUntilDates)
	enc.Write("20b.", f.NeededUntilTimes)
	enc.Write("21.", f.ReportingLocation)
	enc.Write("22.", f.ContactOnArrival)
	enc.Write("23.", f.TravelInfo)
	enc.Write("24a.", f.RequesterName)
	enc.Write("24b.", f.RequesterTitle)
	enc.Write("24c.", f.RequesterContact)
	enc.Write("25a.", f.AgencyApproverName)
	enc.Write("25b.", f.AgencyApproverTitle)
	enc.Write("25c.", f.AgencyApproverContact)
	enc.Write("26a.", f.AgencyApprovedDate)
	enc.Write("26b.", f.AgencyApprovedTime)
	f.StdFooter.EncodeBody(enc)
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return subject, sb.String()
}
