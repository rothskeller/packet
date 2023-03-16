package sheltstat

import (
	_ "embed" // oh, please
	"encoding/hex"

	"github.com/rothskeller/packet/xscmsg/sheltstat"
	"github.com/rothskeller/packet/xscpdf"
)

//go:embed XSC_SheltStat_Fillable_v20190619_p12.pdf
var basePDF []byte

var idFull, _ = hex.DecodeString("52aa01f834820e47a3a070d81cead3b2")
var idP12, _ = hex.DecodeString("a3c1958ebba5a2fc16b557b7a7a84ae1")

func init() {
	xscpdf.RegisterReader(xscpdf.ReaderMap{XSCTag: sheltstat.Tag, PDFID: idFull, Fields: fieldMap})
	xscpdf.RegisterReader(xscpdf.ReaderMap{XSCTag: sheltstat.Tag, PDFID: idP12, Fields: fieldMap})
	xscpdf.RegisterWriter(xscpdf.WriterMap{XSCTag: sheltstat.Tag, BasePDF: basePDF, Fields: fieldMap})
}

var fieldMap = []xscpdf.FieldMap{
	{PDFName: "Address", XSCTag: "33a.", FontSize: 10}, // 3
	{PDFName: "ATC20 Insp", XSCTag: "43c.", Values: []xscpdf.ValueMap{
		{PDFValue: "No", XSCValue: "false"},
		{PDFValue: "No", XSCValue: ""},
		{PDFValue: "Off", XSCValue: "false"},
		{PDFValue: "Yes", XSCValue: "checked"},
	}}, // 51
	{PDFName: "Available Services", XSCTag: "44.", FontSize: 10}, // 10
	{PDFName: "Basic Safety Insp", XSCTag: "43b.", Values: []xscpdf.ValueMap{
		{PDFValue: "No", XSCValue: "false"},
		{PDFValue: "No", XSCValue: ""},
		{PDFValue: "Off", XSCValue: "false"},
		{PDFValue: "Yes", XSCValue: "checked"},
	}}, // 50
	{PDFName: "Capacity", XSCTag: "40a.", FontSize: 10},                 // 7
	{PDFName: "City", XSCTag: "34b.", FontSize: 10},                     // 4
	{PDFName: "Comments", XSCTag: "70.", FontSize: 10},                  // 15
	{PDFName: "Date Created", XSCTag: "1a.", FontSize: 10},              // 17
	{PDFName: "Destination Msg Nbr", XSCTag: "DestMsgNo", FontSize: 10}, // 16
	{PDFName: "Floorplan", XSCTag: "46.", FontSize: 10},                 // 33
	{PDFName: "From Contact Info", XSCTag: "8d.", FontSize: 10},         // 26
	{PDFName: "From ICS Position", XSCTag: "8a.", FontSize: 10},         // 23
	{PDFName: "From Location", XSCTag: "8b.", FontSize: 10},             // 24
	{PDFName: "From Name", XSCTag: "8c.", FontSize: 10},                 // 25
	{PDFName: "Handling", XSCTag: "5.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Immediate", XSCValue: "IMMEDIATE"},
		{PDFValue: "Priority", XSCValue: "PRIORITY"},
		{PDFValue: "Routine", XSCValue: "ROUTINE"},
	}}, // 0
	{PDFName: "Input Freq", XSCTag: "62a.", FontSize: 10},        // 40
	{PDFName: "Input Tone", XSCTag: "62b.", FontSize: 10},        // 41
	{PDFName: "Latitude", XSCTag: "37a.", FontSize: 10},          // 29
	{PDFName: "Longitude", XSCTag: "37b.", FontSize: 10},         // 30
	{PDFName: "Managed By Detail", XSCTag: "50b.", FontSize: 10}, // 11
	{PDFName: "Managed By", XSCTag: "50a.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "American Red Cross", XSCValue: "American Red Cross"},
		{PDFValue: "Private", XSCValue: "Private"},
		{PDFValue: "Community", XSCValue: "Community"},
		{PDFValue: "Government", XSCValue: "Government"},
		{PDFValue: "Other", XSCValue: "Other"},
	}}, // 35
	{PDFName: "Meals", XSCTag: "41.", FontSize: 10},                 // 31
	{PDFName: "MOU", XSCTag: "45.", FontSize: 10},                   // 32
	{PDFName: "NSS Number", XSCTag: "42.", FontSize: 10},            // 9
	{PDFName: "Occupancy", XSCTag: "40b.", FontSize: 10},            // 8
	{PDFName: "Offset", XSCTag: "62c.", FontSize: 10},               // 14
	{PDFName: "Op Call Sign", XSCTag: "OpCall", FontSize: 10},       // 47
	{PDFName: "Op Date", XSCTag: "OpDate", FontSize: 10},            // 48
	{PDFName: "Op Name", XSCTag: "OpName", FontSize: 10},            // 46
	{PDFName: "Op Time", XSCTag: "OpTime", FontSize: 10},            // 49
	{PDFName: "Origin Msg Nbr Copy", XSCTag: "MsgNo", FontSize: 10}, // 55
	{PDFName: "Origin Msg Nbr", XSCTag: "MsgNo", FontSize: 10},      // 34
	{PDFName: "Output Freq", XSCTag: "63a.", FontSize: 10},          // 42
	{PDFName: "Output Tone", XSCTag: "63b.", FontSize: 10},          // 43
	{PDFName: "Pet Friendly", XSCTag: "43a.", Values: []xscpdf.ValueMap{
		{PDFValue: "No", XSCValue: "false"},
		{PDFValue: "No", XSCValue: ""},
		{PDFValue: "Off", XSCValue: "false"},
		{PDFValue: "Yes", XSCValue: "checked"},
	}}, // 52
	{PDFName: "Pri Contact Phone", XSCTag: "51b.", FontSize: 10}, // 37
	{PDFName: "Pri Contact", XSCTag: "51a.", FontSize: 10},       // 36
	{PDFName: "Relay Rcvd", XSCTag: "OpRelayRcvd", FontSize: 10}, // 44
	{PDFName: "Relay Sent", XSCTag: "OpRelaySent", FontSize: 10}, // 45
	{PDFName: "Remove", XSCTag: "71.", Values: []xscpdf.ValueMap{
		{PDFValue: "No", XSCValue: ""},
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Yes", XSCValue: "checked"},
	}}, // 53
	{PDFName: "Repeater Call Sign", XSCTag: "61.", FontSize: 10}, // 13
	{PDFName: "Report Type", XSCTag: "19.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Update", XSCValue: "Update"},
		{PDFValue: "Complete", XSCValue: "Complete"},
	}}, // 1
	{PDFName: "Sec Contact Phone", XSCTag: "52b.", FontSize: 10}, // 39
	{PDFName: "Sec Contact", XSCTag: "52a.", FontSize: 10},       // 38
	{PDFName: "Shelter Name", XSCTag: "32.", FontSize: 10},       // 2
	{PDFName: "Shelter Status", XSCTag: "31.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Open", XSCValue: "Open"},
		{PDFValue: "Closed", XSCValue: "Closed"},
		{PDFValue: "Full", XSCValue: "Full"},
	}}, // 28
	{PDFName: "Shelter Type", XSCTag: "30.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Type 1", XSCValue: "Type 1"},
		{PDFValue: "Type 2", XSCValue: "Type 2"},
		{PDFValue: "Type 3", XSCValue: "Type 3"},
		{PDFValue: "Type 4", XSCValue: "Type 4"},
	}}, // 27
	{PDFName: "State", XSCTag: "33c.", FontSize: 10},             // 5
	{PDFName: "Tactical Call Sign", XSCTag: "60.", FontSize: 10}, // 12
	{PDFName: "Time Created", XSCTag: "1b.", FontSize: 10},       // 18
	{PDFName: "To Contact Info", XSCTag: "7d.", FontSize: 10},    // 22
	{PDFName: "To ICS Position", XSCTag: "7a.", FontSize: 10},    // 19
	{PDFName: "To Location", XSCTag: "7b.", FontSize: 10},        // 20
	{PDFName: "To Name", XSCTag: "7c.", FontSize: 10},            // 21
	{PDFName: "Zip", XSCTag: "33d.", FontSize: 10},               // 6
}
