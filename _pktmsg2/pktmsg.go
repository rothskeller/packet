// Package pktmsg handles encoding and decoding packet messages.  It understands
// RFC-4155 and RFC-5322 email encoding, PackItForms form encoding, and
// Outpost-specific feature encodings.
package pktmsg

import (
	"time"
)

// MessageType describes a message type.
type MessageType struct {
	// Tag is the string that identifies the message type.
	Tag string
	// Name is the English name of the message type, in prose case.
	Name string
	// Article is either "a" or "an", whichever is appropriate as the
	// indefinite article preceding Name.
	Article string
}

// ReceivedInfo contains extra metadata that apply only to received messages.
type ReceivedInfo struct {
	// Flags is a bitmask of flags describing the received message.
	Flags MessageFlag
	// BBS is the name of the BBS from which the message was retrieved.
	BBS string
	// CallSign is the call sign with which we connected to the BBS to
	// retrieve the message.
	CallSign string
	// Arrived is the time when the message arrived at the BBS, if known.
	// It is the zero time if the arrival time is not known.
	Arrived time.Time
	// Retrieved is the time when the message was retrieved from the BBS.
	Retrieved time.Time
	// Area is the name of the message area where the message was retrieved.
	// This will be an empty string for messages retrieved from CallSign's
	// home area.  Otherwise, it may be either the call sign of a shared
	// mailbox (e.g. "XSCPERM") or a category@distribution string (e.g.
	// "XND@ALLXSC").
	Area string
	// ReturnAddress is the return address for the message.
	ReturnAddress string
}

// Message represents a packet message.  A message essentially consists of a
// type and a set of Fields.  Some fields are common across all or many message
// types, and are identified with a FieldKey (which see).
type Message interface {
	// Type returns the message type.
	Type() MessageType
	// Received returns information about the reception of the message, or
	// nil if it is an outgoing message.
	Received() *ReceivedInfo
	// FieldCount returns the number of fields in the message type.  It is
	// used as the bounds around iterative calls to FieldByIndex.
	FieldCount() int
	// FieldByIndex returns the field at the specified index in the ordered
	// list of fields in the message.  This is the order of human
	// presentation, which may not be identical to the order of encoding.
	// The index is zero-based.  If the index is out of range, nil is
	// returned.
	FieldByIndex(int) Field
	// FieldByTag returns the message field with the specified tag, or nil
	// if there is none.
	FieldByTag(string) Field
	// FieldByKey returns the message field with the specified key, or nil
	// if there is none.
	FieldByKey(FieldKey) Field
	// FieldValue returns the value of the field with the specified tag, or
	// an empty string if there is none.
	FieldValue(string) Value
	// KeyValue returns the value of the field with the specified key, or an
	// empty string if there is none.
	KeyValue(FieldKey) Value
	// Calculate recalculates the values of all calculated fields in the
	// message.
	Calculate()
	// Validate validates the values of all fields in the message, returning
	// a list of problem strings that is empty if the message is valid.  If
	// pifo is true, only those problems that would be flagged by
	// PackItForms are returned.
	Validate(pifo bool) []string
	// Encode returns the RFC-5322-encoded form of the message, suitable for
	// saving on disk.
	Encode() string
	// EncodeTx returns the message, encoded as needed for transmission via
	// JNOS.  Specifically, it returns a list of destination addresses, an
	// encoded subject line, and an encoded body.
	EncodeTx() (to []string, subject, body string)
}

// MessageFlag contains a flag, or a bitmask of flags, describing a Message.
type MessageFlag uint8

const (
	// NotPlainText indicates that the message was not entirely plain text.
	NotPlainText MessageFlag = 1 << iota
	// AutoResponse indicates that the message appears to be a bounce
	// message or other auto-responder message.
	AutoResponse
	// RequestDeliveryReceipt indicates that the recipient should send a
	// delivery receipt to the sender when the message is delivered.
	RequestDeliveryReceipt
	// RequestReadReceipt indicates that the recipient should send a read
	// receipt to the sender when the message has been read by a human.
	RequestReadReceipt
	// OutpostUrgent indicates that the Outpost messaging software should
	// display this message as an urgent message.
	OutpostUrgent
)

// A FieldKey is an identifier of a well-known field that is constant across all
// messages containing that field, even if it has different tags in different
// message types.
type FieldKey string

// Values for FieldKey.  Generally, these are all of the fields that
// non-message-type-specific code needs to interact with.
const (
	// FToAddrs is the comma-separated list of destination addresses for the
	// message.  It will exist in all messages, and be set in all messages
	// with the possible exception of draft outgoing messages.
	FToAddrs FieldKey = "TO"
	// FFromAddr is the sender address for the message.  It will exist in
	// all messages, and be set in all messages that have been transmitted
	// or that are ready to be transmitted.
	FFromAddr FieldKey = "FROM"
	// FSent is the date the message was sent, in time.RFC3339 format.  It
	// will exist in all messages, and be set in all messages that have been
	// transmitted (and only those).
	FSent FieldKey = "SENT"
	// FReceived is the date the message was received, in time.RFC3339
	// format.  It will exist only in messages that have been received, and
	// will always have a value in those messages.
	FReceived FieldKey = "RECEIVED"
	// FRecipient is the address of the mailbox where the message was
	// received.  It will exist only in messages that have been received,
	// and will always have a value in those messages.  Its value will
	// generally be one of the addresses in FToAddrs.
	FRecipient FieldKey = "RECIPIENT"
	// FOriginMsgNo is the origin message number field.  It is set by code
	// generating a new outgoing message, and read by code generating
	// subject lines.
	FOriginMsgNo FieldKey = "ORIGIN_MESSAGE_NUMBER"
	// FDestinationMsgNo is the destination (receiver) message number field.
	// Code that is receiving messages will set this to a local message
	// number.
	FDestinationMsgNo FieldKey = "DESTINATION_MESSAGE_NUMBER"
	// FHandling is the handling order for the message.  It gets used in
	// generating subject lines.
	FHandling FieldKey = "HANDLING"
	// FToICSPosition is the To ICS Position field.
	FToICSPosition FieldKey = "TO_ICS_POSITION"
	// FToLocation is the To Location field.
	FToLocation FieldKey = "TO_LOCATION"
	// FSubject is the field containing the subject of the message.  Note
	// that some encodings place additional information in the "Subject:"
	// header of the RFC-5322-encoded message; that additional information
	// is not part of this field.
	FSubject FieldKey = "SUBJECT"
	// FReference is the Reference field.  It is the field that contains the
	// origin message ID of the message to which the instant message is a
	// reply.
	FReference FieldKey = "REFERENCE"
	// FBody is the field into which default message body text is placed.
	// It is generally the most prominent or general-purpose multi-line text
	// field in the message.
	FBody FieldKey = "BODY"
	// FOpCall is the operator call sign field.  It gets set by code that
	// creates a new outgoing message, or by code receiving a message.
	FOpCall FieldKey = "OPERATOR_CALL_SIGN"
	// FOpName is the operator name field.  It gets set by code that
	// creates a new outgoing message, or by code receiving a message.
	FOpName FieldKey = "OPERATOR_NAME"
	// FTacCall is the tactical call sign field.  It gets set by code that
	// creates a new outgoing message.
	FTacCall FieldKey = "TACTICAL_CALL_SIGN"
	// FTacName is the tactical name field.  It gets set by code that
	// creates a new outgoing message.
	FTacName FieldKey = "TACTICAL_NAME"
	// FOpDate is the transmission date (for outgoing messages) or the
	// reception date (for incoming messages).  It gets set when a message
	// is sent or received.  It is equivalent to the date part of FSent or
	// FReceived, but encoded in MM/DD/YYYY format.
	FOpDate FieldKey = "OPERATOR_DATE"
	// FOpTime is the transmission time (for outgoing messages) or the
	// reception date (for incoming messages).  It gets set when a message
	// is sent or received.  It is equivalent to the time part of FSent or
	// FReceived, but encoded in HH:MM format.
	FOpTime FieldKey = "OPERATOR_TIME"
)
