package pktmsg

import (
	"net/mail"
	"strings"
	"time"
)

// now is a function that returns the current time.  It can be overridden by
// tests.
var now = time.Now

// baseRetrieved extends baseRx and adds non-persistent fields related to the
// retrieval of the received message.
type baseRetrieved struct {
	*baseRx
	returnAddr   returnAddrField
	bbsRxDate    bbsRxDateField
	notPlainText notPlainTextField
}

// parseBaseRetrieved parses a message that was just retrieved from JNOS.
func parseBaseRetrieved(om *outpostMessage, bbs, area, envelope string, h mail.Header, notplain bool) Message {
	var rx = baseRx{outpostMessage: om}

	// Fill in the baseRx fields and create a baseRetrieved from it.
	rx.rxBBS = field(bbs)
	rx.rxArea = field(area)
	rx.rxDate = dateField{field(now().Format(time.RFC1123Z))}
	var m = baseRetrieved{baseRx: &rx}
	if notplain {
		m.notPlainText = field("true")
	}
	// Parse the envelope "From " line if any.
	var hadEnvelope bool
	if envelope != "" {
		hadEnvelope = true
		envelope = envelope[5:] // skip "From "
		if idx := strings.IndexByte(envelope, ' '); idx >= 0 {
			m.returnAddr = field(envelope[:idx])
			envelope = envelope[idx+1:]
		} else {
			m.returnAddr = field(envelope)
			envelope = ""
		}
	}
	if envelope != "" {
		// Looks like there's a timestamp on the envelope line.
		// RFC-4155 says it should be a ctime-style timestamp in UTC.
		// We'll try parsing it as that way, but we'll treat it as local
		// time because that's what JNOS BBSes do, and those are our
		// primary source of messages to parse.
		if t, err := time.ParseInLocation(time.ANSIC, envelope, time.Local); err == nil {
			m.bbsRxDate = dateField{field(t.Format(time.RFC1123Z))}
		}
	}
	// Compute the return address if there wasn't one on the From line.
	if !hadEnvelope {
		var line string
		if line = h.Get("Return-Path"); line == "" {
			if line = h.Get("Reply-To"); line == "" {
				if line = h.Get("Sender"); line == "" {
					line = h.Get("From")
				}
			}
		}
		// Most of those sources can have a name comment in the address,
		// which we don't want.  Also, From can have more than one
		// address in it, and we only want the first.
		if addrs, err := mail.ParseAddressList(line); err == nil && len(addrs) > 0 {
			m.returnAddr = field(addrs[0].Address)
		}
	}
	// If we didn't get a BBS Rx date from the envelope, get it from the
	// Received header.
	if m.bbsRxDate.Time().IsZero() {
		_, date, _ := strings.Cut(h.Get("Received"), ";")
		if t, err := mail.ParseDate(strings.TrimSpace(date)); err == nil {
			m.bbsRxDate = dateField{field(t.Format(time.RFC1123Z))}
		}
	}
	return &m
}

// Accessors for baseRetrieved fields.
func (m *baseRetrieved) BBSRxDate() BBSRxDateField {
	if m.bbsRxDate.Value() != "" {
		return m.bbsRxDate
	}
	return nil
}
func (m *baseRetrieved) NotPlainText() NotPlainTextField {
	if m.notPlainText.Value() != "" {
		return m.notPlainText
	}
	return nil
}
func (m *baseRetrieved) ReturnAddr() ReturnAddrField { return m.returnAddr }

// The BBSRxDateField holds the time the message was received by the BBS.  It is
// present only on instantly-received messages, and then only when the data are
// available.  It is not persisted in local storage.
type BBSRxDateField interface {
	Field
	// Time returns the time encoded in the field.
	Time() time.Time
}
type bbsRxDateField = dateField

// The NotPlainTextField indicates, by its presence, that the instantly-received
// message was not in plain text.  It is not persisted in local storage.
type NotPlainTextField interface{ Field }
type notPlainTextField = field

// The ReturnAddrField holds the return address of the message.  It is present
// only on instantly-received messages; it is not persisted in local storage.
type ReturnAddrField interface{ Field }
type returnAddrField = field
