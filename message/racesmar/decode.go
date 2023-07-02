package racesmar

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/message/common"
)

func decode(subject, body string) (f *RACESMAR) {
	if idx := strings.Index(body, "form-oa-mutual-aid-request"); idx < 0 {
		return nil
	}
	form := common.DecodePIFO(body)
	switch {
	case form == nil:
		return nil
	case form.HTMLIdent == "form-oa-mutual-aid-request.html" && form.FormVersion == "1.6":
		break
	case form.HTMLIdent == "form-oa-mutual-aid-request-v2.html" && (form.FormVersion == "2.1" || form.FormVersion == "2.3"):
		break
	default:
		return nil
	}
	f = new(RACESMAR)
	f.PIFOVersion = form.PIFOVersion
	f.FormVersion = form.FormVersion
	f.StdFields.Decode(form.TaggedValues)
	f.AgencyName = form.TaggedValues["15."]
	f.EventName = form.TaggedValues["16a."]
	f.EventNumber = form.TaggedValues["16b."]
	f.Assignment = form.TaggedValues["17."]
	if f.FormVersion == "1.6" {
		f.Resources[0].Qty = form.TaggedValues["18a."]
		f.Resources[0].RolePos = form.TaggedValues["18b."]
		f.Resources[0].PreferredType = form.TaggedValues["18c."]
		f.Resources[0].MinimumType = form.TaggedValues["18d."]
	} else {
		for i := range f.Resources {
			f.Resources[i].Qty = form.TaggedValues[fmt.Sprintf("18.%da.", i+1)]
			f.Resources[i].RolePos = form.TaggedValues[fmt.Sprintf("18.%db.", i+1)]
			f.Resources[i].PreferredType = form.TaggedValues[fmt.Sprintf("18.%dc.", i+1)]
			f.Resources[i].MinimumType = form.TaggedValues[fmt.Sprintf("18.%dd.", i+1)]
			if f.FormVersion == "2.3" {
				f.Resources[i].Role = form.TaggedValues[fmt.Sprintf("18.%de.", i+1)]
				f.Resources[i].Position = form.TaggedValues[fmt.Sprintf("18.%df.", i+1)]
			}
		}
	}
	f.RequestedArrivalDates = form.TaggedValues["19a."]
	f.RequestedArrivalTimes = form.TaggedValues["19b."]
	f.NeededUntilDates = form.TaggedValues["20a."]
	f.NeededUntilTimes = form.TaggedValues["20b."]
	f.ReportingLocation = form.TaggedValues["21."]
	f.ContactOnArrival = form.TaggedValues["22."]
	f.TravelInfo = form.TaggedValues["23."]
	f.RequestedByName = form.TaggedValues["24a."]
	f.RequestedByTitle = form.TaggedValues["24b."]
	f.RequestedByContact = form.TaggedValues["24c."]
	f.ApprovedByName = form.TaggedValues["25a."]
	f.ApprovedByTitle = form.TaggedValues["25b."]
	f.ApprovedByContact = form.TaggedValues["25c."]
	f.ApprovedByDate = form.TaggedValues["26a."]
	f.ApprovedByTime = form.TaggedValues["26b."]
	return f
}
