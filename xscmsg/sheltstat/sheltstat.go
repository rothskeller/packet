package sheltstat

import (
	"fmt"
	"regexp"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscform"
)

// Tag identifies shelter status forms.
const (
	Tag = "SheltStat"
)

var (
	frequencyRE       = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	frequencyOffsetRE = regexp.MustCompile(`^[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+|[-+]$`)
)

func init() {
	xscmsg.RegisterCreate(Tag, create)
	xscmsg.RegisterType(recognize)
}

func create() xscmsg.Message {
	return xscform.CreateForm(formtype22, makeFields22())
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.Message {
	if form == nil || form.FormType != formtype22.HTML {
		return nil
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.2") {
		return xscform.AdoptForm(formtype22, makeFields22(), msg, form)
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.0") {
		return xscform.AdoptForm(formtype20, makeFields20(), msg, form)
	}
	return nil
}

var formtype22 = &xscmsg.MessageType{
	Tag:     Tag,
	Name:    "OA jurisdiction status form",
	Article: "an",
	HTML:    "form-oa-shelter-status.html",
	Version: "2.2",
}

func makeFields22() []xscmsg.Field {
	return []xscmsg.Field{
		xscform.FOriginMessageNumber(),
		xscform.FDestinationMessageNumber(),
		xscform.FMessageDate(),
		xscform.FMessageTime(),
		&handlingField{ChoicesField: *xscform.FHandling().(*xscform.ChoicesField)},
		&toICSPositionField{*xscform.FToICSPosition().(*xscform.Field)},
		xscform.FFromICSPosition(),
		&toLocationField{*xscform.FToLocation().(*xscform.Field)},
		xscform.FFromLocation(),
		xscform.FToName(),
		xscform.FFromName(),
		xscform.FToContact(),
		xscform.FFromContact(),
		&xscform.ChoicesField{Field: *xscform.NewField(reportTypeID, true), Choices: reportTypeChoices},
		xscform.NewField(shelterNameID, true),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(shelterTypeID, false), Choices: shelterTypeChoices}},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(shelterStatusID, false), Choices: shelterStatusChoices}},
		&requiredForCompleteField{Field: *xscform.NewField(shelterAddressID, false)},
		&shelterCityCodeField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(shelterCityCodeID, false), Choices: shelterCityChoices}},
		&requiredForCompleteField{Field: *xscform.NewField(shelterCityID22, false)},
		&requiredForCompleteField{Field: *xscform.NewField(shelterStateID, false)},
		xscform.NewField(shelterZipID, false),
		&xscform.RealNumberField{Field: *xscform.NewField(latitudeID, false)},
		&xscform.RealNumberField{Field: *xscform.NewField(longitudeID, false)},
		&requiredForCompleteCardinalNumberField{CardinalNumberField: xscform.CardinalNumberField{Field: *xscform.NewField(capacityID, false)}},
		&requiredForCompleteCardinalNumberField{CardinalNumberField: xscform.CardinalNumberField{Field: *xscform.NewField(occupancyID, false)}},
		xscform.NewField(mealsID, false),
		xscform.NewField(nssID, false),
		&xscform.ChoicesField{Field: *xscform.NewField(petFriendlyID, false), Choices: checkedFalseChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(basicSafetyID, false), Choices: checkedFalseChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(atc20ID, false), Choices: checkedFalseChoices},
		xscform.NewField(availableServicesID, false),
		xscform.NewField(mouID, false),
		xscform.NewField(floorPlanID, false),
		&managedByCodeField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(managedByCodeID, false), Choices: managedByChoices}},
		&requiredForCompleteField{Field: *xscform.NewField(managedByID22, false)},
		xscform.NewField(managedByDetailID, false),
		&requiredForCompleteField{Field: *xscform.NewField(primaryContactID, false)},
		&requiredForCompletePhoneNumberField{PhoneNumberField: xscform.PhoneNumberField{Field: *xscform.NewField(primaryPhoneID, false)}},
		xscform.NewField(secondaryContactID, false),
		&xscform.PhoneNumberField{Field: *xscform.NewField(secondaryPhoneID, false)},
		xscform.NewField(tacticalCallID, false),
		&xscform.CallSignField{Field: *xscform.NewField(repeaterCallID, false)},
		&frequencyField{Field: *xscform.NewField(repeaterInputID, false)},
		xscform.NewField(repeaterInputToneID, false),
		&frequencyField{Field: *xscform.NewField(repeaterOutputID, false)},
		xscform.NewField(repeaterOutputToneID, false),
		&frequencyOffsetField{Field: *xscform.NewField(repeaterOffsetID, false)},
		xscform.NewField(commentsID, false),
		&xscform.BooleanField{Field: *xscform.NewField(removeFromActiveListID, false)},
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		xscform.FOpName(),
		xscform.FOpCall(),
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var formtype20 = &xscmsg.MessageType{
	Tag:     Tag,
	Name:    "OA jurisdiction status form",
	Article: "an",
	HTML:    "form-oa-shelter-status.html",
	Version: "2.1",
}

func makeFields20() []xscmsg.Field {
	return []xscmsg.Field{
		xscform.FOriginMessageNumber(),
		xscform.FDestinationMessageNumber(),
		xscform.FMessageDate(),
		xscform.FMessageTime(),
		&handlingField{ChoicesField: *xscform.FHandling().(*xscform.ChoicesField)},
		&toICSPositionField{*xscform.FToICSPosition().(*xscform.Field)},
		xscform.FFromICSPosition(),
		&toLocationField{*xscform.FToLocation().(*xscform.Field)},
		xscform.FFromLocation(),
		xscform.FToName(),
		xscform.FFromName(),
		xscform.FToContact(),
		xscform.FFromContact(),
		&xscform.ChoicesField{Field: *xscform.NewField(reportTypeID, true), Choices: reportTypeChoices},
		xscform.NewField(shelterNameID, true),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(shelterTypeID, false), Choices: shelterTypeChoices}},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(shelterStatusID, false), Choices: shelterStatusChoices}},
		&requiredForCompleteField{Field: *xscform.NewField(shelterAddressID, false)},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(shelterCityID20, false), Choices: shelterCityChoices}},
		&requiredForCompleteField{Field: *xscform.NewField(shelterStateID, false)},
		xscform.NewField(shelterZipID, false),
		&xscform.RealNumberField{Field: *xscform.NewField(latitudeID, false)},
		&xscform.RealNumberField{Field: *xscform.NewField(longitudeID, false)},
		&requiredForCompleteCardinalNumberField{CardinalNumberField: xscform.CardinalNumberField{Field: *xscform.NewField(capacityID, false)}},
		&requiredForCompleteCardinalNumberField{CardinalNumberField: xscform.CardinalNumberField{Field: *xscform.NewField(occupancyID, false)}},
		xscform.NewField(mealsID, false),
		xscform.NewField(nssID, false),
		&xscform.ChoicesField{Field: *xscform.NewField(petFriendlyID, false), Choices: checkedFalseChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(basicSafetyID, false), Choices: checkedFalseChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(atc20ID, false), Choices: checkedFalseChoices},
		xscform.NewField(availableServicesID, false),
		xscform.NewField(mouID, false),
		xscform.NewField(floorPlanID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(managedByID20, false), Choices: managedByChoices}},
		xscform.NewField(managedByDetailID, false),
		&requiredForCompleteField{Field: *xscform.NewField(primaryContactID, false)},
		&requiredForCompletePhoneNumberField{PhoneNumberField: xscform.PhoneNumberField{Field: *xscform.NewField(primaryPhoneID, false)}},
		xscform.NewField(secondaryContactID, false),
		&xscform.PhoneNumberField{Field: *xscform.NewField(secondaryPhoneID, false)},
		xscform.NewField(tacticalCallID, false),
		&xscform.CallSignField{Field: *xscform.NewField(repeaterCallID, false)},
		&frequencyField{Field: *xscform.NewField(repeaterInputID, false)},
		xscform.NewField(repeaterInputToneID, false),
		&frequencyField{Field: *xscform.NewField(repeaterOutputID, false)},
		xscform.NewField(repeaterOutputToneID, false),
		&frequencyOffsetField{Field: *xscform.NewField(repeaterOffsetID, false)},
		xscform.NewField(commentsID, false),
		&xscform.BooleanField{Field: *xscform.NewField(removeFromActiveListID, false)},
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		xscform.FOpName(),
		xscform.FOpCall(),
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var (
	reportTypeID = &xscmsg.FieldID{
		Tag:        "19.",
		Annotation: "report-type",
		Label:      "Report Type",
		Comment:    "required: Update, Complete",
	}
	reportTypeChoices = []string{"Update", "Complete"}
	shelterNameID     = &xscmsg.FieldID{
		Tag:        "32.",
		Annotation: "shelter-name",
		Label:      "Shelter Name",
		Comment:    "required",
		Canonical:  xscmsg.FSubject,
	}
	shelterTypeID = &xscmsg.FieldID{
		Tag:        "30.",
		Annotation: "shelter-type",
		Label:      "Shelter Type",
		Comment:    "required-for-complete: Type 1, Type 2, Type 3, Type 4",
	}
	shelterTypeChoices = []string{"Type 1", "Type 2", "Type 3", "Type 4"}
	shelterStatusID    = &xscmsg.FieldID{
		Tag:        "31.",
		Annotation: "shelter-status",
		Label:      "Status",
		Comment:    "required-for-complete: Open, Closed, Full",
	}
	shelterStatusChoices = []string{"Open", "Closed", "Full"}
	shelterAddressID     = &xscmsg.FieldID{
		Tag:        "33a.",
		Annotation: "shelter-address",
		Label:      "Address",
		Comment:    "required-for-complete",
	}
	shelterCityCodeID = &xscmsg.FieldID{
		Tag:        "33b.",
		Annotation: "shelter-city-code",
		Comment:    "required-for-complete",
		ReadOnly:   true,
	}
	shelterCityID20 = &xscmsg.FieldID{
		Tag:        "33b.",
		Annotation: "shelter-city",
		Label:      "City",
		Comment:    "required-for-complete: Campbell, Cupertino, Gilroy, Los Altos, Los Altos Hills, Los Gatos, Milpitas, Monte Sereno, Morgan Hill, Mountain View, Palo Alto, San Jose, Santa Clara, Saratoga, Sunnyvale, Unincorporated",
	}
	shelterCityChoices = []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"}
	shelterCityID22    = &xscmsg.FieldID{
		Tag:        "34b.",
		Annotation: "shelter-city",
		Label:      "City",
		Comment:    "required-for-complete: Campbell, Cupertino, Gilroy, Los Altos, Los Altos Hills, Los Gatos, Milpitas, Monte Sereno, Morgan Hill, Mountain View, Palo Alto, San Jose, Santa Clara, Saratoga, Sunnyvale, ...",
	}
	shelterStateID = &xscmsg.FieldID{
		Tag:        "33c.",
		Annotation: "shelter-state",
		Label:      "State",
	}
	shelterZipID = &xscmsg.FieldID{
		Tag:        "33d.",
		Annotation: "shelter-zip",
		Label:      "Zip",
	}
	latitudeID = &xscmsg.FieldID{
		Tag:        "37a.",
		Annotation: "latitude",
		Label:      "Latitude",
		Comment:    "real-number",
	}
	longitudeID = &xscmsg.FieldID{
		Tag:        "37b.",
		Annotation: "longitude",
		Label:      "Longitude",
		Comment:    "real-number",
	}
	capacityID = &xscmsg.FieldID{
		Tag:        "40a.",
		Annotation: "capacity",
		Label:      "Capacity",
		Comment:    "cardinal-number required-for-complete",
	}
	occupancyID = &xscmsg.FieldID{
		Tag:        "40b.",
		Annotation: "occupancy",
		Label:      "Occupancy",
		Comment:    "cardinal-number required-for-complete",
	}
	mealsID = &xscmsg.FieldID{
		Tag:        "41.",
		Annotation: "meals",
		Label:      "Meals Served (last 24 hours)",
	}
	nssID = &xscmsg.FieldID{
		Tag:        "42.",
		Annotation: "NSS",
		Label:      "NSS Number",
	}
	petFriendlyID = &xscmsg.FieldID{
		Tag:        "43a.",
		Annotation: "pet-friendly",
		Label:      "Pet Friendly",
		Comment:    "checked, false",
	}
	basicSafetyID = &xscmsg.FieldID{
		Tag:        "43b.",
		Annotation: "basic-safety",
		Label:      "Basic Safety Inspection",
		Comment:    "checked, false",
	}
	atc20ID = &xscmsg.FieldID{
		Tag:        "43c.",
		Annotation: "ATC-20",
		Label:      "ATC 20 Inspection",
		Comment:    "checked, false",
	}
	checkedFalseChoices = []string{"checked", "false"}
	availableServicesID = &xscmsg.FieldID{
		Tag:        "44.",
		Annotation: "available-services",
		Label:      "Available Services",
	}
	mouID = &xscmsg.FieldID{
		Tag:        "45.",
		Annotation: "MOU",
		Label:      "",
	}
	floorPlanID = &xscmsg.FieldID{
		Tag:        "46.",
		Annotation: "floor-plan",
		Label:      "Floor Plan",
	}
	managedByCodeID = &xscmsg.FieldID{
		Tag:        "50a.",
		Annotation: "managed-by-code",
		ReadOnly:   true,
	}
	managedByID20 = &xscmsg.FieldID{
		Tag:        "50a.",
		Annotation: "managed-by",
		Label:      "Managed By",
		Comment:    "required-for-complete: American Red Cross, Private, Community, Government, Other",
	}
	managedByChoices = []string{"American Red Cross", "Private", "Community", "Government", "Other"}
	managedByID22    = &xscmsg.FieldID{
		Tag:        "49a.",
		Annotation: "managed-by",
		Label:      "Managed By",
		Comment:    "required-for-complete: American Red Cross, Private, Community, Government, ...",
	}
	managedByDetailID = &xscmsg.FieldID{
		Tag:        "50b.",
		Annotation: "managed-by-detail",
		Label:      "Managed By Detail",
	}
	primaryContactID = &xscmsg.FieldID{
		Tag:        "51a.",
		Annotation: "primary-contact",
		Label:      "Primary Contact",
		Comment:    "required-for-complete",
	}
	primaryPhoneID = &xscmsg.FieldID{
		Tag:        "51b.",
		Annotation: "primary-phone",
		Label:      "Primary Contact Phone",
		Comment:    "phone-number required-for-complete",
	}
	secondaryContactID = &xscmsg.FieldID{
		Tag:        "52a.",
		Annotation: "secondary-contact",
		Label:      "Secondary Contact",
	}
	secondaryPhoneID = &xscmsg.FieldID{
		Tag:        "52b.",
		Annotation: "secondary-phone",
		Label:      "Secondary Contact Phone",
		Comment:    "phone-number",
	}
	tacticalCallID = &xscmsg.FieldID{
		Tag:        "60.",
		Annotation: "tactical-call",
		Label:      "Tactical Call Sign",
	}
	repeaterCallID = &xscmsg.FieldID{
		Tag:        "61.",
		Annotation: "repeater-call",
		Label:      "Repeater Call Sign",
		Comment:    "call-sign",
	}
	repeaterInputID = &xscmsg.FieldID{
		Tag:        "62a.",
		Annotation: "repeater-input",
		Label:      "Repeater Input (MHz)",
		Comment:    "frequency",
	}
	repeaterInputToneID = &xscmsg.FieldID{
		Tag:        "62b.",
		Annotation: "repeater-input-tone",
		Label:      "Repeater Input Tone or Code",
	}
	repeaterOutputID = &xscmsg.FieldID{
		Tag:        "63a.",
		Annotation: "repeater-output",
		Label:      "Repeater Output (MHz)",
		Comment:    "frequency",
	}
	repeaterOutputToneID = &xscmsg.FieldID{
		Tag:        "63b.",
		Annotation: "repeater-output-tone",
		Label:      "Repeater Output Tone or Code",
	}
	repeaterOffsetID = &xscmsg.FieldID{
		Tag:        "62c.",
		Annotation: "repeater-offset",
		Label:      "Repeater Offset (MHz or \"+\" or \"-\" for standard)",
		Comment:    "frequency-offset",
	}
	commentsID = &xscmsg.FieldID{
		Tag:        "70.",
		Annotation: "comments",
		Label:      "Comments",
	}
	removeFromActiveListID = &xscmsg.FieldID{
		Tag:        "71.",
		Annotation: "remove-from-active-list",
		Label:      "Remove from List",
		Comment:    "boolean",
	}
)

type handlingField struct{ xscform.ChoicesField }

func (f *handlingField) Default() string { return "PRIORITY" }

type toICSPositionField struct{ xscform.Field }

func (f *toICSPositionField) Default() string { return "Mass Care and Shelter Unit" }

type toLocationField struct{ xscform.Field }

func (f *toLocationField) Default() string { return "County EOC" }

type shelterCityCodeField struct{ xscform.ChoicesField }

func (f *shelterCityCodeField) Validate(msg xscmsg.Message, strict bool) string {
	if f.Get() == "" {
		j := msg.Field("34b.").Get()
		if j != "" {
			for _, choice := range shelterCityChoices {
				if j == choice {
					f.Set(j)
					break
				}
			}
			if f.Get() == "" {
				f.Set(shelterCityChoices[len(shelterCityChoices)-1])
			}
		}
	}
	return f.ChoicesField.Validate(msg, strict)
}

type managedByCodeField struct{ xscform.ChoicesField }

func (f *managedByCodeField) Validate(msg xscmsg.Message, strict bool) string {
	if f.Get() == "" {
		j := msg.Field("49a.").Get()
		if j != "" {
			for _, choice := range managedByChoices {
				if j == choice {
					f.Set(j)
					break
				}
			}
			if f.Get() == "" {
				f.Set(managedByChoices[len(managedByChoices)-1])
			}
		}
	}
	return f.ChoicesField.Validate(msg, strict)
}

// frequencyField is a field with a frequency value.
type frequencyField struct{ xscform.Field }

// Validate ensures that the value is a properly-formatted frequency value.
func (f *frequencyField) Validate(msg xscmsg.Message, strict bool) string {
	if v := f.Get(); v != "" && !frequencyRE.MatchString(v) {
		return fmt.Sprintf("%q is not a valid frequency value for field %q", v, f.ID().Tag)
	}
	return ""
}

// frequencyOffsetField is a field with a frequency offset value.
type frequencyOffsetField struct{ xscform.Field }

// Validate ensures that the value is a properly-formatted frequencyOffset value.
func (f *frequencyOffsetField) Validate(msg xscmsg.Message, strict bool) string {
	if v := f.Get(); v != "" && !frequencyOffsetRE.MatchString(v) {
		return fmt.Sprintf("%q is not a valid frequency offset value for field %q", v, f.ID().Tag)
	}
	return ""
}

type requiredForCompleteField struct{ xscform.Field }

func (f *requiredForCompleteField) Validate(msg xscmsg.Message, strict bool) string {
	if err := validateRequiredIfComplete(f, msg); err != "" {
		return err
	}
	return f.Field.Validate(msg, strict)
}

type requiredForCompleteCardinalNumberField struct{ xscform.CardinalNumberField }

func (f *requiredForCompleteCardinalNumberField) Validate(msg xscmsg.Message, strict bool) string {
	if err := validateRequiredIfComplete(f, msg); err != "" {
		return err
	}
	return f.CardinalNumberField.Validate(msg, strict)
}

type requiredForCompleteChoicesField struct{ xscform.ChoicesField }

func (f *requiredForCompleteChoicesField) Validate(msg xscmsg.Message, strict bool) string {
	if err := validateRequiredIfComplete(f, msg); err != "" {
		return err
	}
	return f.ChoicesField.Validate(msg, strict)
}

type requiredForCompletePhoneNumberField struct{ xscform.PhoneNumberField }

func (f *requiredForCompletePhoneNumberField) Validate(msg xscmsg.Message, strict bool) string {
	if err := validateRequiredIfComplete(f, msg); err != "" {
		return err
	}
	return f.PhoneNumberField.Validate(msg, strict)
}

func validateRequiredIfComplete(f xscmsg.Field, msg xscmsg.Message) string {
	if f.Get() == "" && msg.Field("19.").Get() == "Complete" {
		return fmt.Sprintf("field %q needs a value when field \"19.\" is \"Complete\"", f.ID().Tag)
	}
	return ""
}
