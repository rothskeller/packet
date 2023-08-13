package message

// This file contains the BaseMessage implementation of Message.RenderPDF.  It
// also defines the PDFMapper interface (value of Field.PDFMap) and provides
// several implementations of it.

import (
	"errors"
	"os"

	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

// ErrNotSupported is the error returned if RenderPDF is called on a message
// with a type that does not support PDF rendering.
var ErrNotSupported = errors.New("message type does not support PDF rendering, or program was not built with -tags packetpdf")

// RenderPDF renders the message as a PDF file with the specified filename,
// overwriting any existing file with that name.
func (bm *BaseMessage) RenderPDF(filename string) (err error) {
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
