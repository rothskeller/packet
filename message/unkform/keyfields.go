package unkform

// GetHandling returns the handling order of the message.
func (f *UnknownForm) GetHandling() string {
	return f.Handling
}

// SetHandling sets the handling order of the message.
func (f *UnknownForm) SetHandling(handling string) {
	f.Handling = handling
}

// GetOriginID returns the origin message ID of the message.
func (f *UnknownForm) GetOriginID() string {
	return f.OriginMsgID
}

// SetOriginID sets the origin message ID of the message.
func (f *UnknownForm) SetOriginID(id string) {
	f.OriginMsgID = id
}

// GetSubject returns the value of the subject field.
func (f *UnknownForm) GetSubject() string {
	return f.Subject
}
