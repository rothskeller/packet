package jurisstat

import (
	_ "embed" // oh, please
	"encoding/hex"

	"github.com/rothskeller/packet/xscmsg/jurisstat"
	"github.com/rothskeller/packet/xscpdf"
)

//go:embed XSC_JurisStat_Fillable_v20190528_p123.pdf
var basePDF []byte

var idFull, _ = hex.DecodeString("19515d102a4c184c8f9092ec0e13116f")
var idP123, _ = hex.DecodeString("1bbf43d608388df7414a8053bb3e4e8a")

func init() {
	xscpdf.RegisterReader(xscpdf.ReaderMap{XSCTag: jurisstat.Tag, PDFID: idFull, Fields: fieldMap})
	xscpdf.RegisterReader(xscpdf.ReaderMap{XSCTag: jurisstat.Tag, PDFID: idP123, Fields: fieldMap})
	xscpdf.RegisterWriter(xscpdf.WriterMap{XSCTag: jurisstat.Tag, BasePDF: basePDF, Fields: fieldMap})
}

var fieldMap = []xscpdf.FieldMap{
	{PDFName: "Activation", XSCTag: "35.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Normal", XSCValue: "Normal"},
		{PDFValue: "Duty Officer", XSCValue: "Duty Officer"},
		{PDFValue: "Monitor", XSCValue: "Monitor"},
		{PDFValue: "Partial", XSCValue: "Partial"},
		{PDFValue: "Full", XSCValue: "Full"},
	}},
	{PDFName: "Animal Comment", XSCTag: "56.1.", FontSize: 10},
	{PDFName: "Animal", XSCTag: "56.0.", Values: statusValues},
	{PDFName: "Attachment", XSCTag: "99.", FontSize: 10},
	{PDFName: "Casualties Comment", XSCTag: "46.1.", FontSize: 10},
	{PDFName: "Casualties", XSCTag: "46.0.", Values: statusValues},
	{PDFName: "Civil Comment", XSCTag: "55.1.", FontSize: 10},
	{PDFName: "Civil", XSCTag: "55.0.", Values: statusValues},
	{PDFName: "Comm Comment", XSCTag: "41.1.", FontSize: 10},
	{PDFName: "Communications", XSCTag: "41.0.", Values: statusValues},
	{PDFName: "Date Created", XSCTag: "1a.", FontSize: 10},
	{PDFName: "Debris Comment", XSCTag: "42.1.", FontSize: 10},
	{PDFName: "Debris", XSCTag: "42.0.", Values: statusValues},
	{PDFName: "Destination Msg Nbr", XSCTag: "DestMsgNo", FontSize: 10},
	{PDFName: "Em Svcs Comment", XSCTag: "45.1.", FontSize: 10},
	{PDFName: "Em Svcs", XSCTag: "45.0.", Values: statusValues},
	{PDFName: "EOC Close Date", XSCTag: "38.", FontSize: 10},
	{PDFName: "EOC Close Time", XSCTag: "39.", FontSize: 10},
	{PDFName: "EOC Fax", XSCTag: "24.", FontSize: 10},
	{PDFName: "EOC Open Date", XSCTag: "36.", FontSize: 10},
	{PDFName: "EOC Open Time", XSCTag: "37.", FontSize: 10},
	{PDFName: "EOC Open", XSCTag: "34.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Unknown", XSCValue: "Unknown"},
		{PDFValue: "Yes", XSCValue: "Yes"},
		{PDFValue: "No", XSCValue: "No"},
	}},
	{PDFName: "EOC Phone", XSCTag: "23.", FontSize: 10},
	{PDFName: "Flood Comment", XSCTag: "43.1.", FontSize: 10},
	{PDFName: "Flooding", XSCTag: "43.0.", Values: statusValues},
	{PDFName: "From Contact Info", XSCTag: "8d.", FontSize: 10},
	{PDFName: "From ICS Position", XSCTag: "8a.", FontSize: 10},
	{PDFName: "From Location", XSCTag: "8b.", FontSize: 10},
	{PDFName: "From Name", XSCTag: "8c.", FontSize: 10},
	{PDFName: "Handling", XSCTag: "5.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Immediate", XSCValue: "IMMEDIATE"},
		{PDFValue: "Priority", XSCValue: "PRIORITY"},
		{PDFValue: "Routine", XSCValue: "ROUTINE"},
	}},
	{PDFName: "Hazmat Comment", XSCTag: "44.1.", FontSize: 10},
	{PDFName: "Hazmat", XSCTag: "44.0.", Values: statusValues},
	{PDFName: "Infra Pwr Comment", XSCTag: "49.1.", FontSize: 10},
	{PDFName: "Infra Pwr", XSCTag: "49.0.", Values: statusValues},
	{PDFName: "Infra Sewer Comment", XSCTag: "51.1.", FontSize: 10},
	{PDFName: "Infra Sewer", XSCTag: "51.0.", Values: statusValues},
	{PDFName: "Infra Water Comment", XSCTag: "50.1.", FontSize: 10},
	{PDFName: "Infra Water", XSCTag: "50.0.", Values: statusValues},
	{PDFName: "Jurisdiction Name", XSCTag: "22.", FontSize: 10},
	{PDFName: "Office Close Date", XSCTag: "32.", FontSize: 10},
	{PDFName: "Office Close Time", XSCTag: "33.", FontSize: 10},
	{PDFName: "Office Open Date", XSCTag: "30.", FontSize: 10},
	{PDFName: "Office Open Time", XSCTag: "31.", FontSize: 10},
	{PDFName: "Office Status", XSCTag: "29.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Unknown", XSCValue: "Unknown"},
		{PDFValue: "Open", XSCValue: "Open"},
		{PDFValue: "Closed", XSCValue: "Closed"},
	}},
	{PDFName: "Op Call Sign", XSCTag: "OpCall", FontSize: 10},
	{PDFName: "Op Date", XSCTag: "OpDate", FontSize: 10},
	{PDFName: "Op Name", XSCTag: "OpName", FontSize: 10},
	{PDFName: "Op Time", XSCTag: "OpTime", FontSize: 10},
	{PDFName: "Origin Msg Nbr Copy", XSCTag: "MsgNo", FontSize: 10},
	{PDFName: "Origin Msg Nbr", XSCTag: "MsgNo", FontSize: 10},
	{PDFName: "Pri EM Contact Name", XSCTag: "25.", FontSize: 10},
	{PDFName: "Pri EM Contact Phone", XSCTag: "26.", FontSize: 10},
	{PDFName: "Relay Rcvd", XSCTag: "OpRelayRcvd", FontSize: 10},
	{PDFName: "Relay Sent", XSCTag: "OpRelaySent", FontSize: 10},
	{PDFName: "Report Type", XSCTag: "19.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Update", XSCValue: "Update"},
		{PDFValue: "Complete", XSCValue: "Complete"},
	}},
	{PDFName: "SAR Comment", XSCTag: "52.1.", FontSize: 10},
	{PDFName: "SAR", XSCTag: "52.0.", Values: statusValues},
	{PDFName: "Sec EM Contact Name", XSCTag: "27.", FontSize: 10},
	{PDFName: "Sec EM Contact Phone", XSCTag: "28.", FontSize: 10},
	{PDFName: "State of Emergency", XSCTag: "40.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Unknown", XSCValue: "Unknown"},
		{PDFValue: "Yes", XSCValue: "Yes"},
		{PDFValue: "No", XSCValue: "No"},
	}},
	{PDFName: "Time Created", XSCTag: "1b.", FontSize: 10},
	{PDFName: "To Contact Info", XSCTag: "7d.", FontSize: 10},
	{PDFName: "To ICS Position", XSCTag: "7a.", FontSize: 10},
	{PDFName: "To Location", XSCTag: "7b.", FontSize: 10},
	{PDFName: "To Name", XSCTag: "7c.", FontSize: 10},
	{PDFName: "Trans Bridges Comment", XSCTag: "54.1.", FontSize: 10},
	{PDFName: "Trans Bridges", XSCTag: "54.0.", Values: statusValues},
	{PDFName: "Trans Roads Comment", XSCTag: "53.1.", FontSize: 10},
	{PDFName: "Trans Roads", XSCTag: "53.0.", Values: statusValues},
	{PDFName: "Util Elec Comment", XSCTag: "48.1.", FontSize: 10},
	{PDFName: "Util Elec", XSCTag: "48.0.", Values: statusValues},
	{PDFName: "Util Gas Comment", XSCTag: "47.1.", FontSize: 10},
	{PDFName: "Util Gas", XSCTag: "47.0.", Values: statusValues},
}

var statusValues = []xscpdf.ValueMap{
	{PDFValue: "Off", XSCValue: ""},
	{PDFValue: "Unknown", XSCValue: "Unknown"},
	{PDFValue: "Normal", XSCValue: "Normal"},
	{PDFValue: "Problem", XSCValue: "Problem"},
	{PDFValue: "Failure", XSCValue: "Failure"},
	{PDFValue: "Delayed", XSCValue: "Delayed"},
	{PDFValue: "Closed", XSCValue: "Closed"},
	{PDFValue: "Early Out", XSCValue: "Early Out"},
}
