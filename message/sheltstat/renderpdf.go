//xgo:build packetpdf

package sheltstat

import (
	_ "embed" // oh, please
	"os"

	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

//go:embed XSC_SheltStat_Fillable_v20190619_p12.pdf
var basePDF []byte

// RenderPDF renders the message as a PDF file with the specified filename,
// overwriting any existing file with that name.  Note that the program needs to
// be built with the packet-pdf build tag in order to include these methods.
func (f *SheltStat) RenderPDF(filename string) (err error) {
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
	err = updatePDFField(pdf, "Address", f.ShelterAddress, err)
	err = updatePDFField(pdf, "ATC20 Insp", translate(atc20InspectionMap, f.ATC20Inspection), err)
	err = updatePDFField(pdf, "Available Services", f.AvailableServices, err)
	err = updatePDFField(pdf, "Basic Safety Insp", translate(basicSafetyInspectionMap, f.BasicSafetyInspection), err)
	err = updatePDFField(pdf, "Capacity", f.Capacity, err)
	err = updatePDFField(pdf, "City", f.ShelterCity, err)
	err = updatePDFField(pdf, "Comments", f.Comments, err)
	err = updatePDFField(pdf, "Date Created", f.MessageDate, err)
	err = updatePDFField(pdf, "Destination Msg Nbr", f.DestinationMsgID, err)
	err = updatePDFField(pdf, "Floorplan", f.FloorPlan, err)
	err = updatePDFField(pdf, "From Contact Info", f.FromContact, err)
	err = updatePDFField(pdf, "From ICS Position", f.FromICSPosition, err)
	err = updatePDFField(pdf, "From Location", f.FromLocation, err)
	err = updatePDFField(pdf, "From Name", f.FromName, err)
	err = updatePDFField(pdf, "Handling", translate(handlingMap, f.Handling), err)
	err = updatePDFField(pdf, "Input Freq", f.RepeaterInput, err)
	err = updatePDFField(pdf, "Input Tone", f.RepeaterInputTone, err)
	err = updatePDFField(pdf, "Latitude", f.Latitude, err)
	err = updatePDFField(pdf, "Longitude", f.Longitude, err)
	err = updatePDFField(pdf, "Managed By Detail", f.ManagedByDetail, err)
	err = updatePDFField(pdf, "Managed By", translate(managedByMap, f.ManagedBy), err)
	err = updatePDFField(pdf, "Meals", f.MealsServed, err)
	err = updatePDFField(pdf, "MOU", f.MOU, err)
	err = updatePDFField(pdf, "NSS Number", f.NSSNumber, err)
	err = updatePDFField(pdf, "Occupancy", f.Occupancy, err)
	err = updatePDFField(pdf, "Offset", f.RepeaterOffset, err)
	err = updatePDFField(pdf, "Op Call Sign", f.OpCall, err)
	err = updatePDFField(pdf, "Op Date", f.OpDate, err)
	err = updatePDFField(pdf, "Op Name", f.OpName, err)
	err = updatePDFField(pdf, "Op Time", f.OpTime, err)
	err = updatePDFField(pdf, "Origin Msg Nbr Copy", f.OriginMsgID, err)
	err = updatePDFField(pdf, "Origin Msg Nbr", f.OriginMsgID, err)
	err = updatePDFField(pdf, "Output Freq", f.RepeaterOutput, err)
	err = updatePDFField(pdf, "Output Tone", f.RepeaterOutputTone, err)
	err = updatePDFField(pdf, "Pet Friendly", translate(petFriendlyMap, f.PetFriendly), err)
	err = updatePDFField(pdf, "Pri Contact Phone", f.PrimaryPhone, err)
	err = updatePDFField(pdf, "Pri Contact", f.PrimaryContact, err)
	err = updatePDFField(pdf, "Relay Rcvd", f.OpRelayRcvd, err)
	err = updatePDFField(pdf, "Relay Sent", f.OpRelaySent, err)
	err = updatePDFField(pdf, "Remove", translate(removeFromListMap, f.RemoveFromList), err)
	err = updatePDFField(pdf, "Repeater Call Sign", f.RepeaterCallSign, err)
	err = updatePDFField(pdf, "Report Type", translate(reportTypeMap, f.ReportType), err)
	err = updatePDFField(pdf, "Sec Contact Phone", f.SecondaryPhone, err)
	err = updatePDFField(pdf, "Sec Contact", f.SecondaryContact, err)
	err = updatePDFField(pdf, "Shelter Name", f.ShelterName, err)
	err = updatePDFField(pdf, "Shelter Status", translate(shelterStatusMap, f.ShelterStatus), err)
	err = updatePDFField(pdf, "Shelter Type", translate(shelterTypeMap, f.ShelterType), err)
	err = updatePDFField(pdf, "State", f.ShelterState, err)
	err = updatePDFField(pdf, "Tactical Call Sign", f.TacticalCallSign, err)
	err = updatePDFField(pdf, "Time Created", f.MessageTime, err)
	err = updatePDFField(pdf, "To Contact Info", f.ToContact, err)
	err = updatePDFField(pdf, "To ICS Position", f.ToICSPosition, err)
	err = updatePDFField(pdf, "To Location", f.ToLocation, err)
	err = updatePDFField(pdf, "To Name", f.ToName, err)
	err = updatePDFField(pdf, "Zip", f.ShelterZip, err)
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
	return pdfform.SetField(pdf, fieldname, value, 10)
}

func translate(m map[string]string, v string) string {
	if tv, ok := m[v]; ok {
		return tv
	}
	return v
}

var atc20InspectionMap = map[string]string{
	"":        "Off",
	"false":   "No",
	"checked": "Yes",
}
var basicSafetyInspectionMap = map[string]string{
	"":        "Off",
	"false":   "No",
	"checked": "Yes",
}
var handlingMap = map[string]string{
	"":          "Off",
	"PRIORITY":  "Priority",
	"ROUTINE":   "Routine",
	"IMMEDIATE": "Immediate",
}
var managedByMap = map[string]string{
	"": "Off",
}
var petFriendlyMap = map[string]string{
	"":        "Off",
	"false":   "No",
	"checked": "Yes",
}
var removeFromListMap = map[string]string{
	"":        "No",
	"checked": "Yes",
}
var reportTypeMap = map[string]string{
	"": "Off",
}
var shelterStatusMap = map[string]string{
	"": "Off",
}
var shelterTypeMap = map[string]string{
	"": "Off",
}
