package message

import (
	"iter"
	"os"
	"strings"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/pdf/v2"
)

// MType is the interface satisfied by all messages types.
type MType interface {
	// Recognize detects whether the argument message is of the type
	// described by this MType, and if so, calls its SetType method to
	// assign this MType to it.
	Recognize(Message)
	// CreateTag returns the tag(s) used to identify this message type on a
	// "new" command line.  Multiple tags may be returned, separated by
	// spaces.  Version numbers may be included in a tag; see the "new"
	// command documentation for details.  Message types that cannot be
	// created with "new" return an empty string.
	CreateTag() string
	// Name returns the name of the message type, as a phrase in lower case
	// (other than acronyms) starting with "a " or "an ".
	Name() string
	// Validate validates the contents of the message and returns any
	// problems.  If pifo is true, it returns only problems that
	// PackItForms would raise; otherwise the checks may be more extensive.
	Validate(m Message, pifo bool) error
	// Fields returns an iterator on the set of message fields.
	Fields() iter.Seq[field.Field]
	// RenderPDF creates a PDF representation of the message in the
	// specified file.  If copyname is not empty, it is placed in the
	// footer of each page.  The returned error may be a Warning, showing a
	// non-fatal rendering issue.
	RenderPDF(m Message, filename, copyname string) error
}

// Warning wraps a non-fatal "error".
type Warning struct{ Err error }

func (w Warning) Error() string { return w.Err.Error() }
func (w Warning) Unwrap() error { return w.Err }

//-----------------------------------------------------------------------------

// mtypes is the list of registered message types.
type registeredMType struct {
	mt       MType
	fallback bool
}

var mtypes []registeredMType

// RegisterType registers a message type.
func RegisterType(mt MType) {
	mtypes = append(mtypes, registeredMType{mt: mt})
}

// RegisterFallbackType registers a message type whose recognizer is guaranteed
// to be called after all types registered with RegisterType, even if they are
// registered later.
func RegisterFallbackType(mt MType) {
	mtypes = append(mtypes, registeredMType{mt: mt, fallback: true})
}

// SetType determines the type of a Message and sets its Type property to the
// first registered type whose Recognize method recognizes it, or otherwise to
// PlainMessage.
func SetType(m Message) {
	if m.Type() != nil {
		return
	}
	for mt := range AllTypes() {
		mt.Recognize(m)
		if m.Type() != nil {
			return
		}
	}
	m.SetType(PlainMessage)
}

// AllTypes returns an iterator on all of the registered message types.
func AllTypes() iter.Seq[MType] {
	return func(yield func(MType) bool) {
		for _, mt := range mtypes {
			if !mt.fallback && !yield(mt.mt) {
				return
			}
		}
		for _, mt := range mtypes {
			if mt.fallback && !yield(mt.mt) {
				return
			}
		}
		yield(PlainMessage)
	}
}

//-----------------------------------------------------------------------------

// BaseMType is the common core implementation for all message types.
type BaseMType struct {
	createTag string
	name      string
	fields    []field.Field
}

// NewBaseMType returns a new BaseMType with the specified details to be
// returned by Name and CreateTag methods.
func NewBaseMType(name, createTag string) *BaseMType {
	return &BaseMType{name: name, createTag: createTag}
}

// CreateTag returns the tag(s) used to identify this message type on a "new"
// command line.  Multiple tags may be returned, separated by spaces.  Version
// numbers may be included in a tag; see the "new" command documentation for
// details.  Message types that cannot be created with "new" return an empty
// string.
func (t *BaseMType) CreateTag() string {
	return t.createTag
}

// Name returns the name of the message type, as a phrase in lower case (other
// than acronyms) starting with "a " or "an ".
func (t *BaseMType) Name() string { return t.name }

// Validate validates the contents of the message and returns any problems.
// If pifo is true, it returns only problems that PackItForms would raise;
// otherwise the checks may be more extensive.
func (t *BaseMType) Validate(m Message, pifo bool) (err error) {
	for f := range t.Fields() {
		err = errors.Join(err, f.Validate(m, f, pifo))
	}
	return err
}

// AddField adds fields to a form definition.  It is usually called by the
// form definition's initFunc.  The arguments can be either field.Field
// implementations or *field.FieldFactory instances.
func (t *BaseMType) AddField(fs ...any) {
	for _, f := range fs {
		switch f := f.(type) {
		case field.Field:
			t.fields = append(t.fields, f)
		case *field.FieldFactory:
			t.fields = append(t.fields, f.MakeField())
		default:
			panic("AddField arguments must be field.Field or *field.FieldFactory")
		}
	}
}

// Fields returns an iterator on the set of message fields.
func (t *BaseMType) Fields() iter.Seq[field.Field] {
	return func(yield func(field.Field) bool) {
		for _, f := range t.fields {
			if !yield(f) {
				return
			}
			for _, cf := range f.Children() {
				if !yield(cf) {
					return
				}
			}
		}
	}
}

// RenderPDF creates a PDF representation of the message in the specified file.
func (t *BaseMType) RenderPDF(m Message, filename, copyname string) (err error) {
	const (
		margin              = 48
		headingFont         = "Helvetica-Bold"
		headingFontSize     = 14
		metadataFont        = "Helvetica"
		metadataFontSize    = 12
		metadataLineSpacing = 14
		metadataLabelFont   = "Helvetica-Bold"
		metadataLabelWidth  = 60
		bodyFont            = "Courier"
		bodyFontSize        = 12
		footerFont          = "Helvetica"
		footerFontSize      = 12
		timestampFormat     = "Monday, January 2, 2006 at 15:04:05"
	)
	var (
		fh       *os.File
		ph       *pdf.PDF
		avail    pdf.Rectangle
		from     string
		date     string
		rcvd     string
		body     string
		warnings pdf.ErrTextRendering
		page     = 1
	)
	if fh, err = os.Create(filename); err != nil {
		return err
	}
	ph = pdf.New(fh)
	ph.Info["Title"] = m.Subject().EncodedSubject()
	ph.Info["Producer"] = "https://github.com/rothskeller/packet"
	ph.AddPage(pdf.USLetterPortrait)
	avail = pdf.USLetterPortrait
	// Set margins.
	avail.LLX, avail.LLY = avail.LLX+margin, avail.LLY+margin
	avail.URX, avail.URY = avail.URX-margin, avail.URY-margin
	// Room for footer.
	avail.LLY += 2 * footerFontSize
	// Add a banner.
	(pdf.Text{String: plainPDFBanner(m.Type().Name()), Rectangle: avail, Font: headingFont, FontSize: headingFontSize, Align: "lt"}).Draw(ph)
	avail.URY -= 2 * headingFontSize
	if m, ok := m.(interface{ From() string }); ok {
		from = m.From()
	}
	if m, ok := m.(interface{ Date() time.Time }); ok {
		date = m.Date().Format(timestampFormat)
	}
	if m, ok := m.(interface{ RxDate() time.Time }); ok {
		rcvd = m.RxDate().Format(timestampFormat)
	}
	// TODO: add receipt information
	if from != "" {
		(pdf.Text{String: "From", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		(pdf.Text{String: from, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, avail.URY-metadataLineSpacing
	}
	if m.To() != "" {
		(pdf.Text{String: "To", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		(pdf.Text{String: m.To(), Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, avail.URY-metadataLineSpacing
	}
	if sub := m.Subject().EncodedSubject(); sub != "" {
		(pdf.Text{String: "Subject", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		(pdf.Text{String: sub, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, avail.URY-metadataLineSpacing
	}
	if date != "" {
		(pdf.Text{String: "Date", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		(pdf.Text{String: date, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, avail.URY-metadataLineSpacing
	}
	if rcvd != "" {
		(pdf.Text{String: "Received", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		(pdf.Text{String: rcvd, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, avail.URY-metadataLineSpacing
	}
	avail.URY -= metadataLineSpacing
	body = m.Body().EncodedBody()
	for body != "" {
		t := pdf.Text{String: body, Page: page, Rectangle: avail, Font: bodyFont, FontSize: bodyFontSize, Align: "lt", Wrap: true, Clip: true}
		fits, overflow, _, w := t.WrapText()
		if w != nil {
			etr := w.(pdf.ErrTextRendering)
			warnings.Merge(etr)
		}
		t.String = fits
		t.Draw(ph)
		if body = overflow; body != "" {
			ph.AddPage(pdf.USLetterPortrait)
			page++
			avail.LLX, avail.LLY = avail.LLX+margin, avail.LLY+margin
			avail.URX, avail.URY = avail.URX-margin, avail.URY-margin
			avail.LLY += 2 * footerFontSize
		}
	}
	// TODO: page footers
	if err = ph.Write(); err != nil {
		return err
	}
	if err = fh.Close(); err != nil {
		return err
	}
	return warnings.AsError()
}

func plainPDFBanner(typeName string) string {
	_, typeName, _ = strings.Cut(typeName, " ")
	return strings.ToUpper(typeName)
}
