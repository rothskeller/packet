package xscform

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// CreateForm creates a new Message with the specified form definition and
// fields, filling in the defaults.
func CreateForm(mtype *xscmsg.MessageType, fields []*xscmsg.FieldDef) *xscmsg.Message {
	var m = xscmsg.Message{Type: mtype, Fields: make([]*xscmsg.Field, len(fields))}
	for i, fd := range fields {
		var f = xscmsg.Field{Def: fd}
		f.Value = f.Default()
		m.Fields[i] = &f
	}
	return &m
}

// AdoptForm returns a new Message for the specified message.
func AdoptForm(mtype *xscmsg.MessageType, fields []*xscmsg.FieldDef, msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	var m = xscmsg.Message{Type: mtype, Fields: make([]*xscmsg.Field, len(fields)), RawMessage: msg, RawForm: form}
	for i, fd := range fields {
		m.Fields[i] = &xscmsg.Field{Def: fd}
	}
	for _, f := range form.Fields {
		nf := m.Field(f.Tag)
		if nf == nil {
			if f.Value == "" {
				continue // ignore unknown fields with no value
			}
			nf = &xscmsg.Field{Def: &xscmsg.FieldDef{
				Tag: f.Tag, Label: f.Tag,
				Validators: []xscmsg.Validator{ValidateUnknownField},
			}}
			m.Fields = append(m.Fields, nf)
		}
		nf.Value = f.Value
	}
	return &m
}

// EncodeSubject returns the encoded subject of the message.  It is generally
// not called directly; rather, it is installed as the MessageType.SubjectFunc
// for a form-based message type.
func EncodeSubject(m *xscmsg.Message) string {
	ho, _ := xscmsg.ParseHandlingOrder(m.KeyField(xscmsg.FHandling).Value)
	omsgno := m.KeyField(xscmsg.FOriginMsgNo).Value
	subject := m.KeyField(xscmsg.FSubject).Value
	return xscmsg.EncodeSubject(omsgno, ho, m.Type.Tag, subject)
}

// EncodeBody returns the encoded body of the message.  EncodeBody is generally
// not called directly; rather, it is installed as the MessageType.BodyFunc for
// a form-based message type.
func EncodeBody(m *xscmsg.Message) string {
	var form = pktmsg.Form{
		FormType:    m.Type.HTML,
		FormVersion: m.Type.Version,
	}
	if m.RawForm != nil {
		form.PIFOVersion, form.FormType, form.FormVersion =
			m.RawForm.PIFOVersion, m.RawForm.FormType, m.RawForm.FormVersion
	}
	for _, f := range m.Fields {
		if f.Value != "" {
			form.Fields = append(form.Fields, pktmsg.FormField{Tag: f.Def.Tag, Value: f.Value})
		}
	}
	return form.Encode()
}
