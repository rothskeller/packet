package message_test

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/ics213"
)

func TestRoundTrip(t *testing.T) {
	message.Register(&ics213.Type)
	ics213 := message.Create("ICS213").(*ics213.ICS213)
	ics213.PIFOVersion = "3.9"
	ics213.Subject = "Subject"
	ics213.Message = "Message"
	subject := ics213.EncodeSubject()
	body := ics213.EncodeBody()
	back := message.Decode(subject, body)
	if !reflect.DeepEqual(ics213, back) {
		spew.Dump(ics213)
		spew.Dump(back)
		t.Fatal("mismatch")
	}
}
