package ics213

import "github.com/rothskeller/packet/message"

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (f *ICS213) KeyFields() *message.KeyFields {
	return &message.KeyFields{
		PIFOVersion:   f.PIFOVersion,
		FormVersion:   f.FormVersion,
		OriginMsgID:   f.OriginMsgID,
		Handling:      f.Handling,
		ToICSPosition: f.ToICSPosition,
		ToLocation:    f.ToLocation,
		Subject:       f.Subject,
		SubjectLabel:  "Subject",
		OpCall:        f.OpCall,
	}
}
