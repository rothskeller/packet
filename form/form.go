package form

import (
	"fmt"
	"io/fs"
	"iter"
	"os"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/subject"
	"github.com/rothskeller/pdf/v2"
)

type FormType struct {
	*formdef.FormDef
}

var _ message.MType = (*FormType)(nil)

// CreateTag returns the tag and version number for creating the form.
func (ft FormType) CreateTag() string { return strings.Join(ft.CreateTags, " ") }

// Name returns the name of the message type, as a phrase in lower case (other
// than acronyms) starting with "a " or "an ".
func (ft FormType) Name() string { return ft.IndefName }

// Validate validates the contents of the message and returns any
// problems.  If pifo is true, it returns only problems that
// PackItForms would raise; otherwise the checks may be more extensive.
func (ft FormType) Validate(m message.Message, pifo bool) error {
	panic("not implemented") // TODO: Implement
}

// Fields returns an iterator on the set of message fields.
func (ft FormType) Fields() iter.Seq[field.Field] {
	panic("not implemented") // TODO: Implement
}

// RenderPDF renders the form in PDF format.
func (ft FormType) RenderPDF(m message.Message, filename, copyname string) (err error) {
	var (
		srcFH    fs.File
		outFH    *os.File
		src      *pdf.PDF
		out      *pdf.PDF
		imp      *pdf.Importer
		warnings error
		maxPage  int
		msgID    string
		body     = m.Body().(*FormBody)
	)
	if ft.PDFFile == "" {
		return message.RenderPlainPDF(m, filename, copyname)
	}
	if outFH, err = os.Create(filename); err != nil {
		return err
	}
	defer func() {
		switch err.(type) {
		case nil, message.Warning:
			// nothing
		default:
			os.Remove(filename)
		}
	}()
	defer outFH.Close()
	// Open the template PDF file.
	if srcFH, err = ft.FormFS.Open(ft.PDFFile); err != nil {
		return err
	}
	defer srcFH.Close()
	if src, err = pdf.Open(srcFH.(pdf.Reader)); err != nil {
		return err
	}
	out = pdf.New(outFH)
	// Set the title of the PDF to the subject line of the message.
	out.Info["Title"] = m.Subject().EncodedSubject()
	out.Info["Producer"] = "https://github.com/rothskeller/packet"
	// Get the highest page number that we have to write something on.  We
	// only import that many pages.  (This allows the EOC-213RR and
	// RACES-MAR forms to include the recipient-only pages only on the
	// receiving side.)
	for fd := range ft.AllFields() {
	RENDERER1:
		for _, pr := range fd.PDF {
			for _, cond := range pr.Conditions {
				value := body.Field(cond.Field)
				if (cond.Set && value == "") || (!cond.Set && value != cond.Value) {
					continue RENDERER1
				}
			}
			switch shape := pr.Renderer.(type) {
			case formdef.CircleRenderer:
				maxPage = max(maxPage, shape.Page)
			case formdef.CrossRenderer:
				maxPage = max(maxPage, shape.Page)
			case formdef.TextRenderer:
				maxPage = max(maxPage, shape.Page)
			}
		}
	}
	// Copy pages from the source to the new PDF.
	if imp, err = out.NewImporter(src); err != nil {
		return err
	}
	for i := 1; i <= maxPage; i++ {
		if err = imp.ImportPage(i, i); err != nil {
			return err
		}
	}
	// Loop through the fields in the definition.
	for fd := range ft.AllFields() {
		// Save the message ID if we see it.
		if fd.Common == "originMessageID" {
			msgID = body.Field(fd.Tag)
		}
		// Loop through the PDF renderers in the field.
	RENDERER2:
		for _, pr := range fd.PDF {
			var (
				value string
				etr   pdf.ErrTextRendering
			)
			// Test the conditions.
			for _, cond := range pr.Conditions {
				value := body.Field(cond.Field)
				if (cond.Set && value == "") || (!cond.Set && value != cond.Value) {
					continue RENDERER2
				}
			}
			// Compute the value to be displayed.
			if fd.Type == "static" {
				value = fd.Value
			} else {
				value = body.Field(fd.Tag)
			}
			// Apply the renderer.
			err = pr.Renderer.Draw(out, value)
			if errors.As(err, &etr) {
				warnings = errors.Join(warnings, fmt.Errorf("field %s: %s", fd.Tag, err))
				err = nil
			}
			if err != nil {
				return err
			}
		}
	}
	// If we didn't find a message ID in the body, check the subject line.
	if msgID == "" {
		if s, ok := m.Subject().(*subject.SCCoSubject); ok {
			msgID = s.SubjectMessageID()
		}
	}
	// Write the page footers.
	if err = message.RenderPDFFooters(out, msgID, copyname); err != nil {
		return err
	}
	// Write the new PDF.
	if err = out.Write(); err != nil {
		return err
	}
	if err = outFH.Close(); err != nil {
		return err
	}
	if warnings != nil {
		return message.Warning{Err: warnings}
	}
	return nil
}

// recognize examines a message to see if it belongs to this form type, and if
// so, sets its MType to this form type and returns true.  Otherwise it returns
// false.
func (ft FormType) Recognize(m message.Message) {
	var body *FormBody

	// If it doesn't have a form body, it's clearly not our form.
	if body, _ = m.Body().(*FormBody); body == nil {
		return
	}
	// Check the addon, filename, and version number.
	if body.addonName != ft.AddonName || body.formHTML != ft.HTMLName || body.formVersion != ft.Version {
		return
	}
	// It's our form.
	m.SetType(ft)
}
