package pktmsg

// This file defines TxDeliveryReceipt, RxDeliveryReceipt, TxReadReceipt, and
// RxReadReceipt.

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// A TxDeliveryReceipt is an outgoing delivery receipt message.
type TxDeliveryReceipt struct {
	TxBase
	DeliveredTo      string
	DeliveredSubject string
	LocalMessageID   string
	DeliveredTime    time.Time
}

// Encode returns the encoded subject line and body of the message.
func (dr *TxDeliveryReceipt) Encode() (subject, body string, err error) {
	var fdt = dr.DeliveredTime.Format("2006-01-02 15:04:05")

	if dr.DeliveredTo == "" || dr.DeliveredSubject == "" || dr.LocalMessageID == "" || dr.DeliveredTime.IsZero() {
		return "", "", ErrIncomplete
	}
	dr.SubjectLine = "DELIVERED: " + dr.DeliveredSubject
	dr.Body = fmt.Sprintf(
		"!LMI!%s!DR!%s\nYour Message\nTo: %s\nSubject: %s\nwas delivered on %s\nRecipient's Local Message ID: %s\n",
		dr.LocalMessageID, fdt, dr.DeliveredTo, dr.DeliveredSubject, fdt, dr.LocalMessageID,
	)
	return dr.TxBase.Encode()
}

//------------------------------------------------------------------------------

// An RxDeliveryReceipt is a received delivery receipt message.
type RxDeliveryReceipt struct {
	RxBase
	DeliveredSubject string
	LocalMessageID   string
	DeliveredTime    time.Time
}

// TypeCode returns the machine-readable code for the message type.
func (*RxDeliveryReceipt) TypeCode() string { return "DELIVERED" }

// TypeName returns the human-reading name of the message type.
func (*RxDeliveryReceipt) TypeName() string { return "delivery receipt" }

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (*RxDeliveryReceipt) TypeArticle() string { return "a" }

// deliveryReceiptRE matches the first line of a delivery receipt message.  Its
// substrings are the local message ID and the delivery time.
var deliveryReceiptRE = regexp.MustCompile(`^!LMI!([^!]+)!DR!(20\d\d-\d\d-\d\d \d\d:\d\d:\d\d)\n`)

// parseRxDeliveryReceipt examines an RxBase to see if it is a delivery receipt,
// and if so, wraps it in an RxDeliveryReceipt and returns it.  If it is not, it
// returns nil.
func parseRxDeliveryReceipt(b *RxBase) *RxDeliveryReceipt {
	var dr RxDeliveryReceipt

	if b.ParseError != "" || b.ReturnAddress == "" || b.NotPlainText {
		return nil
	}
	if sl := b.SubjectLine; strings.HasPrefix(sl, "DELIVERED: ") {
		dr.DeliveredSubject = sl[11:]
	} else {
		return nil
	}
	if match := deliveryReceiptRE.FindStringSubmatch(b.Body); match != nil {
		dr.LocalMessageID = match[1]
		dr.DeliveredTime, _ = time.ParseInLocation("2006-01-02 15:04:05", match[2], time.Local)
	} else {
		return nil
	}
	dr.RxBase = *b
	return &dr
}

//------------------------------------------------------------------------------

// A TxReadReceipt is an outgoing read receipt message.
type TxReadReceipt struct {
	TxBase
	ReadTo      string
	ReadSubject string
	ReadTime    time.Time
}

// Encode returns the encoded subject line and body of the message.
func (rr *TxReadReceipt) Encode() (subject, body string, err error) {
	var frt = rr.ReadTime.Format("2006-01-02 15:04:05")

	if rr.ReadTo == "" || rr.ReadSubject == "" || rr.ReadTime.IsZero() {
		return "", "", ErrIncomplete
	}
	rr.SubjectLine = "READ: " + rr.ReadSubject
	rr.Body = fmt.Sprintf("!RR!%s\nYour Message\n\nTo: %s\nSubject: %s\n\nwas read on %s\n",
		frt, rr.ReadTo, rr.ReadSubject, frt,
	)
	return rr.TxBase.Encode()
}

//------------------------------------------------------------------------------

// An RxReadReceipt is a received read receipt message.
type RxReadReceipt struct {
	RxBase
	ReadSubject string
	ReadTime    time.Time
}

// TypeCode returns the machine-readable code for the message type.
func (*RxReadReceipt) TypeCode() string { return "READ" }

// TypeName returns the human-reading name of the message type.
func (*RxReadReceipt) TypeName() string { return "read receipt" }

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (*RxReadReceipt) TypeArticle() string { return "a" }

// readReceiptRE matches the first line of a read receipt message.  Its
// substring is the read time.
var readReceiptRE = regexp.MustCompile(`^!RR!(20\d\d-\d\d-\d\d \d\d:\d\d:\d\d)\n`)

// parseRxReadReceipt examines an RxBase to see if it is a delivery receipt,
// and if so, wraps it in an RxReadReceipt and returns it.  If it is not, it
// returns nil.
func parseRxReadReceipt(b *RxBase) *RxReadReceipt {
	var rr RxReadReceipt

	if b.ParseError != "" || b.ReturnAddress == "" || b.NotPlainText {
		return nil
	}
	if sl := b.SubjectLine; strings.HasPrefix(sl, "READ: ") {
		rr.ReadSubject = sl[6:]
	} else {
		return nil
	}
	if match := readReceiptRE.FindStringSubmatch(b.Body); match != nil {
		rr.ReadTime, _ = time.ParseInLocation("2006-01-02 15:04:05", match[1], time.Local)
	} else {
		return nil
	}
	rr.RxBase = *b
	return &rr
}
