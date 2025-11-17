// Package formdef contains the FormDef type and related sub-structs, and the
// code to read and write them from *.form files.
package formdef

import (
	"io/fs"
	"iter"

	"github.com/rothskeller/pdf/v2"
	"k8s.io/apimachinery/pkg/util/sets"
)

// A FormDef defines a PackItForms form.  It is version-specific; each version
// of a form has its own FormDef.
type FormDef struct {
	// FormFS is the filesystem containing the HTMLFile and PDFFile.
	FormFS fs.FS
	// AddonName is the Outpost addon name under which this form is defined.
	AddonName string
	// HTMLName is the name (identifier) of this form as seen in its encoded
	// representation.  While it looks like an HTML file name (form-*.html),
	// no such file has to exist.
	HTMLName string
	// Version is the version number of this version of the form.
	Version string
	// Title is the name of the form, in title case.
	Title string
	// Sort is the key used to place this form in a list with other forms.
	// It is often the same as Title, but doesn't have to be.
	Sort string
	// IndefName is the name of the form preceded by the proper indefinite
	// article ("a" or "an"), and rendered in prose case.
	IndefName string
	// CreateTags is a list of tag words that can be used to identify this
	// version of this form when creating a new message.  It is empty if new
	// messages cannot be created with this version of this form.
	CreateTags []string
	// SubjectTag is the tag word that should be used to identify the form
	// on the subject line of a message.
	SubjectTag string
	// HTMLFile is the filename of the HTML file, within FormFS, used to
	// create or edit messages containing this version of this form.  It is
	// empty if new messages cannot be created with this version of this
	// form.
	HTMLFile string
	// PDFFile is the filename of the PDF file template, within FormFS,
	// used to render this version of this form in PDF format.  If may be
	// empty for a plain text rendering.
	PDFFile string
	// Fields is a list of fields of the form.
	Fields []*FieldDef
}

// AllFields returns an iterator on all fields of the form, including children
// of aggregate fields.
func (def *FormDef) AllFields() iter.Seq[*FieldDef] {
	return func(yield func(*FieldDef) bool) {
		for _, fd := range def.Fields {
			if !yield(fd) {
				return
			}
			for _, cd := range fd.Children {
				if !yield(cd) {
					return
				}
			}
		}
	}
}

// A FieldDef defines one field of a form.
type FieldDef struct {
	// Parent is the parent of this field, if this is a child field.
	Parent *FieldDef
	// Children is a list of subsidiary fields.
	Children []*FieldDef
	// Label is the label for the field.  It should be short (definitely no
	// more than 40 characters), and should not rely on surrounding field
	// labels for context.
	Label string
	// Tag is the tag string that identifies the field in the PackItForms
	// encoding of the message.  It is empty for fields that do not get
	// encoded.
	Tag string
	// Common is the tag string that identifies the field as one of the
	// common fields.  It allows software to address those fields even when
	// their PackItForms tags vary.
	Common string
	// Type identifies the type of field, which usually means identifying
	// the type of the information stored in the field.
	Type string
	// Value is the static value of the field, for fields of type `static`.
	// For fields of type `join`, it is the format string used to join the
	// subsidiary field values.  For most other field types, it is the default
	// value of the field when creating a new form.
	Value string
	// Choices is the set of allowed values for the field, for fields of
	// type `restricted`.
	Choices []Choice
	// Presence indicates whether the field is allowed or required.  See
	// also PresenceCond.
	Presence Presence
	// PresenceCond describes a conditional change to Presence.  It is
	// optional.
	PresenceCond []*PresenceCond
	// EditWidth is the width in characters of the input control for editing
	// the field.  It is zero for non-editable fields.
	EditWidth int
	// EditHelp is the help text describing the content and validation rules
	// for the field.  It is empty for non-editable fields.
	EditHelp string
	// ChildLabel is the label used to identify this field in the context of
	// the set of other fields that are children of the same parent.
	ChildLabel string
	// CompareMethod is the method used to compare values of this field
	// across messages.  By default, the compare method for the field type
	// is used.
	CompareMethod string
	// PDF is the set of PDF rendering instructions for this field.
	PDF []PDFFieldRenderer
}

// A Choice is a pair of raw and human representations of the same value, and
// optionally a condition to enable the choice.
type Choice struct {
	Raw       string
	Human     string
	CondField string
	CondValue string
}

// A PresenceCond describes a conditional change to a field's Presence based on
// the value of another field.  If OtherField's value matches OtherValue, or
// OtherField's value is non-empty and OtherValue is empty, the instant field's
// Presence is set to this structure's Presence value.
type PresenceCond struct {
	OtherField string
	OtherValue string
	Presence   Presence
}

type Presence string

const (
	Required Presence = "required"
	Optional Presence = "optional"
	Blocked  Presence = "blocked"
)

// A PDFFieldRenderer specifies how to render a field onto a PDF.
type PDFFieldRenderer struct {
	// Conditions is a list of conditions that must be satisfied for this
	// renderer to be used.
	Conditions []PDFCondition
	// Renderer is the actual renderer.
	Renderer interface{ Draw(*pdf.PDF, string) error }
}

// A PDFCondition represents a conditional governing a PDF field renderer.
// The renderer will render only if the condition is satisfied.
type PDFCondition struct {
	// Field is the tag of the field being tested in the conditional.
	Field string
	// Value is the value that the Field must have to satisfy the condition.
	Value string
	// If Set is true, any non-empty value of Field satisfies the condition,
	// and Value is disregarded.
	Set bool
}

// A BoxRenderer draws a filled rectangle (mainly used to overwrite static page
// numbering so we can replace it with our own).
type BoxRenderer struct{ *pdf.Box }

func (r BoxRenderer) Draw(p *pdf.PDF, _ string) error {
	return r.Box.Draw(p)
}

// A CircleRenderer draws a filled circle (i.e., fills in a radio button).
type CircleRenderer struct{ *pdf.Circle }

func (r CircleRenderer) Draw(p *pdf.PDF, _ string) error {
	return r.Circle.Draw(p)
}

// A CrossRenderer draws an X (i.e., fills in a checkbox).
type CrossRenderer struct{ *pdf.Cross }

func (r CrossRenderer) Draw(p *pdf.PDF, _ string) error {
	return r.Cross.Draw(p)
}

// A TextRenderer draws a text string.
type TextRenderer struct{ *pdf.Text }

func (r TextRenderer) Draw(p *pdf.PDF, v string) (err error) {
	// Make a copy of the Text structure so we can change it without
	// concurrency issues.
	var t = *r.Text

	if t.String == "" {
		t.String = v
	}
	return t.Draw(p)
}

// CommonTags contains the allowed values for Field.Common.
var CommonTags = sets.New(
	"defaultBody",
	"destinationMessageID",
	"formDate",
	"fromICSPosition",
	"fromLocation",
	"handling",
	"messageDate",
	"messageSummary",
	"messageTime",
	"operatorCall",
	"operatorDate",
	"operatorName",
	"operatorTime",
	"originMessageID",
	"reference",
	"tacticalCall",
	"tacticalName",
	"toICSPosition",
	"toLocation",
)
