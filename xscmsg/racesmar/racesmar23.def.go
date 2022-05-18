// Code generated by extract-pifo-fields. DO NOT EDIT.

package racesmar

import "github.com/rothskeller/packet/xscmsg/internal/xscform"

var racesmar23 = &xscform.FormDefinition{
	HTML:                   "form-oa-mutual-aid-request-v2.html",
	Tag:                    "RACES-MAR",
	Name:                   "RACES mutual aid request form",
	Article:                "a",
	Version:                "2.3",
	OriginNumberField:      "MsgNo",
	DestinationNumberField: "",
	HandlingOrderField:     "5.",
	SubjectField:           "15.",
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
			Tag: "15.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "16a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "16b.",
		},
		{
			Tag: "17.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "18.1a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateCardinalNumber},
		},
		{
			Tag: "18.1e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "18.1f.",
		},
		{
			Tag: "18.1b.",
		},
		{
			Tag: "18.1c.",
		},
		{
			Tag: "18.1d.",
		},
		{
			Tag: "18.2a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "18.2e.",
		},
		{
			Tag: "18.2f.",
		},
		{
			Tag: "18.2b.",
		},
		{
			Tag: "18.2c.",
		},
		{
			Tag: "18.2d.",
		},
		{
			Tag: "18.3a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "18.3e.",
		},
		{
			Tag: "18.3f.",
		},
		{
			Tag: "18.3b.",
		},
		{
			Tag: "18.3c.",
		},
		{
			Tag: "18.3d.",
		},
		{
			Tag: "18.4a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "18.4e.",
		},
		{
			Tag: "18.4f.",
		},
		{
			Tag: "18.4b.",
		},
		{
			Tag: "18.4c.",
		},
		{
			Tag: "18.4d.",
		},
		{
			Tag: "18.5a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "18.5e.",
		},
		{
			Tag: "18.5f.",
		},
		{
			Tag: "18.5b.",
		},
		{
			Tag: "18.5c.",
		},
		{
			Tag: "18.5d.",
		},
		{
			Tag: "19a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "19b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "20a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "20b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "21.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "22.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "23.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "24a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "24b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "24c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "25a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "25b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "25c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "26a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateDate},
			Default: "«date»",
		},
		{
			Tag: "26b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateTime},
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
		"15.": "agency",
		"16a.": "event-name",
		"16b.": "event-number",
		"17.": "assignment",
		"18.1a.": "resources-qty",
		"18.1b.": "resources-role",
		"18.1c.": "preferred-type",
		"18.1d.": "minimum-type",
		"18.1e.": "role",
		"18.1f.": "position",
		"18.2a.": "resources-qty",
		"18.2b.": "resources-role",
		"18.2c.": "preferred-type",
		"18.2d.": "minimum-type",
		"18.2e.": "role",
		"18.2f.": "position",
		"18.3a.": "resources-qty",
		"18.3b.": "resources-role",
		"18.3c.": "preferred-type",
		"18.3d.": "minimum-type",
		"18.3e.": "role",
		"18.3f.": "position",
		"18.4a.": "resources-qty",
		"18.4b.": "resources-role",
		"18.4c.": "preferred-type",
		"18.4d.": "minimum-type",
		"18.4e.": "role",
		"18.4f.": "position",
		"18.5a.": "resources-qty",
		"18.5b.": "resources-role",
		"18.5c.": "preferred-type",
		"18.5d.": "minimum-type",
		"18.5e.": "role",
		"18.5f.": "position",
		"19a.": "arrival-dates",
		"19b.": "arrival-times",
		"1a.": "date",
		"1b.": "time",
		"20a.": "needed-dates",
		"20b.": "needed-times",
		"21.": "reporting-location",
		"22.": "contact-on-arrival",
		"23.": "travel-info",
		"24a.": "requester-name",
		"24b.": "requester-title",
		"24c.": "requester-contact",
		"25a.": "agency-approver-name",
		"25b.": "agency-approver-title",
		"25c.": "agency-approver-contact",
		"26a.": "agency-approved-date",
		"26b.": "agency-approved-time",
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
		"15.": "required",
		"16a.": "required",
		"17.": "required",
		"18.1a.": "required cardinal-number",
		"18.1e.": "required",
		"18.2a.": "cardinal-number",
		"18.3a.": "cardinal-number",
		"18.4a.": "cardinal-number",
		"18.5a.": "cardinal-number",
		"19a.": "required",
		"19b.": "required",
		"1a.": "required date",
		"1b.": "required time",
		"20a.": "required",
		"20b.": "required",
		"21.": "required",
		"22.": "required",
		"23.": "required",
		"24a.": "required",
		"24b.": "required",
		"24c.": "required",
		"25a.": "required",
		"25b.": "required",
		"25c.": "required",
		"26a.": "required date",
		"26b.": "required time",
		"5.": "required: IMMEDIATE, PRIORITY, ROUTINE",
		"7a.": "RACES Chief Radio Officer, RACES Unit, Operations Section, ...",
		"7b.": "required",
		"8a.": "required",
		"8b.": "required",
		"MsgNo": "required message-number",
		"OpCall": "required call-sign",
		"OpDate": "date",
		"OpTime": "time",
	},
}
