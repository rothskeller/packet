package delivrcpt

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

const (
	deliveryReceiptTag  = "DELIVERED"
	deliveryReceiptType = "delivery-receipt.html"
)

func init() {
	xscmsg.RegisterType(Create, Recognize)
}

// Create creates a new message of the type identified by the supplied tag.  If
// the tag is not recognized by this package, Create returns nil.
func Create(tag string) xscmsg.XSCMessage {
	if tag == deliveryReceiptTag {
		return new(DeliveryReceipt)
	}
	return nil
}

// deliveryReceiptRE matches the first lines of a delivery receipt message.  Its
// substrings are the local message ID, the delivery time, and the To address.
var deliveryReceiptRE = regexp.MustCompile(`^!LMI!([^!]+)!DR!(.+)\n.*\nTo: (.+)`)

// Recognize examines the supplied Message to see if it is of the type defined
// in this package.  If so, it returns the appropriate XSCMessage implementation
// wrapping it.  If not, it returns nil.
func Recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.XSCMessage {
	var dr = DeliveryReceipt{msg: msg}
	if form != nil {
		if form.FormType == deliveryReceiptType {
			dr.DeliveredTo = form.Get("DeliveredTo")
			dr.DeliveredSubject = form.Get("DeliveredSubject")
			dr.LocalMessageID = form.Get("LocalMessageID")
			dr.DeliveredTime = form.Get("DeliveredTime")
			return &dr
		}
		return nil
	}
	if subject := msg.Header.Get("Subject"); strings.HasPrefix(subject, "DELIVERED: ") {
		dr.DeliveredSubject = subject[11:]
	} else {
		return nil
	}
	if match := deliveryReceiptRE.FindStringSubmatch(msg.Body); match != nil {
		dr.DeliveredTo = match[3]
		dr.LocalMessageID = match[1]
		dr.DeliveredTime = match[2]
	}
	return &dr
}

// DeliveryReceipt is a delivery receipt message, acknowledging delivery of a
// previous message.
type DeliveryReceipt struct {
	msg              *pktmsg.Message
	DeliveredTo      string
	DeliveredSubject string
	LocalMessageID   string
	DeliveredTime    string
}

// TypeTag returns the tag string used to identify the message type.
func (dr *DeliveryReceipt) TypeTag() string { return deliveryReceiptTag }

// TypeName returns the English name of the message type.
func (dr *DeliveryReceipt) TypeName() string { return "delivery receipt" }

// TypeArticle returns the indefinite article ("a" or "an") to be used
// preceding the TypeName, in a sentence that needs one.
func (dr *DeliveryReceipt) TypeArticle() string { return "a" }

// Validate ensures that the contents of the message are correct.  It returns a
// list of problems, which is empty if the message is fine.  If strict is true,
// the message must be exactly correct; otherwise, some trivial issues are
// corrected and not reported.
func (dr *DeliveryReceipt) Validate(strict bool) (problems []string) {
	if dr.DeliveredTo == "" {
		problems = append(problems, "DeliveredTo is not specified")
	}
	if dr.DeliveredSubject == "" {
		problems = append(problems, "DeliveredSubject is not specified")
	}
	if dr.DeliveredTime == "" {
		problems = append(problems, "DeliveredTime is not specified")
	}
	if dr.LocalMessageID == "" {
		problems = append(problems, "LocalMessageID is not specified")
	}
	return problems
}

// Message returns the encoded message.  If human is true, it is encoded for
// human reading or editing; if false, it is encoded for transmission.  If the
// XSCMessage was originally created by a call to Recognize, the Message
// structure passed to it is updated and reused; otherwise, a new Message
// structure is created and filled in.
func (dr *DeliveryReceipt) Message(human bool) *pktmsg.Message {
	if dr.msg == nil {
		dr.msg = pktmsg.New()
	}
	dr.msg.Header.Set("Subject", "DELIVERED: "+dr.DeliveredSubject)
	if human {
		form := &pktmsg.Form{
			FormType: deliveryReceiptType, PIFOVersion: "0", FormVersion: "0",
			Fields: []pktmsg.FormField{
				{Tag: "DeliveredTo", Value: dr.DeliveredTo},
				{Tag: "DeliveredSubject", Value: dr.DeliveredSubject},
				{Tag: "DeliveredTime", Value: dr.DeliveredTime},
				{Tag: "LocalMessageID", Value: dr.LocalMessageID},
			},
		}
		dr.msg.Body = form.Encode(nil, nil, true)
	} else {
		dr.msg.Body = fmt.Sprintf(
			"!LMI!%s!DR!%s\nYour Message\nTo: %s\nSubject: %s\nwas delivered on %s\nRecipient's Local Message ID: %s\n",
			dr.LocalMessageID, dr.DeliveredTime, dr.DeliveredTo, dr.DeliveredSubject, dr.DeliveredTime, dr.LocalMessageID,
		)
	}
	return dr.msg
}

// OriginNumber returns the origin message number of the message, if any.
func (*DeliveryReceipt) OriginNumber() string { return "" }

// SetOriginNumber sets the origin message number of the message, if the message
// type supports that.
func (*DeliveryReceipt) SetOriginNumber(string) {}

// DestinationNumber returns the destination message number of the message, if
// any.
func (*DeliveryReceipt) DestinationNumber() string { return "" }

// SetDestinationNumber sets the destination message number of the message, if
// the message type supports that.
func (*DeliveryReceipt) SetDestinationNumber(string) {}

// HandlingOrder returns the message handling order, if any, in both string and
// parsed forms.
func (*DeliveryReceipt) HandlingOrder() (string, xscmsg.HandlingOrder) { return "", 0 }

// SetHandlingOrder sets the message handling order, if the message type
// supports that.
func (*DeliveryReceipt) SetHandlingOrder(xscmsg.HandlingOrder) {}

// Operator returns the operator name and call sign from the message, if any.
func (*DeliveryReceipt) Operator() (string, string) { return "", "" }

// SetOperator sets the operator name and call sign in the message, if the
// message type supports that.
func (*DeliveryReceipt) SetOperator(string, string) {}

// ActionTime returns the time that the outgoing message was sent or the
// incoming message was received, if the message has that information.  It
// returns them in string form, and also in parsed form if they are in standard
// format.
func (*DeliveryReceipt) ActionTime() (string, string, time.Time) { return "", "", time.Time{} }

// SetActionTime sets the time that the outgoing message was sent or the income
// message was received, if the message type supports that.
func (*DeliveryReceipt) SetActionTime(time.Time) {}
