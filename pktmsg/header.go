package pktmsg

import (
	"net/mail"
	"regexp"
	"strings"
	"time"
)

// ReturnAddress returns the return address of the message, taken from the
// envelope From line, or failing that the Return-Path header, the Reply-To
// header, the Sender header, or the From header.
func (msg *Message) ReturnAddress() string {
	var (
		header []string
		ok     bool
	)
	if msg.EnvelopeAddress != "" {
		return msg.EnvelopeAddress
	}

	if header, ok = msg.Header["Return-Path"]; !ok {
		if header, ok = msg.Header["Reply-To"]; !ok {
			if header, ok = msg.Header["Sender"]; !ok {
				header = msg.Header["From"]
			}
		}
	}
	if len(header) != 0 {
		if al, err := mail.ParseAddressList(header[0]); err == nil && len(al) != 0 {
			return al[0].Address
		}
	}
	return ""
}

// To returns all of the destination addresses of the message, combining the
// To:, Cc:, and Bcc: headers.  Comments on the addresses are removed.
func (msg *Message) To() (to []string) {
	for _, hdr := range []string{"To", "Cc", "Bcc"} {
		for _, list := range msg.Header[hdr] {
			addrs, _ := mail.ParseAddressList(list)
			for _, addr := range addrs {
				to = append(to, addr.Address)
			}
		}
	}
	return to
}

// Date returns the date the message was received, taken from the envelope From
// line, or failing that the Received header, or failing that the Date header.
func (msg *Message) Date() (date time.Time) {
	if !msg.EnvelopeDate.IsZero() {
		return msg.EnvelopeDate
	}
	received := msg.Header.Get("Received")
	if semi := strings.LastIndexByte(received, ';'); semi >= 0 {
		date, _ = mail.ParseDate(received[semi+1:])
	}
	if date.IsZero() {
		date, _ = mail.ParseDate(msg.Header.Get("Date"))
	}
	return date
}

// fromCallSignRE extracts the fromCallSign from the return address.  It looks
// for a call sign at the start of the string, followed either by an @ or the
// end of the string.  It is not case-sensitive.  The substring returned is the
// call sign.
var fromCallSignRE = regexp.MustCompile(`(?i)^([AKNW][A-Z]?[0-9][A-Z]{1,3})(?:@|$)`)

// FromCallSign extracts a call sign from the message return address.
func (msg *Message) FromCallSign() string {
	if match := fromCallSignRE.FindStringSubmatch(msg.ReturnAddress()); match != nil {
		return strings.ToUpper(match[1])
	}
	return ""
}

// fromBBSRE extracts the fromBBS from the return address.  It looks for an @,
// followed by a call sign, optionally followed by ".ampr.org", at the end of
// the string.  It is not case-sensitive.  The substring returned is the call
// sign (i.e., the BBS name).
var fromBBSRE = regexp.MustCompile(`(?i)@([AKNW][A-Z]?[0-9][A-Z]{1,3})(?:\.ampr\.org)?$`)

// FromBBS extracts the sending BBS from the message return address.
func (msg *Message) FromBBS() string {
	if match := fromBBSRE.FindStringSubmatch(msg.ReturnAddress()); match != nil {
		return strings.ToUpper(match[1])
	}
	return ""
}
