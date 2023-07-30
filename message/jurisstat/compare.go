package jurisstat

import (
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// Compare compares two messages.  It returns a score indicating how
// closely they match, and the detailed comparisons of each field in the
// message.  The comparison is not symmetric:  the receiver of the call
// is the "expected" message and the argument is the "actual" message.
func (exp *JurisStat) Compare(actual message.Message) (score int, outOf int, fields []*message.CompareField) {
	var (
		act *JurisStat
		ok  bool
	)
	if act, ok = actual.(*JurisStat); !ok {
		return 0, 1, []*message.CompareField{{
			Label: "Message Type",
			Score: 0, OutOf: 1,
			Expected:     Type.Name,
			ExpectedMask: strings.Repeat("*", len(Type.Name)),
			Actual:       actual.Type().Name,
			ActualMask:   strings.Repeat("*", len(actual.Type().Name)),
		}}
	}
	fields = append(exp.StdFields.Compare(&act.StdFields),
		common.CompareExact("Report Type", exp.ReportType, act.ReportType),
		common.CompareText("Jurisdiction", exp.Jurisdiction, act.Jurisdiction),
		common.ComparePhoneNumber("EOC Phone", exp.EOCPhone, act.EOCPhone),
		common.ComparePhoneNumber("EOC Fax", exp.EOCFax, act.EOCFax),
		common.CompareText("Primary EM Contact Name", exp.PriEMContactName, act.PriEMContactName),
		common.ComparePhoneNumber("Primary EM Contact Phone", exp.PriEMContactPhone, act.PriEMContactPhone),
		common.CompareText("Secondary EM Contact Name", exp.SecEMContactName, act.SecEMContactName),
		common.ComparePhoneNumber("Secondary EM Contact Phone", exp.SecEMContactPhone, act.SecEMContactPhone),
		common.CompareExact("Govt. Office Status", exp.OfficeStatus, act.OfficeStatus),
		common.CompareDate("Govt. Office Expected to Open Date", exp.GovExpectedOpenDate, act.GovExpectedOpenDate),
		common.CompareTime("Govt. Office Expected to Open Time", exp.GovExpectedOpenTime, act.GovExpectedOpenTime),
		common.CompareDate("Govt. Office Expected to Close Date", exp.GovExpectedCloseDate, act.GovExpectedCloseDate),
		common.CompareTime("Govt. Office Expected to Close Time", exp.GovExpectedCloseTime, act.GovExpectedCloseTime),
		common.CompareExact("EOC Open", exp.EOCOpen, act.EOCOpen),
		common.CompareExact("EOC Activation Level", exp.EOCActivationLevel, act.EOCActivationLevel),
		common.CompareDate("EOC Expected to Open Date", exp.EOCExpectedOpenDate, act.EOCExpectedOpenDate),
		common.CompareTime("EOC Expected to Open Time", exp.EOCExpectedOpenTime, act.EOCExpectedOpenTime),
		common.CompareDate("EOC Expected to Close Date", exp.EOCExpectedCloseDate, act.EOCExpectedCloseDate),
		common.CompareTime("EOC Expected to Close Time", exp.EOCExpectedCloseTime, act.EOCExpectedCloseTime),
		common.CompareExact("State of Emergency", exp.StateOfEmergency, act.StateOfEmergency),
		common.CompareText("How SOE Sent", exp.HowSOESent, act.HowSOESent),
		common.CompareExact("Communications", exp.Communications, act.Communications),
		common.CompareText("Communications: Comments", exp.CommunicationsComments, act.CommunicationsComments),
		common.CompareExact("Debris", exp.Debris, act.Debris),
		common.CompareText("Debris: Comments", exp.DebrisComments, act.DebrisComments),
		common.CompareExact("Flooding", exp.Flooding, act.Flooding),
		common.CompareText("Flooding: Comments", exp.FloodingComments, act.FloodingComments),
		common.CompareExact("Hazmat", exp.Hazmat, act.Hazmat),
		common.CompareText("Hazmat: Comments", exp.HazmatComments, act.HazmatComments),
		common.CompareExact("Emergency Services", exp.EmergencyServices, act.EmergencyServices),
		common.CompareText("Emergency Services: Comments", exp.EmergencyServicesComments, act.EmergencyServicesComments),
		common.CompareExact("Casualties", exp.Casualties, act.Casualties),
		common.CompareText("Casualties: Comments", exp.CasualtiesComments, act.CasualtiesComments),
		common.CompareExact("Utilities (Gas)", exp.UtilitiesGas, act.UtilitiesGas),
		common.CompareText("Utilities (Gas): Comments", exp.UtilitiesGasComments, act.UtilitiesGasComments),
		common.CompareExact("Utilities (Electric)", exp.UtilitiesElectric, act.UtilitiesElectric),
		common.CompareText("Utilities (Electric): Comments", exp.UtilitiesElectricComments, act.UtilitiesElectricComments),
		common.CompareExact("Infrastructure (Power)", exp.InfrastructurePower, act.InfrastructurePower),
		common.CompareText("Infrastructure (Power): Comments", exp.InfrastructurePowerComments, act.InfrastructurePowerComments),
		common.CompareExact("Infrastructure (Water)", exp.InfrastructureWater, act.InfrastructureWater),
		common.CompareText("Infrastructure (Water): Comments", exp.InfrastructureWaterComments, act.InfrastructureWaterComments),
		common.CompareExact("Infrastructure (Sewer)", exp.InfrastructureSewer, act.InfrastructureSewer),
		common.CompareText("Infrastructure (Sewer): Comments", exp.InfrastructureSewerComments, act.InfrastructureSewerComments),
		common.CompareExact("Search and Rescue", exp.SearchAndRescue, act.SearchAndRescue),
		common.CompareText("Search and Rescue: Comments", exp.SearchAndRescueComments, act.SearchAndRescueComments),
		common.CompareExact("Transportation (Roads)", exp.TransportationRoads, act.TransportationRoads),
		common.CompareText("Transportation (Roads): Comments", exp.TransportationRoadsComments, act.TransportationRoadsComments),
		common.CompareExact("Transportation (Bridges)", exp.TransportationBridges, act.TransportationBridges),
		common.CompareText("Transportation (Bridges): Comments", exp.TransportationBridgesComments, act.TransportationBridgesComments),
		common.CompareExact("Civil Unrest", exp.CivilUnrest, act.CivilUnrest),
		common.CompareText("Civil Unrest: Comments", exp.CivilUnrestComments, act.CivilUnrestComments),
		common.CompareExact("Animal Issues", exp.AnimalIssues, act.AnimalIssues),
		common.CompareText("Animal Issues: Comments", exp.AnimalIssuesComments, act.AnimalIssuesComments),
	)
	return common.ConsolidateCompareFields(fields)
}
