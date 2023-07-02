package xscmsg

import (
	"strings"

	"github.com/rothskeller/packet/xscmsg/forms/pifo"
	"github.com/rothskeller/packet/xscmsg/forms/versions"
	"github.com/rothskeller/packet/xscmsg/forms/xscsubj"
)

// JurisStat form metadata:
const (
	JurisStatTag     = "JurisStat"
	jurisStatTag21   = "MuniStat"
	JurisStatHTML    = "form-oa-muni-status.html"
	JurisStatVersion = "2.2"
)

// JurisStat holds an OA jurisdiction status form.
type JurisStat struct {
	StdHeader
	ReportType                    string
	JurisdictionCode              string // added in 2.2
	Jurisdiction                  string
	EocPhone                      string
	EocFax                        string
	PriEmContactName              string
	PriEmContactPhone             string
	SecEmContactName              string
	SecEmContactPhone             string
	OfficeStatus                  string
	GovExpectedOpenDate           string
	GovExpectedOpenTime           string
	GovExpectedCloseDate          string
	GovExpectedCloseTime          string
	EocOpen                       string
	EocActivationLevel            string
	EocExpectedOpenDate           string
	EocExpectedOpenTime           string
	EocExpectedCloseDate          string
	EocExpectedCloseTime          string
	StateOfEmergency              string
	HowToSendSOE                  string // removed in 2.1
	HowSOESent                    string
	Communications                string
	CommunicationsComments        string
	Debris                        string
	DebrisComments                string
	Flooding                      string
	FloodingComments              string
	Hazmat                        string
	HazmatComments                string
	EmergencyServices             string
	EmergencyServicesComments     string
	Casualties                    string
	CasualtiesComments            string
	UtilitiesGas                  string
	UtilitiesGasComments          string
	UtilitiesElectric             string
	UtilitiesElectricComments     string
	InfrastructurePower           string
	InfrastructurePowerComments   string
	InfrastructureWater           string
	InfrastructureWaterComments   string
	InfrastructureSewer           string
	InfrastructureSewerComments   string
	SearchAndRescue               string
	SearchAndRescueComments       string
	TransportationRoads           string
	TransportationRoadsComments   string
	TransportationBridges         string
	TransportationBridgesComments string
	CivilUnrest                   string
	CivilUnrestComments           string
	AnimalIssues                  string
	AnimalIssuesComments          string
	StdFooter
}

// DecodeJurisStat decodes the supplied form if it is a JurisStat form.  It
// returns the decoded form and strings describing any non-fatal decoding
// problems.  It returns nil, nil if the form is not a JurisStat form or has an
// unknown version.
func DecodeJurisStat(form *pifo.Form) (f *JurisStat, problems []string) {
	if form.HTMLIdent != JurisStatHTML {
		return nil, nil
	}
	switch form.FormVersion {
	case "2.0", "2.1", "2.2":
		break
	default:
		return nil, nil
	}
	f = new(JurisStat)
	f.FormVersion = form.FormVersion
	f.StdHeader.PullTags(form.TaggedValues)
	f.ReportType = PullTag(form.TaggedValues, "19.")
	if f.FormVersion == "2.2" {
		f.JurisdictionCode = PullTag(form.TaggedValues, "21.")
		f.Jurisdiction = PullTag(form.TaggedValues, "22.")
	} else {
		f.Jurisdiction = PullTag(form.TaggedValues, "21.")
	}
	f.EocPhone = PullTag(form.TaggedValues, "23.")
	f.EocFax = PullTag(form.TaggedValues, "24.")
	f.PriEmContactName = PullTag(form.TaggedValues, "25.")
	f.PriEmContactPhone = PullTag(form.TaggedValues, "26.")
	f.SecEmContactName = PullTag(form.TaggedValues, "27.")
	f.SecEmContactPhone = PullTag(form.TaggedValues, "28.")
	f.OfficeStatus = PullTag(form.TaggedValues, "29.")
	f.GovExpectedOpenDate = PullTag(form.TaggedValues, "30.")
	f.GovExpectedOpenTime = PullTag(form.TaggedValues, "31.")
	f.GovExpectedCloseDate = PullTag(form.TaggedValues, "32.")
	f.GovExpectedCloseTime = PullTag(form.TaggedValues, "33.")
	f.EocOpen = PullTag(form.TaggedValues, "34.")
	f.EocActivationLevel = PullTag(form.TaggedValues, "35.")
	f.EocExpectedOpenDate = PullTag(form.TaggedValues, "36.")
	f.EocExpectedOpenTime = PullTag(form.TaggedValues, "37.")
	f.EocExpectedCloseDate = PullTag(form.TaggedValues, "38.")
	f.EocExpectedCloseTime = PullTag(form.TaggedValues, "39.")
	f.StateOfEmergency = PullTag(form.TaggedValues, "40.")
	if f.FormVersion == "2.0" {
		f.HowToSendSOE = PullTag(form.TaggedValues, "98.")
	}
	f.HowSOESent = PullTag(form.TaggedValues, "99.")
	f.Communications = PullTag(form.TaggedValues, "41.0.")
	f.CommunicationsComments = PullTag(form.TaggedValues, "41.1.")
	f.Debris = PullTag(form.TaggedValues, "42.0.")
	f.DebrisComments = PullTag(form.TaggedValues, "42.1.")
	f.Flooding = PullTag(form.TaggedValues, "43.0.")
	f.FloodingComments = PullTag(form.TaggedValues, "43.1.")
	f.Hazmat = PullTag(form.TaggedValues, "44.0.")
	f.HazmatComments = PullTag(form.TaggedValues, "44.1.")
	f.EmergencyServices = PullTag(form.TaggedValues, "45.0.")
	f.EmergencyServicesComments = PullTag(form.TaggedValues, "45.1.")
	f.Casualties = PullTag(form.TaggedValues, "46.0.")
	f.CasualtiesComments = PullTag(form.TaggedValues, "46.1.")
	f.UtilitiesGas = PullTag(form.TaggedValues, "47.0.")
	f.UtilitiesGasComments = PullTag(form.TaggedValues, "47.1.")
	f.UtilitiesElectric = PullTag(form.TaggedValues, "48.0.")
	f.UtilitiesElectricComments = PullTag(form.TaggedValues, "48.1.")
	f.InfrastructurePower = PullTag(form.TaggedValues, "49.0.")
	f.InfrastructurePowerComments = PullTag(form.TaggedValues, "49.1.")
	f.InfrastructureWater = PullTag(form.TaggedValues, "50.0.")
	f.InfrastructureWaterComments = PullTag(form.TaggedValues, "50.1.")
	f.InfrastructureSewer = PullTag(form.TaggedValues, "51.0.")
	f.InfrastructureSewerComments = PullTag(form.TaggedValues, "51.1.")
	f.SearchAndRescue = PullTag(form.TaggedValues, "52.0.")
	f.SearchAndRescueComments = PullTag(form.TaggedValues, "52.1.")
	f.TransportationRoads = PullTag(form.TaggedValues, "53.0.")
	f.TransportationRoadsComments = PullTag(form.TaggedValues, "53.1.")
	f.TransportationBridges = PullTag(form.TaggedValues, "54.0.")
	f.TransportationBridgesComments = PullTag(form.TaggedValues, "54.1.")
	f.CivilUnrest = PullTag(form.TaggedValues, "55.0.")
	f.CivilUnrestComments = PullTag(form.TaggedValues, "55.1.")
	f.AnimalIssues = PullTag(form.TaggedValues, "56.0.")
	f.AnimalIssuesComments = PullTag(form.TaggedValues, "56.1.")
	f.StdFooter.PullTags(form.TaggedValues)
	return f, LeftoverTagProblems(JurisStatTag, form.FormVersion, form.TaggedValues)
}

// Encode encodes the message contents.
func (f *JurisStat) Encode() (subject, body string) {
	var (
		sb  strings.Builder
		enc *pifo.Encoder
	)
	if versions.Older(f.FormVersion, "2.2") {
		subject = xscsubj.Encode(f.OriginMsgID, f.Handling, jurisStatTag21, f.Jurisdiction)
	} else {
		subject = xscsubj.Encode(f.OriginMsgID, f.Handling, JurisStatTag, f.Jurisdiction)
	}
	if f.FormVersion == "" {
		f.FormVersion = "2.2"
	}
	enc = pifo.NewEncoder(&sb, JurisStatHTML, f.FormVersion)
	f.StdHeader.EncodeBody(enc)
	enc.Write("19.", f.ReportType)
	if f.FormVersion == "2.2" {
		enc.Write("21.", f.JurisdictionCode)
		enc.Write("22.", f.Jurisdiction)
	} else {
		enc.Write("21.", f.Jurisdiction)
	}
	enc.Write("23.", f.EocPhone)
	enc.Write("24.", f.EocFax)
	enc.Write("25.", f.PriEmContactName)
	enc.Write("26.", f.PriEmContactPhone)
	enc.Write("27.", f.SecEmContactName)
	enc.Write("28.", f.SecEmContactPhone)
	enc.Write("29.", f.OfficeStatus)
	enc.Write("30.", f.GovExpectedOpenDate)
	enc.Write("31.", f.GovExpectedOpenTime)
	enc.Write("32.", f.GovExpectedCloseDate)
	enc.Write("33.", f.GovExpectedCloseTime)
	enc.Write("34.", f.EocOpen)
	enc.Write("35.", f.EocActivationLevel)
	enc.Write("36.", f.EocExpectedOpenDate)
	enc.Write("37.", f.EocExpectedOpenTime)
	enc.Write("38.", f.EocExpectedCloseDate)
	enc.Write("39.", f.EocExpectedCloseTime)
	enc.Write("40.", f.StateOfEmergency)
	if f.FormVersion == "2.0" {
		enc.Write("98.", f.HowToSendSOE)
	}
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
	f.StdFooter.EncodeBody(enc)
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return subject, sb.String()
}
