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
	var msg Message
	msg.To = "nobody"
	msg.Body = "hello\n"
	msg.OutpostUrgent = true
	var save = msg.Save()
	const expected = "To: nobody\n\n!URG!hello\n"
	if save != expected {
		t.Fail()
	}
}

func TestEncodeOutpostB64(t *testing.T) {
	var msg Message
	msg.To = "nobody"
	msg.Body = "hellö\n"
	var save = msg.Save()
	const expected = "To: nobody\n\n!B64!aGVsbMO2Cg==\n"
	if save != expected {
		t.Fail()
	}
}

func TestEncodeMinimalForm(t *testing.T) {
	var a = TaggedField{Tag: "A"}
	var b = TaggedField{"B", "b"}
	var msg = Message{
		To:           "nobody",
		PIFOVersion:  CurrentPIFOVersion,
		FormHTML:     "tt.html",
		FormVersion:  "2",
		Subject:      "Subject",
		Handling:     "ROUTINE",
		FormTag:      "FTag",
		OriginMsgID:  "AAA-111P",
		TaggedFields: []TaggedField{a, b},
	}
	var save = msg.Save()
	const expected = "To: nobody\nSubject: AAA-111P_R_FTag_Subject\n\n!SCCoPIFO!\n#T: tt.html\n#V: 3.9-2\nB: [b]\n!/ADDON!\n"
	if save != expected {
		t.Fail()
	}
}

func TestEncodeFormWithLineContinuation(t *testing.T) {
	var a = TaggedField{"A", "1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 "}
	var msg = Message{
		To:           "nobody",
		PIFOVersion:  CurrentPIFOVersion,
		FormHTML:     "tt.html",
		FormVersion:  "2",
		TaggedFields: []TaggedField{a},
	}
	var save = msg.Save()
	const expected = "To: nobody\n\n!SCCoPIFO!\n#T: tt.html\n#V: 3.9-2\nA: [1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 123\n4567890 ]\n!/ADDON!\n"
	if save != expected {
		t.Fail()
	}
}
