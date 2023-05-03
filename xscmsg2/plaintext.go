package xscmsg

import (
	"fmt"
	"net/textproto"

	"github.com/rothskeller/packet/pktmsg"
)

// PlainTextTag identifies a plain text message.
const PlainTextTag = "plain"
const plainTextName = "plain text message"

func createPlainText() Message {
	var pt plainText

	pt.fields = []FormField{
		newPTOriginMessageNumberField(pt),
		newPTHandlingField(pt),
		newPTSubjectField(pt),
		newPTBodyField(pt),
	}
	return pt
}

func recognizePlainText(m *pktmsg.Message, _ *pktmsg.Form) Message {
	var pt = createPlainText().(plainText)

	pt.header = m.Header
	if subject := ParseSubject(m.Header.Get("Subject")); subject != nil {
		pt.fields[0].SetValue(Value(subject.MessageNumber))
		pt.fields[1].SetValue(Value(subject.HandlingOrderCode))
		pt.fields[2].SetValue(Value(subject.Subject))
	} else {
		pt.fields[2].SetValue(Value(m.Header.Get("Subject")))
	}
	pt.fields[3].SetValue(Value(m.Body))
	return pt
	// TODO set immutable flag
}

type plainText struct {
	fieldContainer
	header    textproto.MIMEHeader
	immutable bool
}

func (m plainText) TypeTag() string     { return PlainTextTag }
func (m plainText) TypeName() string    { return plainTextName }
func (m plainText) TypeArticle() string { return "a" }
func (m plainText) EncodedSubject() string {
	if m.immutable {
		return m.header.Get("Subject")
	}
	if m.FieldValue("MSGNO") != "" || m.FieldValue("HANDLING") != "" {
		return fmt.Sprintf("%s_%s_%s", m.FieldValue("MSGNO"), m.FieldValue("HANDLING"), m.FieldValue("SUBJECT"))
	}
	return string(m.FieldValue("SUBJECT"))
}
func (m plainText) EncodedBody() string { return string(m.FieldValue("BODY")) }

func newPTOriginMessageNumberField(c FieldContainer) (f FormField) {
	f = BaseField(c, FOriginMsgNo, "MSGNO", "Origin Message Number", 10, 1,
		`This is the message number for this message assigned to it by its origin station.`)
	f = MessageNumberField(true, f)
	f = RequiredField(f)
	return f
}

func newPTHandlingField(c FieldContainer) (f FormField) {
	f = BaseField(c, FHandling, "HANDLING", "Handling Order", 1, 1,
		`This is the handling order with which this message should be processed.`)
	f = ChoicesField(f, Value("IMMEDIATE"), Value("PRIORITY"), Value("ROUTINE"))
	f = RequiredField(f)
	return f
}

func newPTSubjectField(c FieldContainer) (f FormField) {
	f = BaseField(c, FSubject, "SUBJECT", "Subject", 80, 1, `This is the subject of the message.`)
	f = RequiredField(f)
	return f
}

func newPTBodyField(c FieldContainer) (f FormField) {
	f = BaseField(c, FBody, "BODY", "Body", 80, 5, `This is the body of the message.`)
	f = RequiredField(f)
	return f
}
