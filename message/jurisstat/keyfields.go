package jurisstat

import "github.com/rothskeller/packet/message"

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (f *JurisStat) KeyFields() (kf *message.KeyFields) {
	kf = f.StdFields.KeyFields()
	kf.Subject = f.Jurisdiction
	kf.SubjectLabel = "Jurisdiction Name"
	return kf
}
