package jurisstat

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

func decode(subject, body string) (f *JurisStat) {
	if idx := strings.Index(body, "form-oa-muni-status.html"); idx < 0 {
		return nil
	}
	form := common.DecodePIFO(body)
	if form == nil || form.HTMLIdent != "form-oa-muni-status.html" {
		return nil
	}
	switch form.FormVersion {
	case "2.0", "2.1", "2.2":
		break
	default:
		return nil
	}
	f = new(JurisStat)
	f.PIFOVersion = form.PIFOVersion
	f.FormVersion = form.FormVersion
	f.StdFields.Decode(form.TaggedValues)
	f.ReportType = form.TaggedValues["19."]
	if f.FormVersion >= "2.2" {
		f.JurisdictionCode = form.TaggedValues["21."]
		f.Jurisdiction = form.TaggedValues["22."]
	} else {
		f.Jurisdiction = form.TaggedValues["21."]
	}
	f.EOCPhone = form.TaggedValues["23."]
	f.EOCFax = form.TaggedValues["24."]
	f.PriEMContactName = form.TaggedValues["25."]
	f.PriEMContactPhone = form.TaggedValues["26."]
	f.SecEMContactName = form.TaggedValues["27."]
	f.SecEMContactPhone = form.TaggedValues["28."]
	f.OfficeStatus = form.TaggedValues["29."]
	f.GovExpectedOpenDate = form.TaggedValues["30."]
	f.GovExpectedOpenTime = form.TaggedValues["31."]
	f.GovExpectedCloseDate = form.TaggedValues["32."]
	f.GovExpectedCloseTime = form.TaggedValues["33."]
	f.EOCOpen = form.TaggedValues["34."]
	f.EOCActivationLevel = form.TaggedValues["35."]
	f.EOCExpectedOpenDate = form.TaggedValues["36."]
	f.EOCExpectedOpenTime = form.TaggedValues["37."]
	f.EOCExpectedCloseDate = form.TaggedValues["38."]
	f.EOCExpectedCloseTime = form.TaggedValues["39."]
	f.StateOfEmergency = form.TaggedValues["40."]
	f.HowSOESent = form.TaggedValues["99."]
	f.Communications = form.TaggedValues["41.0."]
	f.CommunicationsComments = form.TaggedValues["41.1."]
	f.Debris = form.TaggedValues["42.0."]
	f.DebrisComments = form.TaggedValues["42.1."]
	f.Flooding = form.TaggedValues["43.0."]
	f.FloodingComments = form.TaggedValues["43.1."]
	f.Hazmat = form.TaggedValues["44.0."]
	f.HazmatComments = form.TaggedValues["44.1."]
	f.EmergencyServices = form.TaggedValues["45.0."]
	f.EmergencyServicesComments = form.TaggedValues["45.1."]
	f.Casualties = form.TaggedValues["46.0."]
	f.CasualtiesComments = form.TaggedValues["46.1."]
	f.UtilitiesGas = form.TaggedValues["47.0."]
	f.UtilitiesGasComments = form.TaggedValues["47.1."]
	f.UtilitiesElectric = form.TaggedValues["48.0."]
	f.UtilitiesElectricComments = form.TaggedValues["48.1."]
	f.InfrastructurePower = form.TaggedValues["49.0."]
	f.InfrastructurePowerComments = form.TaggedValues["49.1."]
	f.InfrastructureWater = form.TaggedValues["50.0."]
	f.InfrastructureWaterComments = form.TaggedValues["50.1."]
	f.InfrastructureSewer = form.TaggedValues["51.0."]
	f.InfrastructureSewerComments = form.TaggedValues["51.1."]
	f.SearchAndRescue = form.TaggedValues["52.0."]
	f.SearchAndRescueComments = form.TaggedValues["52.1."]
	f.TransportationRoads = form.TaggedValues["53.0."]
	f.TransportationRoadsComments = form.TaggedValues["53.1."]
	f.TransportationBridges = form.TaggedValues["54.0."]
	f.TransportationBridgesComments = form.TaggedValues["54.1."]
	f.CivilUnrest = form.TaggedValues["55.0."]
	f.CivilUnrestComments = form.TaggedValues["55.1."]
	f.AnimalIssues = form.TaggedValues["56.0."]
	f.AnimalIssuesComments = form.TaggedValues["56.1."]
	return f
}
