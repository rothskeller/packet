package checkin

// GetHandling returns the handling order of the message.
func (m *CheckIn) GetHandling() string {
	return m.Handling
}

// SetHandling sets the handling order of the message.
func (m *CheckIn) SetHandling(handling string) {
	m.Handling = handling
}

// GetOriginID returns the origin message ID of the message.
func (m *CheckIn) GetOriginID() string {
	return m.OriginMsgID
}

// SetOriginID sets the origin message ID of the message.
func (m *CheckIn) SetOriginID(id string) {
	m.OriginMsgID = id
}

// GetSubject returns the value of the subject field.
func (m *CheckIn) GetSubject() string {
	return m.EncodeSubject()
}
