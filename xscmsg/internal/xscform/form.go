package xscform

import (
	"fmt"
	"time"

	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/xscmsg"
)

// A FormDefinition defines the characteristics of a form such that the base
// XSCForm implementation can work with it.  FormDefinitions are generally
// auto-generated from SCCoPIFO HTML files with the extract-pifo-fields build
// tool.
type FormDefinition struct {
	HTML                   string
	Tag                    string
	Name                   string
	Article                string
	Version                string
	OriginNumberField      string
	DestinationNumberField string
	HandlingOrderField     string
	SubjectField           string
	OperatorNameField      string
	OperatorCallField      string
	ActionDateField        string
	ActionTimeField        string
	Fields                 []*FieldDefinition
	Annotations            map[string]string
	Comments               map[string]string
}

// A FieldDefinition is the definition of a single field in a FormDefinition.
type FieldDefinition struct {
	Tag         string
	Values      []string
	Validations []ValidateFunc
	Default     string
}

// A ValidateFunc is a function that validates the value of a field in a form.
// It takes the form, field definition, and value to be validated as input.  It
// returns the value — which may have been corrected — and a problem string.
// The problem string will describe any uncorrectable issue and will otherwise
// be empty.
type ValidateFunc func(xf *XSCForm, fd *FieldDefinition, value string) (newval, problem string)

// CreateForm creates a new XSCForm with the specified form definition, filling
// in the defaults.
func CreateForm(def *FormDefinition) *XSCForm {
	var xf = XSCForm{def: def, form: new(pktmsg.Form)}
	xf.form.FormType = def.HTML
	xf.form.FormVersion = def.Version
	for _, field := range def.Fields {
		var val = field.Default
		if val == "«date»" {
			val = time.Now().Format("01/02/2006")
		}
		xf.form.Set(field.Tag, val)
	}
	return &xf
}

// RecognizeForm returns a new XSCForm for the specified message if it contains
// a form matching the specified form definition.  The form version must be the
// same as or newer than that of the form definition.
func RecognizeForm(def *FormDefinition, msg *pktmsg.Message, form *pktmsg.Form) *XSCForm {
	if form == nil || form.FormType != def.HTML || xscmsg.OlderVersion(form.FormVersion, def.Version) {
		return nil
	}
	return &XSCForm{def: def, msg: msg, form: form}
}

// XSCForm is the base implementation of all XSC form types.
type XSCForm struct {
	def  *FormDefinition
	msg  *pktmsg.Message
	form *pktmsg.Form
}

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
		if strict && !xf.form.Has(fd.Tag) {
			problems = append(problems, fmt.Sprintf("field %q is not present", fd.Tag))
			continue
		}
		seen[fd.Tag] = true
		if len(fd.Validations) == 0 {
			continue
		}
		value := xf.form.Get(fd.Tag)
		for _, validation := range fd.Validations {
			if newval, problem := validation(xf, fd, value); problem != "" {
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
			problems = append(problems, "unrecognized field %q", field.Tag)
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
