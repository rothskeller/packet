package ics213

import "github.com/rothskeller/packet/message"

// GetHandling returns the handling order of the message.
func (f *ICS213) GetHandling() string {
	return f.Handling
}

// SetHandling sets the handling order of the message.
func (f *ICS213) SetHandling(handling string) {
	f.Handling = handling
}

// GetOriginID returns the origin message ID of the message.
func (f *ICS213) GetOriginID() string {
	return f.OriginMsgID
}

// SetOriginID sets the origin message ID of the message.
func (f *ICS213) SetOriginID(id string) {
	f.OriginMsgID = id
}

// GetSubject returns the value of the subject field.
func (f *ICS213) GetSubject() string {
	return f.Subject
}

// SetSubject sets the value of the subject field.
func (f *ICS213) SetSubject(subject string) {
	f.Subject = subject
}

// GetBody gets the value of the primary text field.
func (f *ICS213) GetBody() string {
	return f.Message
}

// SetBody sets the value of the primary text field.
func (f *ICS213) SetBody(body string) {
	f.Message = body
}

// KeyFields returns a structure containing the values of certain key fields of
// the message that are needed for message analysis.
func (f *ICS213) KeyFields() *message.KeyFields {
	return &message.KeyFields{
		PIFOVersion:   f.PIFOVersion,
		FormVersion:   f.FormVersion,
		ToICSPosition: f.ToICSPosition,
		ToLocation:    f.ToLocation,
		OpCall:        f.OpCall,
	}
}
