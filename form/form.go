package form

import (
	"bytes"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/form/htmlop"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
	"github.com/rothskeller/pdf/v2"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type FormType struct {
	*formdef.FormDef
}

var _ message.MType = (*FormType)(nil)

type EditableFormType struct {
	FormType
}

var _ message.EditableMType = EditableFormType{}

// CreateTag returns the tag for creating the form.
func (ft EditableFormType) CreateTag() string { return ft.FormDef.CreateTag }

// CreateKey returns the key for creating the form.
func (ft EditableFormType) CreateKey() string { return ft.FormDef.CreateKey }

// Name returns the name of the message type, as a phrase in lower case (other
// than acronyms) starting with "a " or "an ".
func (ft FormType) Name() string { return ft.IndefName }

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
		msgID = m.Subject().(*FormSubject).SubjectMessageID()
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
	body.def = ft.FormDef
	// If we have a subject line, convert it to a form subject.  Otherwise,
	// manufacture a form subject from the form contents.  (This can happen
	// when Outpost sends us a message body without headers.)
	if m.Subject().EncodedSubject() != "" {
		m.SetSubject(formSubjectFromPlainSubject(m.Subject().(*subject.PlainSubject)))
	} else {
		var msgID, handling, summary string
		for fd := range ft.AllFields() {
			switch fd.Common {
			case field.COriginMessageID:
				msgID = body.Field(fd.Tag)
			case field.CHandling:
				handling = body.Field(fd.Tag)
			case field.CMessageSummary:
				summary = body.Field(fd.Tag)
			}
		}
		subj, _ := NewFormSubject(msgID, handling, ft.SubjectTag, summary)
		m.SetSubject(subj)
	}
	m.SetType(ft)
}

func (ft EditableFormType) Recognize(m message.Message) {
	if ft.FormType.Recognize(m); m.Type() != nil {
		m.SetType(ft)
	}
}

// NewDraft returns a new draft message of this form type.
func (ft EditableFormType) NewDraft() message.Message {
	var urgent bool

	body, _ := NewFormBody(ft.FormDef)
	for f := range ft.AllFields() {
		if f.Tag != "" && f.Value != "" {
			body.SetField(f.Tag, f.Value)
			if f.Common == field.CHandling && f.Value == "IMMEDIATE" {
				urgent = true
			}
		}
	}
	pload := payload.NewOutpostPayload(body)
	pload.SetUrgent(urgent)
	subj, _ := NewFormSubject("", "", ft.SubjectTag, "")
	return message.NewDraftMessage(ft, subj, pload, false)
}

// EditHTML returns the HTML form for editing the message.
func (ft EditableFormType) EditHTML(msg message.Message, vars message.EditHTMLVars) (out []byte, err error) {
	var (
		formFile []byte
		formHTML *html.Node
		bundle   string
		defFile  string
		body     *FormBody
		formBuf  bytes.Buffer
		fields   = make(map[string]string)
		values   = make(url.Values)
	)
	// Read and parse the HTML for the form.
	if formFile, err = fs.ReadFile(ft.FormFS, ft.HTMLFile); err != nil {
		slog.Error("fs.ReadFile", "f", ft.HTMLFile, "err", err)
		return nil, err
	}
	if formHTML, err = html.Parse(bytes.NewReader(formFile)); err != nil {
		slog.Error("html.Parse", "f", ft.HTMLFile, "err", err)
		return nil, fmt.Errorf("%s: %s", ft.HTMLFile, err)
	}
	bundle, _, _ = strings.Cut(ft.HTMLFile, "/")
	// If there is a definitions.html in the same directory, read and parse
	// it too, and prepend it to the form HTML.
	defFile = bundle + "/definitions.html"
	if formFile, err = fs.ReadFile(ft.FormFS, defFile); err == nil {
		var defHTML *html.Node
		if defHTML, err = html.Parse(bytes.NewReader(formFile)); err != nil {
			slog.Error("html.Parse", "f", defFile, "err", err)
			return nil, fmt.Errorf("%s: %s", defFile, err)
		}
		formBody := findBody(formHTML)
		defBody := findBody(defHTML)
		for c := defBody.LastChild; c != nil; c = defBody.LastChild {
			defBody.RemoveChild(c)
			formBody.InsertBefore(c, formBody.FirstChild)
		}
	}
	fields["assets"] = vars.AssetBase
	fields["submit-url"] = vars.SubmitURL
	fields["submit-label"] = vars.SubmitLabel
	fields["save-label"] = vars.SaveLabel
	if vars.ShowAddressFields {
		fields["show-addrs"] = "true"
	}
	fields["addon-name"] = ft.AddonName
	fields["form-html"] = ft.HTMLName
	fields["form-version"] = ft.Version
	if ft.PDFFile != "" {
		_, rmbundle, _ := strings.Cut(ft.PDFFile, "/")
		fields["pdf-url"] = path.Join(vars.AssetBase, rmbundle)
	}
	// Expand the templates in the form HTML, using the supplied fields.
	htmlop.Expand(formHTML, fields)
	// Fill in the form using the fields from the message.
	values.Set("ToAddr", msg.To())
	values.Set("FromAddr", vars.FromAddress)
	body = msg.Body().(*FormBody)
	for f := range ft.AllFields() {
		if f.Tag != "" {
			if v := body.Field(f.Tag); v != "" {
				values.Set(f.Tag, v)
			}
		}
	}
	htmlop.FillForm(formHTML, values)
	// Render and minimize the result.
	htmlop.Minify(&formBuf, formHTML)
	return formBuf.Bytes(), nil
}

func findBody(doc *html.Node) *html.Node {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Body {
			return n
		}
	}
	return nil
}

// EditAssets returns the file system containing the form assets.
func (ft EditableFormType) EditAssets() (assets fs.FS) {
	bundle, _, _ := strings.Cut(ft.HTMLFile, "/")
	assets, _ = fs.Sub(ft.FormFS, bundle)
	return assets
}

// FromPOST translates the HTML response back into a DraftMessage.
func (ft EditableFormType) FromPOST(r *http.Request) (msg message.Message, err error) {
	var (
		body *FormBody
		subj *FormSubject
		payl *payload.OutpostPayload
		dm   *message.DraftMessage
	)
	if ft.RenderBody != "" {
		// This is a message that is edited as a form but rendered as
		// plain text, e.g. a check-in message.  Handled separately.
		return ft.renderFromPOST(r)
	}
	if body, err = NewFormBody(ft.FormDef); err != nil {
		slog.Error("form.NewFormBody", "err", err)
		return nil, err
	}
	payl = payload.NewOutpostPayload(body)
	subj, _ = NewFormSubject("", "", ft.SubjectTag, "") // will give an error, ignored
	dm = message.NewDraftMessage(ft, subj, payl, false)
	for f := range body.Fields() {
		if tag := f.Tag(); tag != "" {
			if val := r.FormValue(tag); val != "" {
				f.SetValue(dm, strings.ReplaceAll(val, "\r", ""))
			}
		}
	}
	dm.SetTo(r.FormValue("ToAddr"))
	return dm, nil
}

// renderFromPOST translates the HTML response into a plain text DraftMessage.
func (ft EditableFormType) renderFromPOST(r *http.Request) (msg message.Message, err error) {
	// For both the subject summary and the body, we use the form fields as
	// variables to substitute into the rendering templates.
	var variables = make(map[string]string)
	r.FormValue("") // ensure form has been parsed
	for k := range r.Form {
		if v := r.FormValue(k); v != "" {
			variables[k] = v
		}
	}
	summary := renderString(ft.RenderSummary, variables)
	bodytext := renderString(ft.RenderBody, variables)
	// Walk through the fields to get the other subject line elements.
	var msgID, handling string
	for f := range ft.AllFields() {
		switch f.Common {
		case field.COriginMessageID:
			msgID = r.FormValue(f.Tag)
		case field.CHandling:
			if f.Tag != "" {
				handling = r.FormValue(f.Tag)
			} else {
				handling = f.Value
			}
		}
	}
	// Build the message.
	b := body.NewPlainBody(bodytext)
	p := payload.NewOutpostPayload(b)
	if handling == "IMMEDIATE" {
		p.SetUrgent(true)
	}
	s, _ := subject.NewPlainSubject(msgID, handling, summary)
	dm := message.NewDraftMessage(message.PlainMessage, s, p, false)
	dm.SetTo(r.FormValue("ToAddr"))
	return dm, nil
}

func renderString(tmpl string, variables map[string]string) string {
	var buf bytes.Buffer

	// First, parse the template.
	doc, _ := htmlop.Parse(strings.NewReader(tmpl))
	if doc == nil {
		return ""
	}
	// Next, apply the variables.
	htmlop.Expand(doc, variables)
	// Finally, render the result.
	htmlop.Minify(&buf, doc)
	return buf.String()
}
