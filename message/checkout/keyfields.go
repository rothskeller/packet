package checkout

// GetHandling returns the handling order of the message.
func (m *CheckOut) GetHandling() string {
	return m.Handling
}

// SetHandling sets the handling order of the message.
func (m *CheckOut) SetHandling(handling string) {
	m.Handling = handling
}

// GetOriginID returns the origin message ID of the message.
func (m *CheckOut) GetOriginID() string {
	return m.OriginMsgID
}

// SetOriginID sets the origin message ID of the message.
func (m *CheckOut) SetOriginID(id string) {
	m.OriginMsgID = id
}

// GetSubject returns the value of the subject field.
func (m *CheckOut) GetSubject() string {
	return m.EncodeSubject()
}
