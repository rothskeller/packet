package pktmsg

// This file defines TxMessage and RxMessage.

import (
	"fmt"
	"regexp"
	"strings"
)

// TxMessage is the foundation for all outgoing messages with human content
// ("real" messages, as opposed to things like automatic receipts).
type TxMessage struct {
	TxBase
	// MessageNumber is the message number to be embedded in the message's
	// Subject line.
	MessageNumber string
	// HandlingOrder is the handling order to be embedded in the message's
	// Subject line.
	HandlingOrder HandlingOrder
	// FormName is the form name to be embedded in the message's Subject
	// line.  It may be empty for plain text messages.
	FormName string
	// Subject is the message subject, to be embedded in the message's
	// Subject line (which also contains the message number, handling order,
	// and possibly form name).
	Subject string
}

// Encode returns the encoded subject line and body of the message.
func (m *TxMessage) Encode() (subject, body string, err error) {
	if m.MessageNumber == "" || m.HandlingOrder == 0 || m.Subject == "" {
		return "", "", ErrIncomplete
	}
	if m.FormName != "" {
		m.SubjectLine = fmt.Sprintf("%s_%s_%s_%s", m.MessageNumber, m.HandlingOrder.Code(), m.FormName, m.Subject)
	} else {
		m.SubjectLine = fmt.Sprintf("%s_%s_%s", m.MessageNumber, m.HandlingOrder.Code(), m.Subject)
	}
	m.OutpostUrgent = m.HandlingOrder == HandlingImmediate
	return m.TxBase.Encode()
}

//------------------------------------------------------------------------------

// RxMessage is the foundation for all received messages with human content
// ("real" messages, as opposed to things like bounce messages and automatic
// receipts).
type RxMessage struct {
	RxBase
	// FromCallSign is the call sign extracted from the return address, if
	// any.
	FromCallSign string
	// FromBBS is the BBS name extracted from the return address, if any.
	FromBBS string
	// MessageNumber is the message number embedded in the message's Subject
	// line.
	MessageNumber string
	// Severity is the message severity embedded in the message's Subject
	// line, if any.
	Severity MessageSeverity
	// SeverityCode is the code for the Severity.
	SeverityCode string
	// HandlingOrder is the handling order embedded in the message's Subject
	// line.
	HandlingOrder HandlingOrder
	// HandlingOrderCode is the code for the HandlingOrder.
	HandlingOrderCode string
	// FormName is the form name embedded in the message's Subject line, if
	// any.
	FormName string
	// Subject is the message subject, i.e., the remainder of the message's
	// Subject line after removal of any message number, severity, handling
	// order, and/or form name).
	Subject string
}

// Message returns a pointer to the RxMessage portion of a message object.  It
// can be used to reach fields of the RxMessage object that are occluded by
// types that embed RxMessage.
func (m *RxMessage) Message() *RxMessage { return m }

// TypeCode returns the machine-readable code for the message type.
func (*RxMessage) TypeCode() string { return "plain" }

// TypeName returns the human-reading name of the message type.
func (*RxMessage) TypeName() string { return "plain text message" }

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (*RxMessage) TypeArticle() string { return "a" }

// fromCallSignRE extracts the fromCallSign from the envelopeFrom.  It looks for
// a call sign at the start of the string, followed either by an @ or the end of
// the string.  It is not case-sensitive.  The substring returned is the call
// sign.
var fromCallSignRE = regexp.MustCompile(`(?i)^([AKNW][A-Z]?[0-9][A-Z]{1,3})(?:@|$)`)

// fromBBSRE extracts the fromBBS from the envelopeFrom.  It looks for an @,
// followed by a call sign, optionally followed by ".ampr.org", at the end of
// the string.  It is not case-sensitive.  The substring returned is the call
// sign (i.e., the BBS name).
var fromBBSRE = regexp.MustCompile(`(?i)@([AKNW][A-Z]?[0-9][A-Z]{1,3})(?:\.ampr\.org)?$`)

// subjectLineRE extracts the message number, severity code, handling order
// code, form name, and subject from the Subject line, assuming it's properly
// formatted.
var subjectLineRE = regexp.MustCompile(`(?i)^([A-Z0-9]+-?[0-9]+[A-Z]?)_(?:([A-Z])/)?([A-Z])_(?:([^_\s]+)_)?([^_\s]+(?:\s.*|$))`)

// parseRxMessage examines an RxBase to see if it is a human-content message,
// and if so, wraps it in an RxMessage and returns it.  If it is not, it returns
// nil.
func parseRxMessage(b *RxBase) *RxMessage {
	var m RxMessage

	if b.ParseError != "" || b.ReturnAddress == "" {
		return nil // Not a human-content message.
	}
	m.RxBase = *b
	// Extract the call sign and BBS name from the return address, if they
	// are there.
	if match := fromCallSignRE.FindStringSubmatch(m.ReturnAddress); match != nil {
		m.FromCallSign = strings.ToUpper(match[1])
	}
	if match := fromBBSRE.FindStringSubmatch(m.ReturnAddress); match != nil {
		m.FromBBS = strings.ToUpper(match[1])
	}
	// If the subject line is properly formatted, extract the pieces from
	// it.  Otherwise, consider the whole thing to be the subject.
	if match := subjectLineRE.FindStringSubmatch(m.SubjectLine); match != nil {
		m.MessageNumber = match[1]
		m.SeverityCode = match[2]
		m.Severity, _ = ParseSeverity(match[2])
		m.HandlingOrderCode = match[3]
		m.HandlingOrder, _ = ParseHandlingOrder(match[3])
		m.FormName = match[4]
		m.Subject = match[5]
	} else {
		m.Subject = m.SubjectLine
	}
	return &m
}
