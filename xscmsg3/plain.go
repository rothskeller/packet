package xscmsg

import (
	"github.com/rothskeller/packet/pktmsg"
)

// PlainTextTag identifies a plain text message.
const PlainTextTag = "plain"
const plainTextName = "plain text message"

var plainTextType = MessageType{PlainTextTag, plainTextName, "a"}

func createPlainText() Message {
	return NewBaseMessage(plainTextType, nil)
}

func recognizePlainText(pm pktmsg.Message) Message {
	return NewBaseMessage(plainTextType, pm)
}
