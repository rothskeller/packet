package basemsg

import (
	"os"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

// TODO defining these functions on BaseMessage means that only messages for
// which they're appropriate can leverage BaseMessage.

// RenderTable renders the message as a set of field label / field value
// pairs, intended for read-only display to a human.
func (bm *BaseMessage) RenderTable() (lvs []message.LabelValue) {
	for _, f := range bm.Fields {
		if value := f.TableValue(f); value != "" {
			lvs = append(lvs, message.LabelValue{Label: f.Label, Value: value})
		}
	}
	return lvs
}

// TableOmit is a TableValue function that causes the field to be
// unconditionally omitted from the table rendering.
func TableOmit(*Field) string { return "" }

// RenderPDF renders the message as a PDF file with the specified filename,
// overwriting any existing file with that name.
func (bm *BaseMessage) RenderPDF(filename string) (err error) {
	var (
		fh  *os.File
		pdf *pdfstruct.PDF
	)
	// First, write the base PDF.
	if fh, err = os.Create(filename); err != nil {
		return err
	}
	defer fh.Close()
	if _, err = fh.Write(bm.PDFBase); err != nil {
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
				pf.Size = bm.PDFFontSize
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
