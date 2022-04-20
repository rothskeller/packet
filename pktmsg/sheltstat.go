package pktmsg

// This file defines TxSheltStatForm and RxSheltStatForm.

import (
	"fmt"
	"strconv"
)

// A TxSheltStatForm is an outgoing PackItForss-encoded message containing an
// SCCo OA Shelter Status form.
type TxSheltStatForm struct {
	TxSCCoForm
	ReportType           string
	ShelterType          string
	ShelterStatus        string
	ShelterName          string
	ShelterAddress       string
	ShelterCity          string
	ShelterState         string
	ShelterZip           string
	Latitude             float64
	Longitude            float64
	Capacity             int
	Occupancy            int
	Meals                string
	NSS                  string
	PetFriendly          string
	BasicSafety          string
	ATC20                string
	AvailableServices    string
	MOU                  string
	FloorPlan            string
	ManagedBy            string
	ManagedByDetail      string
	PrimaryContact       string
	PrimaryPhone         string
	SecondaryContact     string
	SecondaryPhone       string
	TacticalCall         string
	RepeaterCall         string
	RepeaterInput        string
	RepeaterInputTone    string
	RepeaterOutput       string
	RepeaterOutputTone   string
	RepeaterOffset       string
	Comments             string
	RemoveFromActiveList bool
}

var (
	// validReportType is defined in munistat.go
	validShelterType   = map[string]bool{"": true, "Type 1": true, "Type 2": true, "Type 3": true, "Type 4": true}
	validShelterStatus = map[string]bool{"": true, "Open": true, "Closed": true, "Full": true}
	validPetFriendly   = map[string]bool{"": true, "Yes": true, "No": true}
	validBasicSafety   = map[string]bool{"": true, "Yes": true, "No": true}
	validATC20         = map[string]bool{"": true, "Yes": true, "No": true}
	validManagedBy     = map[string]bool{"": true, "American Red Cross": true, "Private": true, "Community": true, "Government": true, "Other": true}
)

// Encode returns the encoded subject line and body of the message.
func (ss *TxSheltStatForm) Encode() (subject, body string, err error) {
	if err = ss.checkHeaderFooterFields(); err != nil {
		return "", "", err
	}
	if ss.Subject != "" {
		return "", "", ErrDontSet
	}
	if ss.ReportType == "" || ss.ShelterName == "" {
		return "", "", ErrIncomplete
	}
	if !validReportType[ss.ReportType] ||
		!validShelterType[ss.ShelterType] ||
		!validShelterStatus[ss.ShelterStatus] ||
		!validPetFriendly[ss.PetFriendly] ||
		!validBasicSafety[ss.BasicSafety] ||
		!validATC20[ss.ATC20] ||
		!validManagedBy[ss.ManagedBy] {
		return "", "", ErrInvalid
	}
	ss.FormName = "SheltStat"
	ss.FormHTML = "form-oa-shelter-status.html"
	ss.FormVersion = "2.1"
	ss.Subject = ss.ShelterName
	ss.encodeHeaderFields()
	ss.SetField("19.", ss.ReportType)
	ss.SetField("30.", ss.ShelterType)
	ss.SetField("31.", ss.ShelterStatus)
	ss.SetField("32.", ss.ShelterName)
	ss.SetField("33a.", ss.ShelterAddress)
	ss.SetField("33b.", ss.ShelterCity)
	ss.SetField("33c.", ss.ShelterState)
	ss.SetField("33d.", ss.ShelterZip)
	ss.SetField("37a.", strconv.FormatFloat(ss.Latitude, 'f', -1, 64))
	ss.SetField("37b.", strconv.FormatFloat(ss.Longitude, 'f', -1, 64))
	if ss.Capacity != 0 {
		ss.SetField("40a.", strconv.Itoa(ss.Capacity))
	} else {
		ss.SetField("40a.", "")
	}
	if ss.Occupancy != 0 {
		ss.SetField("40b.", strconv.Itoa(ss.Occupancy))
	} else {
		ss.SetField("40b.", "")
	}
	ss.SetField("41.", ss.Meals)
	ss.SetField("42.", ss.NSS)
	ss.SetField("43a.", ss.PetFriendly)
	ss.SetField("43b.", ss.BasicSafety)
	ss.SetField("43c.", ss.ATC20)
	ss.SetField("44.", ss.AvailableServices)
	ss.SetField("45.", ss.MOU)
	ss.SetField("46.", ss.FloorPlan)
	ss.SetField("50a.", ss.ManagedBy)
	ss.SetField("50b.", ss.ManagedByDetail)
	ss.SetField("51a.", ss.PrimaryContact)
	ss.SetField("51b.", ss.PrimaryPhone)
	ss.SetField("52a.", ss.SecondaryContact)
	ss.SetField("52b.", ss.SecondaryPhone)
	ss.SetField("60.", ss.TacticalCall)
	ss.SetField("61.", ss.RepeaterCall)
	ss.SetField("62a.", ss.RepeaterInput)
	ss.SetField("62b.", ss.RepeaterInputTone)
	ss.SetField("63a.", ss.RepeaterOutput)
	ss.SetField("63b.", ss.RepeaterOutputTone)
	ss.SetField("62c.", ss.RepeaterOffset)
	ss.SetField("70.", ss.Comments)
	ss.SetField("71.", boolToChecked(ss.RemoveFromActiveList))
	ss.encodeFooterFields()
	return ss.TxSCCoForm.Encode()
}

//------------------------------------------------------------------------------

// An RxSheltStatForm is a received PackItForss-encoded message containing an
// SCCo EOC-213RR form.
type RxSheltStatForm struct {
	RxSCCoForm
	ReportType           string
	ShelterType          string
	ShelterStatus        string
	ShelterName          string
	ShelterAddress       string
	ShelterCity          string
	ShelterState         string
	ShelterZip           string
	Latitude             float64
	Longitude            float64
	Capacity             int
	Occupancy            int
	Meals                string
	NSS                  string
	PetFriendly          string
	BasicSafety          string
	ATC20                string
	AvailableServices    string
	MOU                  string
	FloorPlan            string
	ManagedBy            string
	ManagedByDetail      string
	PrimaryContact       string
	PrimaryPhone         string
	SecondaryContact     string
	SecondaryPhone       string
	TacticalCall         string
	RepeaterCall         string
	RepeaterInput        string
	RepeaterInputTone    string
	RepeaterOutput       string
	RepeaterOutputTone   string
	RepeaterOffset       string
	Comments             string
	RemoveFromActiveList bool
}

// parseRxSheltStatForm examines an RxForm to see if it contains an EOC-213RR
// form, and if so, wraps it in an RxSheltStatForm and returns it.  If it is not,
// it returns nil.
func parseRxSheltStatForm(f *RxForm) *RxSheltStatForm {
	var ss RxSheltStatForm

	if f.FormHTML != "form-oa-shelter-status.html" {
		return nil
	}
	ss.RxSCCoForm.RxForm = *f
	ss.extractHeaderFields()
	ss.ReportType = ss.Fields["19."]
	ss.ShelterType = ss.Fields["30."]
	ss.ShelterStatus = ss.Fields["31."]
	ss.ShelterName = ss.Fields["32."]
	ss.ShelterAddress = ss.Fields["33a."]
	ss.ShelterCity = ss.Fields["33b."]
	ss.ShelterState = ss.Fields["33c."]
	ss.ShelterZip = ss.Fields["33d."]
	ss.Latitude, _ = strconv.ParseFloat(ss.Fields["37a."], 64)
	ss.Longitude, _ = strconv.ParseFloat(ss.Fields["37b."], 64)
	ss.Capacity, _ = strconv.Atoi(ss.Fields["40a."])
	ss.Occupancy, _ = strconv.Atoi(ss.Fields["40b."])
	ss.Meals = ss.Fields["41."]
	ss.NSS = ss.Fields["42."]
	ss.PetFriendly = ss.Fields["43a."]
	ss.BasicSafety = ss.Fields["43b."]
	ss.ATC20 = ss.Fields["43c."]
	ss.AvailableServices = ss.Fields["44."]
	ss.MOU = ss.Fields["45."]
	ss.FloorPlan = ss.Fields["46."]
	ss.ManagedBy = ss.Fields["50a."]
	ss.ManagedByDetail = ss.Fields["50b."]
	ss.PrimaryContact = ss.Fields["51a."]
	ss.PrimaryPhone = ss.Fields["51b."]
	ss.SecondaryContact = ss.Fields["52a."]
	ss.SecondaryPhone = ss.Fields["52b."]
	ss.TacticalCall = ss.Fields["60."]
	ss.RepeaterCall = ss.Fields["61."]
	ss.RepeaterInput = ss.Fields["62a."]
	ss.RepeaterInputTone = ss.Fields["62b."]
	ss.RepeaterOutput = ss.Fields["63a."]
	ss.RepeaterOutputTone = ss.Fields["63b."]
	ss.RepeaterOffset = ss.Fields["62c."]
	ss.Comments = ss.Fields["70."]
	ss.RemoveFromActiveList = ss.Fields["71."] != ""
	ss.extractFooterFields()
	return &ss
}

// Valid returns whether all of the fields of the form have valid values, and
// all required fields are filled in.
func (ss *RxSheltStatForm) Valid() bool {
	return ss.RxSCCoForm.Valid() &&
		validReportType[ss.ReportType] &&
		ss.ShelterName != "" &&
		validShelterType[ss.ShelterType] &&
		validShelterStatus[ss.ShelterStatus] &&
		validPetFriendly[ss.PetFriendly] &&
		validBasicSafety[ss.BasicSafety] &&
		validATC20[ss.ATC20] &&
		validManagedBy[ss.ManagedBy] &&
		(ss.ReportType == "Update" ||
			(ss.ShelterType != "" &&
				ss.ShelterStatus != "" &&
				ss.ShelterAddress != "" &&
				ss.ShelterCity != "" &&
				ss.ShelterState != "" &&
				ss.Capacity != 0 &&
				ss.Occupancy != 0 &&
				ss.ManagedBy != "" &&
				ss.PrimaryContact != "" &&
				ss.PrimaryPhone != ""))
}

// EncodeSubjectLine returns what the subject line should be based on the
// received form contents.
func (ss *RxSheltStatForm) EncodeSubjectLine() string {
	return fmt.Sprintf("%s_%s_SheltStat_%s", ss.OriginMessageNumber, ss.HandlingOrder.Code(), ss.ShelterName)
}

// TypeCode returns the machine-readable code for the message type.
func (*RxSheltStatForm) TypeCode() string { return "SheltStat" }

// TypeName returns the human-reading name of the message type.
func (*RxSheltStatForm) TypeName() string { return "OA Shelter Status form" }

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (*RxSheltStatForm) TypeArticle() string { return "an" }
