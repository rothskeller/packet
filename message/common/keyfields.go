package common

import "github.com/rothskeller/packet/message"

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (s *StdFields) KeyFields() *message.KeyFields {
	return &message.KeyFields{
		PIFOVersion:   s.PIFOVersion,
		FormVersion:   s.FormVersion,
		OriginMsgID:   s.OriginMsgID,
		Handling:      s.Handling,
		ToICSPosition: s.ToICSPosition,
		ToLocation:    s.ToLocation,
		OpCall:        s.OpCall,
	}
}
