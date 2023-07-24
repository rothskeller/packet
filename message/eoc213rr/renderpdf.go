//go:build packetpdf

package eoc213rr

import (
	_ "embed" // oh, please
	"os"

	"github.com/rothskeller/packet/message/common"
	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

//go:embed XSC_EOC-213RR_Fillable_v20170803_with_XSC_RACES_Routing_Slip_Fillable_v20190527.pdf
var basePDF []byte

// RenderPDF renders the message as a PDF file with the specified filename,
// overwriting any existing file with that name.  Note that the program needs to
// be built with the packet-pdf build tag in order to include these methods.
func (f *EOC213RR) RenderPDF(filename string) (err error) {
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
	var withsig string
	if f.WithSignature != "" {
		withsig = "[with signature]"
	}
	err = updatePDFField(pdf, "Origin Msg Nbr", f.OriginMsgID, err)
	err = updatePDFField(pdf, "Destination Msg Nbr", f.DestinationMsgID, err)
	err = updatePDFField(pdf, "Date Created", f.MessageDate, err)
	err = updatePDFField(pdf, "Time Created", f.MessageTime, err)
	err = updatePDFField(pdf, "Handling", translate(handlingMap, f.Handling), err)
	err = updatePDFField(pdf, "To ICS Position", f.ToICSPosition, err)
	err = updatePDFField(pdf, "To Location", f.ToLocation, err)
	err = updatePDFField(pdf, "To Name", f.ToName, err)
	err = updatePDFField(pdf, "To Contact Info", f.ToContact, err)
	err = updatePDFField(pdf, "From ICS Position", f.FromICSPosition, err)
	err = updatePDFField(pdf, "From Location", f.FromLocation, err)
	err = updatePDFField(pdf, "From Name", f.FromName, err)
	err = updatePDFField(pdf, "From Contact Info", f.FromContact, err)
	err = updatePDFField(pdf, "Form Type", "XSC EOC-213RR", err)
	err = updatePDFField(pdf, "Form Topic", f.IncidentName, err)
	err = updatePDFField(pdf, "Relay Rcvd", f.OpRelayRcvd, err)
	err = updatePDFField(pdf, "Relay Sent", f.OpRelaySent, err)
	err = updatePDFField(pdf, "Op Name", f.OpName, err)
	err = updatePDFField(pdf, "Op Call Sign", f.OpCall, err)
	err = updatePDFField(pdf, "Op Date", f.OpDate, err)
	err = updatePDFField(pdf, "Op Time", f.OpTime, err)
	err = updatePDFField(pdf, "1 Incident Name", f.IncidentName, err)
	err = updatePDFField(pdf, "2 Date Initiated", f.DateInitiated, err)
	err = updatePDFField(pdf, "3 Time Initiated", f.TimeInitiated, err)
	err = updatePDFField(pdf, "5 Requested By", f.RequestedBy, err)
	err = updatePDFField(pdf, "6 Prepared by", f.PreparedBy, err)
	err = updatePDFField(pdf, "7 Approved By", common.SmartJoin(f.ApprovedBy, withsig, "\n"), err)
	err = updatePDFField(pdf, "8 QtyUnit", f.QtyUnit, err)
	err = updatePDFField(pdf, "9 Resource Description", f.ResourceDescription, err)
	err = updatePDFField(pdf, "10 Arrival", f.ResourceArrival, err)
	err = updatePDFField(pdf, "11 Priority", translate(priorityMap, f.Priority), err)
	err = updatePDFField(pdf, "12 Estd Cost", f.EstdCost, err)
	err = updatePDFField(pdf, "13 Deliver To", f.DeliverTo, err)
	err = updatePDFField(pdf, "14 Location", f.DeliverToLocation, err)
	err = updatePDFField(pdf, "15 Sub Sugg Sources", f.Substitutes, err)
	err = updatePDFField(pdf, "Equip Oper", translate(checkboxMap, f.EquipmentOperator), err)
	err = updatePDFField(pdf, "Fuel", translate(checkboxMap, f.Fuel), err)
	err = updatePDFField(pdf, "Fuel Type", f.FuelType, err)
	err = updatePDFField(pdf, "Meals", translate(checkboxMap, f.Meals), err)
	err = updatePDFField(pdf, "Water", translate(checkboxMap, f.Water), err)
	err = updatePDFField(pdf, "Lodging", translate(checkboxMap, f.Lodging), err)
	err = updatePDFField(pdf, "Power", translate(checkboxMap, f.Power), err)
	err = updatePDFField(pdf, "Maintenance", translate(checkboxMap, f.Maintenance), err)
	err = updatePDFField(pdf, "Other", translate(checkboxMap, f.Other), err)
	err = updatePDFField(pdf, "17 Special Instructions", f.Instructions, err)
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
	return pdfform.SetField(pdf, fieldname, value, 12)
}

func translate(m map[string]string, v string) string {
	if tv, ok := m[v]; ok {
		return tv
	}
	return v
}

var handlingMap = map[string]string{
	"":          "Off",
	"PRIORITY":  "Priority",
	"ROUTINE":   "Routine",
	"IMMEDIATE": "Immediate",
}
var priorityMap = map[string]string{
	"": "Off",
}
var checkboxMap = map[string]string{
	"":        "Off",
	"false":   "Off",
	"checked": "Yes",
}
