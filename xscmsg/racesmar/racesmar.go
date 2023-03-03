package racesmar

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/xscform"
)

// Tag identifies RACES MAR forms.
const (
	Tag = "RACES-MAR"
)

func init() {
	xscmsg.RegisterCreate(formtype23, create)
	xscmsg.RegisterType(recognize)

	// Our handling, toICSPosition, and toLocation fields are variants of
	// the standard ones, adding default values to them.
	handlingDef = *xscform.HandlingDef
	handlingDef.DefaultValue = "ROUTINE"
	toICSPositionDef = *xscform.ToICSPositionDef
	toICSPositionDef.DefaultValue = "RACES Chief Radio Officer"
	toICSPositionDef.Choices = []string{"Operations Section", "RACES Chief Radio Officer", "RACES Unit"}
	toLocationDef = *xscform.ToLocationDef
	toLocationDef.DefaultValue = "County EOC"
	toLocationDef.Choices = []string{"County EOC"}
}

func create() *xscmsg.Message {
	return xscform.CreateForm(formtype23, fieldDefsV23)
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	if form == nil {
		return nil
	}
	if form.FormType == formtype23.HTML && !xscmsg.OlderVersion(form.FormVersion, "2.3") {
		return xscform.AdoptForm(formtype23, fieldDefsV23, msg, form)
	}
	if form.FormType == formtype23.HTML && !xscmsg.OlderVersion(form.FormVersion, "2.1") {
		return xscform.AdoptForm(formtype23, fieldDefsV21, msg, form)
	}
	if form.FormType == formtype16.HTML && !xscmsg.OlderVersion(form.FormVersion, "1.6") {
		return xscform.AdoptForm(formtype16, fieldDefsV16, msg, form)
	}
	return nil
}

var formtype23 = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "RACES mutual aid request form",
	Article:     "a",
	HTML:        "form-oa-mutual-aid-request-v2.html",
	Version:     "2.3",
	SubjectFunc: xscform.EncodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

var formtype16 = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "RACES mutual aid request form",
	Article:     "a",
	HTML:        "form-oa-mutual-aid-request.html",
	Version:     "1.6",
	SubjectFunc: xscform.EncodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

// fieldDefsV23 replaces the single resourcesRole field for each resource with
// separate role and position fields, plus a computed resourcesRole field.  It
// also adds validation to the preferredType and mininumType fields.
var fieldDefsV23 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// RACES MAR fields
	agencyDef, eventNameDef, eventNumberDef, assignmentDef,
	resourcesQty1Def, role1Def, position1Def, resourcesRole1DefV23, preferredType1DefV23, minimumType1DefV23,
	resourcesQty2Def, role2Def, position2Def, resourcesRole2DefV23, preferredType2DefV23, minimumType2DefV23,
	resourcesQty3Def, role3Def, position3Def, resourcesRole3DefV23, preferredType3DefV23, minimumType3DefV23,
	resourcesQty4Def, role4Def, position4Def, resourcesRole4DefV23, preferredType4DefV23, minimumType4DefV23,
	resourcesQty5Def, role5Def, position5Def, resourcesRole5DefV23, preferredType5DefV23, minimumType5DefV23,
	arrivalDatesDef, arrivalTimesDef, neededDatesDef, neededTimesDef, reportingLocationDef, contactOnArrivalDef,
	travelInfoDef, requesterNameDef, requesterTitleDef, requesterContactDef, agencyApproverNameDef, agencyApproverTitleDef,
	agencyApproverContactDef, agencyApprovedDateDef, agencyApprovedTimeDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

// fieldDefsV21 replaces a single set of resource definition fields with five
// sets of fields.
var fieldDefsV21 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// RACES MAR fields
	agencyDef, eventNameDef, eventNumberDef, assignmentDef,
	resourcesQty1Def, resourcesRole1DefV21, preferredType1DefV21, minimumType1DefV21,
	resourcesQty2Def, resourcesRole2DefV21, preferredType2DefV21, minimumType2DefV21,
	resourcesQty3Def, resourcesRole3DefV21, preferredType3DefV21, minimumType3DefV21,
	resourcesQty4Def, resourcesRole4DefV21, preferredType4DefV21, minimumType4DefV21,
	resourcesQty5Def, resourcesRole5DefV21, preferredType5DefV21, minimumType5DefV21,
	arrivalDatesDef, arrivalTimesDef, neededDatesDef, neededTimesDef, reportingLocationDef, contactOnArrivalDef,
	travelInfoDef, requesterNameDef, requesterTitleDef, requesterContactDef, agencyApproverNameDef, agencyApproverTitleDef,
	agencyApproverContactDef, agencyApprovedDateDef, agencyApprovedTimeDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var fieldDefsV16 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// RACES MAR fields
	agencyDef, eventNameDef, eventNumberDef, assignmentDef,
	resourcesQtyDef, resourcesRoleDef, preferredTypeDef, minimumTypeDef,
	arrivalDatesDef, arrivalTimesDef, neededDatesDef, neededTimesDef, reportingLocationDef, contactOnArrivalDef,
	travelInfoDef, requesterNameDef, requesterTitleDef, requesterContactDef, agencyApproverNameDef, agencyApproverTitleDef,
	agencyApproverContactDef, agencyApprovedDateDef, agencyApprovedTimeDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var (
	handlingDef      = *xscform.HandlingDef      // modified in func init
	toICSPositionDef = *xscform.ToICSPositionDef // modified in func init
	toLocationDef    = *xscform.ToLocationDef    // modified in func init
	agencyDef        = &xscmsg.FieldDef{
		Tag:   "15.",
		Label: "Agency Name",
		Key:   xscmsg.FSubject,
		Flags: xscmsg.Required,
	}
	eventNameDef = &xscmsg.FieldDef{
		Tag:   "16a.",
		Label: "Event / Incident Name",
		Flags: xscmsg.Required,
	}
	eventNumberDef = &xscmsg.FieldDef{
		Tag:        "16b.",
		Label:      "Event / Incident Number",
		Validators: []xscmsg.Validator{},
	}
	assignmentDef = &xscmsg.FieldDef{
		Tag:   "17.",
		Label: "Assignment",
		Key:   xscmsg.FBody,
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	resourcesQtyDef = &xscmsg.FieldDef{
		Tag:        "18a.",
		Label:      "Qty",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
		Flags:      xscmsg.Required,
	}
	resourcesRoleDef = &xscmsg.FieldDef{
		Tag:   "18b.",
		Label: "Role/Position",
		Flags: xscmsg.Required,
	}
	preferredTypeDef = &xscmsg.FieldDef{
		Tag:   "18c.",
		Label: "Preferred Type",
		Flags: xscmsg.Required,
	}
	minimumTypeDef = &xscmsg.FieldDef{
		Tag:   "18d.",
		Label: "Minimum Type",
		Flags: xscmsg.Required,
	}
	resourcesQty1Def = &xscmsg.FieldDef{
		Tag:        "18.1a.",
		Label:      "Resource 1 Qty",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
		Flags:      xscmsg.Required,
	}
	role1Def = &xscmsg.FieldDef{
		Tag:        "18.1e.",
		Label:      "Resource 1 Role",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
		Flags:      xscmsg.Required,
	}
	position1Def = &xscmsg.FieldDef{
		Tag:        "18.1f.",
		Label:      "Resource 1 Position",
		Validators: []xscmsg.Validator{},
	}
	resourcesRole1DefV23 = &xscmsg.FieldDef{
		Tag:        "18.1b.",
		Validators: []xscmsg.Validator{setResourcesRole},
		Flags:      xscmsg.Readonly | xscmsg.Required,
	}
	resourcesRole1DefV21 = &xscmsg.FieldDef{
		Tag:   "18.1b.",
		Label: "Resource 1 Role/Position",
		Flags: xscmsg.Required,
	}
	preferredType1DefV23 = &xscmsg.FieldDef{
		Tag:        "18.1c.",
		Label:      "Resource 1 Preferred Type",
		Validators: []xscmsg.Validator{validateType},
	}
	preferredType1DefV21 = &xscmsg.FieldDef{
		Tag:   "18.1c.",
		Label: "Resource 1 Preferred Type",
		Flags: xscmsg.Required,
	}
	minimumType1DefV23 = &xscmsg.FieldDef{
		Tag:        "18.1d.",
		Label:      "Resource 1 Minimum Type",
		Validators: []xscmsg.Validator{validateType},
	}
	minimumType1DefV21 = &xscmsg.FieldDef{
		Tag:   "18.1d.",
		Label: "Resource 1 Minimum Type",
		Flags: xscmsg.Required,
	}
	resourcesQty2Def = &xscmsg.FieldDef{
		Tag:        "18.2a.",
		Label:      "Resource 2 Qty",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	role2Def = &xscmsg.FieldDef{
		Tag:        "18.2e.",
		Label:      "Resource 2 Role",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
	}
	position2Def = &xscmsg.FieldDef{
		Tag:        "18.2f.",
		Label:      "Resource 2 Position",
		Validators: []xscmsg.Validator{},
	}
	resourcesRole2DefV23 = &xscmsg.FieldDef{
		Tag:        "18.2b.",
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{setResourcesRole},
	}
	resourcesRole2DefV21 = &xscmsg.FieldDef{
		Tag:        "18.2b.",
		Label:      "Resource 2 Role/Position",
		Validators: []xscmsg.Validator{},
	}
	preferredType2DefV23 = &xscmsg.FieldDef{
		Tag:        "18.2c.",
		Label:      "Resource 2 Preferred Type",
		Validators: []xscmsg.Validator{validateType},
	}
	preferredType2DefV21 = &xscmsg.FieldDef{
		Tag:        "18.2c.",
		Label:      "Resource 2 Preferred Type",
		Validators: []xscmsg.Validator{},
	}
	minimumType2DefV23 = &xscmsg.FieldDef{
		Tag:        "18.2d.",
		Label:      "Resource 2 Minimum Type",
		Validators: []xscmsg.Validator{validateType},
	}
	minimumType2DefV21 = &xscmsg.FieldDef{
		Tag:        "18.2d.",
		Label:      "Resource 2 Minimum Type",
		Validators: []xscmsg.Validator{},
	}
	resourcesQty3Def = &xscmsg.FieldDef{
		Tag:        "18.3a.",
		Label:      "Resource 3 Qty",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	role3Def = &xscmsg.FieldDef{
		Tag:        "18.3e.",
		Label:      "Resource 3 Role",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
	}
	position3Def = &xscmsg.FieldDef{
		Tag:        "18.3f.",
		Label:      "Resource 3 Position",
		Validators: []xscmsg.Validator{},
	}
	resourcesRole3DefV23 = &xscmsg.FieldDef{
		Tag:        "18.3b.",
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{setResourcesRole},
	}
	resourcesRole3DefV21 = &xscmsg.FieldDef{
		Tag:        "18.3b.",
		Label:      "Resource 3 Role/Position",
		Validators: []xscmsg.Validator{},
	}
	preferredType3DefV23 = &xscmsg.FieldDef{
		Tag:        "18.3c.",
		Label:      "Resource 3 Preferred Type",
		Validators: []xscmsg.Validator{validateType},
	}
	preferredType3DefV21 = &xscmsg.FieldDef{
		Tag:        "18.3c.",
		Label:      "Resource 3 Preferred Type",
		Validators: []xscmsg.Validator{},
	}
	minimumType3DefV23 = &xscmsg.FieldDef{
		Tag:        "18.3d.",
		Label:      "Resource 3 Minimum Type",
		Validators: []xscmsg.Validator{validateType},
	}
	minimumType3DefV21 = &xscmsg.FieldDef{
		Tag:        "18.3d.",
		Label:      "Resource 3 Minimum Type",
		Validators: []xscmsg.Validator{},
	}
	resourcesQty4Def = &xscmsg.FieldDef{
		Tag:        "18.4a.",
		Label:      "Resource 4 Qty",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	role4Def = &xscmsg.FieldDef{
		Tag:        "18.4e.",
		Label:      "Resource 4 Role",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
	}
	position4Def = &xscmsg.FieldDef{
		Tag:        "18.4f.",
		Label:      "Resource 4 Position",
		Validators: []xscmsg.Validator{},
	}
	resourcesRole4DefV23 = &xscmsg.FieldDef{
		Tag:        "18.4b.",
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{setResourcesRole},
	}
	resourcesRole4DefV21 = &xscmsg.FieldDef{
		Tag:        "18.4b.",
		Label:      "Resource 4 Role/Position",
		Validators: []xscmsg.Validator{},
	}
	preferredType4DefV23 = &xscmsg.FieldDef{
		Tag:        "18.4c.",
		Label:      "Resource 4 Preferred Type",
		Validators: []xscmsg.Validator{validateType},
	}
	preferredType4DefV21 = &xscmsg.FieldDef{
		Tag:        "18.4c.",
		Label:      "Resource 4 Preferred Type",
		Validators: []xscmsg.Validator{},
	}
	minimumType4DefV23 = &xscmsg.FieldDef{
		Tag:        "18.4d.",
		Label:      "Resource 4 Minimum Type",
		Validators: []xscmsg.Validator{validateType},
	}
	minimumType4DefV21 = &xscmsg.FieldDef{
		Tag:        "18.4d.",
		Label:      "Resource 4 Minimum Type",
		Validators: []xscmsg.Validator{},
	}
	resourcesQty5Def = &xscmsg.FieldDef{
		Tag:        "18.5a.",
		Label:      "Resource 5 Qty",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	role5Def = &xscmsg.FieldDef{
		Tag:        "18.5e.",
		Label:      "Resource 5 Role",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
	}
	position5Def = &xscmsg.FieldDef{
		Tag:        "18.5f.",
		Label:      "Resource 5 Position",
		Validators: []xscmsg.Validator{},
	}
	resourcesRole5DefV23 = &xscmsg.FieldDef{
		Tag:        "18.5b.",
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{setResourcesRole},
	}
	resourcesRole5DefV21 = &xscmsg.FieldDef{
		Tag:        "18.5b.",
		Label:      "Resource 5 Role/Position",
		Validators: []xscmsg.Validator{},
	}
	preferredType5DefV23 = &xscmsg.FieldDef{
		Tag:        "18.5c.",
		Label:      "Resource 5 Preferred Type",
		Validators: []xscmsg.Validator{validateType},
	}
	preferredType5DefV21 = &xscmsg.FieldDef{
		Tag:        "18.5c.",
		Label:      "Resource 5 Preferred Type",
		Validators: []xscmsg.Validator{},
	}
	minimumType5DefV23 = &xscmsg.FieldDef{
		Tag:        "18.5d.",
		Label:      "Resource 5 Minimum Type",
		Validators: []xscmsg.Validator{validateType},
	}
	minimumType5DefV21 = &xscmsg.FieldDef{
		Tag:        "18.5d.",
		Label:      "Resource 5 Minimum Type",
		Validators: []xscmsg.Validator{},
	}
	arrivalDatesDef = &xscmsg.FieldDef{
		Tag:   "19a.",
		Label: "Requested Arrival Date(s)",
		Flags: xscmsg.Required,
	}
	arrivalTimesDef = &xscmsg.FieldDef{
		Tag:   "19b.",
		Label: "Requested Arrival Time(s)",
		Flags: xscmsg.Required,
	}
	neededDatesDef = &xscmsg.FieldDef{
		Tag:   "20a.",
		Label: "Needed Until Date(s)",
		Flags: xscmsg.Required,
	}
	neededTimesDef = &xscmsg.FieldDef{
		Tag:   "20b.",
		Label: "Needed Until Time(s)",
		Flags: xscmsg.Required,
	}
	reportingLocationDef = &xscmsg.FieldDef{
		Tag:   "21.",
		Label: "Reporting Location",
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	contactOnArrivalDef = &xscmsg.FieldDef{
		Tag:   "22.",
		Label: "Contact on Arrival",
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	travelInfoDef = &xscmsg.FieldDef{
		Tag:   "23.",
		Label: "Travel Info",
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	requesterNameDef = &xscmsg.FieldDef{
		Tag:   "24a.",
		Label: "Requested By Name",
		Flags: xscmsg.Required,
	}
	requesterTitleDef = &xscmsg.FieldDef{
		Tag:   "24b.",
		Label: "Requested By Title",
		Flags: xscmsg.Required,
	}
	requesterContactDef = &xscmsg.FieldDef{
		Tag:   "24c.",
		Label: "Requested By Contact",
		Flags: xscmsg.Required,
	}
	agencyApproverNameDef = &xscmsg.FieldDef{
		Tag:   "25a.",
		Label: "Approved By Name",
		Flags: xscmsg.Required,
	}
	agencyApproverTitleDef = &xscmsg.FieldDef{
		Tag:   "25b.",
		Label: "Approved By Title",
		Flags: xscmsg.Required,
	}
	agencyApproverContactDef = &xscmsg.FieldDef{
		Tag:   "25c.",
		Label: "Approved By Contact",
		Flags: xscmsg.Required,
	}
	agencyApprovedDateDef = &xscmsg.FieldDef{
		Tag:        "26a.",
		Label:      "Approved By Date",
		Comment:    "MM/DD/YYYY",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
		Flags:      xscmsg.Required,
	}
	agencyApprovedTimeDef = &xscmsg.FieldDef{
		Tag:        "26b.",
		Label:      "Approved By Time",
		Comment:    "HH:MM",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
		Flags:      xscmsg.Required,
	}
)

func setResourcesRole(f *xscmsg.Field, msg *xscmsg.Message, strict bool) string {
	var btag = f.Def.Tag
	var etag = strings.Replace(btag, "b", "e", 1)
	var ftag = strings.Replace(btag, "b", "f", 1)
	var eval = msg.Field(etag).Value
	var fval = msg.Field(ftag).Value
	var bval = eval
	if fval != "" {
		bval += " / " + fval
	}
	if !strict {
		f.Value = bval
	} else if f.Value != bval {
		return fmt.Sprintf("The value of field %q is not consistent with the values of fields %q and %q.", btag, etag, ftag)
	}
	return ""
}

func validateType(f *xscmsg.Field, msg *xscmsg.Message, strict bool) string {
	var tag = f.Def.Tag
	var etag = tag[:len(tag)-2] + "e."
	var eval = msg.Field(etag).Value
	if eval == "" {
		if f.Value != "" {
			return fmt.Sprintf("Field %q must not have a value unless field %q has a value.", tag, etag)
		}
		return ""
	}
	var val = f.Value
	if val == "" || val == "Type IV" || val == "Type V" {
		return ""
	}
	if !strict {
		if strings.EqualFold(val, "Type IV") {
			f.Value = "Type IV"
			return ""
		}
		if strings.EqualFold(val, "Type V") {
			f.Value = "Type V"
			return ""
		}
	}
	if len(val) != 2 || (val[1] != '1' && val[1] != '2' && val[1] != '3') {
		return fmt.Sprintf("%q is not a valid value for field %q.", val, tag)
	}
	tcode := val[:1]
	if !strict {
		tcode = strings.ToUpper(tcode)
	}
	if tcode != "F" && tcode != "N" && tcode != "P" && tcode != "S" {
		return fmt.Sprintf("%q is not a valid value for field %q.", val, tag)
	}
	if tcode != eval[:1] {
		return fmt.Sprintf("The value %q for field %q is not compatible with the value %q for field %q.", val, tag, eval, etag)
	}
	f.Value = tcode + val[1:2]
	return ""
}
