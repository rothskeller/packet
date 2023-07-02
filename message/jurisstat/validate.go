package jurisstat

import "github.com/rothskeller/packet/message/common"

// Validate checks the contents of the message for compliance with rules
// enforced by standard Santa Clara County packet software (Outpost and
// PackItForms).  It returns a list of strings describing problems that those
// programs would flag or block.
func (f *JurisStat) Validate() (problems []string) {
	problems = append(problems, f.StdFields.Validate()...)
	switch f.ReportType {
	case "":
		problems = append(problems, `The "Report Type" field is required.`)
	case "Update", "Complete":
		break
	default:
		problems = append(problems, `The "Report Type" field does not contain a valid report type value.`)
	}
	if f.FormVersion >= "2.2" {
		switch f.JurisdictionCode {
		case "":
			if f.Jurisdiction != "" {
				problems = append(problems, `The hidden, required "Jurisdiction Code" field is not set.`)
			} // otherwise we let the problem with Jurisdiction be reported.
		case "Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated":
			break
		default:
			problems = append(problems, `The hidden "Jurisdiction Code" field does not contain a valid jurisdiction code.`)
		}
		if f.Jurisdiction == "" {
			problems = append(problems, `The "Jurisdiction Name" field is required.`)
		}
	} else {
		switch f.Jurisdiction {
		case "":
			problems = append(problems, `The "Jurisdiction Name" field is required.`)
		case "Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated":
			break
		default:
			problems = append(problems, `The "Jurisdiction Name" field does not contain a valid jurisdiction name.`)
		}
	}
	if f.EOCPhone == "" {
		if f.ReportType == "Complete" {
			problems = append(problems, `The "EOC Phone" field is required when the "Report Type" is "Complete".`)
		}
	} else if !common.PIFOPhoneNumberRE.MatchString(f.EOCPhone) {
		problems = append(problems, `The "EOC Phone" field does not contain a valid phone number.`)
	}
	if f.EOCFax != "" && !common.PIFOPhoneNumberRE.MatchString(f.EOCFax) {
		problems = append(problems, `The "EOC Fax" field does not contain a valid phone number.`)
	}
	if f.PriEMContactName == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Primary EM Contact Name" field is required when the "Report Type" is "Complete".`)
	}
	if f.PriEMContactPhone == "" {
		if f.ReportType == "Complete" {
			problems = append(problems, `The "Primary EM Contact Phone" field is required when the "Report Type" is "Complete".`)
		}
	} else if !common.PIFOPhoneNumberRE.MatchString(f.PriEMContactPhone) {
		problems = append(problems, `The "Primary EM Contact Phone" field does not contain a valid phone number.`)
	}
	if f.SecEMContactPhone != "" && !common.PIFOPhoneNumberRE.MatchString(f.SecEMContactPhone) {
		problems = append(problems, `The "Secondary EM Contact Phone" field does not contain a valid phone number.`)
	}
	switch f.OfficeStatus {
	case "":
		if f.ReportType == "Complete" {
			problems = append(problems, `The "Office Status" field is required when the "Report Type" is "Complete".`)
		}
	case "Unknown", "Open", "Closed":
		break
	default:
		problems = append(problems, `The "Office Status" field does not contain a valid office status value.`)
	}
	if f.GovExpectedOpenDate != "" && !common.PIFODateRE.MatchString(f.GovExpectedOpenDate) {
		problems = append(problems, `The "Government Office Status: Expected to Open Date" field does not contain a valid date.`)
	}
	if f.GovExpectedOpenTime != "" && !common.PIFOTimeRE.MatchString(f.GovExpectedOpenTime) {
		problems = append(problems, `The "Government Office Status: Expected to Open Time" field does not contain a valid time.`)
	}
	if f.GovExpectedCloseDate != "" && !common.PIFODateRE.MatchString(f.GovExpectedCloseDate) {
		problems = append(problems, `The "Government Office Status: Expected to Close Date" field does not contain a valid date.`)
	}
	if f.GovExpectedCloseTime != "" && !common.PIFOTimeRE.MatchString(f.GovExpectedCloseTime) {
		problems = append(problems, `The "Government Office Status: Expected to Close Time" field does not contain a valid time.`)
	}
	switch f.EOCOpen {
	case "":
		if f.ReportType == "Complete" {
			problems = append(problems, `The "EOC Open" field is required when the "Report Type" is "Complete".`)
		}
	case "Unknown", "Yes", "No":
		break
	default:
		problems = append(problems, `The "EOC Open" field does not contain a valid EOC open value.`)
	}
	switch f.EOCActivationLevel {
	case "":
		if f.ReportType == "Complete" {
			problems = append(problems, `The "EOC Status: Activation" field is required when the "Report Type" is "Complete".`)
		}
	case "Normal", "Duty Officer", "Monitor", "Partial", "Full":
		break
	default:
		problems = append(problems, `The "EOC Status: Activation" field does not contain a valid EOC activation value.`)
	}
	if f.EOCExpectedOpenDate != "" && !common.PIFODateRE.MatchString(f.EOCExpectedOpenDate) {
		problems = append(problems, `The "EOC Status: Expected to Open Date" field does not contain a valid date.`)
	}
	if f.EOCExpectedOpenTime != "" && !common.PIFOTimeRE.MatchString(f.EOCExpectedOpenTime) {
		problems = append(problems, `The "EOC Status: Expected to Open Time" field does not contain a valid time.`)
	}
	if f.EOCExpectedCloseDate != "" && !common.PIFODateRE.MatchString(f.EOCExpectedCloseDate) {
		problems = append(problems, `The "EOC Status: Expected to Close Date" field does not contain a valid date.`)
	}
	if f.EOCExpectedCloseTime != "" && !common.PIFOTimeRE.MatchString(f.GovExpectedCloseTime) {
		problems = append(problems, `The "EOC Status: Expected to Close Time" field does not contain a valid time.`)
	}
	switch f.StateOfEmergency {
	case "":
		if f.ReportType == "Complete" {
			problems = append(problems, `The "State of Emergency" field is required when the "Report Type" is "Complete".`)
		}
	case "Unknown", "Yes", "No":
		break
	default:
		problems = append(problems, `The "State of Emergency" field does not contain a valid SOE declaration value.`)
	}
	if f.HowSOESent == "" && f.StateOfEmergency == "Yes" {
		problems = append(problems, `The "Declaration: Attachment" field is required when the "State of Emergency" field is "Yes".`)
	} else if f.HowSOESent != "" && f.StateOfEmergency != "Yes" {
		problems = append(problems, `The "Declaration: Attachment" field is allowed only when the "State of Emergency" field is "Yes".`)
	}
	switch f.Communications {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Communications" field does not contain a valid situation code.`)
	}
	switch f.Debris {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Debris" field does not contain a valid situation code.`)
	}
	switch f.Flooding {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Flooding" field does not contain a valid situation code.`)
	}
	switch f.Hazmat {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Hazmat" field does not contain a valid situation code.`)
	}
	switch f.EmergencyServices {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Emergency Services" field does not contain a valid situation code.`)
	}
	switch f.Casualties {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Casualties" field does not contain a valid situation code.`)
	}
	switch f.UtilitiesGas {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Utilities (Gas)" field does not contain a valid situation code.`)
	}
	switch f.UtilitiesElectric {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Utilities (Electric)" field does not contain a valid situation code.`)
	}
	switch f.InfrastructurePower {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Infrastructure (Power)" field does not contain a valid situation code.`)
	}
	switch f.InfrastructureWater {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Infrastructure (Water Systems)" field does not contain a valid situation code.`)
	}
	switch f.InfrastructureSewer {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Infrastructure (Sewer Systems)" field does not contain a valid situation code.`)
	}
	switch f.SearchAndRescue {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Search and Rescue" field does not contain a valid situation code.`)
	}
	switch f.TransportationRoads {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Transportation (Roads)" field does not contain a valid situation code.`)
	}
	switch f.TransportationBridges {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Transportation (Bridges)" field does not contain a valid situation code.`)
	}
	switch f.CivilUnrest {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Civil Unrest" field does not contain a valid situation code.`)
	}
	switch f.AnimalIssues {
	case "", "Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out":
		break
	default:
		problems = append(problems, `The "Current Situation: Animal Issues" field does not contain a valid situation code.`)
	}
	return problems
}
