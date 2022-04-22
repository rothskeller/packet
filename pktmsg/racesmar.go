package pktmsg

// This file defines TxRACESMARForm and RxRACESMARForm.

import (
	"fmt"
	"strconv"
)

// RACESMARFields contains the form-specific fields of the RACES Mutual Aid
// Request form.
type RACESMARFields struct {
	Agency                string
	EventName             string
	EventNumber           string
	Assignment            string
	ResourcesRequested    []*RACESResource
	ArrivalDates          string
	ArrivalTimes          string
	NeededDates           string
	NeededTimes           string
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
}

// RACESResource describes one requested RACES resource.
type RACESResource struct {
	ResourcesQty  int
	ResourcesRole string
	PreferredType string
	MinimumType   string
}

// A TxRACESMARForm is an outgoing PackItForms-encoded message containing an
// SCCo RACES Mutual Aid Request form.
type TxRACESMARForm struct {
	TxSCCoForm
	RACESMARFields
}

// Encode returns the encoded subject line and body of the message.
func (mar *TxRACESMARForm) Encode() (subject, body string, err error) {
	if err = mar.checkHeaderFooterFields(); err != nil {
		return "", "", err
	}
	if mar.Subject != "" {
		return "", "", ErrDontSet
	}
	if mar.Agency == "" ||
		mar.EventName == "" ||
		mar.Assignment == "" ||
		len(mar.ResourcesRequested) == 0 ||
		mar.ResourcesRequested[0].ResourcesQty == 0 ||
		mar.ResourcesRequested[0].ResourcesRole == "" ||
		mar.ResourcesRequested[0].PreferredType == "" ||
		mar.ResourcesRequested[0].MinimumType == "" ||
		mar.ArrivalDates == "" ||
		mar.ArrivalTimes == "" ||
		mar.NeededDates == "" ||
		mar.NeededTimes == "" ||
		mar.ReportingLocation == "" ||
		mar.ContactOnArrival == "" ||
		mar.RequesterName == "" ||
		mar.RequesterTitle == "" ||
		mar.RequesterContact == "" ||
		mar.AgencyApproverName == "" ||
		mar.AgencyApproverTitle == "" ||
		mar.AgencyApproverContact == "" ||
		mar.AgencyApprovedDate == "" ||
		mar.AgencyApprovedTime == "" {
		return "", "", ErrIncomplete
	}
	if len(mar.ResourcesRequested) > 5 {
		return "", "", ErrInvalid
	}
	mar.FormName = "RACES-MAR"
	mar.FormHTML = "form-oa-mutual-aid-request-v2.html"
	mar.FormVersion = "2.1"
	mar.Subject = mar.Agency
	mar.encodeHeaderFields()
	mar.SetField("15.", mar.Agency)
	mar.SetField("16a.", mar.EventName)
	mar.SetField("16b.", mar.EventNumber)
	mar.SetField("17.", mar.Assignment)
	for i := 0; i < len(mar.ResourcesRequested); i++ {
		mar.SetField(fmt.Sprintf("18.%da.", i+1), intToText(mar.ResourcesRequested[i].ResourcesQty))
		mar.SetField(fmt.Sprintf("18.%db.", i+1), mar.ResourcesRequested[i].ResourcesRole)
		mar.SetField(fmt.Sprintf("18.%dc.", i+1), mar.ResourcesRequested[i].PreferredType)
		mar.SetField(fmt.Sprintf("18.%dd.", i+1), mar.ResourcesRequested[i].MinimumType)
	}
	for i := len(mar.ResourcesRequested); i < 5; i++ {
		mar.SetField(fmt.Sprintf("18.%da.", i+1), "")
		mar.SetField(fmt.Sprintf("18.%db.", i+1), "")
		mar.SetField(fmt.Sprintf("18.%dc.", i+1), "")
		mar.SetField(fmt.Sprintf("18.%dd.", i+1), "")
	}
	mar.SetField("19a.", mar.ArrivalDates)
	mar.SetField("19b.", mar.ArrivalTimes)
	mar.SetField("20a.", mar.NeededDates)
	mar.SetField("20b.", mar.NeededTimes)
	mar.SetField("21.", mar.ReportingLocation)
	mar.SetField("22.", mar.ContactOnArrival)
	mar.SetField("23.", mar.TravelInfo)
	mar.SetField("24a.", mar.RequesterName)
	mar.SetField("24b.", mar.RequesterTitle)
	mar.SetField("24c.", mar.RequesterContact)
	mar.SetField("25a.", mar.AgencyApproverName)
	mar.SetField("25b.", mar.AgencyApproverTitle)
	mar.SetField("25c.", mar.AgencyApproverContact)
	mar.SetField("26a.", mar.AgencyApprovedDate)
	mar.SetField("26b.", mar.AgencyApprovedTime)
	mar.encodeFooterFields()
	return mar.TxSCCoForm.Encode()
}

//------------------------------------------------------------------------------

// An RxRACESMARForm is a received PackItForms-encoded message containing an
// SCCo OA Shelter Status form.
type RxRACESMARForm struct {
	RxSCCoForm
	RACESMARFields
}

// parseRxRACESMARForm examines an RxForm to see if it contains an EOC-213RR
// form, and if so, wraps it in an RxRACESMARForm and returns it.  If it is not,
// it returns nil.
func parseRxRACESMARForm(f *RxForm) *RxRACESMARForm {
	var mar RxRACESMARForm

	if f.FormHTML != "form-oa-mutual-aid-request.html" && f.FormHTML != "form-oa-mutual-aid-request-v2.html" {
		return nil
	}
	mar.RxSCCoForm.RxForm = *f
	mar.extractHeaderFields()
	mar.Agency = mar.Fields["15."]
	mar.EventName = mar.Fields["16a."]
	mar.EventNumber = mar.Fields["16b."]
	mar.Assignment = mar.Fields["17."]
	{ // form v1
		var res RACESResource

		res.ResourcesQty, _ = strconv.Atoi(mar.Fields["18a."])
		res.ResourcesRole = mar.Fields["18b."]
		res.PreferredType = mar.Fields["18c."]
		res.MinimumType = mar.Fields["18d."]
		if res.ResourcesQty != 0 || res.ResourcesRole != "" || res.PreferredType != "" || res.MinimumType != "" {
			mar.ResourcesRequested = append(mar.ResourcesRequested, &res)
		}
	}
	for i := 0; i < 5; i++ { // form v2
		var res RACESResource

		res.ResourcesQty, _ = strconv.Atoi(mar.Fields[fmt.Sprintf("18.%da.", i+1)])
		res.ResourcesRole = mar.Fields[fmt.Sprintf("18.%db.", i+1)]
		res.PreferredType = mar.Fields[fmt.Sprintf("18.%dc.", i+1)]
		res.MinimumType = mar.Fields[fmt.Sprintf("18.%dd.", i+1)]
		if res.ResourcesQty != 0 || res.ResourcesRole != "" || res.PreferredType != "" || res.MinimumType != "" {
			mar.ResourcesRequested = append(mar.ResourcesRequested, &res)
		}
	}
	mar.ArrivalDates = mar.Fields["19a."]
	mar.ArrivalTimes = mar.Fields["19b."]
	mar.NeededDates = mar.Fields["20a."]
	mar.NeededTimes = mar.Fields["20b."]
	mar.ReportingLocation = mar.Fields["21."]
	mar.ContactOnArrival = mar.Fields["22."]
	mar.TravelInfo = mar.Fields["23."]
	mar.RequesterName = mar.Fields["24a."]
	mar.RequesterTitle = mar.Fields["24b."]
	mar.RequesterContact = mar.Fields["24c."]
	mar.AgencyApproverName = mar.Fields["25a."]
	mar.AgencyApproverTitle = mar.Fields["25b."]
	mar.AgencyApproverContact = mar.Fields["25c."]
	mar.AgencyApprovedDate = mar.Fields["26a."]
	mar.AgencyApprovedTime = mar.Fields["26b."]
	mar.extractFooterFields()
	return &mar
}

// Valid returns whether all of the fields of the form have valid values, and
// all required fields are filled in.
func (mar *RxRACESMARForm) Valid() bool {
	return mar.RxSCCoForm.Valid() &&
		mar.Agency != "" &&
		mar.EventName != "" &&
		mar.Assignment != "" &&
		len(mar.ResourcesRequested) != 0 &&
		mar.ResourcesRequested[0].ResourcesQty != 0 &&
		mar.ResourcesRequested[0].ResourcesRole != "" &&
		mar.ResourcesRequested[0].PreferredType != "" &&
		mar.ResourcesRequested[0].MinimumType != "" &&
		mar.ArrivalDates != "" &&
		mar.ArrivalTimes != "" &&
		mar.NeededDates != "" &&
		mar.NeededTimes != "" &&
		mar.ReportingLocation != "" &&
		mar.ContactOnArrival != "" &&
		mar.RequesterName != "" &&
		mar.RequesterTitle != "" &&
		mar.RequesterContact != "" &&
		mar.AgencyApproverName != "" &&
		mar.AgencyApproverTitle != "" &&
		mar.AgencyApproverContact != "" &&
		mar.AgencyApprovedDate != "" &&
		mar.AgencyApprovedTime != ""
}

// EncodeSubjectLine returns what the subject line should be based on the
// received form contents.
func (mar *RxRACESMARForm) EncodeSubjectLine() string {
	return fmt.Sprintf("%s_%s_RACES-MAR_%s", mar.OriginMessageNumber, mar.HandlingOrder.Code(), mar.Agency)
}

// TypeCode returns the machine-readable code for the message type.
func (*RxRACESMARForm) TypeCode() string { return "RACES-MAR" }

// TypeName returns the human-reading name of the message type.
func (*RxRACESMARForm) TypeName() string { return "RACES Mutual Aid Request form" }

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (*RxRACESMARForm) TypeArticle() string { return "a" }
