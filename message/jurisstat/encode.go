package jurisstat

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

// EncodeSubject encodes the message subject.
func (f *JurisStat) EncodeSubject() string {
	if f.FormVersion == "2.2" {
		return common.EncodeSubject(f.OriginMsgID, f.Handling, Type.Tag, f.Jurisdiction)
	}
	return common.EncodeSubject(f.OriginMsgID, f.Handling, "MuniStat", f.Jurisdiction)
}

// EncodeBody encodes the message body.
func (f *JurisStat) EncodeBody() string {
	var (
		sb  strings.Builder
		enc *common.PIFOEncoder
	)
	if f.FormVersion == "" {
		f.FormVersion = "2.2"
	}
	enc = common.NewPIFOEncoder(&sb, "form-oa-muni-status.html", f.FormVersion)
	f.StdFields.EncodeHeader(enc)
	enc.Write("19.", f.ReportType)
	if f.FormVersion == "2.2" {
		enc.Write("21.", f.JurisdictionCode)
		enc.Write("22.", f.Jurisdiction)
	} else {
		enc.Write("21.", f.Jurisdiction)
	}
	enc.Write("23.", f.EOCPhone)
	enc.Write("24.", f.EOCFax)
	enc.Write("25.", f.PriEMContactName)
	enc.Write("26.", f.PriEMContactPhone)
	enc.Write("27.", f.SecEMContactName)
	enc.Write("28.", f.SecEMContactPhone)
	enc.Write("29.", f.OfficeStatus)
	enc.Write("30.", f.GovExpectedOpenDate)
	enc.Write("31.", f.GovExpectedOpenTime)
	enc.Write("32.", f.GovExpectedCloseDate)
	enc.Write("33.", f.GovExpectedCloseTime)
	enc.Write("34.", f.EOCOpen)
	enc.Write("35.", f.EOCActivationLevel)
	enc.Write("36.", f.EOCExpectedOpenDate)
	enc.Write("37.", f.EOCExpectedOpenTime)
	enc.Write("38.", f.EOCExpectedCloseDate)
	enc.Write("39.", f.EOCExpectedCloseTime)
	enc.Write("40.", f.StateOfEmergency)
	enc.Write("99.", f.HowSOESent)
	enc.Write("41.0.", f.Communications)
	enc.Write("41.1.", f.CommunicationsComments)
	enc.Write("42.0.", f.Debris)
	enc.Write("42.1.", f.DebrisComments)
	enc.Write("43.0.", f.Flooding)
	enc.Write("43.1.", f.FloodingComments)
	enc.Write("44.0.", f.Hazmat)
	enc.Write("44.1.", f.HazmatComments)
	enc.Write("45.0.", f.EmergencyServices)
	enc.Write("45.1.", f.EmergencyServicesComments)
	enc.Write("46.0.", f.Casualties)
	enc.Write("46.1.", f.CasualtiesComments)
	enc.Write("47.0.", f.UtilitiesGas)
	enc.Write("47.1.", f.UtilitiesGasComments)
	enc.Write("48.0.", f.UtilitiesElectric)
	enc.Write("48.1.", f.UtilitiesElectricComments)
	enc.Write("49.0.", f.InfrastructurePower)
	enc.Write("49.1.", f.InfrastructurePowerComments)
	enc.Write("50.0.", f.InfrastructureWater)
	enc.Write("50.1.", f.InfrastructureWaterComments)
	enc.Write("51.0.", f.InfrastructureSewer)
	enc.Write("51.1.", f.InfrastructureSewerComments)
	enc.Write("52.0.", f.SearchAndRescue)
	enc.Write("52.1.", f.SearchAndRescueComments)
	enc.Write("53.0.", f.TransportationRoads)
	enc.Write("53.1.", f.TransportationRoadsComments)
	enc.Write("54.0.", f.TransportationBridges)
	enc.Write("54.1.", f.TransportationBridgesComments)
	enc.Write("55.0.", f.CivilUnrest)
	enc.Write("55.1.", f.CivilUnrestComments)
	enc.Write("56.0.", f.AnimalIssues)
	enc.Write("56.1.", f.AnimalIssuesComments)
	f.StdFields.EncodeFooter(enc)
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return sb.String()
}
