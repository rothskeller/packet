package pktmsg

import "time"

// IMessage is the interface satisfied by Message (and, more importantly, by
// all other types that embed a *Message).  The methods are all getters and
// setters on the Message fields.
type IMessage interface {
	IsAutoresponse() bool
	GetBody() string
	SetBody(string)
	GetBBSRxDate() time.Time
	GetFormHTML() string
	GetFormTag() string
	GetFormVersion() string
	GetFrom() string
	GetFromAddrs() []string
	SetFrom(string)
	GetHandling() string
	SetHandling(string)
	IsNotPlainText() bool
	SetOutpostUrgent(bool)
	GetOriginMsgID() string
	SetOriginMsgID(string)
	GetPIFOVersion() string
	IsRequestDeliveryReceipt() bool
	IsRequestReadReceipt() bool
	GetReturnAddr() string
	GetRxArea() string
	GetRxBBS() string
	GetRxDate() string
	GetRxDateTime() time.Time
	GetSentDate() string
	GetSentDateTime() time.Time
	SetSentDate(time.Time)
	GetSeverity() string
	GetSubject() string
	SetSubject(string)
	GetSubjectHeader() string
	GetTo() string
	GetToAddrs() []string
	SetTo(string)
	GetTaggedFields() []TaggedField
	IsForm() bool
	Finalized() bool
	Save() string
	Transmit() ([]string, string, string)
}

func (m Message) IsAutoresponse() bool { return m.Autoresponse }

func (m Message) GetBody() string       { return m.Body }
func (m *Message) SetBody(value string) { m.Body = value }

func (m Message) GetBBSRxDate() time.Time { return m.BBSRxDate.Time() }

func (m Message) GetFormHTML() string { return m.FormHTML }

func (m Message) GetFormTag() string { return m.FormTag }

func (m Message) GetFormVersion() string { return m.FormVersion }

func (m Message) GetFrom() string        { return m.From.String() }
func (m Message) GetFromAddrs() []string { return m.From.Addresses() }
func (m *Message) SetFrom(value string)  { m.From.SetString(value) }

func (m Message) GetHandling() string       { return m.Handling }
func (m *Message) SetHandling(value string) { m.Handling = value }

func (m Message) IsNotPlainText() bool { return m.NotPlainText }

func (m *Message) SetOutpostUrgent(value bool) { m.OutpostUrgent = value }

func (m Message) GetOriginMsgID() string       { return m.OriginMsgID }
func (m *Message) SetOriginMsgID(value string) { m.OriginMsgID = value }

func (m Message) GetPIFOVersion() string { return m.PIFOVersion }

func (m Message) IsRequestDeliveryReceipt() bool { return m.RequestDeliveryReceipt }

func (m Message) IsRequestReadReceipt() bool { return m.RequestReadReceipt }

func (m Message) GetReturnAddr() string { return m.ReturnAddr }

func (m Message) GetRxArea() string { return m.RxArea }

func (m Message) GetRxBBS() string { return m.RxBBS }

func (m Message) GetRxDate() string        { return m.RxDate.String() }
func (m Message) GetRxDateTime() time.Time { return m.RxDate.Time() }

func (m Message) GetSentDate() string          { return m.SentDate.String() }
func (m Message) GetSentDateTime() time.Time   { return m.SentDate.Time() }
func (m *Message) SetSentDate(value time.Time) { m.SentDate.SetTime(value) }

func (m Message) GetSeverity() string { return m.Severity }

func (m Message) GetSubject() string       { return m.Subject }
func (m *Message) SetSubject(value string) { m.Subject = value }

func (m Message) GetSubjectHeader() string { return m.SubjectHeader }

func (m Message) GetTo() string        { return m.To.String() }
func (m Message) GetToAddrs() []string { return m.To.Addresses() }
func (m *Message) SetTo(value string)  { m.To.SetString(value) }

func (m Message) GetTaggedFields() []TaggedField { return m.TaggedFields }
