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
	// From is a comma-separated list of addresses of senders of the
	// message, from the "From:" header.  In the vast majority of cases,
	// there's only one address on the list.
	From string
	// To is a comma-separated list of destination addresses for the message,
	// from the "To:", "Cc:", and "Bcc:" headers (we don't distinguish).
	To string
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
	// DeliveredDate is the date/time when the (sent) message was delivered
	// to its recipient, as pulled from the delivery receipt they sent back
	// to us.  Note that this is a string, not a time.Time, because there is
	// no standard format for it in the delivery receipt.
	DeliveredDate string
	// DeliveredRMI is the message ID assigned to the (sent) message by its
	// recipient, as pulled from the delivery receipt they sent back to us.
	DeliveredRMI string
}

// IsReceived returns whether the message was received (as opposed to sent or
// pending send).
func (env *Envelope) IsReceived() bool {
	return env.ReceivedBBS != ""
}

// IsFinal returns whether the message has gone over the air (either received or
// transmitted).  Note that even if this returns true, the message can still
// be changed in some ways (e.g. to add a destination message number when a
// delivery receipt arrives).
func (env *Envelope) IsFinal() bool {
	return !env.Date.IsZero()
}
