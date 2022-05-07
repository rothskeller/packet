// Code generated by extract-pifo-fields. DO NOT EDIT.

package ics213

import "steve.rothskeller.net/packet/xscmsg/internal/xscform"

var ics213v22 = &xscform.FormDefinition{
	HTML:                   "form-ics213.html",
	Tag:                    "ICS213",
	Name:                   "ICS-213 general message form",
	Article:                "an",
	Version:                "2.2",
	OriginNumberField:      "MsgNo",
	DestinationNumberField: "3.",
	HandlingOrderField:     "5.",
	SubjectField:           "10.",
	OperatorNameField:      "OpName",
	OperatorCallField:      "OpCall",
	ActionDateField:        "OpDate",
	ActionTimeField:        "OpTime",
	Fields:                 []*xscform.FieldDefinition{
		{
			Tag: "MsgNo",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateMessageNumber},
		},
		{
			Tag: "3.",
			Validations: []xscform.ValidateFunc{xscform.ValidateMessageNumber},
		},
		{
			Tag: "1a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateDate},
			Default: "«date»",
		},
		{
			Tag: "1b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateTime},
		},
		{
			Tag: "5.",
			Values: []string{"IMMEDIATE", "PRIORITY", "ROUTINE"},
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateSelect},
		},
		{
			Tag: "6a.",
			Values: []string{"Yes", "No"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
		},
		{
			Tag: "6b.",
			Values: []string{"Yes", "No"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
		},
		{
			Tag: "6d.",
		},
		{
			Tag: "7.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "8.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "9a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "9b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "ToName",
		},
		{
			Tag: "FmName",
		},
		{
			Tag: "ToTel",
		},
		{
			Tag: "FmTel",
		},
		{
			Tag: "10.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "11.",
		},
		{
			Tag: "12.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "OpRelayRcvd",
		},
		{
			Tag: "OpRelaySent",
		},
		{
			Tag: "Rec-Sent",
			Values: []string{"receiver", "sender"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
		},
		{
			Tag: "OpCall",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateCallSign},
		},
		{
			Tag: "Method",
			Values: []string{"Telephone", "Dispatch Center", "EOC Radio", "FAX", "Courier", "Amateur Radio", "Other"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
			Default: "Other",
		},
		{
			Tag: "OpName",
		},
		{
			Tag: "Other",
			Default: "Packet",
		},
		{
			Tag: "OpDate",
			Validations: []xscform.ValidateFunc{xscform.ValidateDate},
		},
		{
			Tag: "OpTime",
			Validations: []xscform.ValidateFunc{xscform.ValidateTime},
		},
	},
	Annotations: map[string]string{
		"10.": "subject",
		"11.": "reference",
		"12.": "message",
		"1a.": "date",
		"1b.": "time",
		"3.": "rxmsgno",
		"5.": "handling",
		"6a.": "take-action",
		"6b.": "reply",
		"6d.": "reply-by",
		"7.": "to-ics-position",
		"8.": "from-ics-position",
		"9a.": "to-location",
		"9b.": "from-location",
	},
	Comments: map[string]string{
		"10.": "required",
		"12.": "required",
		"1a.": "required date",
		"1b.": "required time",
		"3.": "message-number",
		"5.": "required: IMMEDIATE, PRIORITY, ROUTINE",
		"6a.": "Yes, No",
		"6b.": "Yes, No",
		"7.": "required",
		"8.": "required",
		"9a.": "required",
		"9b.": "required",
		"Method": "Telephone, Dispatch Center, EOC Radio, FAX, Courier, Amateur Radio, Other",
		"MsgNo": "required message-number",
		"OpCall": "required call-sign",
		"OpDate": "date",
		"OpTime": "time",
		"Rec-Sent": "receiver, sender",
	},
}
