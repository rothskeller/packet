package form

import (
	"fmt"

	"github.com/rothskeller/packet/v4/form/formdef"
	"github.com/rothskeller/packet/v4/message"
	"github.com/rothskeller/packet/v4/message/msgifc"
)

type ErrUnrecognizedForm string

func (e ErrUnrecognizedForm) Error() string {
	return "The form contained in this message is not recognized.  " + string(e)
}

type ErrNotFormBody string

func (e ErrNotFormBody) Error() string {
	return fmt.Sprintf("A form message must have a form body, not %s.", string(e))
}

// An UnrecognizedForm is a message with a PackItForms-encoded body of a form
// that isn't known to this program (i.e., not recognized by any registered
// Recognizer).
type UnrecognizedForm struct{ *message.BaseMType }

func init() {
	var ft UnrecognizedForm

	ft.BaseMType = message.NewBaseMType("UNKNOWN", "an unrecognized form message")
	message.RegisterFallbackType(ft)
}

// Recognize works differently from other message types.  If it detects an
// unknown form, it manufactures a new message type for it on the fly, and
// assigns it to that message without registering it as an otherwise known type.
func (ft UnrecognizedForm) Recognize(m message.Message) {
	var (
		body *FormBody
	)
	if body, _ = m.Body().(*FormBody); body == nil {
		return // It's not a form.
	}
	body.def = &formdef.FormDef{
		AddonName: body.addonName,
		HTMLName:  body.formHTML,
		Version:   body.formVersion,
		IndefName: fmt.Sprintf("an unknown form %s version %s", body.formHTML, body.formVersion),
	}
	if fs, _ := m.Subject().(*FormSubject); fs != nil && fs.SubjectFormTag() != "" {
		body.def.SubjectTag = fs.SubjectFormTag()
	} else {
		body.def.SubjectTag = "UNKNOWN"
	}
	for _, f := range body.FieldList() {
		body.def.Fields = append(body.def.Fields, &formdef.FieldDef{
			Label: fmt.Sprintf("Field %s", f),
			Tag:   f,
			Type:  "text",
		})
	}
	m.SetType(&FormType{body.def, true})
}

func (ft UnrecognizedForm) CanCompareAgainst(_ msgifc.MType) bool { return false }

/*
// Validate validates an unrecognized form message.
func (m *UnrecognizedForm) Validate(pifo bool) (err error) {
	// The message should have an SCCo-standard subject line.
	switch s := m.Subject().(type) {
	case *subject.SCCoSubject:
		err = errors.Join(
			subject.ErrNoFormTag,
			s.Validate(!m.envelope.Received()),
		)
	case *subject.SCCoFormSubject:
		err = s.Validate(!m.envelope.Received())
	default:
		err = ErrNonStandardSubject
	}
	switch b := m.Body().(type) {
	case *body.FormBody:
		err = errors.Join(err, ErrUnrecognizedForm(fmt.Sprintf("%s/%s/%s", b.AddonName(), b.FormHTML(), b.FormVersion())))
	default:
		err = errors.Join(err, ErrNotFormBody(fmt.Sprintf("%T", b)))
	}
	return err
}
*/
