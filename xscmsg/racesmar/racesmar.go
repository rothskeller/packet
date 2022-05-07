package racesmar

import (
	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/xscmsg"
	"steve.rothskeller.net/packet/xscmsg/internal/xscform"
)

func init() {
	for _, fd := range formDefinitions {
		fd.Name = "RACES mutual aid request form"
		fd.Article = "a"
		fd.Comments["7a."] = "RACES Chief Radio Officer, RACES Unit, Operations Section, ..."
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

// Form is a RACES mutual aid request form.
type Form struct {
	*xscform.XSCForm
}
