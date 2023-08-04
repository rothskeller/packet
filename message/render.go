package message

import (
	"errors"
	"os"

	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

// TableOmit is a TableValue function that causes the field to be
// unconditionally omitted from the table rendering.
func TableOmit(*Field) string { return "" }

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
