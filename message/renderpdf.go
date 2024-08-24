package message

// This file contains the BaseMessage implementation of Message.RenderPDF.  It
// also defines the PDFRenderer interface (value of Field.PDFRenderer) and
// provides several implementations of it.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
	"github.com/rothskeller/pdf/pdftext"
)

// ShowLayout is a global flag that can be set by callers.  When true, generated
// PDFs have a pink background in every field, showing the exact field placement.
var ShowLayout = true

// Warning is a wrapper around an error that makes it semantically a warning
// (operation completed with issues) rather than a true error (operation
// failed).  Callers of a function whose documentation says it can return a
// Warning can test whether it actually did using errors.As.
type Warning struct{ err error }

func (w Warning) Unwrap() error { return w.err }
func (w Warning) Error() string { return w.err.Error() }

// ErrNotSupported is the error returned if RenderPDF is called on a message
// with a type that does not support PDF rendering.
var ErrNotSupported = errors.New("message type does not support PDF rendering, or program was not built with -tags packetpdf")

// RenderPDF renders the message as a PDF file with the specified filename,
// overwriting any existing file with that name.  If the returned error is nil
// or is an instance of Warning, the specified PDF file was created; otherwise,
// no such file is left existing.
func (bm *BaseMessage) RenderPDF(_ *envelope.Envelope, filename string) (err error) {
	if bm.Type.PDFRenderV2 {
		return bm.renderPDFV2(filename)
	} else {
		return bm.renderPDFV1(filename)
	}
}

// RenderPDF renders the message as a PDF file with the specified filename,
// overwriting any existing file with that name.
func (bm *BaseMessage) renderPDFV1(filename string) (err error) {
	var (
		fh  *os.File
		pdf *pdfstruct.PDF
	)
	if bm.Type.PDFBase == nil {
		return ErrNotSupported
	}
	// First, write the base PDF.
	if fh, err = os.Create(filename); err != nil {
		return err
	}
	defer fh.Close()
	if _, err = fh.Write(bm.Type.PDFBase); err != nil {
		os.Remove(filename)
		return err
	}
	// Next, open it as a PDF.
	if pdf, err = pdfstruct.Open(fh); err != nil {
		os.Remove(filename)
		return err
	}
	// Update the fields of the PDF.
	for _, f := range bm.Fields {
		var pdfFields []PDFField
		if f.PDFMap != nil {
			pdfFields = f.PDFMap.RenderPDF(f)
		}
		for _, pf := range pdfFields {
			if pf.Size == 0 {
				pf.Size = bm.Type.PDFFontSize
			}
			if err = pdfform.SetField(pdf, pf.Name, pf.Value, pf.Size); err != nil {
				os.Remove(filename)
				return err
			}
		}
	}
	if err != nil {
		os.Remove(filename)
		return err
	}
	// Write the changes to the PDF.
	if err = pdf.Write(); err != nil {
		os.Remove(filename)
		return err
	}
	return nil
}

func (bm *BaseMessage) renderPDFV2(filename string) (err error) {
	var (
		rdr   io.ReadSeeker
		pdf   *gofpdf.Fpdf
		imp   *gofpdi.Importer
		sizes map[int]map[string]map[string]float64
		warn  Warning
		page  = 1
		nump  = 1
	)
	if bm.Type.PDFBase == nil {
		os.Remove(filename)
		return ErrNotSupported
	}
	// Create the output PDF and the importer from the base PDF.
	rdr = bytes.NewReader(bm.Type.PDFBase)
	pdf = gofpdf.New("P", "pt", "Letter", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.SetMargins(0, 0, 0)
	imp = gofpdi.NewImporter()
	// Walk through each page in the base PDF.
	for page <= nump {
		// Import the page.
		tpl := imp.ImportPageFromStream(pdf, &rdr, page, "/MediaBox")
		if sizes == nil {
			// After importing the first page, we can get the page
			// count and sizes.
			sizes = imp.GetPageSizes()
			nump = len(sizes)
		}
		// Create the page in the output PDF and copy the imported page
		// to it.
		orient := "P"
		w, h := sizes[page]["/MediaBox"]["w"], sizes[page]["/MediaBox"]["h"]
		if w > h {
			orient = "L"
		}
		pdf.AddPageFormat(orient, gofpdf.SizeType{Wd: w, Ht: h})
		imp.UseImportedTemplate(pdf, tpl, 0, 0, w, h)
		// Look for fields that need to be written to the page.
		for _, f := range bm.Fields {
			if f.PDFRenderer != nil {
				if err = f.PDFRenderer.RenderToPDF(f, pdf, page); err != nil {
					if !errors.As(err, &warn) {
						return err
					}
				}
			}
		}
		page++
	}
	// Write the resulting PDF.
	if err = pdf.OutputFileAndClose(filename); err != nil {
		os.Remove(filename)
		return err
	}
	if warn.err != nil {
		return warn
	}
	return nil
}

// PDFRenderer is the interface honored by a Field.PDFRenderer value.
type PDFRenderer interface {
	RenderToPDF(f *Field, pdf *gofpdf.Fpdf, page int) error
}

// PDFMultiRenderer is a PDFRenderer that invokes multiple sub-renderers.
type PDFMultiRenderer []PDFRenderer

func (mr PDFMultiRenderer) RenderToPDF(f *Field, pdf *gofpdf.Fpdf, page int) (err error) {
	for _, r := range mr {
		if rerr := r.RenderToPDF(f, pdf, page); rerr != nil && err == nil {
			err = rerr
		}
	}
	return err
}

// PDFTextRenderer is a PDFRenderer that renders the value of a field as text
// in a rectangle of the PDF page.
type PDFTextRenderer struct {
	// Page is the page number onto which the field should be rendered.
	// For convenience, zero is treated as 1.
	Page int
	// X, Y, W, and H are the position (upper left corner) and dimensions of
	// the box into which to render the text.
	X, Y, W, H float64
	// Style is the styling to apply to the text.
	Style PDFTextStyle
}
type PDFTextStyle = pdftext.Style

func (r *PDFTextRenderer) RenderToPDF(f *Field, pdf *gofpdf.Fpdf, page int) error {
	var (
		fits  bool
		style = pdftext.Style{MinFontSize: 10.0, LineHeight: 1.15, Color: []byte{0, 0, 153}, Wrap: 1}
	)
	if (r.Page == 0 && page != 1) || (r.Page != 0 && r.Page != page) {
		return nil
	}
	style = style.Merge(r.Style)
	if ShowLayout {
		pdf.SetFillColor(255, 0, 0)
		pdf.SetAlpha(0.5, "")
		pdf.Rect(r.X, r.Y, r.W, r.H, "F")
		pdf.SetAlpha(1.0, "")
	}
	if fits = pdftext.Draw(pdf, *f.Value, r.X, r.Y, r.W, r.H, style); !fits {
		return Warning{fmt.Errorf("value of %q does not fit in PDF", f.Label)}
	}
	return nil
}

// PDFStaticTextRenderer is a PDFRenderer that draws a static text string at
// the specified place.
type PDFStaticTextRenderer struct {
	// Page is the page number onto which the string should be rendered.
	// For convenience, zero is treated as 1.
	Page int
	// X, Y, and H are the position (upper left corner) and height of the
	// box into which to render the string.
	X, Y, H float64
	// Text is the static text string to render.
	Text string
}

func (r *PDFStaticTextRenderer) RenderToPDF(_ *Field, pdf *gofpdf.Fpdf, page int) error {
	if (r.Page == 0 && page != 1) || (r.Page != 0 && r.Page != page) {
		return nil
	}
	if ShowLayout {
		pdf.SetFillColor(255, 0, 0)
		pdf.SetAlpha(0.5, "")
		pdf.Rect(r.X, r.Y, 10, r.H, "F")
		pdf.Polygon([]gofpdf.PointType{{X: r.X + 10, Y: r.Y}, {X: r.X + 20, Y: r.Y + r.H/2}, {X: r.X + 10, Y: r.Y + r.H}}, "F")
		pdf.SetAlpha(1.0, "")
	}
	pdftext.Draw(pdf, r.Text, r.X, r.Y, 10000, r.H, pdftext.Style{
		LineHeight: 1.2, Color: []byte{0, 0, 153},
	})
	return nil
}

// PDFRadioRenderer is a PDFRenderer that draws a radio button in the
// appropriate place on a PDF page based on the value of a field.
type PDFRadioRenderer struct {
	// Page is the page number onto which the field should be rendered.
	// For convenience, zero is treated as 1.
	Page int
	// Points is a map from field value to radio button center point.  Each
	// entry is a slice of two numbers: x and then y.
	Points map[string][]float64
	// Radius is the radius of the radio button indicator.  The default
	// radius is 3.
	Radius float64
}

func (r *PDFRadioRenderer) RenderToPDF(f *Field, pdf *gofpdf.Fpdf, page int) error {
	var radius = 3.0

	if (r.Page == 0 && page != 1) || (r.Page != 0 && r.Page != page) {
		return nil
	}
	if r.Radius != 0 {
		radius = r.Radius
	}
	if ShowLayout {
		pdf.SetFillColor(255, 0, 0)
		pdf.SetAlpha(0.5, "")
		for _, pt := range r.Points {
			pdf.Circle(pt[0], pt[1], radius, "F")
		}
		pdf.SetAlpha(1.0, "")
	}
	if *f.Value == "" {
		return nil
	}
	if pt, ok := r.Points[*f.Value]; !ok {
		return Warning{fmt.Errorf("field %q: unknown value %q", f.Label, *f.Value)}
	} else {
		pdf.SetFillColor(0, 0, 153)
		pdf.Circle(pt[0], pt[1], radius, "F")
	}
	return nil
}

// PDFCheckRenderer is a PDFRenderer that draws an "X" in the appropriate place
// on a PDF page based on the value of a field.
type PDFCheckRenderer struct {
	// Page is the page number onto which the field should be rendered.
	// For convenience, zero is treated as 1.
	Page int
	// Points is a map from field value to checkbox top left point.  Each
	// entry is a slice of two numbers: x and then y.
	Points map[string][]float64
	// W and H are the width and height of the checkbox.
	W, H float64
}

func (r *PDFCheckRenderer) RenderToPDF(f *Field, pdf *gofpdf.Fpdf, page int) error {
	if (r.Page == 0 && page != 1) || (r.Page != 0 && r.Page != page) {
		return nil
	}
	if ShowLayout {
		pdf.SetFillColor(255, 0, 0)
		pdf.SetAlpha(0.5, "")
		for _, pt := range r.Points {
			pdf.Rect(pt[0], pt[1], r.W, r.H, "F")
		}
		pdf.SetAlpha(1.0, "")
	}
	if *f.Value == "" {
		return nil
	}
	if pt, ok := r.Points[*f.Value]; !ok {
		return Warning{fmt.Errorf("field %q: unknown value %q", f.Label, *f.Value)}
	} else {
		pdf.SetFillColor(0, 0, 153)
		pdf.SetDrawColor(0, 0, 153)
		pdf.SetLineWidth(r.W / 10)
		pdf.SetLineCapStyle("butt")
		pdf.ClipRect(pt[0], pt[1], r.W, r.H, false)
		pdf.Line(pt[0], pt[1], pt[0]+r.W, pt[1]+r.H)
		pdf.Line(pt[0]+r.W, pt[1], pt[0], pt[1]+r.H)
		pdf.ClipEnd()
	}
	return nil
}

// PDFMapper is the interface honored by a Field.PDFMap value.
type PDFMapper interface {
	RenderPDF(*Field) []PDFField
}

// PDFField describes the rendering of a single PDF field, as returned by the
// PDFMap interface in a Field structure.
type PDFField struct {
	// Name is the name of the field in the PDF.
	Name string
	// Value is the value to be placed in the PDF field.
	Value string
	// Size is the font size to be used for the PDF field; it can be zero
	// if supplied by the PDF template or the Type.PDFSize value.
	Size float64
}

// NoPDFField is a zero implementation of PDFMapper that returns no mappings.
type NoPDFField struct{}

func (NoPDFField) RenderPDF(*Field) []PDFField { return nil }

// PDFName is a simple implementation of PDFMapper.  It maps the field value,
// unedited, to the specified PDF field name, with the default font size.
type PDFName string

func (n PDFName) RenderPDF(f *Field) []PDFField {
	value := *f.Value
	if f.Choices != nil {
		value = f.Choices.ToHuman(value)
	}
	return []PDFField{{Name: string(n), Value: value}}
}

// PDFNameMap is an implementation of PDFMapper.  The first element of the slice
// is the PDF field name.  The remaining elements are (PIFO value, PDF value)
// pairs.  If the current value of the field matches any of the PIFO values in
// the slice, it is mapped to the corresponding PDF value.  Otherwise it is
// rendered unchanged.  All values are rendered with the default font size.
type PDFNameMap []string

func (m PDFNameMap) RenderPDF(f *Field) []PDFField {
	value := *f.Value
	var mapped bool
	for i := 1; i < len(m)-1; i += 2 {
		if value == m[i] {
			value, mapped = m[i+1], true
			break
		}
	}
	if !mapped && f.Choices != nil {
		value = f.Choices.ToHuman(value)
	}
	return []PDFField{{Name: m[0], Value: value}}
}

// PDFMapFunc converts a PDF mapping function into an interface that satisfies
// PDFMapper.
type PDFMapFunc func(*Field) []PDFField

func (fn PDFMapFunc) RenderPDF(f *Field) []PDFField {
	return fn(f)
}
