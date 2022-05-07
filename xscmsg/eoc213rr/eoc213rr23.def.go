// Code generated by extract-pifo-fields. DO NOT EDIT.

package eoc213rr

import "steve.rothskeller.net/packet/xscmsg/internal/xscform"

var eoc213rr23 = &xscform.FormDefinition{
	HTML:                   "form-scco-eoc-213rr.html",
	Tag:                    "EOC213RR",
	Name:                   "EOC-213RR resource request form",
	Article:                "an",
	Version:                "2.3",
	OriginNumberField:      "MsgNo",
	DestinationNumberField: "",
	HandlingOrderField:     "5.",
	SubjectField:           "21.",
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
			Tag: "7a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "8a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "7b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "8b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "7c.",
		},
		{
			Tag: "8c.",
		},
		{
			Tag: "7d.",
		},
		{
			Tag: "8d.",
		},
		{
			Tag: "21.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "22.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateDate},
			Default: "«date»",
		},
		{
			Tag: "23.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateTime},
		},
		{
			Tag: "24.",
		},
		{
			Tag: "25.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "26.",
		},
		{
			Tag: "27.",
		},
		{
			Tag: "28.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "29.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "30.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "31.",
			Values: []string{"Now", "High", "Medium", "Low"},
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateSelect},
		},
		{
			Tag: "32.",
		},
		{
			Tag: "33.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "34.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "35.",
		},
		{
			Tag: "36a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateBoolean},
		},
		{
			Tag: "36b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateBoolean},
		},
		{
			Tag: "36c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateBoolean},
		},
		{
			Tag: "36d.",
		},
		{
			Tag: "36e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateBoolean},
		},
		{
			Tag: "36f.",
			Validations: []xscform.ValidateFunc{xscform.ValidateBoolean},
		},
		{
			Tag: "36g.",
			Validations: []xscform.ValidateFunc{xscform.ValidateBoolean},
		},
		{
			Tag: "36h.",
			Validations: []xscform.ValidateFunc{xscform.ValidateBoolean},
		},
		{
			Tag: "36i.",
			Validations: []xscform.ValidateFunc{xscform.ValidateBoolean},
		},
		{
			Tag: "37.",
		},
		{
			Tag: "OpRelayRcvd",
		},
		{
			Tag: "OpRelaySent",
		},
		{
			Tag: "OpName",
		},
		{
			Tag: "OpCall",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateCallSign},
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
		"1a.": "date",
		"1b.": "time",
		"21.": "incident-name",
		"22.": "date",
		"23.": "time",
		"24.": "tracking-number",
		"25.": "requested-by",
		"26.": "prepared-by",
		"27.": "approved-by",
		"28.": "qty-unit",
		"29.": "resource-description",
		"30.": "resource-arrival",
		"31.": "priority",
		"32.": "resource-priority",
		"33.": "deliver-to",
		"34.": "deliver-to-location",
		"35.": "substitutes",
		"36a.": "equipment-operator",
		"36b.": "lodging",
		"36c.": "fuel",
		"36d.": "fuel-type",
		"36e.": "power",
		"36f.": "meals",
		"36g.": "maintenance",
		"36h.": "water",
		"36i.": "other",
		"37.": "instructions",
		"5.": "handling",
		"7a.": "to-ics-position",
		"7b.": "to-location",
		"7c.": "to-name",
		"7d.": "to-contact",
		"8a.": "from-ics-position",
		"8b.": "from-location",
		"8c.": "from-name",
		"8d.": "from-contact",
	},
	Comments: map[string]string{
		"1a.": "required date",
		"1b.": "required time",
		"21.": "required",
		"22.": "required date",
		"23.": "required time",
		"25.": "required",
		"28.": "required",
		"29.": "required",
		"30.": "required",
		"31.": "required: Now, High, Medium, Low",
		"33.": "required",
		"34.": "required",
		"36a.": "boolean",
		"36b.": "boolean",
		"36c.": "boolean",
		"36e.": "boolean",
		"36f.": "boolean",
		"36g.": "boolean",
		"36h.": "boolean",
		"36i.": "boolean",
		"5.": "required: IMMEDIATE, PRIORITY, ROUTINE",
		"7a.": "required",
		"7b.": "required",
		"8a.": "required",
		"8b.": "required",
		"MsgNo": "required message-number",
		"OpCall": "required call-sign",
		"OpDate": "date",
		"OpTime": "time",
	},
}
