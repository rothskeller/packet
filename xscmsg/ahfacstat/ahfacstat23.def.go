// Code generated by extract-pifo-fields. DO NOT EDIT.

package ahfacstat

import "github.com/rothskeller/packet/xscmsg/internal/xscform"

var ahfacstat23 = &xscform.FormDefinition{
	HTML:                   "form-allied-health-facility-status.html",
	Tag:                    "AHFacStat",
	Name:                   "allied health status report",
	Article:                "an",
	Version:                "2.3",
	OriginNumberField:      "MsgNo",
	DestinationNumberField: "",
	HandlingOrderField:     "5.",
	SubjectField:           "20.",
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
			Tag: "20.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
		{
			Tag: "21.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "22d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateDate},
			Default: "«date»",
		},
		{
			Tag: "22t.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateTime},
		},
		{
			Tag: "23.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "23p.",
			Validations: []xscform.ValidateFunc{xscform.ValidatePhoneNumber, xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "23f.",
			Validations: []xscform.ValidateFunc{xscform.ValidatePhoneNumber},
		},
		{
			Tag: "24.",
		},
		{
			Tag: "25.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "25d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateDate, xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "35.",
			Values: []string{"Fully Functional", "Limited Services", "Impaired/Closed"},
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete, xscform.ValidateSelect},
		},
		{
			Tag: "26a.",
			Values: []string{"Yes", "No"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
		},
		{
			Tag: "26b.",
			Values: []string{"Yes", "No"},
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete, xscform.ValidateSelect},
		},
		{
			Tag: "26c.",
			Values: []string{"Yes", "No"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
		},
		{
			Tag: "26d.",
			Values: []string{"Yes", "No"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
		},
		{
			Tag: "27p.",
			Validations: []xscform.ValidateFunc{xscform.ValidatePhoneNumber, xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "26e.",
			Values: []string{"Yes", "No"},
			Validations: []xscform.ValidateFunc{xscform.ValidateSelect},
		},
		{
			Tag: "27f.",
			Validations: []xscform.ValidateFunc{xscform.ValidatePhoneNumber},
		},
		{
			Tag: "28.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "34.",
		},
		{
			Tag: "28p.",
			Validations: []xscform.ValidateFunc{xscform.ValidatePhoneNumber},
		},
		{
			Tag: "29.",
		},
		{
			Tag: "29p.",
			Validations: []xscform.ValidateFunc{xscform.ValidatePhoneNumber},
		},
		{
			Tag: "29e.",
		},
		{
			Tag: "30.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "30p.",
			Validations: []xscform.ValidateFunc{xscform.ValidatePhoneNumber, xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "40a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "40b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "40c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "40d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "40e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "30e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequiredForComplete},
		},
		{
			Tag: "41a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "41b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "41c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "41d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "41e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "42a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "42b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "42c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "42d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "42e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "31a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "43a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "43b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "43c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "43d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "43e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "31b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "44a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "44b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "44c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "44d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "44e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "31c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "45a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "45b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "45c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "45d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "45e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "33.",
		},
		{
			Tag: "46.",
		},
		{
			Tag: "46a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "46b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "46c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "46d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "46e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "50a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "50b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "50c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "50d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "50e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "51a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "51b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "51c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "51d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "51e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "52a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "52b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "52c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "52d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "52e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "53a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "53b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "53c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "53d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "53e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "54a.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "54b.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "54c.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "54d.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
		},
		{
			Tag: "54e.",
			Validations: []xscform.ValidateFunc{xscform.ValidateCardinalNumber},
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
		"20.": "facility",
		"21.": "facility-type",
		"22d.": "date",
		"22t.": "time",
		"23.": "contact",
		"23f.": "contact-fax",
		"23p.": "contact-phone",
		"24.": "other-contact",
		"25.": "incident",
		"25d.": "incident-date",
		"26a.": "attach-org-chart",
		"26b.": "attach-RR",
		"26c.": "attach-status",
		"26d.": "attach-action-plan",
		"26e.": "attach-directory",
		"27f.": "eoc-fax",
		"27p.": "eoc-phone",
		"28.": "liaison",
		"28p.": "liaison-phone",
		"29.": "info-officer",
		"29e.": "info-officer-email",
		"29p.": "info-officer-phone",
		"30.": "eoc-closed-contact",
		"30e.": "eoc-email",
		"30p.": "eoc-phone",
		"31a.": "to-evacuate",
		"31b.": "injured",
		"31c.": "transfered",
		"33.": "other-care",
		"34.": "summary",
		"35.": "status",
		"40a.": "skilled-nursing-beds-staffed-m",
		"40b.": "skilled-nursing-beds-staffed-f",
		"40c.": "skilled-nursing-beds-vacant-m",
		"40d.": "skilled-nursing-beds-vacant-f",
		"40e.": "skilled-nursing-beds-surge",
		"41a.": "assisted-living-beds-staffed-m",
		"41b.": "assisted-living-beds-staffed-f",
		"41c.": "assisted-living-beds-vacant-m",
		"41d.": "assisted-living-beds-vacant-f",
		"41e.": "assisted-living-beds-surge",
		"42a.": "sub-acute-beds-staffed-m",
		"42b.": "sub-acute-beds-staffed-f",
		"42c.": "sub-acute-beds-vacant-m",
		"42d.": "sub-acute-beds-vacant-f",
		"42e.": "sub-acute-beds-surge",
		"43a.": "alzheimers-beds-staffed-m",
		"43b.": "alzheimers-beds-staffed-f",
		"43c.": "alzheimers-beds-vacant-m",
		"43d.": "alzheimers-beds-vacant-f",
		"43e.": "alzheimers-beds-surge",
		"44a.": "ped-sub-acute-beds-staffed-m",
		"44b.": "ped-sub-acute-beds-staffed-f",
		"44c.": "ped-sub-acute-beds-vacant-m",
		"44d.": "ped-sub-acute-beds-vacant-f",
		"44e.": "ped-sub-acute-beds-surge",
		"45a.": "psychiatric-beds-staffed-m",
		"45b.": "psychiatric-beds-staffed-f",
		"45c.": "psychiatric-beds-vacant-m",
		"45d.": "psychiatric-beds-vacant-f",
		"45e.": "psychiatric-beds-surge",
		"46.": "bed-resource",
		"46a.": "other-care-beds-staffed-m",
		"46b.": "other-care-beds-staffed-f",
		"46c.": "other-care-beds-vacant-m",
		"46d.": "other-care-beds-vacant-f",
		"46e.": "other-care-beds-surge",
		"5.": "handling",
		"50a.": "dialysis-chairs",
		"50b.": "dialysis-vacant-chairs",
		"50c.": "dialysis-front-staff",
		"50d.": "dialysis-support-staff",
		"50e.": "dialysis-providers",
		"51a.": "surgical-chairs",
		"51b.": "surgical-vacant-chairs",
		"51c.": "surgical-front-staff",
		"51d.": "surgical-support-staff",
		"51e.": "surgical-providers",
		"52a.": "clinic-chairs",
		"52b.": "clinic-vacant-chairs",
		"52c.": "clinic-front-staff",
		"52d.": "clinic-support-staff",
		"52e.": "clinic-providers",
		"53a.": "home-health-chairs",
		"53b.": "home-health-vacant-chairs",
		"53c.": "home-health-front-staff",
		"53d.": "home-health-support-staff",
		"53e.": "home-health-providers",
		"54a.": "adulty-day-ctr-chairs",
		"54b.": "adulty-day-ctr-vacant-chairs",
		"54c.": "adulty-day-ctr-front-staff",
		"54d.": "adulty-day-ctr-support-staff",
		"54e.": "adulty-day-ctr-providers",
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
		"20.": "required",
		"21.": "required-for-complete",
		"22d.": "required date",
		"22t.": "required time",
		"23.": "required-for-complete",
		"23f.": "phone-number",
		"23p.": "phone-number required-for-complete",
		"25.": "required-for-complete",
		"25d.": "date required-for-complete",
		"26a.": "Yes, No",
		"26b.": "required-for-complete: Yes, No",
		"26c.": "Yes, No",
		"26d.": "Yes, No",
		"26e.": "Yes, No",
		"27f.": "phone-number",
		"27p.": "phone-number required-for-complete",
		"28.": "required-for-complete",
		"28p.": "phone-number",
		"29p.": "phone-number",
		"30.": "required-for-complete",
		"30e.": "required-for-complete",
		"30p.": "phone-number required-for-complete",
		"31a.": "cardinal-number",
		"31b.": "cardinal-number",
		"31c.": "cardinal-number",
		"35.": "required-for-complete: Fully Functional, Limited Services, Impaired/Closed",
		"40a.": "cardinal-number",
		"40b.": "cardinal-number",
		"40c.": "cardinal-number",
		"40d.": "cardinal-number",
		"40e.": "cardinal-number",
		"41a.": "cardinal-number",
		"41b.": "cardinal-number",
		"41c.": "cardinal-number",
		"41d.": "cardinal-number",
		"41e.": "cardinal-number",
		"42a.": "cardinal-number",
		"42b.": "cardinal-number",
		"42c.": "cardinal-number",
		"42d.": "cardinal-number",
		"42e.": "cardinal-number",
		"43a.": "cardinal-number",
		"43b.": "cardinal-number",
		"43c.": "cardinal-number",
		"43d.": "cardinal-number",
		"43e.": "cardinal-number",
		"44a.": "cardinal-number",
		"44b.": "cardinal-number",
		"44c.": "cardinal-number",
		"44d.": "cardinal-number",
		"44e.": "cardinal-number",
		"45a.": "cardinal-number",
		"45b.": "cardinal-number",
		"45c.": "cardinal-number",
		"45d.": "cardinal-number",
		"45e.": "cardinal-number",
		"46a.": "cardinal-number",
		"46b.": "cardinal-number",
		"46c.": "cardinal-number",
		"46d.": "cardinal-number",
		"46e.": "cardinal-number",
		"5.": "required: IMMEDIATE, PRIORITY, ROUTINE",
		"50a.": "cardinal-number",
		"50b.": "cardinal-number",
		"50c.": "cardinal-number",
		"50d.": "cardinal-number",
		"50e.": "cardinal-number",
		"51a.": "cardinal-number",
		"51b.": "cardinal-number",
		"51c.": "cardinal-number",
		"51d.": "cardinal-number",
		"51e.": "cardinal-number",
		"52a.": "cardinal-number",
		"52b.": "cardinal-number",
		"52c.": "cardinal-number",
		"52d.": "cardinal-number",
		"52e.": "cardinal-number",
		"53a.": "cardinal-number",
		"53b.": "cardinal-number",
		"53c.": "cardinal-number",
		"53d.": "cardinal-number",
		"53e.": "cardinal-number",
		"54a.": "cardinal-number",
		"54b.": "cardinal-number",
		"54c.": "cardinal-number",
		"54d.": "cardinal-number",
		"54e.": "cardinal-number",
		"7a.": "EMS Unit, Public Health Unit, Medical Health Branch, Operations Section, ...",
		"7b.": "MHJOC, County EOC, ...",
		"8a.": "required",
		"8b.": "required",
		"MsgNo": "required message-number",
		"OpCall": "required call-sign",
		"OpDate": "date",
		"OpTime": "time",
	},
}
