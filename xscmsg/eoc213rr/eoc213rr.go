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
	toICSPositionDef.Comment = "required: Planning Section, ..."
	toLocationDef.DefaultValue = "County EOC"
	toLocationDef.Comment = "required: County EOC, ..."
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
		Tag:        "21.",
		Annotation: "incident-name",
		Label:      "1. Incident Name",
		Comment:    "required",
		Key:        xscmsg.FSubject,
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	dateInitiatedDef = &xscmsg.FieldDef{
		Tag:        "22.",
		Annotation: "date",
		Label:      "2. Date Initiated",
		Comment:    "required date",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateDate},
	}
	timeInitiatedDef = &xscmsg.FieldDef{
		Tag:        "23.",
		Annotation: "time",
		Label:      "3. Time Initiated",
		Comment:    "required time",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateTime},
	}
	trackingNumberDef = &xscmsg.FieldDef{
		Tag:        "24.",
		Annotation: "tracking-number",
		Label:      "4. Tracking Number (OA EOC)",
	}
	requestedByDef = &xscmsg.FieldDef{
		Tag:        "25.",
		Annotation: "requested-by",
		Label:      "5. Requested by (name, agency, position, email, phone)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	preparedByDef = &xscmsg.FieldDef{
		Tag:        "26.",
		Annotation: "prepared-by",
		Label:      "6. Prepared by (name, position, email, phone)",
	}
	approvedByDef = &xscmsg.FieldDef{
		Tag:        "27.",
		Annotation: "approved-by",
		Label:      "7. Approved By (name, position, email, phone)",
	}
	qtyUnitDef = &xscmsg.FieldDef{
		Tag:        "28.",
		Annotation: "qty-unit",
		Label:      "8. Qty/Unit",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	resourceDescriptionDef = &xscmsg.FieldDef{
		Tag:        "29.",
		Annotation: "resource-description",
		Label:      "9. Resource Description (kind/type if applicable)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	resourceArrivalDef = &xscmsg.FieldDef{
		Tag:        "30.",
		Annotation: "resource-arrival",
		Label:      "10. Arrival (date/time)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	priorityDef = &xscmsg.FieldDef{
		Tag:        "31.",
		Annotation: "priority",
		Label:      "11. Priority",
		Comment:    "required: Now, High, Medium, Low",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateChoices},
		Choices:    []string{"Now", "High", "Medium", "Low"},
	}
	estdCostDef = &xscmsg.FieldDef{
		Tag:        "32.",
		Annotation: "estd-cost", // it's resource-priority in PackItForms, but that's clearly wrong
		Label:      "12. Est'd Cost",
	}
	deliverToDef = &xscmsg.FieldDef{
		Tag:        "33.",
		Annotation: "deliver-to",
		Label:      "13. Deliver to (name, agency, position, email, phone)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	deliverToLocationDef = &xscmsg.FieldDef{
		Tag:        "34.",
		Annotation: "deliver-to-location",
		Label:      "14. Location (address or lat/long, site type)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	substitutesDef = &xscmsg.FieldDef{
		Tag:        "35.",
		Annotation: "substitutes",
		Label:      "15. Substitute/Suggested Sources (name, phone, website)",
	}
	equipmentOperatorDef = &xscmsg.FieldDef{
		Tag:        "36a.",
		Annotation: "equipment-operator",
		Label:      "16. Supplemental Requirements: Equipment Operator",
		Comment:    "boolean",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	lodgingDef = &xscmsg.FieldDef{
		Tag:        "36b.",
		Annotation: "lodging",
		Label:      "16. Supplemental Requirements: Lodging",
		Comment:    "boolean",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	fuelDef = &xscmsg.FieldDef{
		Tag:        "36c.",
		Annotation: "fuel",
		Label:      "16. Supplemental Requirements: Fuel",
		Comment:    "boolean",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	fuelTypeDef = &xscmsg.FieldDef{
		Tag:        "36d.",
		Annotation: "fuel-type",
		Label:      "16. Supplemental Requirements: Fuel Type",
	}
	powerDef = &xscmsg.FieldDef{
		Tag:        "36e.",
		Annotation: "power",
		Label:      "16. Supplemental Requirements: Power",
		Comment:    "boolean",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	mealsDef = &xscmsg.FieldDef{
		Tag:        "36f.",
		Annotation: "meals",
		Label:      "16. Supplemental Requirements: Meals",
		Comment:    "boolean",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	maintenanceDef = &xscmsg.FieldDef{
		Tag:        "36g.",
		Annotation: "maintenance",
		Label:      "16. Supplemental Requirements: Maintenance",
		Comment:    "boolean",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	waterDef = &xscmsg.FieldDef{
		Tag:        "36h.",
		Annotation: "water",
		Label:      "16. Supplemental Requirements: Water",
		Comment:    "boolean",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	otherDef = &xscmsg.FieldDef{
		Tag:        "36i.",
		Annotation: "other",
		Label:      "16. Supplemental Requirements: Other",
		Comment:    "boolean",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	instructionsDef = &xscmsg.FieldDef{
		Tag:        "37.",
		Annotation: "instructions",
		Label:      "17. Special Instructions",
	}
)
