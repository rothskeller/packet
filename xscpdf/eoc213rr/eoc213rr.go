package eoc213rr

import (
	_ "embed" // oh, please
	"encoding/hex"

	"github.com/rothskeller/packet/xscmsg/eoc213rr"
	"github.com/rothskeller/packet/xscpdf"
)

//go:embed XSC_EOC-213RR_Fillable_v20170803_with_XSC_RACES_Routing_Slip_Fillable_v20190527.pdf
var basePDF []byte

var id, _ = hex.DecodeString("d807e45d22add1680c19a7beba47d7a7")

func init() {
	xscpdf.RegisterReader(xscpdf.ReaderMap{XSCTag: eoc213rr.Tag, PDFID: id, Fields: fieldMap})
	xscpdf.RegisterWriter(xscpdf.WriterMap{XSCTag: eoc213rr.Tag, BasePDF: basePDF, Fields: fieldMap})
}

var fieldMap = []xscpdf.FieldMap{
	{PDFName: "Origin Msg Nbr", XSCTag: "MsgNo", FontSize: 12},          // 1
	{PDFName: "Destination Msg Nbr", XSCTag: "DestMsgNo", FontSize: 12}, // 2
	{PDFName: "Date Created", XSCTag: "1a.", FontSize: 12},              // 3
	{PDFName: "Time Created", XSCTag: "1b.", FontSize: 12},              // 4
	{PDFName: "Handling", XSCTag: "5.", Values: []xscpdf.ValueMap{ // 5
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Immediate", XSCValue: "IMMEDIATE"},
		{PDFValue: "Priority", XSCValue: "PRIORITY"},
		{PDFValue: "Routine", XSCValue: "ROUTINE"},
	}},
	{PDFName: "To ICS Position", XSCTag: "7a.", FontSize: 12},        // 6
	{PDFName: "To Location", XSCTag: "7b.", FontSize: 12},            // 7
	{PDFName: "To Name", XSCTag: "7c.", FontSize: 12},                // 8
	{PDFName: "To Contact Info", XSCTag: "7d.", FontSize: 12},        // 9
	{PDFName: "From ICS Position", XSCTag: "8a.", FontSize: 12},      // 10
	{PDFName: "From Location", XSCTag: "8b.", FontSize: 12},          // 11
	{PDFName: "From Name", XSCTag: "8c.", FontSize: 12},              // 12
	{PDFName: "From Contact Info", XSCTag: "8d.", FontSize: 12},      // 13
	{PDFName: "Form Type", PDFFixed: "XSC EOC-213RR", FontSize: 12},  // 14
	{PDFName: "Form Topic", XSCTag: "21.", FontSize: 12},             // 15
	{PDFName: "Relay Rcvd", XSCTag: "OpRelayRcvd", FontSize: 12},     // 16
	{PDFName: "Relay Sent", XSCTag: "OpRelaySent", FontSize: 12},     // 17
	{PDFName: "Op Name", XSCTag: "OpName", FontSize: 12},             // 18
	{PDFName: "Op Call Sign", XSCTag: "OpCall", FontSize: 12},        // 19
	{PDFName: "Op Date", XSCTag: "OpDate", FontSize: 12},             // 20
	{PDFName: "Op Time", XSCTag: "OpTime", FontSize: 12},             // 21
	{PDFName: "1 Incident Name", XSCTag: "21.", FontSize: 12},        // 22
	{PDFName: "2 Date Initiated", XSCTag: "22.", FontSize: 12},       // 23
	{PDFName: "3 Time Initiated", XSCTag: "23.", FontSize: 12},       // 24
	{PDFName: "5 Requested By", XSCTag: "25.", FontSize: 12},         // 25
	{PDFName: "6 Prepared by", XSCTag: "26.", FontSize: 12},          // 26
	{PDFName: "7 Approved By", XSCTag: "27.", FontSize: 12},          // 27
	{PDFName: "8 QtyUnit", XSCTag: "28.", FontSize: 12},              // 28
	{PDFName: "9 Resource Description", XSCTag: "29.", FontSize: 12}, // 29
	{PDFName: "10 Arrival", XSCTag: "30.", FontSize: 12},             // 30
	{PDFName: "11 Priority", XSCTag: "31.", Values: []xscpdf.ValueMap{ // 31
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Now", XSCValue: "Now"},
		{PDFValue: "High", XSCValue: "High"},
		{PDFValue: "Medium", XSCValue: "Medium"},
		{PDFValue: "Low", XSCValue: "Low"},
	}},
	{PDFName: "12 Estd Cost", XSCTag: "32.", FontSize: 12},               // 32
	{PDFName: "13 Deliver To", XSCTag: "33.", FontSize: 12},              // 33
	{PDFName: "14 Location", XSCTag: "34.", FontSize: 12},                // 34
	{PDFName: "15 Sub Sugg Sources", XSCTag: "35.", FontSize: 12},        // 35
	{PDFName: "Equip Oper", XSCTag: "36a.", Values: xscpdf.CheckboxMap},  // 36
	{PDFName: "Fuel", XSCTag: "36c.", Values: xscpdf.CheckboxMap},        // 37
	{PDFName: "Fuel Type", XSCTag: "36d.", FontSize: 12},                 // 38
	{PDFName: "Meals", XSCTag: "36f.", Values: xscpdf.CheckboxMap},       // 39
	{PDFName: "Water", XSCTag: "36h.", Values: xscpdf.CheckboxMap},       // 40
	{PDFName: "Lodging", XSCTag: "36b.", Values: xscpdf.CheckboxMap},     // 41
	{PDFName: "Power", XSCTag: "36e.", Values: xscpdf.CheckboxMap},       // 42
	{PDFName: "Maintenance", XSCTag: "36g.", Values: xscpdf.CheckboxMap}, // 43
	{PDFName: "Other", XSCTag: "36i.", Values: xscpdf.CheckboxMap},       // 44
	{PDFName: "17 Special Instructions", XSCTag: "37.", FontSize: 12},    // 46
}
