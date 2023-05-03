package xscmsg

import (
	"github.com/rothskeller/packet/pktmsg"
)

// baseMessage is the base implementation of xscmsg.Message.
type baseMessage struct {
	*xscMessage
	mtype     MessageType
	fromAddr  fromAddrField
	toAddrs   toAddrsField
	sentDate  sentDateField
	retrieved *retrievedField
	body      bodyField
}

// NewBaseMessage creates a new Message which handles the message envelope
// fields: From, To, Subject, Date, and Retrieved.  If pm is set, the new
// Message is based on that underlying pktmsg.Message; otherwise, it is based on
// a newly-created, outgoing message.
func NewBaseMessage(mtype MessageType, pm pktmsg.Message) Message {
	var xm *xscMessage
	switch pm := pm.(type) {
	case nil:
		xm = parseSubject(pktmsg.NewMessage())
	case *xscMessage:
		xm = pm
	default:
		xm = parseSubject(pm)
	}
	var m = baseMessage{xscMessage: xm, mtype: mtype}
	// Wrap the fields of the underlying message with xscmsg field
	// definitions that have labels, validation, etc.  Note that subject
	// doesn't need to be wrapped; that was already done by xscMessage.
	m.fromAddr = fromAddrField{xm.FromAddr()}
	m.toAddrs = toAddrsField{xm.ToAddrs()}
	if !pktmsg.Finalized(pm) {
		m.toAddrs.SetValue(m.toAddrs.Value()) // clean it up
	}
	m.sentDate = sentDateField{xm.SentDate()}
	if f := xm.RxBBS(); f != nil {
		m.retrieved = &retrievedField{f, xm.RxArea(), xm.RxDate()}
	}
	if f := xm.Body(); f != nil {
		m.body = bodyField{f}
	}
	return &m
}

// Default implementations for message fields that may not exist.  Overriding
// types may return values for these.
func (baseMessage) DestinationMsgID() pktmsg.SettableField { return nil }
func (baseMessage) MessageDate() pktmsg.SettableField      { return nil }
func (baseMessage) MessageTime() pktmsg.SettableField      { return nil }
func (baseMessage) OpCall() pktmsg.SettableField           { return nil }
func (baseMessage) OpDate() pktmsg.SettableField           { return nil }
func (baseMessage) OpName() pktmsg.SettableField           { return nil }
func (baseMessage) OpTime() pktmsg.SettableField           { return nil }
func (baseMessage) Reference() pktmsg.SettableField        { return nil }
func (baseMessage) TacCall() pktmsg.SettableField          { return nil }
func (baseMessage) TacName() pktmsg.SettableField          { return nil }
func (baseMessage) ToICSPosition() pktmsg.SettableField    { return nil }
func (baseMessage) ToLocation() pktmsg.SettableField       { return nil }

// Implementations for fields in the xscMessage.
func (m *baseMessage) FormTag() pktmsg.Field            { return &m.formtag }
func (m *baseMessage) Handling() pktmsg.SettableField   { return &m.handling }
func (m baseMessage) OriginMsgID() pktmsg.SettableField { return &m.originMsgID }
func (m *baseMessage) RawSubject() pktmsg.Field         { return m.xscMessage.subject }
func (m *baseMessage) Severity() pktmsg.SettableField   { return &m.severity }

func (m *baseMessage) Retrieved() pktmsg.Field { return m.retrieved }

func (m baseMessage) Type() MessageType { return m.mtype }

func (m *baseMessage) EditableFields(fn func(EditableField)) {}
func (m *baseMessage) ViewableFields(fn func(ViewableField)) {}
func (m *baseMessage) Iterate(fn func(pktmsg.Field)) {
	// First emit the visible/editable fields in the desired human order.
	fn(&m.fromAddr)
	fn(&m.toAddrs)
	fn(&m.subject)
	fn(&m.sentDate)
	if m.retrieved != nil {
		fn(m.retrieved)
	}
	fn(&m.body)
	// Then the rest of the underlying fields.
	m.Message.Iterate(func(f pktmsg.Field) {
		if f, ok := f.(pktmsg.KeyedField); ok {
			switch f.Key() {
			case pktmsg.FBody, pktmsg.FFromAddr, pktmsg.FSentDate, pktmsg.FSubject, pktmsg.FToAddrs:
				return // already emitted
			}
		}
		fn(f)
	})
}

func (m *baseMessage) KeyedField(key pktmsg.FieldKey) pktmsg.KeyedField {
	switch key {
	case FRetrieved:
		if m.retrieved != nil {
			return m.retrieved
		}
	case pktmsg.FFromAddr:
		return &m.fromAddr
	case pktmsg.FToAddrs:
		return &m.toAddrs
	case pktmsg.FSubject:
		return &m.subject
	case pktmsg.FSentDate:
		return &m.sentDate
	case pktmsg.FBody:
		return &m.body
	}
	return m.Message.KeyedField(key)
}
