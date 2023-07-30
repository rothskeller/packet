package jurisstat

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
	"golang.org/x/exp/slices"
)

type jurisStatEdit struct {
	common.StdFieldsEdit
	ReportType                    message.EditField
	Jurisdiction                  message.EditField
	EOCPhone                      message.EditField
	EOCFax                        message.EditField
	PriEMContactName              message.EditField
	PriEMContactPhone             message.EditField
	SecEMContactName              message.EditField
	SecEMContactPhone             message.EditField
	OfficeStatus                  message.EditField
	GovExpectedOpen               message.EditField
	GovExpectedClose              message.EditField
	EOCOpen                       message.EditField
	EOCActivationLevel            message.EditField
	EOCExpectedOpen               message.EditField
	EOCExpectedClose              message.EditField
	StateOfEmergency              message.EditField
	HowSOESent                    message.EditField
	Communications                message.EditField
	CommunicationsComments        message.EditField
	Debris                        message.EditField
	DebrisComments                message.EditField
	Flooding                      message.EditField
	FloodingComments              message.EditField
	Hazmat                        message.EditField
	HazmatComments                message.EditField
	EmergencyServices             message.EditField
	EmergencyServicesComments     message.EditField
	Casualties                    message.EditField
	CasualtiesComments            message.EditField
	UtilitiesGas                  message.EditField
	UtilitiesGasComments          message.EditField
	UtilitiesElectric             message.EditField
	UtilitiesElectricComments     message.EditField
	InfrastructurePower           message.EditField
	InfrastructurePowerComments   message.EditField
	InfrastructureWater           message.EditField
	InfrastructureWaterComments   message.EditField
	InfrastructureSewer           message.EditField
	InfrastructureSewerComments   message.EditField
	SearchAndRescue               message.EditField
	SearchAndRescueComments       message.EditField
	TransportationRoads           message.EditField
	TransportationRoadsComments   message.EditField
	TransportationBridges         message.EditField
	TransportationBridgesComments message.EditField
	CivilUnrest                   message.EditField
	CivilUnrestComments           message.EditField
	AnimalIssues                  message.EditField
	AnimalIssuesComments          message.EditField
	fields                        []*message.EditField
}

// EditFields returns the set of editable fields of the message.
func (f *JurisStat) EditFields() []*message.EditField {
	if f.edit == nil {
		f.edit = &jurisStatEdit{
			StdFieldsEdit: common.StdFieldsEditTemplate,
			ReportType: message.EditField{
				Label:   "Report Type",
				Width:   8,
				Choices: []string{"Update", "Complete"},
				Help:    `This indicates whether the form should "Update" the previous status report for the jurisdiction, or whether it is a "Complete" replacement of the previous report.  This field is required.`,
			},
			Jurisdiction: message.EditField{
				Label:   "Jurisdiction",
				Width:   42,
				Choices: []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Help:    `This is the name of the jurisdiction being described by the form.  It is required.`,
			},
			EOCPhone: message.EditField{
				Label: "EOC Phone",
				Width: 34,
				Help:  `This is the phone number of the jurisdiction's Emergency Operations Center (EOC).  It is required when "Report Type" is "Complete".`,
			},
			EOCFax: message.EditField{
				Label: "EOC Fax",
				Width: 37,
				Help:  `This is the fax number of the jurisdiction's Emergency Operations Center (EOC).`,
			},
			PriEMContactName: message.EditField{
				Label: "Primary EM Contact Name",
				Width: 27,
				Help:  `This is the name of the primary emergency manager of the jurisdiction.  It is required when "Report Type" is "Complete".`,
			},
			PriEMContactPhone: message.EditField{
				Label: "Primary EM Contact Phone",
				Width: 26,
				Help:  `This is the phone number of the primary emergency manager of the jurisdiction.  It is required when "Report Type" is "Complete".`,
			},
			SecEMContactName: message.EditField{
				Label: "Secondary EM Contact Name",
				Width: 26,
				Help:  `This is the name of the secondary emergency manager of the jurisdiction.`,
			},
			SecEMContactPhone: message.EditField{
				Label: "Secondary EM Contact Phone",
				Width: 26,
				Help:  `This is the phone number of the secondary emergency manager of the jurisdiction.`,
			},
			OfficeStatus: message.EditField{
				Label:   "Govt. Office Status",
				Width:   7,
				Choices: []string{"Unknown", "Open", "Closed"},
				Help:    `This indicates whether the jurisdiction's regular business offices are open.  It is required when "Report Type" is "Complete".`,
			},
			GovExpectedOpen: message.EditField{
				Label: "Govt. Office Expected to Open",
				Width: 16,
				Hint:  "MM/DD/YYYY HH:MM",
				Help:  `This is the date and time when the jurisdiction's regular business offices are expected to open, in MM/DD/YYYY HH:MM format (24-hour clock).`,
			},
			GovExpectedClose: message.EditField{
				Label: "Govt. Office Expected to Close",
				Width: 16,
				Hint:  "MM/DD/YYYY HH:MM",
				Help:  `This is the date and time when the jurisdiction's regular business offices are expected to close, in MM/DD/YYYY HH:MM format (24-hour clock).`,
			},
			EOCOpen: message.EditField{
				Label:   "EOC Open",
				Width:   7,
				Choices: []string{"Unknown", "Yes", "No"},
				Help:    `This indicates whether the jurisdiction's Emergency Operations Center (EOC) is open.  It is required when "Report Type" is "Complete".`,
			},
			EOCActivationLevel: message.EditField{
				Label:   "EOC Activation Level",
				Width:   12,
				Choices: []string{"Normal", "Duty Officer", "Monitor", "Partial", "Full"},
				Help:    `This indicates the activation level of the jurisdiction's Emergency Operations Center (EOC).  It is required when "Report Type" is "Complete".`,
			},
			EOCExpectedOpen: message.EditField{
				Label: "EOC Expected to Open",
				Width: 16,
				Hint:  "MM/DD/YYYY HH:MM",
				Help:  `This is the date and time when the jurisdiction's Emergency Operations Center (EOC) is expected to open, in MM/DD/YYYY HH:MM format (24-hour clock).`,
			},
			EOCExpectedClose: message.EditField{
				Label: "EOC Expected to Close",
				Width: 16,
				Hint:  "MM/DD/YYYY HH:MM",
				Help:  `This is the date and time when the jurisdiction's Emergency Operations Center (EOC) is expected to close, in MM/DD/YYYY HH:MM format (24-hour clock).`,
			},
			StateOfEmergency: message.EditField{
				Label:   "State of Emergency",
				Width:   7,
				Choices: []string{"Unknown", "Yes", "No"},
				Help:    `This indicates whether the jurisdiction has a declared state of emergency.  It is required when "Report Type" is "Complete".`,
			},
			HowSOESent: message.EditField{
				Label: "How SOE Sent",
				Width: 58,
				Help:  `This describes where and how the jurisdiction's "state of emergency" declaration was delivered.`,
			},
			Communications: message.EditField{
				Label:   "Communications",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to communications.`,
			},
			CommunicationsComments: message.EditField{
				Label: "Communications: Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to communications.`,
			},
			Debris: message.EditField{
				Label:   "Debris",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to debris.`,
			},
			DebrisComments: message.EditField{
				Label: "Debris: Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to debris.`,
			},
			Flooding: message.EditField{
				Label:   "Flooding",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to flooding.`,
			},
			FloodingComments: message.EditField{
				Label: "Flooding: Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to flooding.`,
			},
			Hazmat: message.EditField{
				Label:   "Hazmat",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to hazmat.`,
			},
			HazmatComments: message.EditField{
				Label: "Hazmat: Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to hazmat.`,
			},
			EmergencyServices: message.EditField{
				Label:   "Emergency Services",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to emergency services.`,
			},
			EmergencyServicesComments: message.EditField{
				Label: "Emergency Services: Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to emergency services.`,
			},
			Casualties: message.EditField{
				Label:   "Casualties",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to casualties.`,
			},
			CasualtiesComments: message.EditField{
				Label: "Casualties: Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to casualties.`,
			},
			UtilitiesGas: message.EditField{
				Label:   "Utilities (Gas)",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to utilities (gas).`,
			},
			UtilitiesGasComments: message.EditField{
				Label: "Utilities (Gas): Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to utilities (gas).`,
			},
			UtilitiesElectric: message.EditField{
				Label:   "Utilities (Electric)",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to utilities (electric).`,
			},
			UtilitiesElectricComments: message.EditField{
				Label: "Utilities (Electric): Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to utilities (electric).`,
			},
			InfrastructurePower: message.EditField{
				Label:   "Infrastructure (Power)",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to infrastructure (power).`,
			},
			InfrastructurePowerComments: message.EditField{
				Label: "Infrastructure (Power): Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to infrastructure (power).`,
			},
			InfrastructureWater: message.EditField{
				Label:   "Infrastructure (Water)",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to infrastructure (water).`,
			},
			InfrastructureWaterComments: message.EditField{
				Label: "Infrastructure (Water): Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to infrastructure (water).`,
			},
			InfrastructureSewer: message.EditField{
				Label:   "Infrastructure (Sewer)",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to infrastructure (sewer).`,
			},
			InfrastructureSewerComments: message.EditField{
				Label: "Infrastructure (Sewer): Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to infrastructure (sewer).`,
			},
			SearchAndRescue: message.EditField{
				Label:   "Search and Rescue",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to search and rescue.`,
			},
			SearchAndRescueComments: message.EditField{
				Label: "Search and Rescue: Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to search and rescue.`,
			},
			TransportationRoads: message.EditField{
				Label:   "Transportation (Roads)",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to transportation (roads).`,
			},
			TransportationRoadsComments: message.EditField{
				Label: "Transportation (Roads): Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to transportation (roads).`,
			},
			TransportationBridges: message.EditField{
				Label:   "Transportation (Bridges)",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to transportation (bridges).`,
			},
			TransportationBridgesComments: message.EditField{
				Label: "Transportation (Bridges): Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to transportation (bridges).`,
			},
			CivilUnrest: message.EditField{
				Label:   "Civil Unrest",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to civil unrest.`,
			},
			CivilUnrestComments: message.EditField{
				Label: "Civil Unrest: Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to civil unrest.`,
			},
			AnimalIssues: message.EditField{
				Label:   "Animal Issues",
				Width:   9,
				Choices: []string{"Normal", "Unknown", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
				Help:    `This describes the current situation status with respect to animal issues.`,
			},
			AnimalIssuesComments: message.EditField{
				Label: "Animal Issues: Comments",
				Width: 60, Multiline: true,
				Help: `These are comments on the current situation status with respect to animal issues.`,
			},
		}
		// Set the field list slice.
		f.edit.fields = append(f.edit.StdFieldsEdit.EditFields1(),
			&f.edit.ReportType,
			&f.edit.Jurisdiction,
			&f.edit.EOCPhone,
			&f.edit.EOCFax,
			&f.edit.PriEMContactName,
			&f.edit.PriEMContactPhone,
			&f.edit.SecEMContactName,
			&f.edit.SecEMContactPhone,
			&f.edit.OfficeStatus,
			&f.edit.GovExpectedOpen,
			&f.edit.GovExpectedClose,
			&f.edit.EOCOpen,
			&f.edit.EOCActivationLevel,
			&f.edit.EOCExpectedOpen,
			&f.edit.EOCExpectedClose,
			&f.edit.StateOfEmergency,
			&f.edit.HowSOESent,
			&f.edit.Communications,
			&f.edit.CommunicationsComments,
			&f.edit.Debris,
			&f.edit.DebrisComments,
			&f.edit.Flooding,
			&f.edit.FloodingComments,
			&f.edit.Hazmat,
			&f.edit.HazmatComments,
			&f.edit.EmergencyServices,
			&f.edit.EmergencyServicesComments,
			&f.edit.Casualties,
			&f.edit.CasualtiesComments,
			&f.edit.UtilitiesGas,
			&f.edit.UtilitiesGasComments,
			&f.edit.UtilitiesElectric,
			&f.edit.UtilitiesElectricComments,
			&f.edit.InfrastructurePower,
			&f.edit.InfrastructurePowerComments,
			&f.edit.InfrastructureWater,
			&f.edit.InfrastructureWaterComments,
			&f.edit.InfrastructureSewer,
			&f.edit.InfrastructureSewerComments,
			&f.edit.SearchAndRescue,
			&f.edit.SearchAndRescueComments,
			&f.edit.TransportationRoads,
			&f.edit.TransportationRoadsComments,
			&f.edit.TransportationBridges,
			&f.edit.TransportationBridgesComments,
			&f.edit.CivilUnrest,
			&f.edit.CivilUnrestComments,
			&f.edit.AnimalIssues,
			&f.edit.AnimalIssuesComments,
		)
		f.edit.fields = append(f.edit.fields, f.edit.StdFieldsEdit.EditFields2()...)
		f.toEdit()
		f.validate()
	}
	return f.edit.fields
}

// ApplyEdits applies the revised Values in the EditFields to the
// message.
func (f *JurisStat) ApplyEdits() {
	f.fromEdit()
	f.toEdit()
	f.validate()
}

func (f *JurisStat) fromEdit() {
	f.StdFields.FromEdit(&f.edit.StdFieldsEdit)
	f.ReportType = common.ExpandRestricted(&f.edit.ReportType)
	f.Jurisdiction = common.ExpandRestricted(&f.edit.Jurisdiction)
	if f.Jurisdiction == "" || slices.Contains(f.edit.Jurisdiction.Choices, f.Jurisdiction) {
		f.JurisdictionCode = f.Jurisdiction
	} else {
		f.JurisdictionCode = "Unincorporated"
	}
	f.EOCPhone = common.CleanPhoneNumber(f.edit.EOCPhone.Value)
	f.EOCFax = common.CleanPhoneNumber(f.edit.EOCFax.Value)
	f.PriEMContactName = strings.TrimSpace(f.edit.PriEMContactName.Value)
	f.PriEMContactPhone = common.CleanPhoneNumber(f.edit.PriEMContactPhone.Value)
	f.SecEMContactName = strings.TrimSpace(f.edit.SecEMContactName.Value)
	f.SecEMContactPhone = common.CleanPhoneNumber(f.edit.SecEMContactPhone.Value)
	f.OfficeStatus = common.ExpandRestricted(&f.edit.OfficeStatus)
	f.GovExpectedOpenDate, f.GovExpectedOpenTime = common.CleanDateTime(f.edit.GovExpectedOpen.Value)
	f.GovExpectedCloseDate, f.GovExpectedCloseTime = common.CleanDateTime(f.edit.GovExpectedClose.Value)
	f.EOCOpen = common.ExpandRestricted(&f.edit.EOCOpen)
	f.EOCActivationLevel = common.ExpandRestricted(&f.edit.EOCActivationLevel)
	f.EOCExpectedOpenDate, f.EOCExpectedOpenTime = common.CleanDateTime(f.edit.EOCExpectedOpen.Value)
	f.EOCExpectedCloseDate, f.EOCExpectedCloseTime = common.CleanDateTime(f.edit.EOCExpectedClose.Value)
	f.StateOfEmergency = common.ExpandRestricted(&f.edit.StateOfEmergency)
	f.HowSOESent = strings.TrimSpace(f.edit.HowSOESent.Value)
	f.Communications = common.ExpandRestricted(&f.edit.Communications)
	f.CommunicationsComments = strings.TrimSpace(f.edit.CommunicationsComments.Value)
	f.Debris = common.ExpandRestricted(&f.edit.Debris)
	f.DebrisComments = strings.TrimSpace(f.edit.DebrisComments.Value)
	f.Flooding = common.ExpandRestricted(&f.edit.Flooding)
	f.FloodingComments = strings.TrimSpace(f.edit.FloodingComments.Value)
	f.Hazmat = common.ExpandRestricted(&f.edit.Hazmat)
	f.HazmatComments = strings.TrimSpace(f.edit.HazmatComments.Value)
	f.EmergencyServices = common.ExpandRestricted(&f.edit.EmergencyServices)
	f.EmergencyServicesComments = strings.TrimSpace(f.edit.EmergencyServicesComments.Value)
	f.Casualties = common.ExpandRestricted(&f.edit.Casualties)
	f.CasualtiesComments = strings.TrimSpace(f.edit.CasualtiesComments.Value)
	f.UtilitiesGas = common.ExpandRestricted(&f.edit.UtilitiesGas)
	f.UtilitiesGasComments = strings.TrimSpace(f.edit.UtilitiesGasComments.Value)
	f.UtilitiesElectric = common.ExpandRestricted(&f.edit.UtilitiesElectric)
	f.UtilitiesElectricComments = strings.TrimSpace(f.edit.UtilitiesElectricComments.Value)
	f.InfrastructurePower = common.ExpandRestricted(&f.edit.InfrastructurePower)
	f.InfrastructurePowerComments = strings.TrimSpace(f.edit.InfrastructurePowerComments.Value)
	f.InfrastructureWater = common.ExpandRestricted(&f.edit.InfrastructureWater)
	f.InfrastructureWaterComments = strings.TrimSpace(f.edit.InfrastructureWaterComments.Value)
	f.InfrastructureSewer = common.ExpandRestricted(&f.edit.InfrastructureSewer)
	f.InfrastructureSewerComments = strings.TrimSpace(f.edit.InfrastructureSewerComments.Value)
	f.SearchAndRescue = common.ExpandRestricted(&f.edit.SearchAndRescue)
	f.SearchAndRescueComments = strings.TrimSpace(f.edit.SearchAndRescueComments.Value)
	f.TransportationRoads = common.ExpandRestricted(&f.edit.TransportationRoads)
	f.TransportationRoadsComments = strings.TrimSpace(f.edit.TransportationRoadsComments.Value)
	f.TransportationBridges = common.ExpandRestricted(&f.edit.TransportationBridges)
	f.TransportationBridgesComments = strings.TrimSpace(f.edit.TransportationBridgesComments.Value)
	f.CivilUnrest = common.ExpandRestricted(&f.edit.CivilUnrest)
	f.CivilUnrestComments = strings.TrimSpace(f.edit.CivilUnrestComments.Value)
	f.AnimalIssues = common.ExpandRestricted(&f.edit.AnimalIssues)
	f.AnimalIssuesComments = strings.TrimSpace(f.edit.AnimalIssuesComments.Value)
}

func (f *JurisStat) toEdit() {
	f.StdFields.ToEdit(&f.edit.StdFieldsEdit)
	f.edit.ReportType.Value = f.ReportType
	f.edit.Jurisdiction.Value = f.Jurisdiction
	f.edit.EOCPhone.Value = f.EOCPhone
	f.edit.EOCFax.Value = f.EOCFax
	f.edit.PriEMContactName.Value = f.PriEMContactName
	f.edit.PriEMContactPhone.Value = f.PriEMContactPhone
	f.edit.SecEMContactName.Value = f.SecEMContactName
	f.edit.SecEMContactPhone.Value = f.SecEMContactPhone
	f.edit.OfficeStatus.Value = f.OfficeStatus
	f.edit.GovExpectedOpen.Value = common.SmartJoin(f.GovExpectedOpenDate, f.GovExpectedOpenTime, " ")
	f.edit.GovExpectedClose.Value = common.SmartJoin(f.GovExpectedCloseDate, f.GovExpectedCloseTime, " ")
	f.edit.EOCOpen.Value = f.EOCOpen
	f.edit.EOCActivationLevel.Value = f.EOCActivationLevel
	f.edit.EOCExpectedOpen.Value = common.SmartJoin(f.EOCExpectedOpenDate, f.EOCExpectedOpenTime, " ")
	f.edit.EOCExpectedClose.Value = common.SmartJoin(f.EOCExpectedCloseDate, f.EOCExpectedCloseTime, " ")
	f.edit.StateOfEmergency.Value = f.StateOfEmergency
	f.edit.HowSOESent.Value = f.HowSOESent
	f.edit.Communications.Value = f.Communications
	f.edit.CommunicationsComments.Value = f.CommunicationsComments
	f.edit.Debris.Value = f.Debris
	f.edit.DebrisComments.Value = f.DebrisComments
	f.edit.Flooding.Value = f.Flooding
	f.edit.FloodingComments.Value = f.FloodingComments
	f.edit.Hazmat.Value = f.Hazmat
	f.edit.HazmatComments.Value = f.HazmatComments
	f.edit.EmergencyServices.Value = f.EmergencyServices
	f.edit.EmergencyServicesComments.Value = f.EmergencyServicesComments
	f.edit.Casualties.Value = f.Casualties
	f.edit.CasualtiesComments.Value = f.CasualtiesComments
	f.edit.UtilitiesGas.Value = f.UtilitiesGas
	f.edit.UtilitiesGasComments.Value = f.UtilitiesGasComments
	f.edit.UtilitiesElectric.Value = f.UtilitiesElectric
	f.edit.UtilitiesElectricComments.Value = f.UtilitiesElectricComments
	f.edit.InfrastructurePower.Value = f.InfrastructurePower
	f.edit.InfrastructurePowerComments.Value = f.InfrastructurePowerComments
	f.edit.InfrastructureWater.Value = f.InfrastructureWater
	f.edit.InfrastructureWaterComments.Value = f.InfrastructureWaterComments
	f.edit.InfrastructureSewer.Value = f.InfrastructureSewer
	f.edit.InfrastructureSewerComments.Value = f.InfrastructureSewerComments
	f.edit.SearchAndRescue.Value = f.SearchAndRescue
	f.edit.SearchAndRescueComments.Value = f.SearchAndRescueComments
	f.edit.TransportationRoads.Value = f.TransportationRoads
	f.edit.TransportationRoadsComments.Value = f.TransportationRoadsComments
	f.edit.TransportationBridges.Value = f.TransportationBridges
	f.edit.TransportationBridgesComments.Value = f.TransportationBridgesComments
	f.edit.CivilUnrest.Value = f.CivilUnrest
	f.edit.CivilUnrestComments.Value = f.CivilUnrestComments
	f.edit.AnimalIssues.Value = f.AnimalIssues
	f.edit.AnimalIssuesComments.Value = f.AnimalIssuesComments
}

func (f *JurisStat) validate() {
	f.edit.StdFieldsEdit.Validate()
	if common.ValidateRequired(&f.edit.ReportType) {
		common.ValidateRestricted(&f.edit.ReportType)
	}
	common.ValidateRequired(&f.edit.Jurisdiction)
	validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.EOCPhone)
	validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.PriEMContactName)
	validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.PriEMContactPhone)
	if validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.OfficeStatus) {
		common.ValidateRestricted(&f.edit.OfficeStatus)
	}
	if f.edit.GovExpectedOpen.Value != "" {
		common.ValidateDateTime(&f.edit.GovExpectedOpen)
	} else {
		f.edit.GovExpectedOpen.Problem = ""
	}
	if f.edit.GovExpectedClose.Value != "" {
		common.ValidateDateTime(&f.edit.GovExpectedClose)
	} else {
		f.edit.GovExpectedClose.Problem = ""
	}
	if validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.EOCOpen) {
		common.ValidateRestricted(&f.edit.EOCOpen)
	}
	if validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.EOCActivationLevel) {
		common.ValidateRestricted(&f.edit.EOCActivationLevel)
	}
	if f.edit.EOCExpectedOpen.Value != "" {
		common.ValidateDateTime(&f.edit.EOCExpectedOpen)
	} else {
		f.edit.EOCExpectedOpen.Problem = ""
	}
	if f.edit.EOCExpectedClose.Value != "" {
		common.ValidateDateTime(&f.edit.EOCExpectedClose)
	} else {
		f.edit.EOCExpectedClose.Problem = ""
	}
	if validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.StateOfEmergency) {
		common.ValidateRestricted(&f.edit.StateOfEmergency)
	}
	common.ValidateRestricted(&f.edit.Communications)
	common.ValidateRestricted(&f.edit.Debris)
	common.ValidateRestricted(&f.edit.Flooding)
	common.ValidateRestricted(&f.edit.Hazmat)
	common.ValidateRestricted(&f.edit.EmergencyServices)
	common.ValidateRestricted(&f.edit.Casualties)
	common.ValidateRestricted(&f.edit.UtilitiesGas)
	common.ValidateRestricted(&f.edit.UtilitiesElectric)
	common.ValidateRestricted(&f.edit.InfrastructurePower)
	common.ValidateRestricted(&f.edit.InfrastructureWater)
	common.ValidateRestricted(&f.edit.InfrastructureSewer)
	common.ValidateRestricted(&f.edit.SearchAndRescue)
	common.ValidateRestricted(&f.edit.TransportationRoads)
	common.ValidateRestricted(&f.edit.TransportationBridges)
	common.ValidateRestricted(&f.edit.CivilUnrest)
	common.ValidateRestricted(&f.edit.AnimalIssues)
}

func validateRequiredIfComplete(reportType string, ef *message.EditField) bool {
	if reportType != "Complete" || ef.Value != "" {
		ef.Problem = ""
		return true
	}
	ef.Problem = fmt.Sprintf(`The %q field is required when "Report Type" is "Complete".`, ef.Label)
	return false

}
