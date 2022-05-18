// Package xscmsg handles recognition, parsing, validation, and encoding of all
// of the Santa Clara County (XSC) standard message types.  Each message type
// is represented by its own concrete type; all of them implement the XSCMessage
// interface.
package xscmsg

import (
	"strconv"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
)

// XSCMessage is the interface implemented by all XSC message types.
type XSCMessage interface {
	// TypeTag returns the tag string used to identify the message type.
	TypeTag() string
	// TypeName returns the English name of the message type.  It is a noun
	// phrase in prose case, such as "foo message" or "bar form".
	TypeName() string
	// TypeArticle returns the indefinite article ("a" or "an") to be used
	// preceding the TypeName, in a sentence that needs one.
	TypeArticle() string
	// Message returns the encoded message.  If human is true, it is encoded
	// for human reading or editing; if false, it is encoded for
	// transmission.  If the XSCMessage was originally created by a call to
	// Recognize, the Message structure passed to it is updated and reused;
	// otherwise, a new Message structure is created and filled in.
	Message(human bool) *pktmsg.Message
}

var createFuncs []func(string) XSCMessage
var recognizeFuncs []func(*pktmsg.Message, *pktmsg.Form) XSCMessage

// RegisterType registers a message type to be accessible through the Create and
// Recognize functions.  (It is normally called by an init function in the
// package that implements the message type.  To make a message type
// recognizable, simply import its implementing package.)
func RegisterType(create func(string) XSCMessage, recognize func(*pktmsg.Message, *pktmsg.Form) XSCMessage) {
	createFuncs = append(createFuncs, create)
	recognizeFuncs = append(recognizeFuncs, recognize)
}

// Create creates a new message of the type identified by the supplied tag.  If
// the tag is not recognized, Create returns nil.
func Create(tag string) XSCMessage {
	for _, createFunc := range createFuncs {
		if xsc := createFunc(tag); xsc != nil {
			return xsc
		}
	}
	return nil
}

// Recognize examines the supplied Message to see if it is one of the standard
// XSC messages.  If so, it returns the appropriate XSCMessage implementation
// wrapping it.  If not, it returns nil.  The strict flag indicates whether any
// form embedded in the message should be parsed strictly; see pktmsg.ParseForm
// for details.
func Recognize(msg *pktmsg.Message, strict bool) XSCMessage {
	var form *pktmsg.Form
	if pktmsg.IsForm(msg.Body) {
		form = pktmsg.ParseForm(msg.Body, strict)
	}
	for _, recognizeFunc := range recognizeFuncs {
		if xsc := recognizeFunc(msg, form); xsc != nil {
			return xsc
		}
	}
	return nil
}

// OlderVersion compares two version numbers, and returns true if the first one
// is older than the second one.  Each version number is a dot-separated
// sequence of parts, each of which is compared independently with the
// corresponding part in the other version number.  The parts are compared
// numerically if they parse as integers, and as strings otherwise.  However, a
// part that starts with a digit is always "newer" than a part that doesn't.
// (This results in "undefined" being treated as infinitely old, which is what
// we want.)
func OlderVersion(a, b string) bool {
	aparts := strings.Split(a, ".")
	bparts := strings.Split(b, ".")
	for len(aparts) != 0 && len(bparts) != 0 {
		anum, aerr := strconv.Atoi(aparts[0])
		bnum, berr := strconv.Atoi(bparts[0])
		if aerr == nil && berr == nil {
			if anum < bnum {
				return true
			}
			if anum > bnum {
				return false
			}
		} else if startsWithDigit(aparts[0]) && !startsWithDigit(bparts[0]) {
			return false
		} else if !startsWithDigit(aparts[0]) && startsWithDigit(bparts[0]) {
			return true
		} else {
			if aparts[0] < bparts[0] {
				return true
			}
			if aparts[0] > bparts[0] {
				return false
			}
		}
		aparts = aparts[1:]
		bparts = bparts[1:]
	}
	if len(bparts) != 0 {
		return true
	}
	return false
}
func startsWithDigit(s string) bool {
	return s != "" && s[0] >= '0' && s[0] <= '9'
}
