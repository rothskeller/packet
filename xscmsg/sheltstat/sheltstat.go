package sheltstat

import (
	"fmt"
	"regexp"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscform"
)

// Tag identifies shelter status forms.
const Tag = "SheltStat"

var (
	frequencyRE       = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	frequencyOffsetRE = regexp.MustCompile(`^[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+|[-+]$`)
)

func init() {
	xscmsg.RegisterCreate(Tag, create)
	xscmsg.RegisterType(recognize)

	// Our handling, toICSPosition, and toLocation fields are variants of
	// the standard ones, adding default values to them.
	handlingDef.DefaultValue = "PRIORITY"
	toICSPositionDef.DefaultValue = "Mass Care and Shelter Unit"
	toICSPositionDef.Comment = "required: Mass Care and Shelter Unit, Care and Shelter Branch, Operations Section, ..."
	toLocationDef.DefaultValue = "County EOC"
	toLocationDef.Comment = "required: County EOC, ..."
}

func create() *xscmsg.Message {
	return xscform.CreateForm(formtype, fieldDefsV22)
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	if form == nil || form.FormType != formtype.HTML {
		return nil
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.2") {
		return xscform.AdoptForm(formtype, fieldDefsV22, msg, form)
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.0") {
		return xscform.AdoptForm(formtype, fieldDefsV21, msg, form)
	}
	return nil
}

var formtype = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "OA shelter status form",
	Article:     "an",
	HTML:        "form-oa-shelter-status.html",
	Version:     "2.2",
	SubjectFunc: xscform.EncodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

// fieldDefsV22 adds shelterCityCode and managedByCode, and changes from V21 to
// V22 for shelterCity and managedBy.
var fieldDefsV22 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// Shelter Status fields
	reportTypeDef, shelterNameDef, shelterTypeDef, shelterStatusDef, shelterAddressDef, shelterCityCodeDef, shelterCityDefV22,
	shelterStateDef, shelterZipDef, latitudeDef, longitudeDef, capacityDef, occupancyDef, mealsDef, nssDef, petFriendlyDef,
	basicSafetyDef, atc20Def, availableServicesDef, mouDef, floorPlanDef, managedByCodeDef, managedByDefV22, managedByDetailDef,
	primaryContactDef, primaryPhoneDef, secondaryContactDef, secondaryPhoneDef, tacticalCallDef, repeaterCallDef,
	repeaterInputDef, repeaterInputToneDef, repeaterOutputDef, repeaterOutputToneDef, repeaterOffsetDef, commentsDef,
	removeFromActiveListDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var fieldDefsV21 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// Shelter Status fields
	reportTypeDef, shelterNameDef, shelterTypeDef, shelterStatusDef, shelterAddressDef, shelterCityDefV21,
	shelterStateDef, shelterZipDef, latitudeDef, longitudeDef, capacityDef, occupancyDef, mealsDef, nssDef, petFriendlyDef,
	basicSafetyDef, atc20Def, availableServicesDef, mouDef, floorPlanDef, managedByDefV21, managedByDetailDef,
	primaryContactDef, primaryPhoneDef, secondaryContactDef, secondaryPhoneDef, tacticalCallDef, repeaterCallDef,
	repeaterInputDef, repeaterInputToneDef, repeaterOutputDef, repeaterOutputToneDef, repeaterOffsetDef, commentsDef,
	removeFromActiveListDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var (
	handlingDef      = *xscform.HandlingDef      // modified in func init
	toICSPositionDef = *xscform.ToICSPositionDef // modified in func init
	toLocationDef    = *xscform.ToLocationDef    // modified in func init
	reportTypeDef    = &xscmsg.FieldDef{
		Tag:        "19.",
		Annotation: "report-type",
		Label:      "Report Type",
		Comment:    "required: Update, Complete",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateChoices},
		Choices:    []string{"Update", "Complete"},
	}
	shelterNameDef = &xscmsg.FieldDef{
		Tag:        "32.",
		Annotation: "shelter-name",
		Label:      "Shelter Name",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
		Key:        xscmsg.FSubject,
	}
	shelterTypeDef = &xscmsg.FieldDef{
		Tag:        "30.",
		Annotation: "shelter-type",
		Label:      "Shelter Type",
		Comment:    "required-for-complete: Type 1, Type 2, Type 3, Type 4",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Type 1", "Type 2", "Type 3", "Type 4"},
	}
	shelterStatusDef = &xscmsg.FieldDef{
		Tag:        "31.",
		Annotation: "shelter-status",
		Label:      "Status",
		Comment:    "required-for-complete: Open, Closed, Full",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Open", "Closed", "Full"},
	}
	shelterAddressDef = &xscmsg.FieldDef{
		Tag:        "33a.",
		Annotation: "shelter-address",
		Label:      "Address",
		Comment:    "required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete},
	}
	shelterCityCodeDef = &xscmsg.FieldDef{
		Tag:        "33b.",
		Annotation: "shelter-city-code",
		Comment:    "required-for-complete",
		ReadOnly:   true,
		Validators: []xscmsg.Validator{setShelterCityCode, xscform.ValidateChoices},
		Choices:    []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
	}
	shelterCityDefV21 = &xscmsg.FieldDef{
		Tag:        "33b.",
		Annotation: "shelter-city",
		Label:      "City",
		Comment:    "required-for-complete: Campbell, Cupertino, Gilroy, Los Altos, Los Altos Hills, Los Gatos, Milpitas, Monte Sereno, Morgan Hill, Mountain View, Palo Alto, San Jose, Santa Clara, Saratoga, Sunnyvale, Unincorporated",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
	}
	shelterCityDefV22 = &xscmsg.FieldDef{
		Tag:        "34b.",
		Annotation: "shelter-city",
		Label:      "City",
		Comment:    "required-for-complete: Campbell, Cupertino, Gilroy, Los Altos, Los Altos Hills, Los Gatos, Milpitas, Monte Sereno, Morgan Hill, Mountain View, Palo Alto, San Jose, Santa Clara, Saratoga, Sunnyvale, ...",
		Validators: []xscmsg.Validator{requiredForComplete},
	}
	shelterStateDef = &xscmsg.FieldDef{
		Tag:        "33c.",
		Annotation: "shelter-state",
		Label:      "State",
	}
	shelterZipDef = &xscmsg.FieldDef{
		Tag:        "33d.",
		Annotation: "shelter-zip",
		Label:      "Zip",
	}
	latitudeDef = &xscmsg.FieldDef{
		Tag:        "37a.",
		Annotation: "latitude",
		Label:      "Latitude",
		Comment:    "real-number",
		Validators: []xscmsg.Validator{xscform.ValidateRealNumber},
	}
	longitudeDef = &xscmsg.FieldDef{
		Tag:        "37b.",
		Annotation: "longitude",
		Label:      "Longitude",
		Comment:    "real-number",
		Validators: []xscmsg.Validator{xscform.ValidateRealNumber},
	}
	capacityDef = &xscmsg.FieldDef{
		Tag:        "40a.",
		Annotation: "capacity",
		Label:      "Capacity",
		Comment:    "cardinal-number required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidateCardinalNumber},
	}
	occupancyDef = &xscmsg.FieldDef{
		Tag:        "40b.",
		Annotation: "occupancy",
		Label:      "Occupancy",
		Comment:    "cardinal-number required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidateCardinalNumber},
	}
	mealsDef = &xscmsg.FieldDef{
		Tag:        "41.",
		Annotation: "meals",
		Label:      "Meals Served (last 24 hours)",
	}
	nssDef = &xscmsg.FieldDef{
		Tag:        "42.",
		Annotation: "NSS",
		Label:      "NSS Number",
	}
	petFriendlyDef = &xscmsg.FieldDef{
		Tag:        "43a.",
		Annotation: "pet-friendly",
		Label:      "Pet Friendly",
		Comment:    "checked, false",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"checked", "false"},
	}
	basicSafetyDef = &xscmsg.FieldDef{
		Tag:        "43b.",
		Annotation: "basic-safety",
		Label:      "Basic Safety Inspection",
		Comment:    "checked, false",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"checked", "false"},
	}
	atc20Def = &xscmsg.FieldDef{
		Tag:        "43c.",
		Annotation: "ATC-20",
		Label:      "ATC 20 Inspection",
		Comment:    "checked, false",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"checked", "false"},
	}
	availableServicesDef = &xscmsg.FieldDef{
		Tag:        "44.",
		Annotation: "available-services",
		Label:      "Available Services",
	}
	mouDef = &xscmsg.FieldDef{
		Tag:        "45.",
		Annotation: "MOU",
		Label:      "",
	}
	floorPlanDef = &xscmsg.FieldDef{
		Tag:        "46.",
		Annotation: "floor-plan",
		Label:      "Floor Plan",
	}
	managedByCodeDef = &xscmsg.FieldDef{
		Tag:        "50a.",
		Annotation: "managed-by-code",
		ReadOnly:   true,
		Validators: []xscmsg.Validator{setManagedByCode, xscform.ValidateChoices},
		Choices:    []string{"American Red Cross", "Private", "Community", "Government", "Other"},
	}
	managedByDefV21 = &xscmsg.FieldDef{
		Tag:        "50a.",
		Annotation: "managed-by",
		Label:      "Managed By",
		Comment:    "required-for-complete: American Red Cross, Private, Community, Government, Other",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"American Red Cross", "Private", "Community", "Government", "Other"},
	}
	managedByDefV22 = &xscmsg.FieldDef{
		Tag:        "49a.",
		Annotation: "managed-by",
		Label:      "Managed By",
		Comment:    "required-for-complete: American Red Cross, Private, Community, Government, ...",
		Validators: []xscmsg.Validator{requiredForComplete},
	}
	managedByDetailDef = &xscmsg.FieldDef{
		Tag:        "50b.",
		Annotation: "managed-by-detail",
		Label:      "Managed By Detail",
	}
	primaryContactDef = &xscmsg.FieldDef{
		Tag:        "51a.",
		Annotation: "primary-contact",
		Label:      "Primary Contact",
		Comment:    "required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete},
	}
	primaryPhoneDef = &xscmsg.FieldDef{
		Tag:        "51b.",
		Annotation: "primary-phone",
		Label:      "Primary Contact Phone",
		Comment:    "phone-number required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidatePhoneNumber},
	}
	secondaryContactDef = &xscmsg.FieldDef{
		Tag:        "52a.",
		Annotation: "secondary-contact",
		Label:      "Secondary Contact",
	}
	secondaryPhoneDef = &xscmsg.FieldDef{
		Tag:        "52b.",
		Annotation: "secondary-phone",
		Label:      "Secondary Contact Phone",
		Comment:    "phone-number",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	tacticalCallDef = &xscmsg.FieldDef{
		Tag:        "60.",
		Annotation: "tactical-call",
		Label:      "Tactical Call Sign",
	}
	repeaterCallDef = &xscmsg.FieldDef{
		Tag:        "61.",
		Annotation: "repeater-call",
		Label:      "Repeater Call Sign",
		Comment:    "call-sign",
		Validators: []xscmsg.Validator{xscform.ValidateCallSign},
	}
	repeaterInputDef = &xscmsg.FieldDef{
		Tag:        "62a.",
		Annotation: "repeater-input",
		Label:      "Repeater Input (MHz)",
		Comment:    "frequency",
		Validators: []xscmsg.Validator{validateFrequency},
	}
	repeaterInputToneDef = &xscmsg.FieldDef{
		Tag:        "62b.",
		Annotation: "repeater-input-tone",
		Label:      "Repeater Input Tone or Code",
	}
	repeaterOutputDef = &xscmsg.FieldDef{
		Tag:        "63a.",
		Annotation: "repeater-output",
		Label:      "Repeater Output (MHz)",
		Comment:    "frequency",
		Validators: []xscmsg.Validator{validateFrequency},
	}
	repeaterOutputToneDef = &xscmsg.FieldDef{
		Tag:        "63b.",
		Annotation: "repeater-output-tone",
		Label:      "Repeater Output Tone or Code",
	}
	repeaterOffsetDef = &xscmsg.FieldDef{
		Tag:        "62c.",
		Annotation: "repeater-offset",
		Label:      "Repeater Offset (MHz or \"+\" or \"-\" for standard)",
		Comment:    "frequency-offset",
		Validators: []xscmsg.Validator{validateFrequencyOffset},
	}
	commentsDef = &xscmsg.FieldDef{
		Tag:        "70.",
		Annotation: "comments",
		Label:      "Comments",
	}
	removeFromActiveListDef = &xscmsg.FieldDef{
		Tag:        "71.",
		Annotation: "remove-from-active-list",
		Label:      "Remove from List",
		Comment:    "boolean",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
)

func setShelterCityCode(f *xscmsg.Field, m *xscmsg.Message, strict bool) string {
	if !strict {
		xscform.SetChoiceFieldFromValue(f, m.Field("34b.").Value)
	}
	return ""
}

func setManagedByCode(f *xscmsg.Field, m *xscmsg.Message, strict bool) string {
	if !strict {
		xscform.SetChoiceFieldFromValue(f, m.Field("49a.").Value)
	}
	return ""
}

func validateFrequency(f *xscmsg.Field, _ *xscmsg.Message, _ bool) string {
	if f.Value != "" && !frequencyRE.MatchString(f.Value) {
		return fmt.Sprintf("%q is not a valid frequency value for field %q", f.Value, f.Def.Tag)
	}
	return ""
}

func validateFrequencyOffset(f *xscmsg.Field, _ *xscmsg.Message, _ bool) string {
	if f.Value != "" && !frequencyOffsetRE.MatchString(f.Value) {
		return fmt.Sprintf("%q is not a valid frequency offset value for field %q", f.Value, f.Def.Tag)
	}
	return ""
}

func requiredForComplete(f *xscmsg.Field, m *xscmsg.Message, _ bool) string {
	if m.Field("19.").Value != "Complete" {
		return ""
	}
	if f.Value == "" {
		return fmt.Sprintf("The %q field must have a value when the \"19.\" field is set to \"Complete\".", f.Def.Tag)
	}
	return ""
}
