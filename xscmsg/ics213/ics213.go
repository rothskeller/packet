package ics213

import (
	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/xscmsg"
	"steve.rothskeller.net/packet/xscmsg/internal/xscform"
)

var ics213v20rx, ics213v20tx, ics213v21rx, ics213v21tx *xscform.FormDefinition

func init() {
	// Clean up an ugly annotation in the PIFO HTML encoding.
	ics213v20.Annotations["7."] = "to-ics-position"
	ics213v21.Annotations["7."] = "to-ics-position"
	ics213v22.Annotations["7."] = "to-ics-position"
	// ICS-213 forms prior to version 2.2 have ugly weirdness in the message
	// numbers.  Rather than have a straightforward origin and destination
	// message number, they have three message number fields, and they
	// change which two of those three are used depending on whether we are
	// transmitting or receiving.  So we need to clone the field definitions
	// and tweak them for the two cases.  Thankfully, starting with version
	// 2.2 we no longer have that issue.
	ics213v20rx = makeICS213RxForm(ics213v20)
	ics213v21rx = makeICS213RxForm(ics213v21)
	ics213v20tx = makeICS213TxForm(ics213v20)
	ics213v21tx = makeICS213TxForm(ics213v21)
	// Register the type.
	xscmsg.RegisterType(Create, Recognize)
}

// Create creates a new message of the type identified by the supplied tag.  If
// the tag is not recognized by this package, Create returns nil.
func Create(tag string) xscmsg.XSCMessage {
	if tag == ics213v21tx.Tag {
		return &Form{xscform.CreateForm(ics213v21tx)}
	}
	return nil
}

// Recognize examines the supplied Message to see if it is of the type defined
// in this package.  If so, it returns the appropriate XSCMessage implementation
// wrapping it.  If not, it returns nil.
func Recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.XSCMessage {
	if xf := xscform.RecognizeForm(ics213v22, msg, form); xf != nil {
		return &Form{xf}
	}
	// If it's an older version, we have to guess whether it's transmit-side
	// or receive-side.  If it has field 3. (receiver's message number when
	// sending), we will assume we are the sender and use the transmit-side
	// field definitions.  If it doesn't, we'll assume we are the receiver
	// and use the receive-side field definitions.
	if form != nil && form.Has("3.") {
		if xf := xscform.RecognizeForm(ics213v21tx, msg, form); xf != nil {
			return &Form{xf}
		}
		if xf := xscform.RecognizeForm(ics213v20tx, msg, form); xf != nil {
			return &Form{xf}
		}
	} else {
		if xf := xscform.RecognizeForm(ics213v21rx, msg, form); xf != nil {
			return &Form{xf}
		}
		if xf := xscform.RecognizeForm(ics213v20rx, msg, form); xf != nil {
			return &Form{xf}
		}
	}
	return nil
}

// Form is an ICS-213 general message form.
type Form struct {
	*xscform.XSCForm
}

// makeICS213RxForm makes a receive-side ICS-213 form definition by changing the
// DestinationNumberField to the middle message number field ("MsgNo"), and
// removing the left-side message number field (2.).
func makeICS213RxForm(fd *xscform.FormDefinition) (rx *xscform.FormDefinition) {
	rx = new(xscform.FormDefinition)
	*rx = *fd
	rx.DestinationNumberField = "MsgNo"
	rx.Fields = make([]*xscform.FieldDefinition, len(fd.Fields)-1)
	j := 0
	for _, ff := range fd.Fields {
		if ff.Tag != "2." {
			rx.Fields[j] = ff
			j++
		}
	}
	return rx
}

// makeICS213TxForm makes a transmit-side ICS-213 form definition by changing
// the OriginNumberField to the middle message number field ("MsgNo"), and
// removing the right-side message number field (3.).
func makeICS213TxForm(fd *xscform.FormDefinition) (tx *xscform.FormDefinition) {
	tx = new(xscform.FormDefinition)
	*tx = *fd
	tx.OriginNumberField = "MsgNo"
	tx.Fields = make([]*xscform.FieldDefinition, len(fd.Fields)-1)
	j := 0
	for _, ff := range fd.Fields {
		if ff.Tag != "3." {
			tx.Fields[j] = ff
			j++
		}
	}
	return tx
}
