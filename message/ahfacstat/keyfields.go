package ahfacstat

import "github.com/rothskeller/packet/message"

// GetSubject returns the value of the subject field.
func (f *AHFacStat) GetSubject() string {
	return f.FacilityName
}

// SetSubject sets the value of the subject field.
func (f *AHFacStat) SetSubject(subject string) {
	f.FacilityName = subject
}

// SetBody sets the value of the primary text field.
func (f *AHFacStat) SetBody(body string) {
	f.Summary = body
}

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (f *AHFacStat) KeyFields() (kf *message.KeyFields) {
	kf = f.StdFields.KeyFields()
	return kf
}
