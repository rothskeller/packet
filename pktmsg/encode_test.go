package pktmsg

import (
	"testing"
)

func TestEncodeReceived(t *testing.T) {
	const start = "Received: FROM bbs.ampr.org BY pktmsg.local FOR area;\n\tWed, 01 Dec 2021 08:04:29 +0000\nFrom: <nobody@nowhere>\nTo: <somebody@somewhere>\nSubject: Hello, World\nDate: Wed, 01 Dec 2021 08:04:29 +0000\n\nnothing\n"
	var msg, _ = ParseMessage(start)
	var end = msg.Save()
	if start != end {
		t.Fail()
	}
}

func TestEncodeOutpostFlags(t *testing.T) {
	var msg = NewMessage()
	msg.ToAddrs().SetValue("nobody")
	msg.Body().SetValue("hello\n")
	msg.OutpostFlags().SetUrgent(true)
	var save = msg.Save()
	const expected = "To: nobody\n\n!URG!hello\n"
	if save != expected {
		t.Fail()
	}
}

func TestEncodeOutpostB64(t *testing.T) {
	var msg = NewMessage()
	msg.ToAddrs().SetValue("nobody")
	msg.Body().SetValue("hellö\n")
	var save = msg.Save()
	const expected = "To: nobody\n\n!B64!aGVsbMO2Cg==\n"
	if save != expected {
		t.Fail()
	}
}

func TestEncodeMinimalForm(t *testing.T) {
	var a = taggedField{tag: "A"}
	var b = taggedField{"b", "B"}
	var form = NewForm("tt.html", "2", []*taggedField{&a, &b})
	form.ToAddrs().SetValue("nobody")
	var save = form.Save()
	const expected = "To: nobody\n\n!SCCoPIFO!\n#T: tt.html\n#V: 3.9-2\nB: [b]\n!/ADDON!\n"
	if save != expected {
		t.Fail()
	}
}

func TestEncodeFormWithLineContinuation(t *testing.T) {
	var a = taggedField{"1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 ", "A"}
	var form = NewForm("tt.html", "2", []*taggedField{&a})
	form.ToAddrs().SetValue("nobody")
	var save = form.Save()
	const expected = "To: nobody\n\n!SCCoPIFO!\n#T: tt.html\n#V: 3.9-2\nA: [1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 123\n4567890 ]\n!/ADDON!\n"
	if save != expected {
		t.Fail()
	}
}
