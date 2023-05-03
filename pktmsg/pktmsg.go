// Package pktmsg handles encoding and decoding packet messages.  It understands
// RFC-4155 and RFC-5322 email encoding, Outpost-specific feature encoding, the
// Santa Clara County standard subject line format, and PackItForms form
// encoding.
package pktmsg

import (
	"net/mail"
	"strings"
	"time"
)

// CurrentPIFOVersion is the PIFO version number to use in new messages.
const CurrentPIFOVersion = "3.9"

// Message holds a decoded packet message.
type Message struct {
	// Autoresponse is a flag indicating that the received message was an
	// autoresponse message.  This is a non-persistent field, set only on
	// messages retrieved from JNOS (as opposed to local storage).
	Autoresponse bool
	// Body is the body of the message (after decoding).
	Body string
	// BBSRxDate is the date/time at which the message was received by the
	// BBS from which we retrieved it.  This is a non-persistent field, set
	// only on messages retrieved from JNOS (as opposed to local storage).
	BBSRxDate Date
	// FormHTML is the PackItForms HTML file name for the type of
	// PackItForms form in the message.  It is empty for non-PackItForms
	// messages.
	FormHTML string
	// FormTag holds the tag identifying the type of PackItForms form
	// embedded in the message, if any.
	FormTag string
	// FormVersion is the version number of the PackItForms form in the
	// message.  It is empty for non-PackItForms messages.
	FormVersion string
	// From is the From: header of the message.
	From Addresses
	// Handling holds the message handling order.
	Handling string
	// NotPlainText is a flag indicating that the received message was not
	// in plain text encoding.  This is a non-persistent field, set only for
	// messages retrieved from JNOS (as opposed to local storage).
	NotPlainText bool
	// OutpostUrgent is a flag indicating that Outpost considers this
	// message to be urgent.  (In Santa Clara County this normally
	// correlates with an "IMMEDIATE" handling order.)
	OutpostUrgent bool
	// OriginMsgID holds the message ID assigned by the sender of the
	// message.
	OriginMsgID string
	// PIFOVersion is the version number of the encoding of the PackItForms
	// form in the message.  It is empty for non-PackItForms messages.
	PIFOVersion string
	// RequestDeliveryReceipt is an Outpost flag indicating that the
	// recipient should send a delivery receipt for the message to its
	// sender.
	RequestDeliveryReceipt bool
	// RequestReadReceipt is an Outpost flag indicating that the recipient
	// should send a read receipt for the message to its sender, when the
	// message is first read by a human.
	RequestReadReceipt bool
	// ReturnAddr is the return address for the message.  This is a
	// non-persistent field, set only for messages retrieved from JNOS (as
	// opposed to local storage), and not necessarily all of those.
	ReturnAddr string
	// RxArea is the bulletin area from which the message was retrieved.  It
	// is set only for retrieved bulletin messages.
	RxArea string
	// RxBBS is the name of the BBS from which the message was retrieved.
	// It is set only for retrieved messages.
	RxBBS string
	// RxDate is the date/time at which the message was retrieved from JNOS.
	// It is set only for retrieved messages.
	RxDate Date
	// SentDate is the date/time at which the message was sent (i.e., the
	// Date: header of the message).  It is set only for messages that have
	// gone over the air (incoming or outgoing).
	SentDate Date
	// Severity holds the severity of the situation to which the message
	// pertains.  This is an obsolete field.
	Severity string
	// Subject holds the message subject.
	Subject string
	// SubjectHeader is the unparsed Subject: header of the message.
	SubjectHeader string
	// To is the set of destination addresses for the message, the union of
	// the To:, Cc:, and Bcc: headers.
	To Addresses
	// TaggedFields is the set of PackItForms fields in the message.  It is
	// set only for PackItForms forms.
	TaggedFields []TaggedField
}

// IsForm returns whether the Message is a PackItForms form.
func (m *Message) IsForm() bool {
	return m.PIFOVersion != ""
}

// Finalized returns whether the Message has been finalized, i.e., has been
// transmitted or received over the air.  Such messages may still be modified in
// small ways, e.g. to add a destination message number when a delivery receipt
// is received, but from a human perspective they are no longer changeable.
func (m *Message) Finalized() bool {
	return m.SentDate != ""
}

// Addresses is a string containing a comma-separated list of JNOS mailbox
// names, BBS network addresses, and/or email addresses.
type Addresses string

// String returns the string value of the address list.
func (a Addresses) String() string { return string(a) }

// Address returns the first address on the list, as a bare address with name
// comments and/or angle brackets removed.  If the list is empty, Address
// returns an empty string.
func (a Addresses) Address() string {
	if addrs := a.Addresses(); len(addrs) != 0 {
		return addrs[0]
	}
	return ""
}

// Addresses returns the address list as a slice of strings.  Each element of
// the slice is a bare address, with name comments and/or angle brackets
// removed.
func (a Addresses) Addresses() (addrs []string) {
	addrs = strings.Split(string(a), ",")
	j := 0
	for _, addr := range addrs {
		addr = strings.TrimSpace(addr)
		if pa, err := mail.ParseAddress(addr); err == nil {
			addrs[j] = pa.Address
			j++
		} else if addr != "" {
			addrs[j] = addr
			j++
		}
	}
	addrs = addrs[:j]
	return addrs
}

// SetString sets the string value of the address list.
func (a *Addresses) SetString(value string) { *a = Addresses(value) }

// A Date is an RFC-1123Z-formatted date, as seen in message headers.
type Date string

// String returns the string value of the Date.
func (d Date) String() string { return string(d) }

// Time returns the Time value of the Date.
func (d Date) Time() (t time.Time) {
	t, _ = time.Parse(time.RFC1123Z, string(d))
	return t
}

// SetString sets the string value of the Date.
func (d *Date) SetString(value string) { *d = Date(value) }

// SetTime sets the Time value of the Date.
func (d *Date) SetTime(value time.Time) {
	if value.IsZero() {
		*d = ""
	} else {
		*d = Date(value.Format(time.RFC1123Z))
	}
}

// A TaggedField is a PackItForms field with a tag and associated value.
type TaggedField struct {
	Tag   string
	Value string
}
