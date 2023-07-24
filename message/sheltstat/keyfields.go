package sheltstat

import "github.com/rothskeller/packet/message"

// GetSubject returns the value of the subject field.
func (f *SheltStat) GetSubject() string {
	return f.ShelterName
}

// SetSubject sets the value of the subject field.
func (f *SheltStat) SetSubject(subject string) {
	f.ShelterName = subject
}

// SetBody sets the value of the primary text field.
func (f *SheltStat) SetBody(body string) {
	f.Comments = body
}

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (f *SheltStat) KeyFields() (kf *message.KeyFields) {
	kf = f.StdFields.KeyFields()
	return kf
}
