package jurisstat

import "github.com/rothskeller/packet/message"

// GetSubject returns the value of the subject field.
func (f *JurisStat) GetSubject() string {
	return f.Jurisdiction
}

// SetSubject sets the value of the subject field.
func (f *JurisStat) SetSubject(subject string) {
	f.Jurisdiction = subject
}

// SetBody sets the value of the primary text field.
func (f *JurisStat) SetBody(body string) {
	f.CommunicationsComments = body // not a good choice, but there isn't a better one
}

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (f *JurisStat) KeyFields() (kf *message.KeyFields) {
	kf = f.StdFields.KeyFields()
	return kf
}
