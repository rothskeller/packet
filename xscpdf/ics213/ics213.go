package ics213

import (
	_ "embed" // oh, please
	"encoding/hex"

	"github.com/rothskeller/packet/xscmsg/ics213"
	"github.com/rothskeller/packet/xscpdf"
)

//go:embed ICS-213_SCCo_Message_Form_Fillable_v20220119.pdf
var basePDF []byte

var id, _ = hex.DecodeString("d77f0ce1e7f53f4bc76cf657612be44f")

func init() {
	xscpdf.RegisterReader(xscpdf.ReaderMap{XSCTag: ics213.Tag, PDFID: id, Fields: fieldMap})
	xscpdf.RegisterWriter(xscpdf.WriterMap{XSCTag: ics213.Tag, BasePDF: basePDF, Fields: fieldMap})
}

var fieldMap = []xscpdf.FieldMap{
	{PDFName: "Immediate", XSCTag: "5.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "1", XSCValue: "PRIORITY"},
		{PDFValue: "2", XSCValue: "ROUTINE"},
		{PDFValue: "3", XSCValue: "IMMEDIATE"},
	}},
	{PDFName: "TakeAction", XSCTag: "6a.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "1", XSCValue: "Yes"},
		{PDFValue: "2", XSCValue: "No"},
	}},
	{PDFName: "Reply", XSCTag: "6b.", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Reply-Yes", XSCValue: "Yes"},
		{PDFValue: "Reply-No", XSCValue: "No"},
	}},
	{PDFName: "How: Received", XSCTag: "Rec-Sent", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "0", XSCValue: "sender"},
		{PDFValue: "1", XSCValue: "receiver"},
	}},
	{PDFName: "Telephone", XSCTag: "Method", Values: []xscpdf.ValueMap{
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "1", XSCValue: "Telephone"},
		{PDFValue: "2", XSCValue: "Dispatch Center"},
		{PDFValue: "3", XSCValue: "EOC Radio"},
		{PDFValue: "4", XSCValue: "FAX"},
		{PDFValue: "5", XSCValue: "Courier"},
		{PDFValue: "6", XSCValue: "Amateur Radio"},
		{PDFValue: "7", XSCValue: "Other"},
	}},
	{PDFName: "Origin Msg #", XSCTag: "MsgNo"},
	{PDFName: "Destination Msg#", XSCTag: "3."},
	{PDFName: "FormDate", XSCTag: "1a."},
	{PDFName: "FormTime", XSCTag: "1b."},
	{PDFName: "Reply_2", XSCTag: "6d."},
	{PDFName: "TO ICS Position", XSCTag: "7."},
	{PDFName: "TO ICS Locatoin", XSCTag: "9a."},
	{PDFName: "TO ICS Name", XSCTag: "ToName"},
	{PDFName: "TO ICS Telephone", XSCTag: "ToTel"},
	{PDFName: "From ICS Position", XSCTag: "8."},
	{PDFName: "From ICS Location", XSCTag: "9b."},
	{PDFName: "From ICS Name", XSCTag: "FmName"},
	{PDFName: "From ICS Telephone", XSCTag: "FmTel"},
	{PDFName: "Subject", XSCTag: "10."},
	{PDFName: "Reference", XSCTag: "11."},
	{PDFName: "Message", XSCTag: "12."},
	{PDFName: "Relay Received", XSCTag: "OpRelayRcvd"},
	{PDFName: "Relay Sent", XSCTag: "OpRelaySent"},
	{PDFName: "Operation Call Sign", XSCTag: "OpCall"},
	{PDFName: "Relay Received_2", XSCTag: "OpName"},
	{PDFName: "OperatorDate", XSCTag: "OpDate"},
	{PDFName: "OperatorTime", XSCTag: "OpTime"},
	{PDFName: "OtherText", XSCTag: "Other"},
}
