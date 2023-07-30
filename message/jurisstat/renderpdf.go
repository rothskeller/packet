//go:build packetpdf

package jurisstat

import (
	_ "embed" // oh, please
	"os"

	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

//go:embed XSC_JurisStat_Fillable_v20190528_p123.pdf
var basePDF []byte

// RenderPDF renders the message as a PDF file with the specified filename,
// overwriting any existing file with that name.  Note that the program needs to
// be built with the packet-pdf build tag in order to include these methods.
func (f *JurisStat) RenderPDF(filename string) (err error) {
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
	err = updatePDFField(pdf, "Activation", translate(activationMap, f.EOCActivationLevel), err)
	err = updatePDFField(pdf, "Animal Comment", f.AnimalIssuesComments, err)
	err = updatePDFField(pdf, "Animal", translate(statusMap, f.AnimalIssues), err)
	err = updatePDFField(pdf, "Attachment", f.HowSOESent, err)
	err = updatePDFField(pdf, "Casualties Comment", f.CasualtiesComments, err)
	err = updatePDFField(pdf, "Casualties", translate(statusMap, f.Casualties), err)
	err = updatePDFField(pdf, "Civil Comment", f.CivilUnrestComments, err)
	err = updatePDFField(pdf, "Civil", translate(statusMap, f.CivilUnrest), err)
	err = updatePDFField(pdf, "Comm Comment", f.CommunicationsComments, err)
	err = updatePDFField(pdf, "Communications", translate(statusMap, f.Communications), err)
	err = updatePDFField(pdf, "Date Created", f.MessageDate, err)
	err = updatePDFField(pdf, "Debris Comment", f.DebrisComments, err)
	err = updatePDFField(pdf, "Debris", translate(statusMap, f.Debris), err)
	err = updatePDFField(pdf, "Destination Msg Nbr", f.DestinationMsgID, err)
	err = updatePDFField(pdf, "Em Svcs Comment", f.EmergencyServicesComments, err)
	err = updatePDFField(pdf, "Em Svcs", translate(statusMap, f.EmergencyServices), err)
	err = updatePDFField(pdf, "EOC Close Date", f.EOCExpectedCloseDate, err)
	err = updatePDFField(pdf, "EOC Close Time", f.EOCExpectedCloseTime, err)
	err = updatePDFField(pdf, "EOC Fax", f.EOCFax, err)
	err = updatePDFField(pdf, "EOC Open Date", f.EOCExpectedOpenDate, err)
	err = updatePDFField(pdf, "EOC Open Time", f.EOCExpectedOpenTime, err)
	err = updatePDFField(pdf, "EOC Open", translate(eocOpenMap, f.EOCOpen), err)
	err = updatePDFField(pdf, "EOC Phone", f.EOCPhone, err)
	err = updatePDFField(pdf, "Flood Comment", f.FloodingComments, err)
	err = updatePDFField(pdf, "Flooding", translate(statusMap, f.Flooding), err)
	err = updatePDFField(pdf, "From Contact Info", f.FromContact, err)
	err = updatePDFField(pdf, "From ICS Position", f.FromICSPosition, err)
	err = updatePDFField(pdf, "From Location", f.FromLocation, err)
	err = updatePDFField(pdf, "From Name", f.FromName, err)
	err = updatePDFField(pdf, "Handling", translate(handlingMap, f.Handling), err)
	err = updatePDFField(pdf, "Hazmat Comment", f.HazmatComments, err)
	err = updatePDFField(pdf, "Hazmat", translate(statusMap, f.Hazmat), err)
	err = updatePDFField(pdf, "Infra Pwr Comment", f.InfrastructurePowerComments, err)
	err = updatePDFField(pdf, "Infra Pwr", translate(statusMap, f.InfrastructurePower), err)
	err = updatePDFField(pdf, "Infra Sewer Comment", f.InfrastructureSewerComments, err)
	err = updatePDFField(pdf, "Infra Sewer", translate(statusMap, f.InfrastructureSewer), err)
	err = updatePDFField(pdf, "Infra Water Comment", f.InfrastructureWaterComments, err)
	err = updatePDFField(pdf, "Infra Water", translate(statusMap, f.InfrastructureWater), err)
	err = updatePDFField(pdf, "Jurisdiction Name", f.Jurisdiction, err)
	err = updatePDFField(pdf, "Office Close Date", f.GovExpectedCloseDate, err)
	err = updatePDFField(pdf, "Office Close Time", f.GovExpectedCloseTime, err)
	err = updatePDFField(pdf, "Office Open Date", f.GovExpectedOpenDate, err)
	err = updatePDFField(pdf, "Office Open Time", f.GovExpectedOpenTime, err)
	err = updatePDFField(pdf, "Office Status", translate(officeStatusMap, f.OfficeStatus), err)
	err = updatePDFField(pdf, "Op Call Sign", f.OpCall, err)
	err = updatePDFField(pdf, "Op Date", f.OpDate, err)
	err = updatePDFField(pdf, "Op Name", f.OpName, err)
	err = updatePDFField(pdf, "Op Time", f.OpTime, err)
	err = updatePDFField(pdf, "Origin Msg Nbr Copy", f.OriginMsgID, err)
	err = updatePDFField(pdf, "Origin Msg Nbr", f.OriginMsgID, err)
	err = updatePDFField(pdf, "Pri EM Contact Name", f.PriEMContactName, err)
	err = updatePDFField(pdf, "Pri EM Contact Phone", f.PriEMContactPhone, err)
	err = updatePDFField(pdf, "Relay Rcvd", f.OpRelayRcvd, err)
	err = updatePDFField(pdf, "Relay Sent", f.OpRelaySent, err)
	err = updatePDFField(pdf, "Report Type", translate(reportTypeMap, f.ReportType), err)
	err = updatePDFField(pdf, "SAR Comment", f.SearchAndRescueComments, err)
	err = updatePDFField(pdf, "SAR", translate(statusMap, f.SearchAndRescue), err)
	err = updatePDFField(pdf, "Sec EM Contact Name", f.SecEMContactName, err)
	err = updatePDFField(pdf, "Sec EM Contact Phone", f.SecEMContactPhone, err)
	err = updatePDFField(pdf, "State of Emergency", translate(soeMap, f.StateOfEmergency), err)
	err = updatePDFField(pdf, "Time Created", f.MessageTime, err)
	err = updatePDFField(pdf, "To Contact Info", f.ToContact, err)
	err = updatePDFField(pdf, "To ICS Position", f.ToICSPosition, err)
	err = updatePDFField(pdf, "To Location", f.ToLocation, err)
	err = updatePDFField(pdf, "To Name", f.ToName, err)
	err = updatePDFField(pdf, "Trans Bridges Comment", f.TransportationBridgesComments, err)
	err = updatePDFField(pdf, "Trans Bridges", translate(statusMap, f.TransportationBridges), err)
	err = updatePDFField(pdf, "Trans Roads Comment", f.TransportationRoadsComments, err)
	err = updatePDFField(pdf, "Trans Roads", translate(statusMap, f.TransportationRoads), err)
	err = updatePDFField(pdf, "Util Elec Comment", f.UtilitiesElectricComments, err)
	err = updatePDFField(pdf, "Util Elec", translate(statusMap, f.UtilitiesElectric), err)
	err = updatePDFField(pdf, "Util Gas Comment", f.UtilitiesGasComments, err)
	err = updatePDFField(pdf, "Util Gas", translate(statusMap, f.UtilitiesGas), err)
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

var activationMap = map[string]string{
	"": "Off",
}
var eocOpenMap = map[string]string{
	"": "Off",
}
var handlingMap = map[string]string{
	"":          "Off",
	"PRIORITY":  "Priority",
	"ROUTINE":   "Routine",
	"IMMEDIATE": "Immediate",
}
var officeStatusMap = map[string]string{
	"": "Off",
}
var reportTypeMap = map[string]string{
	"": "Off",
}
var soeMap = map[string]string{
	"": "Off",
}
var statusMap = map[string]string{
	"": "Off",
}
