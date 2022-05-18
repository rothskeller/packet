package readrcpt

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

const (
	readReceiptTag  = "READ"
	readReceiptType = "read-receipt.html"
)

func init() {
	xscmsg.RegisterType(Create, Recognize)
}

// Create creates a new message of the type identified by the supplied tag.  If
// the tag is not recognized by this package, Create returns nil.
func Create(tag string) xscmsg.XSCMessage {
	if tag == readReceiptTag {
		return new(ReadReceipt)
	}
	return nil
}

// readReceiptRE matches the first lines of a read receipt message.  Its
// substrings are the read time and the To address.
var readReceiptRE = regexp.MustCompile(`^!RR!(.+)\n.*\n\nTo: (.+)`)

// Recognize examines the supplied Message to see if it is of the type defined
// in this package.  If so, it returns the appropriate XSCMessage implementation
// wrapping it.  If not, it returns nil.
func Recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.XSCMessage {
	var rr = ReadReceipt{msg: msg}
	if form != nil {
		if form.FormType == readReceiptType {
			rr.ReadTo = form.Get("ReadTo")
			rr.ReadSubject = form.Get("ReadSubject")
			rr.ReadTime = form.Get("ReadTime")
			return &rr
		}
		return nil
	}
	if subject := msg.Header.Get("Subject"); strings.HasPrefix(subject, "READ: ") {
		rr.ReadSubject = subject[6:]
	} else {
		return nil
	}
	if match := readReceiptRE.FindStringSubmatch(msg.Body); match != nil {
		rr.ReadTo = match[2]
		rr.ReadTime = match[1]
	}
	return &rr
}

// ReadReceipt is a read receipt message, acknowledging read of a
// previous message.
type ReadReceipt struct {
	msg         *pktmsg.Message
	ReadTo      string
	ReadSubject string
	ReadTime    string
}

// TypeTag returns the tag string used to identify the message type.
func (rr *ReadReceipt) TypeTag() string { return readReceiptTag }

// TypeName returns the English name of the message type.
func (rr *ReadReceipt) TypeName() string { return "read receipt" }

// TypeArticle returns the indefinite article ("a" or "an") to be used
// preceding the TypeName, in a sentence that needs one.
func (rr *ReadReceipt) TypeArticle() string { return "a" }

// Validate ensures that the contents of the message are correct.  It returns a
// list of problems, which is empty if the message is fine.  If strict is true,
// the message must be exactly correct; otherwise, some trivial issues are
// corrected and not reported.
func (rr *ReadReceipt) Validate(strict bool) (problems []string) {
	if rr.ReadTo == "" {
		problems = append(problems, "ReadTo is not specified")
	}
	if rr.ReadSubject == "" {
		problems = append(problems, "ReadSubject is not specified")
	}
	if rr.ReadTime == "" {
		problems = append(problems, "ReadTime is not specified")
	}
	return problems
}

// Message returns the encoded message.  If human is true, it is encoded for
// human reading or editing; if false, it is encoded for transmission.  If the
// XSCMessage was originally created by a call to Recognize, the Message
// structure passed to it is updated and reused; otherwise, a new Message
// structure is created and filled in.
func (rr *ReadReceipt) Message(human bool) *pktmsg.Message {
	if rr.msg == nil {
		rr.msg = pktmsg.New()
	}
	rr.msg.Header.Set("Subject", "READ: "+rr.ReadSubject)
	if human {
		form := &pktmsg.Form{
			FormType: readReceiptType, PIFOVersion: "0", FormVersion: "0",
			Fields: []pktmsg.FormField{
				{Tag: "ReadTo", Value: rr.ReadTo},
				{Tag: "ReadSubject", Value: rr.ReadSubject},
				{Tag: "ReadTime", Value: rr.ReadTime},
			},
		}
		rr.msg.Body = form.Encode(nil, nil, true)
	} else {
		rr.msg.Body = fmt.Sprintf("!RR!%s\nYour Message\n\nTo: %s\nSubject: %s\n\nwas read on %s\n",
			rr.ReadTime, rr.ReadTo, rr.ReadSubject, rr.ReadTime,
		)
	}
	return rr.msg
}

// OriginNumber returns the origin message number of the message, if any.
func (*ReadReceipt) OriginNumber() string { return "" }

// SetOriginNumber sets the origin message number of the message, if the message
// type supports that.
func (*ReadReceipt) SetOriginNumber(string) {}

// DestinationNumber returns the destination message number of the message, if
// any.
func (*ReadReceipt) DestinationNumber() string { return "" }

// SetDestinationNumber sets the destination message number of the message, if
// the message type supports that.
func (*ReadReceipt) SetDestinationNumber(string) {}

// HandlingOrder returns the message handling order, if any, in both string and
// parsed forms.
func (*ReadReceipt) HandlingOrder() (string, xscmsg.HandlingOrder) { return "", 0 }

// SetHandlingOrder sets the message handling order, if the message type
// supports that.
func (*ReadReceipt) SetHandlingOrder(xscmsg.HandlingOrder) {}

// Operator returns the operator name and call sign from the message, if any.
func (*ReadReceipt) Operator() (string, string) { return "", "" }

// SetOperator sets the operator name and call sign in the message, if the
// message type supports that.
func (*ReadReceipt) SetOperator(string, string) {}

// ActionTime returns the time that the outgoing message was sent or the
// incoming message was received, if the message has that information.  It
// returns them in string form, and also in parsed form if they are in standard
// format.
func (*ReadReceipt) ActionTime() (string, string, time.Time) { return "", "", time.Time{} }

// SetActionTime sets the time that the outgoing message was sent or the income
// message was received, if the message type supports that.
func (*ReadReceipt) SetActionTime(time.Time) {}
