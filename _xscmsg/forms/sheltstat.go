package xscmsg

import (
	"strings"

	"github.com/rothskeller/packet/xscmsg/forms/pifo"
	"github.com/rothskeller/packet/xscmsg/forms/xscsubj"
)

// SheltStat form metadata:
const (
	SheltStatTag     = "SheltStat"
	SheltStatHTML    = "form-oa-shelter-status.html"
	SheltStatVersion = "2.2"
)

// SheltStat holds an OA shelter status form.
type SheltStat struct {
	StdHeader
	ReportType            string
	ShelterName           string
	ShelterType           string
	ShelterStatus         string
	ShelterAddress        string
	ShelterCityCode       string // added in v2.2
	ShelterCity           string
	ShelterState          string
	ShelterZip            string
	Latitude              string
	Longitude             string
	Capacity              string
	Occupancy             string
	MealsServed           string
	NSSNumber             string
	PetFriendly           string
	BasicSafetyInspection string
	ATC20Inspection       string
	AvailableServices     string
	MOU                   string
	FloorPlan             string
	ManagedByCode         string // added in v2.2
	ManagedBy             string
	ManagedByDetail       string
	PrimaryContact        string
	PrimaryPhone          string
	SecondaryContact      string
	SecondaryPhone        string
	TacticalCallSign      string
	RepeaterCallSign      string
	RepeaterInput         string
	RepeaterInputTone     string
	RepeaterOutput        string
	RepeaterOutputTone    string
	RepeaterOffset        string
	Comments              string
	RemoveFromList        string
	StdFooter
}

// DecodeSheltStat decodes the supplied form if it is a SheltStat form.  It
// returns the decoded form and strings describing any non-fatal decoding
// problems.  It returns nil, nil if the form is not a SheltStat form or has an
// unknown version.
func DecodeSheltStat(form *pifo.Form) (f *SheltStat, problems []string) {
	if form.HTMLIdent != SheltStatHTML {
		return nil, nil
	}
	switch form.FormVersion {
	case "2.0", "2.1", "2.2":
		break
	default:
		return nil, nil
	}
	f = new(SheltStat)
	f.FormVersion = form.FormVersion
	f.StdHeader.PullTags(form.TaggedValues)
	f.ReportType = PullTag(form.TaggedValues, "19.")
	f.ShelterName = PullTag(form.TaggedValues, "32.")
	f.ShelterType = PullTag(form.TaggedValues, "30.")
	f.ShelterStatus = PullTag(form.TaggedValues, "31.")
	f.ShelterAddress = PullTag(form.TaggedValues, "33a.")
	if f.FormVersion == "2.2" {
		f.ShelterCityCode = PullTag(form.TaggedValues, "33b.")
		f.ShelterCity = PullTag(form.TaggedValues, "34b.")
	} else {
		f.ShelterCity = PullTag(form.TaggedValues, "33b.")
	}
	f.ShelterState = PullTag(form.TaggedValues, "33c.")
	f.ShelterZip = PullTag(form.TaggedValues, "33d.")
	f.Latitude = PullTag(form.TaggedValues, "37a.")
	f.Longitude = PullTag(form.TaggedValues, "37b.")
	f.Capacity = PullTag(form.TaggedValues, "40a.")
	f.Occupancy = PullTag(form.TaggedValues, "40b.")
	f.MealsServed = PullTag(form.TaggedValues, "41.")
	f.NSSNumber = PullTag(form.TaggedValues, "42.")
	f.PetFriendly = PullTag(form.TaggedValues, "43a.")
	f.BasicSafetyInspection = PullTag(form.TaggedValues, "43b.")
	f.ATC20Inspection = PullTag(form.TaggedValues, "43c.")
	f.AvailableServices = PullTag(form.TaggedValues, "44.")
	f.MOU = PullTag(form.TaggedValues, "45.")
	f.FloorPlan = PullTag(form.TaggedValues, "46.")
	if f.FormVersion == "2.2" {
		f.ManagedByCode = PullTag(form.TaggedValues, "50a.")
		f.ManagedBy = PullTag(form.TaggedValues, "49a.")
	} else {
		f.ManagedBy = PullTag(form.TaggedValues, "50a.")
	}
	f.ManagedByDetail = PullTag(form.TaggedValues, "50b.")
	f.PrimaryContact = PullTag(form.TaggedValues, "51a.")
	f.PrimaryPhone = PullTag(form.TaggedValues, "51b.")
	f.SecondaryContact = PullTag(form.TaggedValues, "52a.")
	f.SecondaryPhone = PullTag(form.TaggedValues, "52b.")
	f.TacticalCallSign = PullTag(form.TaggedValues, "60.")
	f.RepeaterCallSign = PullTag(form.TaggedValues, "61.")
	f.RepeaterInput = PullTag(form.TaggedValues, "62a.")
	f.RepeaterInputTone = PullTag(form.TaggedValues, "62b.")
	f.RepeaterOutput = PullTag(form.TaggedValues, "63a.")
	f.RepeaterOutputTone = PullTag(form.TaggedValues, "63b.")
	f.RepeaterOffset = PullTag(form.TaggedValues, "62c.")
	f.Comments = PullTag(form.TaggedValues, "70.")
	f.RemoveFromList = PullTag(form.TaggedValues, "71.")
	f.StdFooter.PullTags(form.TaggedValues)
	return f, LeftoverTagProblems(SheltStatTag, form.FormVersion, form.TaggedValues)
}

// Encode encodes the message contents.
func (f *SheltStat) Encode() (subject, body string) {
	var (
		sb  strings.Builder
		enc *pifo.Encoder
	)
	subject = xscsubj.Encode(f.OriginMsgID, f.Handling, SheltStatTag, f.ShelterName)
	if f.FormVersion == "" {
		f.FormVersion = "2.2"
	}
	enc = pifo.NewEncoder(&sb, SheltStatHTML, f.FormVersion)
	f.StdHeader.EncodeBody(enc)
	enc.Write("19.", f.ReportType)
	enc.Write("32.", f.ShelterName)
	enc.Write("30.", f.ShelterType)
	enc.Write("31.", f.ShelterStatus)
	enc.Write("33a.", f.ShelterAddress)
	if f.FormVersion == "2.2" {
		enc.Write("33b.", f.ShelterCityCode)
		enc.Write("34b.", f.ShelterCity)
	} else {
		enc.Write("33b.", f.ShelterCity)
	}
	enc.Write("33c.", f.ShelterState)
	enc.Write("33d.", f.ShelterZip)
	enc.Write("37a.", f.Latitude)
	enc.Write("37b.", f.Longitude)
	enc.Write("40a.", f.Capacity)
	enc.Write("40b.", f.Occupancy)
	enc.Write("41.", f.MealsServed)
	enc.Write("42.", f.NSSNumber)
	enc.Write("43a.", f.PetFriendly)
	enc.Write("43b.", f.BasicSafetyInspection)
	enc.Write("43c.", f.ATC20Inspection)
	enc.Write("44.", f.AvailableServices)
	enc.Write("45.", f.MOU)
	enc.Write("46.", f.FloorPlan)
	if f.FormVersion == "2.2" {
		enc.Write("50a.", f.ManagedByCode)
		enc.Write("49a.", f.ManagedBy)
	} else {
		enc.Write("50a.", f.ManagedBy)
	}
	enc.Write("50b.", f.ManagedByDetail)
	enc.Write("51a.", f.PrimaryContact)
	enc.Write("51b.", f.PrimaryPhone)
	enc.Write("52a.", f.SecondaryContact)
	enc.Write("52b.", f.SecondaryPhone)
	enc.Write("60.", f.TacticalCallSign)
	enc.Write("61.", f.RepeaterCallSign)
	enc.Write("62a.", f.RepeaterInput)
	enc.Write("62b.", f.RepeaterInputTone)
	enc.Write("63a.", f.RepeaterOutput)
	enc.Write("63b.", f.RepeaterOutputTone)
	enc.Write("62c.", f.RepeaterOffset)
	enc.Write("70.", f.Comments)
	enc.Write("71.", f.RemoveFromList)
	f.StdFooter.EncodeBody(enc)
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return subject, sb.String()
}
