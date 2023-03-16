package racesmar

import (
	_ "embed" // oh, please
	"encoding/hex"

	"github.com/rothskeller/packet/xscmsg/racesmar"
	"github.com/rothskeller/packet/xscpdf"
)

//go:embed XSC_RACES_MA_Req_Fillable_v20220129_p12.pdf
var basePDF []byte

var idFull, _ = hex.DecodeString("e13ec09cfa3d89595da86324c4097b7b")
var idP12, _ = hex.DecodeString("e3f731ea0af1bb1e3a40124bc9526c20")

func init() {
	xscpdf.RegisterReader(xscpdf.ReaderMap{XSCTag: racesmar.Tag, PDFID: idFull, Fields: fieldMap})
	xscpdf.RegisterReader(xscpdf.ReaderMap{XSCTag: racesmar.Tag, PDFID: idP12, Fields: fieldMap})
	xscpdf.RegisterWriter(xscpdf.WriterMap{XSCTag: racesmar.Tag, BasePDF: basePDF, Fields: fieldMap})
}

var fieldMap = []xscpdf.FieldMap{
	{PDFName: "Immediate", XSCTag: "5.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "1", XSCValue: "IMMEDIATE"},
		{PDFValue: "2", XSCValue: "PRIORITY"},
		{PDFValue: "3", XSCValue: "ROUTINE"},
	}}, // PDFFields[0]
	{PDFName: "ApprovedDate_2", XSCTag: "MsgNo"},     // PDFFields[1]
	{PDFName: "OriginMsg", XSCTag: "MsgNo"},          // PDFFields[2]
	{PDFName: "DestinationMsg", XSCTag: "DestMsgNo"}, // PDFFields[3]
	{PDFName: "FormDate", XSCTag: "1a."},             // PDFFields[4]
	{PDFName: "FormTime", XSCTag: "1b."},             // PDFFields[5]
	{PDFName: "ToICSPosition", XSCTag: "7a."},        // PDFFields[6]
	{PDFName: "ToICSLocation", XSCTag: "7b."},        // PDFFields[7]
	{PDFName: "ToICSPosition_3", XSCTag: "7c."},      // PDFFields[8]
	{PDFName: "ToICSPosition_4", XSCTag: "7d."},      // PDFFields[9]
	{PDFName: "ToICSPosition_2", XSCTag: "8a."},      // PDF Fields[10]
	{PDFName: "ToICSLocation_2", XSCTag: "8b."},      // PDF Fields[11]
	{PDFName: "ToICSPosition_5", XSCTag: "8c."},      // PDF Fields[12]
	{PDFName: "ToICSPosition_6", XSCTag: "8d."},      // PDF Fields[13]
	{PDFName: "ToICSPosition_7", XSCTag: "15."},      // PDF Fields[14]
	{PDFName: "ToICSPosition_8", XSCTag: "16a."},     // PDF Fields[15]
	{PDFName: "Nbr", XSCTag: "16b."},                 // PDF Fields[16]
	{PDFName: "Assignment", XSCTag: "17."},           // PDF Fields[17]
	{PDFName: "Qty1", XSCTag: "18.1a."},              // PDF Fields[18]
	{PDFName: "Qty2", XSCTag: "18.2a."},              // PDF Fields[19]
	{PDFName: "Qty3", XSCTag: "18.3a."},              // PDF Fields[20]
	{PDFName: "Qty4", XSCTag: "18.4a."},              // PDF Fields[21]
	{PDFName: "Qty5", XSCTag: "18.5a."},              // PDF Fields[22]
	{PDFName: "Role1", XSCTag: "18.1e."},             // PDF Fields[23]
	{PDFName: "Role2", XSCTag: "18.2e."},             // PDF Fields[24]
	{PDFName: "Role3", XSCTag: "18.3e."},             // PDF Fields[25]
	{PDFName: "Role4", XSCTag: "18.4e."},             // PDF Fields[26]
	{PDFName: "Role5", XSCTag: "18.5e."},             // PDF Fields[27]
	{PDFName: "Position1", XSCTag: "18.1f."},         // PDF Fields[28]
	{PDFName: "Position2", XSCTag: "18.2f."},         // PDF Fields[29]
	{PDFName: "Position3", XSCTag: "18.3f."},         // PDF Fields[30]
	{PDFName: "Position4", XSCTag: "18.4f."},         // PDF Fields[31]
	{PDFName: "Position5", XSCTag: "18.5f."},         // PDF Fields[32]
	{PDFName: "Pref1", XSCTag: "18.1c."},             // PDF Fields[33]
	{PDFName: "Pref2", XSCTag: "18.2c."},             // PDF Fields[34]
	{PDFName: "Pref3", XSCTag: "18.3c."},             // PDF Fields[35]
	{PDFName: "Pref4", XSCTag: "18.4c."},             // PDF Fields[36]
	{PDFName: "Pref5", XSCTag: "18.5c."},             // PDF Fields[37]
	{PDFName: "Min1", XSCTag: "18.1d."},              // PDF Fields[38]
	{PDFName: "Min2", XSCTag: "18.2d."},              // PDF Fields[39]
	{PDFName: "Min3", XSCTag: "18.3d."},              // PDF Fields[40]
	{PDFName: "Min4", XSCTag: "18.4d."},              // PDF Fields[41]
	{PDFName: "Min5", XSCTag: "18.5d."},              // PDF Fields[42]
	{PDFName: "ReqArriveDates", XSCTag: "19a."},      // PDF Fields[43]
	{PDFName: "NeedUntilDates", XSCTag: "20a."},      // PDF Fields[44]
	{PDFName: "ReqArriveTimes", XSCTag: "19b."},      // PDF Fields[45]
	{PDFName: "NeedUntilTimes", XSCTag: "20b."},      // PDF Fields[46]
	{PDFName: "ReportingLocation", XSCTag: "21."},    // PDF Fields[47]
	{PDFName: "ContactOnArrival", XSCTag: "22."},     // PDF Fields[49]
	{PDFName: "TravelInfo", XSCTag: "23."},           // PDF Fields[50]
	{PDFName: "RequestedName", XSCTag: "24a."},       // PDF Fields[51]
	{PDFName: "RequestedTitle", XSCTag: "24b."},      // PDF Fields[52]
	{PDFName: "ApprovedTitle", XSCTag: "25b."},       // PDF Fields[53]
	{PDFName: "RequestedContact", XSCTag: "24c."},    // PDF Fields[54]
	{PDFName: "ApprovedName", XSCTag: "25a."},        // PDF Fields[55]
	{PDFName: "ApprovedContact", XSCTag: "25c."},     // PDF Fields[56]
	{PDFName: "ApprovedTime", XSCTag: "26b."},        // PDF Fields[58]
	{PDFName: "ApprovedDate", XSCTag: "26a."},        // PDF Fields[59]
	{PDFName: "OperatorRcvd", XSCTag: "OpRelayRcvd"}, // PDF Fields[60]
	{PDFName: "OperatorName", XSCTag: "OpName"},      // PDF Fields[61]
	{PDFName: "OperatorCallSign", XSCTag: "OpCall"},  // PDF Fields[62]
	{PDFName: "OperatorSent", XSCTag: "OpRelaySent"}, // PDF Fields[63]
	{PDFName: "OperatorTime", XSCTag: "OpTime"},      // PDF Fields[64]
	{PDFName: "OperatorDate", XSCTag: "OpDate"},      // PDF Fields[65]
}
