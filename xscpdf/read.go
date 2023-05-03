// Package xscpdf contains all of the data needed to read XSC-standard message
// content from fillable PDFs, or encode message content into fillable PDFs.
// It's a separate package from xscmsg because it embeds the actual PDF files,
// which are large, and because it introduces a dependency on the PDF libraries
// that may not be needed in some applications.
package xscpdf

import (
	"bytes"

	"github.com/rothskeller/packet/typedmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

// FieldMap is a mapping between an xscmsg.Message field and a PDF form field.
type FieldMap struct {
	// PDFName is the name of the fillable field in the PDF.
	PDFName string
	// XSCTag is the Tag of the field in the xscmsg.
	XSCTag string
	// Values, if set, is a mapping between values in the PDF and values in
	// the xscmsg.
	Values []ValueMap
	// FromXSC is a hook function to get the value of a PDF field when a
	// simple mapping won't suffice.
	FromXSC func(xscmsg.IMessage) string
	// FromPDF is a hook function to get the value of an XSC field when a
	// simple mapping won't suffice.
	FromPDF func(map[string]string) string
	// FontSize is the font size to be used for the field in the PDF.  Some
	// fields have a defined font size in the PDF, in which case this isn't
	// used.  But for those that don't, we need this.
	FontSize float64
}

// ValueMap is a mapping between a field value of an xscmsg.Message field and
// the corresponding value of the PDF form field.
type ValueMap struct {
	PDFValue string
	XSCValue string
}

// CheckboxMap is the set of Values for a FieldMap for a field that is a
// checkbox in both the PDF and the XSC form.
var CheckboxMap = []ValueMap{
	{PDFValue: "Off", XSCValue: ""},
	{PDFValue: "Yes", XSCValue: "checked"},
}

// A ReaderMap gives the xscmsg tag and field map associated with a PDF.
type ReaderMap struct {
	XSCTag string
	PDFID  []byte
	Fields []FieldMap
}

// readers is the list of registered readers.
var readers []ReaderMap

// RegisterReader registers a reader function for a particular PDF file.
func RegisterReader(rm ReaderMap) {
	readers = append(readers, rm)
}

// PDFToMessage converts a PDF file into a message.  The PDF file must be one of
// the known message PDFs, with form fields filled in; those values will be
// carried over into the newly created message.  If the PDF file is not known or
// any other problem occurs, PDFToMessage returns nil.
func PDFToMessage(pdf *pdfstruct.PDF) xscmsg.IMessage {
	var id []byte
	if a, ok := pdf.Info["ID"].(pdfstruct.Array); ok {
		if len(a) == 2 {
			if i, ok := a[0].([]byte); ok {
				id = i
			}
		}
	}
	if id == nil {
		return nil
	}
	for _, rm := range readers {
		if bytes.Equal(rm.PDFID, id) {
			msg := typedmsg.Create(rm.XSCTag).(xscmsg.IMessage)
			readFields(pdf, msg, rm.Fields)
			return msg
		}
	}
	return nil
}

// readFields reads fields from the PDF and applies them to the form according
// to the map.
func readFields(pdf *pdfstruct.PDF, msg xscmsg.IMessage, fmaps []FieldMap) {
	var (
		pfields map[string]string
		err     error
	)
	if pfields, err = pdfform.GetFields(pdf); err != nil {
		return
	}
	for pf, pv := range pfields {
		for _, fm := range fmaps {
			if pf == fm.PDFName {
				if fm.XSCTag == "" {
					continue
				}
				if fm.FromPDF != nil {
					pv = fm.FromPDF(pfields)
				} else if len(fm.Values) != 0 {
					for _, vm := range fm.Values {
						if pv == vm.PDFValue {
							pv = vm.XSCValue
							break
						}
					}
				}
				for _, f := range msg.GetTaggedFields() {
					if f.Tag == fm.XSCTag {
						f.Value = pv
						break
					}
				}
				break
			}
		}
	}
}
