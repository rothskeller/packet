package sheltstat

import (
	"fmt"
	"regexp"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/xscform"
)

// Tag identifies shelter status forms.
const Tag = "SheltStat"

var (
	frequencyRE       = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	frequencyOffsetRE = regexp.MustCompile(`^[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+|[-+]$`)
)

func init() {
	xscmsg.RegisterCreate(formtype, create)
	xscmsg.RegisterType(recognize)

	// Our handling, toICSPosition, and toLocation fields are variants of
	// the standard ones, adding default values to them.
	handlingDef.DefaultValue = "PRIORITY"
	toICSPositionDef.DefaultValue = "Mass Care and Shelter Unit"
	toICSPositionDef.Choices = []string{"Care and Shelter Branch", "Mass Care and Shelter Unit", "Operations Section"}
	toLocationDef.DefaultValue = "County EOC"
	toLocationDef.Choices = []string{"County EOC"}
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
	Name:        "OA Shelter Status form",
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
		Key:        xscmsg.FComplete,
		Label:      "Report Type",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Update", "Complete"},
		Flags:      xscmsg.Required,
	}
	shelterNameDef = &xscmsg.FieldDef{
		Tag:   "32.",
		Label: "Shelter Name",
		Flags: xscmsg.Required,
		Key:   xscmsg.FSubject,
	}
	shelterTypeDef = &xscmsg.FieldDef{
		Tag:        "30.",
		Label:      "Shelter Type",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Type 1", "Type 2", "Type 3", "Type 4"},
		Flags:      xscmsg.RequiredForComplete,
	}
	shelterStatusDef = &xscmsg.FieldDef{
		Tag:        "31.",
		Label:      "Status",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Open", "Closed", "Full"},
		Flags:      xscmsg.RequiredForComplete,
	}
	shelterAddressDef = &xscmsg.FieldDef{
		Tag:   "33a.",
		Label: "Address",
		Flags: xscmsg.RequiredForComplete,
	}
	shelterCityCodeDef = &xscmsg.FieldDef{
		Tag:        "33b.",
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{setShelterCityCode, xscform.ValidateChoices},
		Choices:    []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
	}
	shelterCityDefV21 = &xscmsg.FieldDef{
		Tag:        "33b.",
		Label:      "City",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
		Flags:      xscmsg.RequiredForComplete,
	}
	shelterCityDefV22 = &xscmsg.FieldDef{
		Tag:   "34b.",
		Label: "City",
		Flags: xscmsg.RequiredForComplete,
	}
	shelterStateDef = &xscmsg.FieldDef{
		Tag:   "33c.",
		Label: "State",
	}
	shelterZipDef = &xscmsg.FieldDef{
		Tag:   "33d.",
		Label: "Zip",
	}
	latitudeDef = &xscmsg.FieldDef{
		Tag:        "37a.",
		Label:      "Latitude",
		Comment:    "number",
		Validators: []xscmsg.Validator{xscform.ValidateRealNumber},
	}
	longitudeDef = &xscmsg.FieldDef{
		Tag:        "37b.",
		Label:      "Longitude",
		Comment:    "number",
		Validators: []xscmsg.Validator{xscform.ValidateRealNumber},
	}
	capacityDef = &xscmsg.FieldDef{
		Tag:        "40a.",
		Label:      "Capacity",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
		Flags:      xscmsg.RequiredForComplete,
	}
	occupancyDef = &xscmsg.FieldDef{
		Tag:        "40b.",
		Label:      "Occupancy",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
		Flags:      xscmsg.RequiredForComplete,
	}
	mealsDef = &xscmsg.FieldDef{
		Tag:   "41.",
		Label: "Meals Served (last 24 hours)",
	}
	nssDef = &xscmsg.FieldDef{
		Tag:   "42.",
		Label: "NSS Number",
	}
	petFriendlyDef = &xscmsg.FieldDef{
		Tag:        "43a.",
		Label:      "Pet Friendly",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"checked", "false"},
	}
	basicSafetyDef = &xscmsg.FieldDef{
		Tag:        "43b.",
		Label:      "Basic Safety Inspection",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"checked", "false"},
	}
	atc20Def = &xscmsg.FieldDef{
		Tag:        "43c.",
		Label:      "ATC 20 Inspection",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"checked", "false"},
	}
	availableServicesDef = &xscmsg.FieldDef{
		Tag:   "44.",
		Label: "Available Services",
		Flags: xscmsg.Multiline,
	}
	mouDef = &xscmsg.FieldDef{
		Tag:   "45.",
		Label: "",
	}
	floorPlanDef = &xscmsg.FieldDef{
		Tag:   "46.",
		Label: "Floor Plan",
	}
	managedByCodeDef = &xscmsg.FieldDef{
		Tag:        "50a.",
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{setManagedByCode, xscform.ValidateChoices},
		Choices:    []string{"American Red Cross", "Private", "Community", "Government", "Other"},
	}
	managedByDefV21 = &xscmsg.FieldDef{
		Tag:        "50a.",
		Label:      "Managed By",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"American Red Cross", "Private", "Community", "Government", "Other"},
		Flags:      xscmsg.RequiredForComplete,
	}
	managedByDefV22 = &xscmsg.FieldDef{
		Tag:   "49a.",
		Label: "Managed By",
		Flags: xscmsg.RequiredForComplete,
	}
	managedByDetailDef = &xscmsg.FieldDef{
		Tag:   "50b.",
		Label: "Managed By Detail",
	}
	primaryContactDef = &xscmsg.FieldDef{
		Tag:   "51a.",
		Label: "Primary Contact",
		Flags: xscmsg.RequiredForComplete,
	}
	primaryPhoneDef = &xscmsg.FieldDef{
		Tag:        "51b.",
		Label:      "Primary Contact Phone",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
		Flags:      xscmsg.RequiredForComplete,
	}
	secondaryContactDef = &xscmsg.FieldDef{
		Tag:   "52a.",
		Label: "Secondary Contact",
	}
	secondaryPhoneDef = &xscmsg.FieldDef{
		Tag:        "52b.",
		Label:      "Secondary Contact Phone",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	tacticalCallDef = &xscmsg.FieldDef{
		Tag:   "60.",
		Label: "Tactical Call Sign",
	}
	repeaterCallDef = &xscmsg.FieldDef{
		Tag:        "61.",
		Label:      "Repeater Call Sign",
		Validators: []xscmsg.Validator{xscform.ValidateCallSign},
	}
	repeaterInputDef = &xscmsg.FieldDef{
		Tag:        "62a.",
		Label:      "Repeater Input (MHz)",
		Validators: []xscmsg.Validator{validateFrequency},
	}
	repeaterInputToneDef = &xscmsg.FieldDef{
		Tag:   "62b.",
		Label: "Repeater Input Tone or Code",
	}
	repeaterOutputDef = &xscmsg.FieldDef{
		Tag:        "63a.",
		Label:      "Repeater Output (MHz)",
		Validators: []xscmsg.Validator{validateFrequency},
	}
	repeaterOutputToneDef = &xscmsg.FieldDef{
		Tag:   "63b.",
		Label: "Repeater Output Tone or Code",
	}
	repeaterOffsetDef = &xscmsg.FieldDef{
		Tag:        "62c.",
		Label:      "Repeater Offset",
		Comment:    "+, -, number",
		Validators: []xscmsg.Validator{validateFrequencyOffset},
	}
	commentsDef = &xscmsg.FieldDef{
		Tag:   "70.",
		Label: "Comments",
		Key:   xscmsg.FBody,
		Flags: xscmsg.Multiline,
	}
	removeFromActiveListDef = &xscmsg.FieldDef{
		Tag:        "71.",
		Label:      "Remove from List",
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
