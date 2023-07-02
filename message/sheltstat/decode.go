package sheltstat

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

func decode(subject, body string) (f *SheltStat) {
	if idx := strings.Index(body, "form-oa-shelter-status.html"); idx < 0 {
		return nil
	}
	form := common.DecodePIFO(body)
	if form == nil || form.HTMLIdent != "form-oa-shelter-status.html" {
		return nil
	}
	switch form.FormVersion {
	case "2.0", "2.1", "2.2":
		break
	default:
		return nil
	}
	f = new(SheltStat)
	f.PIFOVersion = form.PIFOVersion
	f.FormVersion = form.FormVersion
	f.StdFields.Decode(form.TaggedValues)
	f.ReportType = form.TaggedValues["19."]
	f.ShelterName = form.TaggedValues["32."]
	f.ShelterType = form.TaggedValues["30."]
	f.ShelterStatus = form.TaggedValues["31."]
	f.ShelterAddress = form.TaggedValues["33a."]
	if f.FormVersion == "2.2" {
		f.ShelterCityCode = form.TaggedValues["33b."]
		f.ShelterCity = form.TaggedValues["34b."]
	} else {
		f.ShelterCity = form.TaggedValues["33b."]
	}
	f.ShelterState = form.TaggedValues["33c."]
	f.ShelterZip = form.TaggedValues["33d."]
	f.Latitude = form.TaggedValues["37a."]
	f.Longitude = form.TaggedValues["37b."]
	f.Capacity = form.TaggedValues["40a."]
	f.Occupancy = form.TaggedValues["40b."]
	f.MealsServed = form.TaggedValues["41."]
	f.NSSNumber = form.TaggedValues["42."]
	f.PetFriendly = form.TaggedValues["43a."]
	f.BasicSafetyInspection = form.TaggedValues["43b."]
	f.ATC20Inspection = form.TaggedValues["43c."]
	f.AvailableServices = form.TaggedValues["44."]
	f.MOU = form.TaggedValues["45."]
	f.FloorPlan = form.TaggedValues["46."]
	if f.FormVersion == "2.2" {
		f.ManagedByCode = form.TaggedValues["50a."]
		f.ManagedBy = form.TaggedValues["49a."]
	} else {
		f.ManagedBy = form.TaggedValues["50a."]
	}
	f.ManagedByDetail = form.TaggedValues["50b."]
	f.PrimaryContact = form.TaggedValues["51a."]
	f.PrimaryPhone = form.TaggedValues["51b."]
	f.SecondaryContact = form.TaggedValues["52a."]
	f.SecondaryPhone = form.TaggedValues["52b."]
	f.TacticalCallSign = form.TaggedValues["60."]
	f.RepeaterCallSign = form.TaggedValues["61."]
	f.RepeaterInput = form.TaggedValues["62a."]
	f.RepeaterInputTone = form.TaggedValues["62b."]
	f.RepeaterOutput = form.TaggedValues["63a."]
	f.RepeaterOutputTone = form.TaggedValues["63b."]
	f.RepeaterOffset = form.TaggedValues["62c."]
	f.Comments = form.TaggedValues["70."]
	f.RemoveFromList = form.TaggedValues["71."]
	return f
}
