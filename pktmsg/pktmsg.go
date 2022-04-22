// Package pktmsg handles encoding and decoding packet messages.  It understands
// RFC-4155 and RFC-5322 email encoding, SCCo-standard subject line encoding,
// PackItForms form encoding, and Outpost-specific feature encodings.
//
// This package exports two parallel hierarchies of message-related types.  The
// types starting with "Rx" are for messages that have been received and are
// being parsed.  The types starting with "Tx" are for messages that are being
// built and will be sent.
//
//     RxBase                    TxBase                     (base.go)
//       RxDeliveryReceipt         TxDeliveryReceipt        (receipts.go)
//       RxReadReceipt             TxReadReceipt            (receipts.go)
//       RxMessage                 TxMessage                (message.go)
//         RxForm                    TxForm                 (form.go)
//           RxICS213Form              TxICS213Form         (ics213.go)
//           RxSCCoForm                TxSCCoForm           (sccoform.go)
//             RxAHFacStatForm           TxAHFacStatForm    (ahfacstat.go)
//             RxEOC213RRForm            TxEOC213RRForm     (eoc213rr.go)
//             RxMuniStatForm            TxMuniStatForm     (munistat.go)
//             RxSheltStatForm           TxSheltStatForm    (sheltstat.go)
//
// All message types embed either RxBase or TxBase.  Messages that can't be
// parsed or have no return address are represented as a bare RxBase.
//
// Message types with human content also embed RxMessage or TxMessage.  Plain
// text messages are represented by a bare RxMessage or TxMessage.
//
// Messages containing a PackItForms-encoded form also embed RxForm or TxForm.
// If the form is not one that has its own message type it is represented by a
// bare RxForm or TxForm.
//
// Forms with standard SCCo header and footer fields also embed RxSCCoForm or
// TxSCCoForm.
//
// Received message objects (the "Rx" types) are created by calling ParseMessage
// on the encoded message text.  Outgoing messages are built by creating a
// message object of the appropriate type, filling it its fields, and then
// calling its Encode() method.
package pktmsg

// ParseMessage parses an encoded message and returns an object representing it.
func ParseMessage(rawmsg string) ParsedMessage {
	var base = parseRxBase(rawmsg)
	if dr := parseRxDeliveryReceipt(base); dr != nil {
		return dr
	}
	if rr := parseRxReadReceipt(base); rr != nil {
		return rr
	}
	var msg = parseRxMessage(base)
	if msg == nil {
		return base
	}
	var form = parseRxForm(msg)
	if form == nil {
		return msg
	}
	if ics213 := parseRxICS213Form(form); ics213 != nil {
		return ics213
	}
	if ahfacstat := parseRxAHFacStatForm(form); ahfacstat != nil {
		return ahfacstat
	}
	if eoc213rr := parseRxEOC213RRForm(form); eoc213rr != nil {
		return eoc213rr
	}
	if munistat := parseRxMuniStatForm(form); munistat != nil {
		return munistat
	}
	if racesmar := parseRxRACESMARForm(form); racesmar != nil {
		return racesmar
	}
	if sheltstat := parseRxSheltStatForm(form); sheltstat != nil {
		return sheltstat
	}
	return form
}
