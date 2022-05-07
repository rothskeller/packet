package sheltstat

import (
	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/xscmsg"
	"steve.rothskeller.net/packet/xscmsg/internal/xscform"
)

func init() {
	for _, fd := range formDefinitions {
		fd.Name = "OA shelter status form"
		fd.Article = "an"
		fd.Comments["7a."] = "Mass Care and Shelter Unit, Care and Shelter Branch, Operations Section, ..."
		if fd.FindField("34b.") != nil {
			ff := fd.FindField("33b.")
			ff.Default = "(computed)"
			ff.Values = sheltstat21.FindField("33b.").Values
			ff.ComputedFromField = "34b."
			ff.Validations = []xscform.ValidateFunc{xscform.ValidateComputedChoice}
			fd.Annotations["33b."] = "shelter-city-code"
		}
		if fd.FindField("49a.") != nil {
			ff := fd.FindField("50a.")
			ff.Default = "(computed)"
			ff.Values = sheltstat21.FindField("50a.").Values
			ff.ComputedFromField = "49a."
			ff.Validations = []xscform.ValidateFunc{xscform.ValidateComputedChoice}
			fd.Annotations["50a."] = "managed-by-code"
		}
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

// Form is an OA shelter status form.
type Form struct {
	*xscform.XSCForm
}
