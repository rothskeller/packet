// Code generated by extract-pifo-fields. DO NOT EDIT.

package sheltstat

import "steve.rothskeller.net/packet/xscmsg/internal/xscform"

var sheltstat22 = &xscform.FormDefinition{
	HTML:                   "form-oa-shelter-status.html",
	Tag:                    "SheltStat",
	Name:                   "shelter status",
	Article:                "a",
	Version:                "2.2",
	OriginNumberField:      "MsgNo",
	DestinationNumberField: "",
	HandlingOrderField:     "5.",
	SubjectField:           "32.",
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
			Tag: "19.",
			Values: []string{"Update", "Complete"},
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateSelect},
		},
		{
			Tag: "32.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "30.",
			Values: []string{"Type 1", "Type 2", "Type 3", "Type 4"},
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete, xscform.ValidateSelect},
		},
		{
			Tag: "31.",
			Values: []string{"Open", "Closed", "Full"},
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete, xscform.ValidateSelect},
		},
		{
			Tag: "33a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "33b.",
		},
		{
			Tag: "34b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "33c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "33d.",
		},
		{
			Tag: "37a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRealNumber},
		},
		{
			Tag: "37b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRealNumber},
		},
		{
			Tag: "40a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber, xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "40b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber, xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "41.",
		},
		{
			Tag: "42.",
		},
		{
			Tag: "43a.",
			Values: []string{"checked", "false"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
		},
		{
			Tag: "43b.",
			Values: []string{"checked", "false"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
		},
		{
			Tag: "43c.",
			Values: []string{"checked", "false"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
		},
		{
			Tag: "44.",
		},
		{
			Tag: "45.",
		},
		{
			Tag: "46.",
		},
		{
			Tag: "50a.",
		},
		{
			Tag: "49a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "50b.",
		},
		{
			Tag: "51a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "51b.",
			Validations: []xscform.ValidateFunc{xscform.ValidatePhoneNumber, xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "52a.",
		},
		{
			Tag: "52b.",
			Validations: []xscform.ValidateFunc{xscform.ValidatePhoneNumber},
		},
		{
			Tag: "60.",
		},
		{
			Tag: "61.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCallSign},
		},
		{
			Tag: "62a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateFrequency},
		},
		{
			Tag: "62b.",
		},
		{
			Tag: "63a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateFrequency},
		},
		{
			Tag: "63b.",
		},
		{
			Tag: "62c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateFrequencyOffset},
		},
		{
			Tag: "70.",
		},
		{
			Tag: "71.",
			Validations: []xscform.ValidateFunc{xscform.ValidateBoolean},
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
		"19.": "report-type",
		"1a.": "date",
		"1b.": "time",
		"30.": "shelter-type",
		"31.": "shelter-status",
		"32.": "shelter-name",
		"33a.": "shelter-address",
		"33b.": "shelter-city",
		"33c.": "shelter-state",
		"33d.": "shelter-zip",
		"34b.": "shelter-city",
		"37a.": "latitude",
		"37b.": "longitude",
		"40a.": "capacity",
		"40b.": "occupancy",
		"41.": "meals",
		"42.": "NSS",
		"43a.": "pet-friendly",
		"43b.": "basic-safety",
		"43c.": "ATC-20",
		"44.": "available-services",
		"45.": "MOU",
		"46.": "floor-plan",
		"49a.": "managed-by",
		"5.": "handling",
		"50a.": "managed-by",
		"50b.": "managed-by-detail",
		"51a.": "primary-contact",
		"51b.": "primary-phone",
		"52a.": "secondary-contact",
		"52b.": "secondary-phone",
		"60.": "tactical-call",
		"61.": "repeater-call",
		"62a.": "repeater-input",
		"62b.": "repeater-input-tone",
		"62c.": "repeater-offset",
		"63a.": "repeater-output",
		"63b.": "repeater-output-tone",
		"70.": "comments",
		"71.": "remove-from-active-list",
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
		"19.": "required: Update, Complete",
		"1a.": "required date",
		"1b.": "required time",
		"30.": "required-for-complete: Type 1, Type 2, Type 3, Type 4",
		"31.": "required-for-complete: Open, Closed, Full",
		"32.": "required",
		"33a.": "required-for-complete",
		"33c.": "required-for-complete",
		"34b.": "required-for-complete",
		"37a.": "real-number",
		"37b.": "real-number",
		"40a.": "cardinal-number required-for-complete",
		"40b.": "cardinal-number required-for-complete",
		"43a.": "checked, false",
		"43b.": "checked, false",
		"43c.": "checked, false",
		"49a.": "required-for-complete",
		"5.": "required: IMMEDIATE, PRIORITY, ROUTINE",
		"51a.": "required-for-complete",
		"51b.": "phone-number required-for-complete",
		"52b.": "phone-number",
		"61.": "call-sign",
		"62a.": "frequency",
		"62c.": "frequency-offset",
		"63a.": "frequency",
		"71.": "boolean",
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