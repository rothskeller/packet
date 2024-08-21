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
	"slices"
	"strings"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/pdf/pdffont"
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

// PDFTextRenderer is a PDFRenderer that renders the value of a field as text
// in a rectangle of the PDF page.
type PDFTextRenderer struct {
	// Page is the page number onto which the field should be rendered.
	// For convenience, zero is treated as 1.
	Page int
	// X, Y, W, and H are the position (upper left corner) and dimensions of
	// the box into which to render the text.
	X, Y, W, H float64
	// Font is the name of the font to use to render the text.  It must be
	// one of the fonts supported by pdffont.Measure.  It defaults to
	// Helvetica.
	Font string
	// MaxFontSize is the largest font that can be used to render the text.
	// It will be used unless the text doesn't fit at that size.
	MaxFontSize float64
	// MinFontSize is the smallest font that can be used to render the text.
	// It will be used only if the text doesn't fit at a larger size.
	MinFontSize float64
	// VAlign indicates how the text will be vertically aligned:  "top",
	// "center", or "baseline".  ("baseline" is like "center" except it
	// centers using the bounding box of all ASCII characters rather than
	// only the ones used in the string.  This allows multiple fields with
	// the same vertical extent to have aligned baselines.)  The default is
	// "center".
	VAlign string
}

func (r *PDFTextRenderer) RenderToPDF(f *Field, pdf *gofpdf.Fpdf, page int) error {
	var (
		wrapped  []string
		fits     bool
		top      float64
		fontSize = 12.0
	)
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Font == "" {
		r.Font = "Helvetica"
	}
	if r.MaxFontSize == 0 {
		r.MaxFontSize = 12.0
	}
	if r.MinFontSize == 0 {
		r.MinFontSize = 10.0
	}
	if r.VAlign == "" || (r.VAlign != "top" && r.VAlign != "baseline") {
		r.VAlign = "center"
	}
	if page != r.Page {
		return nil
	}
	if ShowLayout {
		pdf.SetFillColor(255, 0, 0)
		pdf.SetAlpha(0.5, "")
		pdf.Rect(r.X, r.Y, r.W, r.H, "F")
		pdf.SetAlpha(1.0, "")
	}
	if *f.Value == "" {
		return nil
	}
	for {
		if wrapped, fits = r.fitText(*f.Value, fontSize); fits {
			break
		}
		if fontSize-0.5 < r.MinFontSize {
			break
		}
		fontSize -= 0.5
	}
	if r.VAlign == "top" || !fits {
		_, habove, _ := pdffont.Measure(wrapped[0], r.Font, fontSize)
		top = r.Y + habove
	} else if r.VAlign == "baseline" {
		habove, hbelow := pdffont.FontMetrics(r.Font, fontSize)
		height := float64(len(wrapped)-1)*fontSize*1.2 + habove + hbelow
		top = r.Y + (r.H-height)/2 + habove
	} else { // r.VAlign == "center"
		_, habove, _ := pdffont.Measure(wrapped[0], r.Font, fontSize)
		_, _, hbelow := pdffont.Measure(wrapped[len(wrapped)-1], r.Font, fontSize)
		height := float64(len(wrapped)-1)*fontSize*1.2 + habove + hbelow
		top = r.Y + (r.H-height)/2 + habove
	}
	pdf.SetFont(r.Font, "", fontSize)
	pdf.SetTextColor(0, 0, 153)
	for _, line := range wrapped {
		pdf.Text(r.X, top, line)
		top += fontSize * 1.2
	}
	if !fits {
		return Warning{fmt.Errorf("value of %q does not fit in PDF", f.Label)}
	}
	return nil
}

// fitText determines whether the text value fits in the renderer box at the
// specified font size, and how it got word-wrapped in order to fit.
func (r *PDFTextRenderer) fitText(value string, fontSize float64) (wrapped []string, fits bool) {
	// Streamline special case of empty string.
	if value == "" {
		return nil, true
	}
	// Start by assuming it will fit, until we find out otherwise.
	var height = r.H
	fits = true
	// Break the string up into lines and handle each one separately.
	var lines = strings.Split(value, "\n")
	for i := 0; i < len(lines); i++ {
		var stop = len(lines[i])
		for {
			// Measure the line to see if it fits.
			if w, _, _ := pdffont.Measure(lines[i][:stop], r.Font, fontSize); w > r.W {
				// It doesn't fit.  Is there a non-initial run
				// of spaces in it, such that we can word-wrap?
				if idx := strings.LastIndexByte(lines[i][:stop], ' '); idx > 0 {
					for ; idx > 0 && lines[i][idx-1] == ' '; idx-- {
					}
					if idx > 0 {
						// Yes.  Stop the line at that
						// point and try again.
						stop = idx
						continue
					}
				}
				// Can't word wrap (any further).  The whole
				// value will not fit.  We'll accept truncating
				// this line, but we'll still continue
				// word-wrapping the rest of the lines to do the
				// best we can.
				fits = false
			}
			// Remove the line's vertical size from bbox.
			height -= fontSize * 1.2
			// If we had to take a tail off the line to word wrap,
			// put that into the slice as the next line, and remove
			// it from the current line.
			var rest int
			for rest = stop; rest < len(lines[i]) && lines[i][rest] == ' '; rest++ {
			}
			if rest < len(lines[i]) {
				lines = slices.Insert(lines, i+1, lines[i][rest:])
			}
			if stop < len(lines[i]) {
				lines[i] = lines[i][:stop]
			}
			// Move on to the next line.
			break
		}
	}
	// We don't need leading after the last line.
	height += fontSize * 0.2
	// Did the value fit vertically?
	if height < 0 {
		fits = false
	}
	// Return the result.
	return lines, fits
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
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Radius == 0 {
		r.Radius = 3
	}
	if page != r.Page {
		return nil
	}
	if ShowLayout {
		pdf.SetFillColor(255, 0, 0)
		pdf.SetAlpha(0.5, "")
		for _, pt := range r.Points {
			pdf.Circle(pt[0], pt[1], r.Radius, "F")
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
		pdf.Circle(pt[0], pt[1], r.Radius, "F")
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
