package plaintext

// GetHandling returns the handling order of the message.
func (m *PlainText) GetHandling() string {
	return m.Handling
}

// SetHandling sets the handling order of the message.
func (m *PlainText) SetHandling(handling string) {
	m.Handling = handling
}

// GetOriginID returns the origin message ID of the message.
func (m *PlainText) GetOriginID() string {
	return m.OriginMsgID
}

// SetOriginID sets the origin message ID of the message.
func (m *PlainText) SetOriginID(id string) {
	m.OriginMsgID = id
}

// GetSubject returns the value of the subject field.
func (m *PlainText) GetSubject() string {
	return m.Subject
}

// SetSubject sets the value of the subject field.
func (m *PlainText) SetSubject(subject string) {
	m.Subject = subject
}

// GetBody gets the value of the primary text field.
func (m *PlainText) GetBody() string {
	return m.Body
}

// SetBody sets the value of the primary text field.
func (m *PlainText) SetBody(body string) {
	m.Body = body
}
