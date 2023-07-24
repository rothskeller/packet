package eoc213rr

import "github.com/rothskeller/packet/message"

// GetSubject returns the value of the subject field.
func (f *EOC213RR) GetSubject() string {
	return f.IncidentName
}

// SetSubject sets the value of the subject field.
func (f *EOC213RR) SetSubject(subject string) {
	f.IncidentName = subject
}

// SetBody sets the value of the primary text field.
func (f *EOC213RR) SetBody(body string) {
	f.Instructions = body
}

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (f *EOC213RR) KeyFields() (kf *message.KeyFields) {
	kf = f.StdFields.KeyFields()
	return kf
}
