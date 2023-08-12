package envelope

import (
	"testing"
)

func TestEncodeReceived(t *testing.T) {
	const start = "Received: FROM bbs.ampr.org BY pktmsg.local FOR area;\n\tWed, 01 Dec 2021 08:04:29 +0000\nFrom: Nobody <nobody@nowhere>\nTo: Somebody <somebody@somewhere>\nSubject: Hello, World\nDate: Wed, 01 Dec 2021 08:04:29 +0000\n\nnothing\n"
	var env, body, _ = ParseSaved(start)
	var end = env.RenderSaved(body)
	if start != end {
		t.Fatalf("actual:\n%s\nexpected:\n%s\n", end, start)
	}
}

func TestEncodeOutpostFlags(t *testing.T) {
	var env Envelope
	env.To = []string{"nobody"}
	env.OutpostUrgent = true
	var save = env.RenderSaved("hello\n")
	const expected = "To: nobody\n\n!URG!hello\n"
	if save != expected {
		t.Fail()
	}
}

func TestEncodeOutpostB64(t *testing.T) {
	var env Envelope
	env.To = []string{"nobody"}
	var save = env.RenderSaved("hellö\n")
	const expected = "To: nobody\n\n!B64!aGVsbMO2Cg==\n"
	if save != expected {
		t.Fail()
	}
}
