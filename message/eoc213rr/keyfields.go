package eoc213rr

import "github.com/rothskeller/packet/message"

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (f *EOC213RR) KeyFields() (kf *message.KeyFields) {
	kf = f.StdFields.KeyFields()
	kf.Subject = f.IncidentName
	kf.SubjectLabel = "Incident Name"
	return kf
}
