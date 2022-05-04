package checkin

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

var checkin = &xscform.FormDefinition{
	HTML:              "check-in.html", // fake, not really
	Tag:               "Check-In",
	Name:              "check-in message",
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
	if tag == checkin.Tag {
		return &CheckIn{xscform.CreateForm(checkin)}
	}
	return nil
}

var checkInRE = regexp.MustCompile(`(?i)^Check-In\s+([A-Z][A-Z0-9]{2,5})\s*,(.*)(?:\n([AKNW][A-Z0-9]{2,5})\s*,(.*))?`)

// Recognize examines the supplied Message to see if it is of the type defined
// in this package.  If so, it returns the appropriate XSCMessage implementation
// wrapping it.  If not, it returns nil.
func Recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.XSCMessage {
	if xf := xscform.RecognizeForm(checkin, msg, form); xf != nil {
		return &CheckIn{xf}
	}
	if subject := xscmsg.ParseSubject(msg.Header.Get("Subject")); subject != nil && subject.FormTag == "" {
		if strings.HasPrefix(strings.ToLower(subject.Subject), "check-in ") {
			var ci = &CheckIn{xscform.CreateForm(checkin)}
			ci.SetOriginNumber(subject.MessageNumber)
			if match := checkInRE.FindStringSubmatch(msg.Body); match != nil {
				if match[3] != "" {
					ci.SetTactical(strings.TrimSpace(match[2]), match[1])
					ci.SetOperator(strings.TrimSpace(match[4]), match[3])
				} else {
					ci.SetOperator(strings.TrimSpace(match[2]), match[1])
				}
			}
			return ci
		}
	}
	return nil
}

// CheckIn is a check-in message.
type CheckIn struct {
	form *xscform.XSCForm
}

// TypeTag returns the tag string used to identify the message type.
func (ci *CheckIn) TypeTag() string { return checkin.Tag }

// TypeName returns the English name of the message type.  It is a noun
// phrase in prose case, such as "foo message" or "bar form".
func (ci *CheckIn) TypeName() string { return checkin.Name }

// TypeArticle returns the indefinite article ("a" or "an") to be used
// preceding the TypeName, in a sentence that needs one.
func (ci *CheckIn) TypeArticle() string { return checkin.Article }

// Validate ensures that the contents of the message are correct.  It returns a
// list of problems, which is empty if the message is fine.  If strict is true,
// the message must be exactly correct; otherwise, some trivial issues are
// corrected and not reported.
func (ci *CheckIn) Validate(strict bool) (problems []string) {
	problems = ci.form.Validate(strict)
	if (ci.Get("TacCall") == "") != (ci.Get("TacName") == "") {
		problems = append(problems, "must set both or neither of TacCall and TacName")
	}
	return problems
}

// Message returns the encoded message.  If human is true, it is encoded for
// human reading or editing; if false, it is encoded for transmission.  If the
// XSCMessage was originally created by a call to Recognize, the Message
// structure passed to it is updated and reused; otherwise, a new Message
// structure is created and filled in.
func (ci *CheckIn) Message(human bool) (msg *pktmsg.Message) {
	opname, opcall := ci.Operator()
	tacname, taccall := ci.Tactical()
	if human {
		msg = ci.form.Message(human)
	} else {
		msg = pktmsg.New()
		if taccall != "" {
			msg.Body = fmt.Sprintf("Check-In %s, %s\n%s, %s\n", taccall, tacname, opcall, opname)
		} else {
			msg.Body = fmt.Sprintf("Check-In %s, %s", opcall, opname)
		}
	}
	msg.Header.Set("Subject", ci.EncodeSubject())
	return msg
}

// EncodeSubject returns the encoded subject line of the message based on its
// contents.
func (ci *CheckIn) EncodeSubject() string {
	var priname, pricall string
	tacname, taccall := ci.Tactical()
	if taccall != "" {
		priname, pricall = tacname, taccall
	} else {
		priname, pricall = ci.Operator()
	}
	return xscmsg.EncodeSubject(ci.OriginNumber(), xscmsg.HandlingRoutine, "",
		fmt.Sprintf("Check-In %s, %s", pricall, priname))
}

// Get returns the value of the named form field.
func (ci *CheckIn) Get(s string) string { return ci.form.Get(s) }

// Set sets the value of the named form field
func (ci *CheckIn) Set(name, value string) { ci.form.Set(name, value) }

// OriginNumber returns the origin message number of the message, if any.
func (ci *CheckIn) OriginNumber() string { return ci.form.OriginNumber() }

// SetOriginNumber sets the origin message number of the message, if the message
// type supports that.
func (ci *CheckIn) SetOriginNumber(s string) { ci.form.SetOriginNumber(s) }

// Operator returns the operator name and call sign from the message.
func (ci *CheckIn) Operator() (name, callSign string) { return ci.form.Operator() }

// SetOperator sets the operator name and call sign in the message.
func (ci *CheckIn) SetOperator(name, callSign string) { ci.form.SetOperator(name, callSign) }

// Tactical returns the tactical name and call sign from the message.
func (ci *CheckIn) Tactical() (name, callSign string) { return ci.Get("TacName"), ci.Get("TacCall") }

// SetTactical sets the tactical name and call sign in the message.
func (ci *CheckIn) SetTactical(name, callSign string) {
	ci.Set("TacName", name)
	ci.Set("TacCall", callSign)
}
