package racesmar

import "github.com/rothskeller/packet/message"

// GetSubject returns the value of the subject field.
func (f *RACESMAR) GetSubject() string {
	return f.AgencyName
}

// SetSubject sets the value of the subject field.
func (f *RACESMAR) SetSubject(subject string) {
	f.AgencyName = subject
}

// SetBody sets the value of the primary text field.
func (f *RACESMAR) SetBody(body string) {
	f.Assignment = body
}

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (f *RACESMAR) KeyFields() (kf *message.KeyFields) {
	kf = f.StdFields.KeyFields()
	return kf
}
