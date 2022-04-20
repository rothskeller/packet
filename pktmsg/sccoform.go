package pktmsg

// This file defines TxSCCoForm and RxSCCoForm.

import (
	"time"
)

// TxSCCoForm defines the common header and footer fields used by SCCo forms.
type TxSCCoForm struct {
	TxForm
	OriginMessageNumber string
	DateTime            time.Time
	ToICSPosition       string
	ToLocation          string
	ToName              string
	ToContactInfo       string
	FromICSPosition     string
	FromLocation        string
	FromName            string
	FromContactInfo     string
	RelayReceived       string
	RelaySent           string
	OperatorCallSign    string
	OperatorName        string
	OperatorDateTime    time.Time
}

func (s *TxSCCoForm) checkHeaderFooterFields() (err error) {
	if s.MessageNumber != "" {
		return ErrDontSet
	}
	if s.OriginMessageNumber == "" ||
		s.DateTime.IsZero() ||
		s.HandlingOrder == 0 ||
		s.ToICSPosition == "" ||
		s.ToLocation == "" ||
		s.FromICSPosition == "" ||
		s.FromLocation == "" ||
		s.OperatorCallSign == "" ||
		s.OperatorName == "" ||
		s.OperatorDateTime.IsZero() {
		return ErrIncomplete
	}
	return nil
}

// encodeHeaderFields encodes the header field values into the form fields.
func (s *TxSCCoForm) encodeHeaderFields() (err error) {
	s.TxForm.MessageNumber = s.OriginMessageNumber
	s.SetField("MsgNo", s.OriginMessageNumber)
	s.SetField("1a.", s.DateTime.Format("01/02/2006"))
	s.SetField("1b.", s.DateTime.Format("15:04"))
	s.SetField("5.", s.HandlingOrder.String())
	s.SetField("7a.", s.ToICSPosition)
	s.SetField("8a.", s.FromICSPosition)
	s.SetField("7b.", s.ToLocation)
	s.SetField("8b.", s.FromLocation)
	s.SetField("7c.", s.ToName)
	s.SetField("8c.", s.FromName)
	s.SetField("7d.", s.ToContactInfo)
	s.SetField("8d.", s.FromContactInfo)
	return nil
}

// encodeFooterFields encodes the footer field values into the form fields.
func (s *TxSCCoForm) encodeFooterFields() {
	s.SetField("OpRelayRcvd", s.RelayReceived)
	s.SetField("OpRelaySent", s.RelaySent)
	s.SetField("OpName", s.OperatorName)
	s.SetField("OpCall", s.OperatorCallSign)
	s.SetField("OpDate", s.OperatorDateTime.Format("01/02/2006"))
	s.SetField("OpTime", s.OperatorDateTime.Format("15:04"))
}

// boolToChecked returns "checked" if the Boolean is true and "" otherwise.
func boolToChecked(b bool) string {
	if b {
		return "checked"
	}
	return ""
}

//------------------------------------------------------------------------------

// RxSCCoForm defines the common header and footer fields used by SCCo forms.
type RxSCCoForm struct {
	RxForm
	OriginMessageNumber string
	Date                string
	Time                string
	DateTime            time.Time
	HandlingOrder       HandlingOrder
	ToICSPosition       string
	ToLocation          string
	ToName              string
	ToContactInfo       string
	FromICSPosition     string
	FromLocation        string
	FromName            string
	FromContactInfo     string
	RelayReceived       string
	RelaySent           string
	OperatorCallSign    string
	OperatorName        string
	OperatorDate        string
	OperatorTime        string
	OperatorDateTime    time.Time
}

// SCCoForm returns a pointer to the RxSCCoForm portion of a message object.  It
// can be used to reach fields of the RxSCCoForm object that are occluded by
// types that embed RxSCCoForm.
func (s *RxSCCoForm) SCCoForm() *RxSCCoForm { return s }

// Valid returns whether all of the fields of the form have valid values, and
// all required fields are filled in.
func (s *RxSCCoForm) Valid() bool {
	return (s.OriginMessageNumber != "" &&
		!s.DateTime.IsZero() &&
		s.HandlingOrder != 0 &&
		s.ToICSPosition != "" &&
		s.ToLocation != "" &&
		s.FromICSPosition != "" &&
		s.FromLocation != "" &&
		s.OperatorCallSign != "" &&
		s.OperatorName != "" &&
		!s.OperatorDateTime.IsZero())
}

// extractHeaderFields extracts the header field values from the form.
func (s *RxSCCoForm) extractHeaderFields() {
	s.OriginMessageNumber = s.Fields["MsgNo"]
	s.Date = s.Fields["1a."]
	s.Time = s.Fields["1b."]
	s.DateTime = dateTimeParse(s.Date, s.Time)
	s.HandlingOrder, _ = ParseHandlingOrder(s.Fields["5."])
	s.ToICSPosition = s.Fields["7a."]
	s.ToLocation = s.Fields["7b."]
	s.ToName = s.Fields["7c."]
	s.ToContactInfo = s.Fields["7d."]
	s.FromICSPosition = s.Fields["8a."]
	s.FromLocation = s.Fields["8b."]
	s.FromName = s.Fields["8c."]
	s.FromContactInfo = s.Fields["8d."]
}

// extractFooterFields extracts the footer field values from the form.
func (s *RxSCCoForm) extractFooterFields() {
	s.RelayReceived = s.Fields["OpRelayRcvd"]
	s.RelaySent = s.Fields["OpRelaySent"]
	s.OperatorCallSign = s.Fields["OpCall"]
	s.OperatorName = s.Fields["OpName"]
	s.OperatorDate = s.Fields["OpDate"]
	s.OperatorTime = s.Fields["OpTime"]
	s.OperatorDateTime = dateTimeParse(s.OperatorDate, s.OperatorTime)
}

// dateTimeParse attempts to parse a date and time.
func dateTimeParse(d, t string) (dt time.Time) {
	var err error

	if dt, err = time.ParseInLocation("1/2/2006 15:04", d+" "+t, time.Local); err == nil {
		return dt
	}
	dt, _ = time.ParseInLocation("1/2/2006 1504", d+" "+t, time.Local)
	return dt
}
