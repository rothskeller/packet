package main

var cities = []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"}
var managedBy = []string{"American Red Cross", "Private", "Community", "Government", "Other"}

// applyFixups applies manual fixups to the parsed form definitions.  This
// corrects oddities in the PackItForms HTML files, adds helpful comments, and
// addresses fields (e.g., combo boxes) that our parser can't handle.
func applyFixups(fd *formDefinition) {
	switch fd.Tag {
	case "AHFacStat":
		fd.Comments["7a."] = "EMS Unit, Public Health Unit, Medical Health Branch, Operations Section, ..."
		fd.Comments["7b."] = "MHJOC, County EOC, ..."
		fd.Annotations["30."] = "eoc-closed-contact"
		fd.Annotations["31a."] = "to-evacuate"
		fd.Annotations["40a."] = "skilled-nursing-beds-staffed-m"
		fd.Annotations["40b."] = "skilled-nursing-beds-staffed-f"
		fd.Annotations["40c."] = "skilled-nursing-beds-vacant-m"
		fd.Annotations["40d."] = "skilled-nursing-beds-vacant-f"
		fd.Annotations["40e."] = "skilled-nursing-beds-surge"
		fd.Annotations["41a."] = "assisted-living-beds-staffed-m"
		fd.Annotations["41b."] = "assisted-living-beds-staffed-f"
		fd.Annotations["41c."] = "assisted-living-beds-vacant-m"
		fd.Annotations["41d."] = "assisted-living-beds-vacant-f"
		fd.Annotations["41e."] = "assisted-living-beds-surge"
		fd.Annotations["42a."] = "sub-acute-beds-staffed-m"
		fd.Annotations["42b."] = "sub-acute-beds-staffed-f"
		fd.Annotations["42c."] = "sub-acute-beds-vacant-m"
		fd.Annotations["42d."] = "sub-acute-beds-vacant-f"
		fd.Annotations["42e."] = "sub-acute-beds-surge"
		fd.Annotations["43a."] = "alzheimers-beds-staffed-m"
		fd.Annotations["43b."] = "alzheimers-beds-staffed-f"
		fd.Annotations["43c."] = "alzheimers-beds-vacant-m"
		fd.Annotations["43d."] = "alzheimers-beds-vacant-f"
		fd.Annotations["43e."] = "alzheimers-beds-surge"
		fd.Annotations["44a."] = "ped-sub-acute-beds-staffed-m"
		fd.Annotations["44b."] = "ped-sub-acute-beds-staffed-f"
		fd.Annotations["44c."] = "ped-sub-acute-beds-vacant-m"
		fd.Annotations["44d."] = "ped-sub-acute-beds-vacant-f"
		fd.Annotations["44e."] = "ped-sub-acute-beds-surge"
		fd.Annotations["45a."] = "psychiatric-beds-staffed-m"
		fd.Annotations["45b."] = "psychiatric-beds-staffed-f"
		fd.Annotations["45c."] = "psychiatric-beds-vacant-m"
		fd.Annotations["45d."] = "psychiatric-beds-vacant-f"
		fd.Annotations["45e."] = "psychiatric-beds-surge"
		fd.Annotations["46a."] = "other-care-beds-staffed-m"
		fd.Annotations["46b."] = "other-care-beds-staffed-f"
		fd.Annotations["46c."] = "other-care-beds-vacant-m"
		fd.Annotations["46d."] = "other-care-beds-vacant-f"
		fd.Annotations["46e."] = "other-care-beds-surge"
		fd.Annotations["50a."] = "dialysis-chairs"
		fd.Annotations["50b."] = "dialysis-vacant-chairs"
		fd.Annotations["50c."] = "dialysis-front-staff"
		fd.Annotations["50d."] = "dialysis-support-staff"
		fd.Annotations["50e."] = "dialysis-providers"
		fd.Annotations["51a."] = "surgical-chairs"
		fd.Annotations["51b."] = "surgical-vacant-chairs"
		fd.Annotations["51c."] = "surgical-front-staff"
		fd.Annotations["51d."] = "surgical-support-staff"
		fd.Annotations["51e."] = "surgical-providers"
		fd.Annotations["52a."] = "clinic-chairs"
		fd.Annotations["52b."] = "clinic-vacant-chairs"
		fd.Annotations["52c."] = "clinic-front-staff"
		fd.Annotations["52d."] = "clinic-support-staff"
		fd.Annotations["52e."] = "clinic-providers"
		fd.Annotations["53a."] = "home-health-chairs"
		fd.Annotations["53b."] = "home-health-vacant-chairs"
		fd.Annotations["53c."] = "home-health-front-staff"
		fd.Annotations["53d."] = "home-health-support-staff"
		fd.Annotations["53e."] = "home-health-providers"
		fd.Annotations["54a."] = "adulty-day-ctr-chairs"
		fd.Annotations["54b."] = "adulty-day-ctr-vacant-chairs"
		fd.Annotations["54c."] = "adulty-day-ctr-front-staff"
		fd.Annotations["54d."] = "adulty-day-ctr-support-staff"
		fd.Annotations["54e."] = "adulty-day-ctr-providers"
	case "EOC213RR":
		fd.Name = "EOC-213RR resource request form"
		fd.Article = "an"
	case "ICS213":
		fd.Name = "ICS-213 general message form"
		fd.Article = "an"
		fd.Annotations["7."] = "to-ics-position"
	case "JurisStat":
		fd.Name = "OA jurisdiction status form"
		fd.Article = "an"
		fd.Comments["7a."] = "Situation Analysis Unit, Planning Section, ..."
		if fd.findField("22.") != nil {
			ff := fd.findField("21.")
			ff.Default = "(computed)"
			ff.Values = cities
			ff.ComputedFromField = "22."
			ff.Validations = []string{"computed-choice"}
			fd.Annotations["21."] = "jurisdiction-code"
		}
	case "MuniStat":
		fd.Name = "OA municipal status form"
		fd.Article = "an"
		fd.Comments["7a."] = "Situation Analysis Unit, Planning Section, ..."
		if fd.findField("22.") != nil {
			ff := fd.findField("21.")
			ff.Default = "(computed)"
			ff.Values = cities
			ff.ComputedFromField = "22."
			ff.Validations = []string{"computed-choice"}
			fd.Annotations["21."] = "jurisdiction-code"
		}
	case "RACES-MAR":
		fd.Name = "RACES mutual aid request form"
		fd.Article = "a"
		fd.Comments["7a."] = "RACES Chief Radio Officer, RACES Unit, Operations Section, ..."
	case "SheltStat":
		fd.Name = "OA shelter status form"
		fd.Article = "an"
		fd.Comments["7a."] = "Mass Care and Shelter Unit, Care and Shelter Branch, Operations Section, ..."
		if fd.findField("34b.") != nil {
			ff := fd.findField("33b.")
			ff.Default = "(computed)"
			ff.Values = cities
			ff.ComputedFromField = "34b."
			ff.Validations = []string{"computed-choice"}
			fd.Annotations["33b."] = "shelter-city-code"
		}
		if fd.findField("49a.") != nil {
			ff := fd.findField("50a.")
			ff.Default = "(computed)"
			ff.Values = managedBy
			ff.ComputedFromField = "49a."
			ff.Validations = []string{"computed-choice"}
			fd.Annotations["50a."] = "managed-by-code"
		}
	}
}

func (fd *formDefinition) findField(tag string) *fieldDefinition {
	for _, ff := range fd.Fields {
		if ff.Tag == tag {
			return ff
		}
	}
	return nil
}
