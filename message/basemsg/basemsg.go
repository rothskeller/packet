// Package basemsg provides the common code underlying all packet message types.
package basemsg

import (
	"github.com/rothskeller/packet/message"
)

// BaseMessage is the type underlying all packet messages, providing their
// shared functionality.  Every message, regardless of type, embeds a
// BaseMessage and provides access to it through the Base() method on the
// message type.
type BaseMessage struct {
	// MessageType is the type definition for the message type.
	MessageType *message.Type
	// PIFOVersion is the PIFO version used to encode the message.  It is
	// set only for messages with PIFO encoding.
	PIFOVersion string
	// FormVersion identifies the form and version in the message.  It is
	// set only for messages with PIFO encoding.
	Form *FormVersion
	// Fields is an ordered list of fields in the message.  This is the
	// core of the shared message functionality:  most operations are
	// implemented by iterating through these fields.
	Fields []*Field
	// PDFBase is the form-fillable PDF file whose fields will be filled in
	// to create a PDF rendering of the message.  It is nil if PDF rendering
	// is not supported.
	// TODO: move to message type.
	PDFBase []byte
	// PDFFontSize is default font size for the fillable fields in the PDF
	// file.  It can be zero if all of the fields already have assigned
	// sizes in the PDF file.  TODO: move to message type.
	PDFFontSize float64

	// Pointers to key fields.

	// FOriginMsgID points to the value of the Origin Message ID field.  It
	// is nil for message types that do not have that field.
	FOriginMsgID *string
	// FHandling points to the value of the Handling field.  It
	// is nil for message types that do not have that field.
	FHandling *string
	// FSubject points to the value of the field of the message that will
	// get propagated to the message's subject line.  It is nil for message
	// types that do not have any such field.
	FSubject *string
	// RestrictedSubject is a flag indicating that the FSubject field allows
	// only certain restricted values.  (It will not be populated with the
	// subject of a message being replied to, unless that message is of the
	// same type.)
	RestrictedSubject bool
	// FToICSPosition points to the value of the To ICS Position field.  It
	// is nil for message types that do not have that field.
	FToICSPosition *string
	// FToLocation points to the value of the To Location field.  It
	// is nil for message types that do not have that field.
	FToLocation *string
	// FReportType points to the value of the Report Type field.  It is nil
	// for message types that do not have that field.
	FReportType *string
	// FOpCall points to the value of the Operator Call Sign field.  It
	// is nil for message types that do not have that field.
	FOpCall *string
	// FOpName points to the value of the Operator Name field.  It
	// is nil for message types that do not have that field.
	FOpName *string
	// FBody points to the value of the most prominent, or first, multi-line
	// text field of the message.  It is nil for message types that do not
	// have any such field.
	FBody *string

	// Local storage.
	editFields   []*message.EditField
	editFieldMap map[*Field]*message.EditField
}

// A FormVersion identifies a version of a PackItForms form.
type FormVersion struct {
	// HTML is the HTML filename that identifies the type of form.
	HTML string
	// Version is the version number of the form.
	Version string
	// Tag is the tag for the form type, which goes in the subject line.
	Tag string
	// FieldOrder is an ordered list of field tags.  When provided,
	// generated forms list the fields in that order.  (Should any fields be
	// omitted from the list, they are put at the start of the form, before
	// those in the list.)
	FieldOrder []string
}

// Base returns the BaseMessage structure for the message.
func (bm *BaseMessage) Base() *BaseMessage { return bm }

// Type returns the message type definition.
func (bm *BaseMessage) Type() *message.Type { return bm.MessageType }

// Presence is a enumeration indicating whether a field is allowed or required.
type Presence uint8

// Values for Presence:
const (
	// PresenceNotAllowed means a value for this field is not allowed.
	// (This is generally because some parent field is not set, or some
	// conflicting field is set.)
	PresenceNotAllowed Presence = iota
	// PresenceOptional means a value for this field is allowed but not
	// required.
	PresenceOptional
	// PresenceRequired means a value for this field is required.
	PresenceRequired
)

func NotAllowed() (Presence, string) { return PresenceNotAllowed, "" }
func Optional() (Presence, string)   { return PresenceOptional, "" }
func Required() (Presence, string)   { return PresenceRequired, "" }

// A Field describes a single field within a message.  Generally, a message type
// has one Field for each field in the PackItForms encoding of the message, plus
// occasionally other Fields for special purposes (e.g. aggregating the
// underlying fields for display or editing).
type Field struct {
	// Label is the name of the field, as it is displayed to the user.  It
	// should be short, definitely no more than 40 characters.
	Label string
	// Value is a pointer to where the value of the field is stored.  Not
	// all fields have a stored value, so this pointer may be nil.
	Value *string
	// Choices is a set of recommended or allowed values for the field.
	// (Whether other values are allowed is up to the validation functions.)
	Choices ChoiceMapper
	// Presence is a function that returns whether the field is allowed or
	// required.  The function may optionally return a reason, which is
	// interpolated into validation problem strings when needed.
	Presence func() (Presence, string)
	// PIFOTag is the tag for this field in a PackItForms encoding.  If this
	// is empty, the field will not be rendered in PackItForms encoding nor
	// populated from PackItForms decoding.
	PIFOTag string
	// PIFOValid checks the value of the field against the restrictions
	// enforced by the PackItForms software.  It returns a problem
	// description if the value is one that PackItForms would reject, and an
	// empty string otherwise.
	PIFOValid func(*Field) string
	// Compare compares an expected value of this field against an actual
	// value of this field, and returns a description of the comparison.  To
	// disable comparison for a field, set this to CompareNone.
	Compare func(label, exp, act string) *message.CompareField
	// PDFMap is the mapper that tells how to render this field into a
	// form-fillable PDF file.
	PDFMap PDFMapper
	// TableValue returns the value of this field when rendered in flat text
	// table form.  To omit a field from the table rendering, set this to
	// TableOmit.
	TableValue func(*Field) string
	// EditWidth is the width in characters of the input control for this
	// field.  It should correspond to the number of characters that will
	// fit in the PDF rendering of the field, if applicable.
	EditWidth int
	// Multiline indicates that this field can contain multiple lines, i.e.,
	// can contain newline characters.
	Multiline bool
	// EditHelp is the help text for the form field, describing its contents
	// and its validity rules.  If this is empty, the field is not editable.
	EditHelp string
	// EditHint is a short string giving a model for the field value (e.g.,
	// "MM/DD/YYYY" for a date field).  It is optional, and will only be
	// displayed if there is room for it.
	EditHint string
	// EditValue returns the editable representation of the value of the
	// field.
	EditValue func(*Field) string
	// EditApply stores the supplied edited value into the field, revising
	// it if need be to convert from human to internal (PIFO) form.
	EditApply func(*Field, string)
	// EditValid checks the value of the field and returns a problem
	// description, or an empty string if there are no problems.
	EditValid func(*Field) string
	// EditSkip returns whether the field should be skipped while editing
	// the message (e.g., no entry in this field is valid because of the
	// value of some earlier field).
	EditSkip func(*Field) bool
}
