// Package envelope contains all the knowledge about how message content (forms,
// plain text messages, receipts, etc.) are wrapped in envelopes for sending,
// receiving, and storage.
package envelope

import "time"

// An Envelope structure contains the envelope of a packet message.
type Envelope struct {
	// RxBBS is the name of the BBS from which the message was retrieved.
	// It is set only for received messages.
	ReceivedBBS string
	// RxArea is the bulletin area from which the message was retrieved.  It
	// is set only for received bulletin messages.
	ReceivedArea string
	// RxDate is the date/time at which the message was retrieved from JNOS.
	// It is set only for received messages.
	ReceivedDate time.Time
	// Autoresponse is a flag indicating that the received message was an
	// autoresponse message.  This is a non-persistent field, set only on
	// messages retrieved from JNOS (as opposed to local storage).
	Autoresponse bool
	// ReturnAddr is the return address for the message.  This is a
	// non-persistent field, set only for messages retrieved from JNOS (as
	// opposed to local storage), and not necessarily all of those.
	ReturnAddr string
	// BBSReceivedDate is the date/time at which the message was received by
	// the BBS from which we retrieved it.  This is a non-persistent field,
	// set only on messages retrieved from JNOS (as opposed to local
	// storage).
	BBSReceivedDate time.Time
	// From is the sender of the message, from the From: header.  Note that
	// while From: headers are allowed by RFC-5322 to have more than one
	// address, only the first one is stored in this field.
	From string
	// To is the set of destination addresses for the message, the union of
	// the To:, Cc:, and Bcc: headers.
	To []string
	// Date is the date/time at which the message was sent, from the Date:
	// header.  It is set only for messages that have gone over the air
	// (received or transmitted).
	Date time.Time
	// SubjectLine is the subject line of the message.
	SubjectLine string
	// NotPlainText is a flag indicating that the received message was not
	// in plain text encoding.  This is a non-persistent field, set only for
	// messages retrieved from JNOS (as opposed to local storage).
	NotPlainText bool
	// OutpostUrgent is a flag indicating that Outpost considers this
	// message to be urgent.  (In Santa Clara County this normally
	// correlates with an "IMMEDIATE" handling order.)
	OutpostUrgent bool
	// RequestDeliveryReceipt is an Outpost flag indicating that the
	// recipient should send a delivery receipt for the message to its
	// sender.
	RequestDeliveryReceipt bool
	// RequestReadReceipt is an Outpost flag indicating that the recipient
	// should send a read receipt for the message to its sender, when the
	// message is first read by a human.
	RequestReadReceipt bool
	// ReadyToSend is a flag indicating that the (outgoing, untransmitted)
	// message is ready to be sent.  When false, the message is a draft.
	// This field is ignored for received or transmitted messages.
	ReadyToSend bool
}

// IsReceived returns whether the message was received (as opposed to sent or
// pending send).
func (e *Envelope) IsReceived() bool {
	return e.ReceivedBBS != ""
}

// IsFinal returns whether the message has gone over the air (either received or
// transmitted).  Note that even if this returns true, the message can still
// be changed in some ways (e.g. to add a destination message number when a
// delivery receipt arrives).
func (e *Envelope) IsFinal() bool {
	return !e.Date.IsZero()
}