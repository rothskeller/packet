package eoc213rr

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/xscform"
)

// Tag identifies EOC-213RR forms.
const Tag = "EOC213RR"

func init() {
	xscmsg.RegisterCreate(formtype, create)
	xscmsg.RegisterType(recognize)

	// Our toICSPosition and toLocation fields are variants of the standard
	// ones, adding default values to them.
	toICSPositionDef.DefaultValue = "Planning Section"
	toICSPositionDef.Choices = []string{"Planning Section"}
	toLocationDef.DefaultValue = "County EOC"
	toLocationDef.Choices = []string{"County EOC"}
}

func create() *xscmsg.Message {
	return xscform.CreateForm(formtype, fieldDefs)
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	if form == nil || form.FormType != formtype.HTML || xscmsg.OlderVersion(form.FormVersion, "2.0") {
		return nil
	}
	return xscform.AdoptForm(formtype, fieldDefs, msg, form)
}

var formtype = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "EOC-213RR resource request form",
	Article:     "an",
	HTML:        "form-scco-eoc-213rr.html",
	Version:     "2.3",
	SubjectFunc: xscform.EncodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

var fieldDefs = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	xscform.HandlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef,
	xscform.ToNameDef, xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// EOC-213RR fields
	incidentNameDef, dateInitiatedDef, timeInitiatedDef, trackingNumberDef, requestedByDef, preparedByDef, approvedByDef,
	qtyUnitDef, resourceDescriptionDef, resourceArrivalDef, priorityDef, estdCostDef, deliverToDef, deliverToLocationDef,
	substitutesDef, equipmentOperatorDef, lodgingDef, fuelDef, fuelTypeDef, powerDef, mealsDef, maintenanceDef, waterDef,
	otherDef, instructionsDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var (
	toICSPositionDef = *xscform.ToICSPositionDef // modified in func init
	toLocationDef    = *xscform.ToLocationDef    // modified in func init
	incidentNameDef  = &xscmsg.FieldDef{
		Tag:   "21.",
		Label: "1. Incident Name",
		Key:   xscmsg.FSubject,
		Flags: xscmsg.Required,
	}
	dateInitiatedDef = &xscmsg.FieldDef{
		Tag:        "22.",
		Label:      "2. Date Initiated",
		Comment:    "MM/DD/YYYY",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
		Flags:      xscmsg.Required,
	}
	timeInitiatedDef = &xscmsg.FieldDef{
		Tag:        "23.",
		Label:      "3. Time Initiated",
		Comment:    "HH:MM",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
		Flags:      xscmsg.Required,
	}
	trackingNumberDef = &xscmsg.FieldDef{
		Tag:   "24.",
		Label: "4. Tracking Number (OA EOC)",
		Flags: xscmsg.Readonly,
	}
	requestedByDef = &xscmsg.FieldDef{
		Tag:   "25.",
		Label: "5. Requested by",
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	preparedByDef = &xscmsg.FieldDef{
		Tag:   "26.",
		Label: "6. Prepared by",
		Flags: xscmsg.Multiline,
	}
	approvedByDef = &xscmsg.FieldDef{
		Tag:   "27.",
		Label: "7. Approved By",
		Flags: xscmsg.Multiline,
	}
	qtyUnitDef = &xscmsg.FieldDef{
		Tag:   "28.",
		Label: "8. Qty/Unit",
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	resourceDescriptionDef = &xscmsg.FieldDef{
		Tag:   "29.",
		Label: "9. Resource Description",
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	resourceArrivalDef = &xscmsg.FieldDef{
		Tag:   "30.",
		Label: "10. Arrival (date/time)",
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	priorityDef = &xscmsg.FieldDef{
		Tag:        "31.",
		Label:      "11. Priority",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Now", "High", "Medium", "Low"},
		Flags:      xscmsg.Required,
	}
	estdCostDef = &xscmsg.FieldDef{
		Tag:   "32.",
		Label: "12. Est'd Cost",
		Flags: xscmsg.Multiline,
	}
	deliverToDef = &xscmsg.FieldDef{
		Tag:   "33.",
		Label: "13. Deliver to",
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	deliverToLocationDef = &xscmsg.FieldDef{
		Tag:   "34.",
		Label: "14. Location",
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	substitutesDef = &xscmsg.FieldDef{
		Tag:   "35.",
		Label: "15. Substitute/Suggested Sources",
		Flags: xscmsg.Multiline,
	}
	equipmentOperatorDef = &xscmsg.FieldDef{
		Tag:        "36a.",
		Label:      "16. Supplemental Requirements: Equipment Operator",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	lodgingDef = &xscmsg.FieldDef{
		Tag:        "36b.",
		Label:      "16. Supplemental Requirements: Lodging",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	fuelDef = &xscmsg.FieldDef{
		Tag:        "36c.",
		Label:      "16. Supplemental Requirements: Fuel",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	fuelTypeDef = &xscmsg.FieldDef{
		Tag:   "36d.",
		Label: "16. Supplemental Requirements: Fuel Type",
	}
	powerDef = &xscmsg.FieldDef{
		Tag:        "36e.",
		Label:      "16. Supplemental Requirements: Power",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	mealsDef = &xscmsg.FieldDef{
		Tag:        "36f.",
		Label:      "16. Supplemental Requirements: Meals",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	maintenanceDef = &xscmsg.FieldDef{
		Tag:        "36g.",
		Label:      "16. Supplemental Requirements: Maintenance",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	waterDef = &xscmsg.FieldDef{
		Tag:        "36h.",
		Label:      "16. Supplemental Requirements: Water",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	otherDef = &xscmsg.FieldDef{
		Tag:        "36i.",
		Label:      "16. Supplemental Requirements: Other",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	instructionsDef = &xscmsg.FieldDef{
		Tag:   "37.",
		Label: "17. Special Instructions",
		Key:   xscmsg.FBody,
		Flags: xscmsg.Multiline,
	}
)
