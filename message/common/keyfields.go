package common

import "github.com/rothskeller/packet/message"

// GetHandling returns the handling order of the message.
func (s *StdFields) GetHandling() string {
	return s.Handling
}

// SetHandling sets the handling order of the message.
func (s *StdFields) SetHandling(handling string) {
	s.Handling = handling
}

// GetOriginID returns the origin message ID of the message.
func (s *StdFields) GetOriginID() string {
	return s.OriginMsgID
}

// SetOriginID sets the origin message ID of the message.
func (s *StdFields) SetOriginID(id string) {
	s.OriginMsgID = id
}

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (s *StdFields) KeyFields() *message.KeyFields {
	return &message.KeyFields{
		PIFOVersion:   s.PIFOVersion,
		FormVersion:   s.FormVersion,
		ToICSPosition: s.ToICSPosition,
		ToLocation:    s.ToLocation,
		OpCall:        s.OpCall,
	}
}
