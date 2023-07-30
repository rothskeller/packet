package jurisstat

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// RenderTable renders the message as a set of field label / field value
// pairs, intended for read-only display to a human.
func (f *JurisStat) RenderTable() (lvs []message.LabelValue) {
	lvs = append(f.StdFields.RenderTable1(), []message.LabelValue{
		{Label: "Report Type", Value: f.ReportType},
		{Label: "Jurisdiction", Value: f.Jurisdiction},
		{Label: "EOC Phone", Value: f.EOCPhone},
		{Label: "EOC Fax", Value: f.EOCFax},
		{Label: "Primary EM Contact Name", Value: f.PriEMContactName},
		{Label: "Primary EM Contact Phone", Value: f.PriEMContactPhone},
		{Label: "Secondary EM Contact Name", Value: f.SecEMContactName},
		{Label: "Secondary EM Contact Phone", Value: f.SecEMContactPhone},
		{Label: "Govt. Office Status", Value: f.OfficeStatus},
		{Label: "Govt. Office Expected to Open", Value: common.SmartJoin(f.GovExpectedOpenDate, f.GovExpectedOpenTime, " ")},
		{Label: "Govt. Office Expected to Close", Value: common.SmartJoin(f.GovExpectedCloseDate, f.GovExpectedCloseTime, " ")},
		{Label: "EOC Open", Value: f.EOCOpen},
		{Label: "EOC Activation Level", Value: f.EOCActivationLevel},
		{Label: "EOC Expected to Open", Value: common.SmartJoin(f.EOCExpectedOpenDate, f.EOCExpectedOpenTime, " ")},
		{Label: "EOC Expected to Close", Value: common.SmartJoin(f.EOCExpectedCloseDate, f.EOCExpectedCloseTime, " ")},
		{Label: "State Of Emergency", Value: f.StateOfEmergency},
		{Label: "How SOE Sent", Value: f.HowSOESent},
		{Label: "Communications", Value: f.Communications},
		{Label: "Communications: Comments", Value: f.CommunicationsComments},
		{Label: "Debris", Value: f.Debris},
		{Label: "Debris: Comments", Value: f.DebrisComments},
		{Label: "Flooding", Value: f.Flooding},
		{Label: "Flooding: Comments", Value: f.FloodingComments},
		{Label: "Hazmat", Value: f.Hazmat},
		{Label: "Hazmat: Comments", Value: f.HazmatComments},
		{Label: "Emergency Services", Value: f.EmergencyServices},
		{Label: "Emergency Services: Comments", Value: f.EmergencyServicesComments},
		{Label: "Casualties", Value: f.Casualties},
		{Label: "Casualties: Comments", Value: f.CasualtiesComments},
		{Label: "Utilities (Gas)", Value: f.UtilitiesGas},
		{Label: "Utilities (Gas): Comments", Value: f.UtilitiesGasComments},
		{Label: "Utilities (Electric)", Value: f.UtilitiesElectric},
		{Label: "Utilities (Electric): Comments", Value: f.UtilitiesElectricComments},
		{Label: "Infrastructure (Power)", Value: f.InfrastructurePower},
		{Label: "Infrastructure (Power): Comments", Value: f.InfrastructurePowerComments},
		{Label: "Infrastructure (Water)", Value: f.InfrastructureWater},
		{Label: "Infrastructure (Water): Comments", Value: f.InfrastructureWaterComments},
		{Label: "Infrastructure (Sewer)", Value: f.InfrastructureSewer},
		{Label: "Infrastructure (Sewer): Comments", Value: f.InfrastructureSewerComments},
		{Label: "Search and Rescue", Value: f.SearchAndRescue},
		{Label: "Search and Rescue: Comments", Value: f.SearchAndRescueComments},
		{Label: "Transportation (Roads)", Value: f.TransportationRoads},
		{Label: "Transportation (Roads): Comments", Value: f.TransportationRoadsComments},
		{Label: "Transportation (Bridges)", Value: f.TransportationBridges},
		{Label: "Transportation (Bridges): Comments", Value: f.TransportationBridgesComments},
		{Label: "Civil Unrest", Value: f.CivilUnrest},
		{Label: "Civil Unrest: Comments", Value: f.CivilUnrestComments},
		{Label: "Animal Issues", Value: f.AnimalIssues},
		{Label: "Animal Issues: Comments", Value: f.AnimalIssuesComments},
	}...)
	return append(lvs, f.StdFields.RenderTable2()...)
}
