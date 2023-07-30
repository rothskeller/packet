package jurisstat

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for an OA jurisdiction status form.
var Type = message.Type{
	Tag:     "JurisStat",
	Name:    "OA jurisdiction status form",
	Article: "an",
	Create:  New,
	Decode:  decode,
}

// OldType is the previous type definition for an OA jurisdiction status form.
var OldType = message.Type{
	Tag:     "MuniStat",
	Name:    "OA municipal status form",
	Article: "an",
	Create:  nil,
	Decode:  decode,
}

// JurisStat holds an OA jurisdiction status form.
type JurisStat struct {
	common.StdFields
	ReportType                    string
	JurisdictionCode              string // added in 2.2
	Jurisdiction                  string
	EOCPhone                      string
	EOCFax                        string
	PriEMContactName              string
	PriEMContactPhone             string
	SecEMContactName              string
	SecEMContactPhone             string
	OfficeStatus                  string
	GovExpectedOpenDate           string
	GovExpectedOpenTime           string
	GovExpectedCloseDate          string
	GovExpectedCloseTime          string
	EOCOpen                       string
	EOCActivationLevel            string
	EOCExpectedOpenDate           string
	EOCExpectedOpenTime           string
	EOCExpectedCloseDate          string
	EOCExpectedCloseTime          string
	StateOfEmergency              string
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
	edit                          *jurisStatEdit
}

// Type returns the message type definition.
func (*JurisStat) Type() *message.Type { return &Type }
