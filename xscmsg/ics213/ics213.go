package ics213

import (
	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/xscmsg"
	"steve.rothskeller.net/packet/xscmsg/internal/xscform"
)

func init() {
	for _, fd := range formDefinitions {
		fd.Name = "ICS-213 general message form"
		fd.Article = "an"
		fd.Annotations["7."] = "to-ics-position"
	}
	xscmsg.RegisterType(Create, Recognize)
}

// Create creates a new message of the type identified by the supplied tag.  If
// the tag is not recognized by this package, Create returns nil.
func Create(tag string) xscmsg.XSCMessage {
	for _, fd := range formDefinitions {
		if tag == fd.Tag && fd != ics213v22 {
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

// Form is an ICS-213 general message form.
type Form struct {
	*xscform.XSCForm
}

// EncodeSubject returns the encoded subject line of the message based on its
// contents.
func (f *Form) EncodeSubject() string {
	// We have to override the XSCForm implementation because we want to use
	// the message number from MsgNo, which is not marked as the origin
	// message number in the field definitions.
	ho, _ := xscmsg.ParseHandlingOrder(f.Get("5."))
	return xscmsg.EncodeSubject(f.Get("MsgNo"), ho, "ICS213", f.Get("10."))
}

// OriginNumber returns the origin message number of the message, if any.
func (f *Form) OriginNumber() string {
	if msgnum := f.Get("2."); msgnum != "" {
		return msgnum
	}
	return f.Get("MsgNo")
}

// SetOriginNumber sets the originmessage number of the message, if the message
// type supports that.
func (f *Form) SetOriginNumber(msgnum string) { f.Set("MsgNo", msgnum) }

// DestinationNumber returns the destination message number of the message, if
// any.
func (f *Form) DestinationNumber() string {
	if f.Get("2.") != "" {
		return f.Get("MsgNo")
	}
	return f.Get("3.")
}

// SetDestinationNumber sets the destination message number of the message, if
// the message type supports that.
func (f *Form) SetDestinationNumber(msgnum string) {
	if f.Get("2.") != "" {
		f.Set("MsgNo", msgnum)
	} else {
		f.Set("3.", msgnum)
	}
}

// Routing returns the To ICS Position and To Location fields of the form, if
// it has them.
func (f *Form) Routing() (pos, loc string) {
	return f.Get("7."), f.Get("8.")
}
