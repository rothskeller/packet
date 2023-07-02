//go:build packetpdf

package ics213

import (
	_ "embed" // oh, please
	"os"

	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

//go:embed ICS-213_SCCo_Message_Form_Fillable_v20220119_p1.pdf
var basePDF []byte

// RenderPDF renders the message as a PDF file with the specified
// filename, overwriting any existing file with that name.  Note that
// the program needs to be built with the packet-pdf build tag in order
// to include these methods.
func (f *ICS213) RenderPDF(filename string) (err error) {
	var (
		fh  *os.File
		pdf *pdfstruct.PDF
	)
	// First, write the base PDF.
	if fh, err = os.Create(filename); err != nil {
		return err
	}
	defer fh.Close()
	if _, err = fh.Write(basePDF); err != nil {
		os.Remove(filename)
		return err
	}
	// Next, open it as a PDF.
	if pdf, err = pdfstruct.Open(fh); err != nil {
		os.Remove(filename)
		return err
	}
	// Update the fields of the PDF.
	err = updatePDFField(pdf, "Immediate", translate(handlingMap, f.Handling), err)
	err = updatePDFField(pdf, "TakeAction", translate(takeActionMap, f.TakeAction), err)
	err = updatePDFField(pdf, "Reply", translate(replyMap, f.Reply), err)
	err = updatePDFField(pdf, "How: Received", translate(recSentMap, f.ReceivedSent), err)
	err = updatePDFField(pdf, "Telephone", translate(methodMap, f.TxMethod), err)
	err = updatePDFField(pdf, "Origin Msg #", f.OriginMsgID, err)
	err = updatePDFField(pdf, "Destination Msg#", f.DestinationMsgID, err)
	err = updatePDFField(pdf, "FormDate", f.Date, err)
	err = updatePDFField(pdf, "FormTime", f.Time, err)
	err = updatePDFField(pdf, "Reply_2", f.ReplyBy, err)
	err = updatePDFField(pdf, "TO ICS Position", f.ToICSPosition, err)
	err = updatePDFField(pdf, "TO ICS Locatoin", f.ToLocation, err)
	err = updatePDFField(pdf, "TO ICS Name", f.ToName, err)
	err = updatePDFField(pdf, "TO ICS Telephone", f.ToTelephone, err)
	err = updatePDFField(pdf, "From ICS Position", f.FromICSPosition, err)
	err = updatePDFField(pdf, "From ICS Location", f.FromLocation, err)
	err = updatePDFField(pdf, "From ICS Name", f.FromName, err)
	err = updatePDFField(pdf, "From ICS Telephone", f.FromTelephone, err)
	err = updatePDFField(pdf, "Subject", f.Subject, err)
	err = updatePDFField(pdf, "Reference", f.Reference, err)
	err = updatePDFField(pdf, "Message", f.Message, err)
	err = updatePDFField(pdf, "Relay Received", f.OpRelayRcvd, err)
	err = updatePDFField(pdf, "Relay Sent", f.OpRelaySent, err)
	err = updatePDFField(pdf, "Operation Call Sign", f.OpCall, err)
	err = updatePDFField(pdf, "Relay Received_2", f.OpName, err)
	err = updatePDFField(pdf, "OperatorDate", f.OpDate, err)
	err = updatePDFField(pdf, "OperatorTime", f.OpTime, err)
	err = updatePDFField(pdf, "OtherText", f.OtherMethod, err)
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

func updatePDFField(pdf *pdfstruct.PDF, fieldname, value string, err error) error {
	if err != nil {
		return err
	}
	return pdfform.SetField(pdf, fieldname, value, 0)
}

func translate(m map[string]string, v string) string {
	if tv, ok := m[v]; ok {
		return tv
	}
	return v
}

var handlingMap = map[string]string{
	"":          "Off",
	"PRIORITY":  "1",
	"ROUTINE":   "2",
	"IMMEDIATE": "3",
}
var takeActionMap = map[string]string{
	"":    "Off",
	"Yes": "1",
	"No":  "2",
}
var replyMap = map[string]string{
	"":    "Off",
	"Yes": "Reply-Yes",
	"No":  "Reply-No",
}
var recSentMap = map[string]string{
	"":         "Off",
	"sender":   "0",
	"receiver": "1",
}
var methodMap = map[string]string{
	"":                "Off",
	"Telephone":       "1",
	"Dispatch Center": "2",
	"EOC Radio":       "3",
	"FAX":             "4",
	"Courier":         "5",
	"Amateur Radio":   "6",
	"Other":           "7",
}
