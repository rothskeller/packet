// Package message handles encoding and decoding (but not interpretation) of
// messages and their constituent parts.  The structure of a Message looks like
// this:
//
//	Message {
//	    Headers
//	    MType
//	    Subject
//	    Payload {
//	        Outpost Flags
//	        Body
//	    }
//	}
//
// A Message represents a packet message with all associated detail.  There are
// four implementations of the Message interface, representing different states
// of message transmission and the different message headers associated with
// each.  They are DraftMessage, SentMessage, ReceivedMessage, and
// JustReceivedMessage.
//
// All Message implementations contain an MType, short for message type, which
// provides the interpretation of the message and supports actions on it.  All
// messages of the same type share the same MType object; its methods take the
// specific Message to act on as a parameter.  There are three message types
// defined in this package group: message.PlainMessage, receipt.DeliveryReceipt,
// and receipt.ReadReceipt.  Additional message types can be defined elsewhere
// and registered with this package.
//
// All Message implementations contain a subject.Subject, which represents the
// subject line of the message.  There is one implementations of the Subject
// interface defined in this package group:  subject.PlainSubject.  Additional
// implementations can be defined elsewhere.
//
// All Message implementations contain a payload.Payload, which handles content
// transfer encoding wrappers around the message body.  There is one
// implementation of the Payload interface defined in this package group:
// payload.OutpostPayload, which handles !XXX! Outpost body flags as well as
// base64 and quoted-printable transfer encodings.  Additional implementations
// can be defined elsewhere and registered with the payload package.
//
// All Payload implementations contain a body.Body, which represents the actual
// body of the message after all wrappers have been removed.  There are three
// implementations of the Body interface defined in this package:
// body.PlainBody, receipt.DeliveryReceiptBody, and receipt.ReadReceiptBody.
// Additional implementations can be defined elsewhere and registered with the
// body package.
//
// There are three supported ways to create a Message object:
//  1. Call message.Read with the filename of a file containing a message.
//     This will return a DraftMessage, SentMessage, or ReceivedMessage, as
//     appropriate.
//  2. Call message.NewJustReceivedMessage with an encoded message string,
//     just retrieved from a BBS, and details of how it was retrieved.
//  3. Call message.NewDraftMessage with the details of a draft message.
//
// Various operations are possible on a Message object:
//   - To write a message to a file, pass it to message.Write.
//   - To validate the correctness of a message, call its Validate method.
package message
