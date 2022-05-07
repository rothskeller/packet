package ahfacstat

import (
	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/xscmsg"
	"steve.rothskeller.net/packet/xscmsg/internal/xscform"
)

func init() {
	for _, fd := range formDefinitions {
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
	}
	xscmsg.RegisterType(Create, Recognize)
}

// Create creates a new message of the type identified by the supplied tag.  If
// the tag is not recognized by this package, Create returns nil.
func Create(tag string) xscmsg.XSCMessage {
	for _, fd := range formDefinitions {
		if tag == fd.Tag {
			return &Form{xscform.CreateForm(fd)}
		}
	}
	return nil
}

// Recognize examines the supplied Message to see if it is of the type defined
// in this package.  If so, it returns the appropriate XSCMessage implementation
// wrapping it.  If not, it returns nil.
func Recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.XSCMessage {
	for _, fd := range formDefinitions {
		if xf := xscform.RecognizeForm(fd, msg, form); xf != nil {
			return &Form{xf}
		}
	}
	return nil
}

// Form is an allied health facility status report form.
type Form struct {
	*xscform.XSCForm
}
