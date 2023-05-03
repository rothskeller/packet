package pktmsg

// PlainTextTag is the tag that identifies plain text messages.
const PlainTextTag = "plain"

var plainTextType = MessageType{
	Tag:     PlainTextTag,
	Name:    "plain text message",
	Article: "a",
}

// createPlainText creates a new plain text message.
func createPlainText() Message {
	var m PlainText
	var fields = []Field{
		newPTOriginMessageNumberField(m),
		newPTHandlingField(m),
		newPTSubjectField(m),
		newPTBodyField(m),
	}
	m.BaseMessage = NewBaseMessage(plainTextType, nil, fields)
	return m
}

// recognizePlainText recognizes any message as a plain text message.
func recognizePlainText(raw, bbs, callsign, area string) Message {
	var m PlainText

}

// PlainText is a plain text message.
type PlainText struct{ BaseMessage }

func (p PlainText) Encode() string                                {}
func (p PlainText) EncodeTx() (to []string, subject, body string) {}

func newPTOriginMessageNumberField(c FieldContainer) (f Field) {
	f = BaseField(c, FOriginMsgNo, "MSGNO", "Origin Message Number", 10, 1,
		`This is the message number for this message assigned to it by its origin station.`)
	f = MessageNumberField(true, f)
	f = RequiredField(f)
	return f
}

func newPTHandlingField(c FieldContainer) (f Field) {
	f = BaseField(c, FHandling, "HANDLING", "Handling Order", 1, 1,
		`This is the handling order with which this message should be processed.`)
	f = ChoicesField(f, Value("IMMEDIATE"), Value("PRIORITY"), Value("ROUTINE"))
	f = RequiredField(f)
	return f
}

func newPTSubjectField(c FieldContainer) (f Field) {
	f = BaseField(c, FSubject, "SUBJECT", "Subject", 80, 1, `This is the subject of the message.`)
	f = RequiredField(f)
	return f
}

func newPTBodyField(c FieldContainer) (f Field) {
	f = BaseField(c, FBody, "BODY", "Body", 80, 5, `This is the body of the message.`)
	f = RequiredField(f)
	return f
}
