package xscform

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

var fieldAnnotations = make(map[*xscmsg.MessageType]map[string]string)
var fieldComments = make(map[*xscmsg.MessageType]map[string]string)

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

// EncodeBody returns the encoded body of the message.  If human is true, it is
// encoded for human reading or editing; if false, it is encoded for
// transmission.  EncodeBody is generally not called directly; rather, it is
// installed as the MessageType.BodyFunc for a form-based message type.
func EncodeBody(m *xscmsg.Message, human bool) string {
	var form = pktmsg.Form{
		FormType:    m.Type.HTML,
		FormVersion: m.Type.Version,
	}
	if m.RawForm != nil {
		form.PIFOVersion, form.FormType, form.FormVersion =
			m.RawForm.PIFOVersion, m.RawForm.FormType, m.RawForm.FormVersion
	}
	for _, f := range m.Fields {
		if f.Value != "" || (human && !f.Def.ReadOnly) {
			form.Fields = append(form.Fields, pktmsg.FormField{Tag: f.Def.Tag, Value: f.Value})
		}
	}
	if human {
		annotations, comments := generateFieldAnnotationsAndComments(m.Type, m.Fields)
		return form.Encode(annotations, comments, true)
	}
	return form.Encode(nil, nil, false)
}

// generateFieldAnnotationAndComments generates (or returns cached) maps from
// field tag to field annotation and to field comment for the specified message
// type.
func generateFieldAnnotationsAndComments(mtype *xscmsg.MessageType, fields []*xscmsg.Field) (m, c map[string]string) {
	if m = fieldAnnotations[mtype]; m != nil {
		return m, fieldComments[mtype]
	}
	m = make(map[string]string)
	c = make(map[string]string)
	for _, f := range fields {
		if f.Def.Annotation != "" {
			m[f.Def.Tag] = f.Def.Annotation
		}
		if f.Def.Comment != "" {
			c[f.Def.Tag] = f.Def.Comment
		}
	}
	fieldAnnotations[mtype] = m
	fieldComments[mtype] = c
	return m, c
}
