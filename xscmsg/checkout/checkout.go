package checkout

import (
	"fmt"
	"regexp"
	"strings"

	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/xscmsg"
	"steve.rothskeller.net/packet/xscmsg/internal/xscform"
)

func init() {
	xscmsg.RegisterType(Create, Recognize)
}

var checkout = &xscform.FormDefinition{
	HTML:              "check-out.html", // fake, not really
	Tag:               "Check-Out",
	Name:              "check-out message",
	Article:           "a",
	Version:           "1.0",
	OriginNumberField: "MsgNo",
	OperatorNameField: "OpName",
	OperatorCallField: "OpCall",
	Fields: []*xscform.FieldDefinition{
		{
			Tag:         "MsgNo",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateMessageNumber},
		},
		{Tag: "TacCall"},
		{Tag: "TacName"},
		{
			Tag:         "OpCall",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired, xscform.ValidateCallSign},
		},
		{
			Tag:         "OpName",
			Validations: []xscform.ValidateFunc{xscform.ValidateRequired},
		},
	},
}

// Create creates a new message of the type identified by the supplied tag.  If
// the tag is not recognized by this package, Create returns nil.
func Create(tag string) xscmsg.XSCMessage {
	if tag == checkout.Tag {
		return &CheckOut{xscform.CreateForm(checkout)}
	}
	return nil
}

var checkOutRE = regexp.MustCompile(`(?i)^Check-Out\s+([A-Z][A-Z0-9]{2,5})\s*,(.*)(?:\n([AKNW][A-Z0-9]{2,5})\s*,(.*))?`)

// Recognize examines the supplied Message to see if it is of the type defined
// in this package.  If so, it returns the appropriate XSCMessage implementation
// wrapping it.  If not, it returns nil.
func Recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.XSCMessage {
	if xf := xscform.RecognizeForm(checkout, msg, form); xf != nil {
		return &CheckOut{xf}
	}
	if subject := xscmsg.ParseSubject(msg.Header.Get("Subject")); subject != nil && subject.FormTag == "" {
		if strings.HasPrefix(strings.ToLower(subject.Subject), "check-out ") {
			var co = &CheckOut{xscform.CreateForm(checkout)}
			co.SetOriginNumber(subject.MessageNumber)
			if match := checkOutRE.FindStringSubmatch(msg.Body); match != nil {
				if match[3] != "" {
					co.SetTactical(strings.TrimSpace(match[2]), match[1])
					co.SetOperator(strings.TrimSpace(match[4]), match[3])
				} else {
					co.SetOperator(strings.TrimSpace(match[2]), match[1])
				}
			}
			return co
		}
	}
	return nil
}

// CheckOut is a check-out message.
type CheckOut struct {
	form *xscform.XSCForm
}

// TypeTag returns the tag string used to identify the message type.
func (co *CheckOut) TypeTag() string { return checkout.Tag }

// TypeName returns the English name of the message type.  It is a noun
// phrase in prose case, such as "foo message" or "bar form".
func (co *CheckOut) TypeName() string { return checkout.Name }

// TypeArticle returns the indefinite article ("a" or "an") to be used
// preceding the TypeName, in a sentence that needs one.
func (co *CheckOut) TypeArticle() string { return checkout.Article }

// Validate ensures that the contents of the message are correct.  It returns a
// list of problems, which is empty if the message is fine.  If strict is true,
// the message must be exactly correct; otherwise, some trivial issues are
// corrected and not reported.
func (co *CheckOut) Validate(strict bool) (problems []string) {
	problems = co.form.Validate(strict)
	if (co.Get("TacCall") == "") != (co.Get("TacName") == "") {
		problems = append(problems, "must set both or neither of TacCall and TacName")
	}
	return problems
}

// Message returns the encoded message.  If human is true, it is encoded for
// human reading or editing; if false, it is encoded for transmission.  If the
// XSCMessage was originally created by a call to Recognize, the Message
// structure passed to it is updated and reused; otherwise, a new Message
// structure is created and filled in.
func (co *CheckOut) Message(human bool) (msg *pktmsg.Message) {
	opname, opcall := co.Operator()
	tacname, taccall := co.Tactical()
	if human {
		msg = co.form.Message(human)
	} else {
		msg = pktmsg.New()
		if taccall != "" {
			msg.Body = fmt.Sprintf("Check-Out %s, %s\n%s, %s\n", taccall, tacname, opcall, opname)
		} else {
			msg.Body = fmt.Sprintf("Check-Out %s, %s", opcall, opname)
		}
	}
	msg.Header.Set("Subject", co.EncodeSubject())
	return msg
}

// EncodeSubject returns the encoded subject line of the message based on its
// contents.
func (co *CheckOut) EncodeSubject() string {
	var priname, pricall string
	tacname, taccall := co.Tactical()
	if taccall != "" {
		priname, pricall = tacname, taccall
	} else {
		priname, pricall = co.Operator()
	}
	return xscmsg.EncodeSubject(co.OriginNumber(), xscmsg.HandlingRoutine, "",
		fmt.Sprintf("Check-Out %s, %s", pricall, priname))
}

// Get returns the value of the named form field.
func (co *CheckOut) Get(s string) string { return co.form.Get(s) }

// Set sets the value of the named form field
func (co *CheckOut) Set(name, value string) { co.form.Set(name, value) }

// OriginNumber returns the origin message number of the message, if any.
func (co *CheckOut) OriginNumber() string { return co.form.OriginNumber() }

// SetOriginNumber sets the origin message number of the message, if the message
// type supports that.
func (co *CheckOut) SetOriginNumber(s string) { co.form.SetOriginNumber(s) }

// Operator returns the operator name and call sign from the message.
func (co *CheckOut) Operator() (name, callSign string) { return co.form.Operator() }

// SetOperator sets the operator name and call sign in the message.
func (co *CheckOut) SetOperator(name, callSign string) { co.form.SetOperator(name, callSign) }

// Tactical returns the tactical name and call sign from the message.
func (co *CheckOut) Tactical() (name, callSign string) { return co.Get("TacName"), co.Get("TacCall") }

// SetTactical sets the tactical name and call sign in the message.
func (co *CheckOut) SetTactical(name, callSign string) {
	co.Set("TacName", name)
	co.Set("TacCall", callSign)
}
