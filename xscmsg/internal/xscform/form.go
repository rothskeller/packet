package xscform

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

var fieldAnnotations = make(map[*xscmsg.MessageType]map[string]string)
var fieldComments = make(map[*xscmsg.MessageType]map[string]string)

// CreateForm creates a new XSCForm with the specified form definition, filling
// in the defaults.
func CreateForm(mtype *xscmsg.MessageType, fields []xscmsg.Field) *XSCForm {
	var xf = XSCForm{mtype: mtype, fields: fields, form: new(pktmsg.Form)}
	xf.form.FormType = mtype.HTML
	xf.form.FormVersion = mtype.Version
	for _, field := range xf.fields {
		if val := field.Default(); val != "" {
			field.Set(val)
		}
	}
	return &xf
}

// AdoptForm returns a new XSCForm for the specified message.
func AdoptForm(mtype *xscmsg.MessageType, fields []xscmsg.Field, msg *pktmsg.Message, form *pktmsg.Form) *XSCForm {
	var xf = XSCForm{mtype: mtype, fields: fields, msg: msg, form: form}
	for _, f := range form.Fields {
		nf := xf.Field(f.Tag)
		if nf == nil {
			if f.Value == "" {
				continue // ignore unknown fields with no value
			}
			nf = &unknownField{tag: f.Tag}
			xf.fields = append(xf.fields, nf)
		}
		nf.Set(f.Value)
	}
	return &xf
}

// XSCForm is the base implementation of all XSC form types.
type XSCForm struct {
	mtype  *xscmsg.MessageType
	fields []xscmsg.Field
	msg    *pktmsg.Message
	form   *pktmsg.Form
}

// Type returns the message type definition.
func (xf *XSCForm) Type() *xscmsg.MessageType { return xf.mtype }

// Fields returns the list of fields in the message.
func (xf *XSCForm) Fields() []xscmsg.Field { return xf.fields }

// Field returns the field with the specified name, or nil if there is
// no such field.
func (xf *XSCForm) Field(name string) xscmsg.Field {
	for _, f := range xf.fields {
		id := f.ID()
		if id.Tag == name || id.Canonical == name {
			return f
		}
	}
	return nil
}

// Validate ensures that the contents of the message are correct.  It
// returns a list of problems, which is empty if the message is fine.
// If strict is true, the message must be exactly correct; otherwise,
// some trivial issues are corrected and not reported.
func (xf *XSCForm) Validate(strict bool) (problems []string) {
	for _, f := range xf.fields {
		if err := f.Validate(xf, strict); err != "" {
			problems = append(problems, err)
		}
	}
	return problems
}

// EncodeSubject returns the encoded subject of the message.
func (xf *XSCForm) EncodeSubject() string {
	ho, _ := xscmsg.ParseHandlingOrder(xf.Field(xscmsg.FHandling).Get())
	omsgno := xf.Field(xscmsg.FOriginMsgNo).Get()
	subject := xf.Field(xscmsg.FSubject).Get()
	return xscmsg.EncodeSubject(omsgno, ho, xf.mtype.Tag, subject)
}

// EncodeBody returns the encoded body of the message.  If human is true, it is
// encoded for human reading or editing; if false, it is encoded for
// transmission.
func (xf *XSCForm) EncodeBody(human bool) string {
	xf.form.Fields = xf.form.Fields[:0]
	for _, f := range xf.fields {
		value := f.Get()
		if value != "" || (human && !f.ID().ReadOnly) {
			xf.form.Fields = append(xf.form.Fields, pktmsg.FormField{Tag: f.ID().Tag, Value: value})
		}
	}
	if human {
		annotations, comments := generateFieldAnnotationsAndComments(xf.mtype, xf.fields)
		return xf.form.Encode(annotations, comments, true)
	}
	return xf.form.Encode(nil, nil, false)
}

// generateFieldAnnotationAndComments generates (or returns cached) maps from
// field tag to field annotation and to field comment for the specified message
// type.
func generateFieldAnnotationsAndComments(mtype *xscmsg.MessageType, fields []xscmsg.Field) (m, c map[string]string) {
	if m = fieldAnnotations[mtype]; m != nil {
		return m, fieldComments[mtype]
	}
	m = make(map[string]string)
	c = make(map[string]string)
	for _, f := range fields {
		id := f.ID()
		if id.Annotation != "" {
			m[id.Tag] = id.Annotation
		}
		if id.Comment != "" {
			c[id.Tag] = id.Comment
		}
	}
	fieldAnnotations[mtype] = m
	fieldComments[mtype] = c
	return m, c
}

/*
// TypeTag returns the tag string used to identify the message type.
func (xf *XSCForm) TypeTag() string { return xf.def.Tag }

// TypeName returns the English name of the message type.  It is a noun phrase
// in prose case, such as "foo message" or "bar form".
func (xf *XSCForm) TypeName() string { return xf.def.Name }

// TypeArticle returns the indefinite article ("a" or "an") to be used preceding
// the TypeName, in a sentence that needs one.
func (xf *XSCForm) TypeArticle() string { return xf.def.Article }

// Validate ensures that the contents of the message are correct.  It returns a
// list of problems, which is empty if the message is fine.  If strict is true,
// the message must be exactly correct; otherwise, some trivial issues are
// corrected and not reported.
func (xf *XSCForm) Validate(strict bool) (problems []string) {
	var seen = make(map[string]bool)
	for _, fd := range xf.def.Fields {
		seen[fd.Tag] = true
		if len(fd.Validations) == 0 {
			continue
		}
		value := xf.form.Get(fd.Tag)
		for _, validation := range fd.Validations {
			if newval, problem := validation(xf, fd, value, strict); problem != "" {
				problems = append(problems, problem)
			} else if newval != value && strict {
				problems = append(problems, fmt.Sprintf("%q is not a valid value for field %q", value, fd.Tag))
			} else if newval != value {
				xf.form.Set(fd.Tag, newval)
			}
		}
	}
	for _, field := range xf.form.Fields {
		if !seen[field.Tag] {
			problems = append(problems, fmt.Sprintf("unrecognized field %q", field.Tag))
		}
	}
	return problems
}

// Message returns the encoded message.  If human is true, it is encoded for
// human reading or editing; if false, it is encoded for transmission.  If the
// XSCMessage was originally created by a call to Recognize, the Message
// structure passed to it is updated and reused; only its Body and Subject and
// Subject header are changed.  Otherwise, a new Message structure is created
// and filled in.
func (xf *XSCForm) Message(human bool) *pktmsg.Message {
	if xf.msg == nil {
		xf.msg = pktmsg.New()
	}
	if human {
		xf.form.EncodeToMessage(xf.msg, xf.def.Annotations, xf.def.Comments, true)
	} else {
		xf.form.EncodeToMessage(xf.msg, nil, nil, false)
	}
	xf.msg.Header.Set("Subject", xf.EncodeSubject())
	return xf.msg
}

// EncodeSubject returns the encoded subject line of the message based on its
// contents.
func (xf *XSCForm) EncodeSubject() string {
	ho, _ := xscmsg.ParseHandlingOrder(xf.form.Get(xf.def.HandlingOrderField))
	return xscmsg.EncodeSubject(xf.form.Get(xf.def.OriginNumberField), ho, xf.def.Tag,
		xf.form.Get(xf.def.SubjectField),
	)
}

// OriginNumber returns the origin message number of the message, if any.
func (xf *XSCForm) OriginNumber() string { return xf.form.Get(xf.def.OriginNumberField) }

// SetOriginNumber sets the originmessage number of the message, if the message
// type supports that.
func (xf *XSCForm) SetOriginNumber(msgnum string) {
	if field := xf.def.OriginNumberField; field != "" {
		xf.form.Set(field, msgnum)
	}
}

// DestinationNumber returns the destination message number of the message, if
// any.
func (xf *XSCForm) DestinationNumber() string { return xf.form.Get(xf.def.DestinationNumberField) }

// SetDestinationNumber sets the destination message number of the message, if
// the message type supports that.
func (xf *XSCForm) SetDestinationNumber(msgnum string) {
	if field := xf.def.DestinationNumberField; field != "" {
		xf.form.Set(field, msgnum)
	}
}

// HandlingOrder returns the message handling order, if any, in both string and
// parsed forms.
func (xf *XSCForm) HandlingOrder() (s string, ho xscmsg.HandlingOrder) {
	s = xf.form.Get(xf.def.HandlingOrderField)
	ho, _ = xscmsg.ParseHandlingOrder(s)
	return s, ho
}

// SetHandlingOrder sets the message handling order, if the message type
// supports that.
func (xf *XSCForm) SetHandlingOrder(ho xscmsg.HandlingOrder) {
	if field := xf.def.HandlingOrderField; field != "" {
		xf.form.Set(field, ho.String())
	}
}

// Routing returns the To ICS Position and To Location fields of the form, if
// it has them.
func (xf *XSCForm) Routing() (pos, loc string) {
	return xf.form.Get("7a."), xf.form.Get("7b.")
}

// Operator returns the operator name and call sign from the message, if any.
func (xf *XSCForm) Operator() (name string, callSign string) {
	return xf.form.Get(xf.def.OperatorNameField), xf.form.Get(xf.def.OperatorCallField)
}

// SetOperator sets the operator name and call sign in the message, if the
// message type supports that.
func (xf *XSCForm) SetOperator(name string, callSign string) {
	if field := xf.def.OperatorNameField; field != "" {
		xf.form.Set(field, name)
	}
	if field := xf.def.OperatorCallField; field != "" {
		xf.form.Set(field, callSign)
	}
}

// ActionTime returns the time that the outgoing message was sent or the
// incoming message was received, if the message has that information.  It
// returns them in string form, and also in parsed form if they are in standard
// format.
func (xf *XSCForm) ActionTime() (actdate string, acttime string, ts time.Time) {
	actdate = xf.form.Get(xf.def.ActionDateField)
	acttime = xf.form.Get(xf.def.ActionTimeField)
	ts, _ = time.ParseInLocation("1/2/2006 15:04", actdate+" "+acttime, time.Local)
	return
}

// SetActionTime sets the time that the outgoing message was sent or the income
// message was received, if the message type supports that.
func (xf *XSCForm) SetActionTime(ts time.Time) {
	if field := xf.def.ActionDateField; field != "" {
		xf.form.Set(field, ts.Format("01/02/2006"))
	}
	if field := xf.def.ActionTimeField; field != "" {
		xf.form.Set(field, ts.Format("15:04"))
	}
}

// Form returns the underlying pktmsg.Form
func (xf *XSCForm) Form() *pktmsg.Form { return xf.form }

// Get retrieves the value of a form field.
func (xf *XSCForm) Get(tag string) string { return xf.form.Get(tag) }

// Set sets the value of a form field.
func (xf *XSCForm) Set(tag, value string) { xf.form.Set(tag, value) }

// FindField finds the specified field definition in the form definition.
func (fd *FormDefinition) FindField(tag string) *FieldDefinition {
	for _, ff := range fd.Fields {
		if ff.Tag == tag {
			return ff
		}
	}
	return nil
}
*/
