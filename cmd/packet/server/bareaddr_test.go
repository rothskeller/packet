package server

import (
	"strings"
	"testing"

	"github.com/rothskeller/packet/v4/incident"
)

func TestBareAddressWarning(t *testing.T) {
	inc := &incident.Incident{Config: &incident.Config{ConnectBBS: "W3XSC"}}
	for _, tc := range []struct {
		to   string
		want string // substring of the warning, or "" for none
	}{
		{"PKTMON", "PKTMON has no @BBS part"},
		{"pktmon", "pktmon has no @BBS part"},
		{"PKTMON@W2XSC", ""},
		{"PKTMON@w2xsc.#nca.ca.usa.noam", ""},
		{"nobody@example.com", ""},
		{"KK6QRA@W1XSC, PKTMON", "PKTMON has no @BBS part"},
		{"XSCPERM", ""}, // a bulletin address
		{"", ""},
	} {
		got := bareAddressWarning(inc, tc.to)
		if tc.want == "" && got != "" {
			t.Errorf("%q: unexpected warning %q", tc.to, got)
		} else if tc.want != "" && !strings.Contains(got, tc.want) {
			t.Errorf("%q: warning %q lacks %q", tc.to, got, tc.want)
		}
	}
	if got := bareAddressWarning(inc, "PKTMON"); !strings.Contains(got, "(W3XSC)") {
		t.Errorf("warning doesn't name the BBS we connect to: %q", got)
	}
	if got := bareAddressWarning(nil, "PKTMON"); got == "" {
		t.Error("no warning without an incident")
	}
}
