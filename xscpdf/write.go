package xscpdf

import (
	"errors"
	"fmt"
	"os"

	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

// A WriterMap gives the base PDF file and field map associated with a message
// type tag.
type WriterMap struct {
	XSCTag  string
	BasePDF []byte
	Fields  []FieldMap
}

// ErrNoWriter is the error returned when trying to write a PDF for a form that
// has no writer.
var ErrNoWriter = errors.New("no writer registered for this message type")

// writers is the list of registered writers.
var writers []WriterMap

// RegisterWriter registers a writer for a message type.
func RegisterWriter(wm WriterMap) {
	writers = append(writers, wm)
}

// MessageToPDF creates a PDF with the specified filename from the specified
// message.
func MessageToPDF(m xscmsg.IMessage, filename string) (err error) {
	for _, wm := range writers {
		if wm.XSCTag == m.Type().Tag {
			return messageToPDF(m, filename, wm)
		}
	}
	return ErrNoWriter
}

// messageToPDF creates a PDF for the specified message.
func messageToPDF(m xscmsg.IMessage, filename string, wm WriterMap) (err error) {
	var (
		fh  *os.File
		pdf *pdfstruct.PDF
	)
	// First, write the base PDF.
	if fh, err = os.Create(filename); err != nil {
		return err
	}
	defer fh.Close()
	if _, err = fh.Write(wm.BasePDF); err != nil {
		os.Remove(filename)
		return err
	}
	// Next, open it as a PDF.
	if pdf, err = pdfstruct.Open(fh); err != nil {
		os.Remove(filename)
		return err
	}
	// Walk through the mapped fields and update the PDF.
	for _, fm := range wm.Fields {
		var v string

		if fm.PDFName == "" {
			continue
		}
		if fm.FromXSC != nil {
			v = fm.FromXSC(m)
		} else {
			for _, f := range m.GetTaggedFields() {
				if f.Tag == fm.XSCTag {
					v = f.Value
					break
				}
			}
			if len(fm.Values) != 0 {
				for _, vm := range fm.Values {
					if vm.XSCValue == v {
						v = vm.PDFValue
						break
					}
				}
			}
		}
		if err = pdfform.SetField(pdf, fm.PDFName, v, fm.FontSize); err != nil { // TODO
			os.Remove(filename)
			return fmt.Errorf("field %s: %s", fm.XSCTag, err)
		}
	}
	// Write the changes to the PDF.
	if err = pdf.Write(); err != nil {
		os.Remove(filename)
		return err
	}
	return nil
}
