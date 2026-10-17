package tnc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuiltinsValid(t *testing.T) {
	for name, def := range tncs {
		if err := def.Validate(); err != nil {
			t.Errorf("builtin %q: %v", name, err)
		}
	}
}

func TestFetchOneValidation(t *testing.T) {
	const fields = `"CommandPrompt": "cmd:", "ConnectedMessage": "*** CONNECTED", "DisconnectedMessage": "*** DISCONNECTED", "MyCallCommand": "MY", "ConnectCommand": "C"`
	cases := []struct {
		name string // TNC name in the file
		json string
		ok   bool
	}{
		{"Test Good", `{"Name": "Test Good", ` + fields + `, "EchoTimeout": 300, "ReplyTimeout": 0}`, true},
		{"", `{"Name": "", ` + fields + `}`, false},
		{"Test No Prompt", `{"Name": "Test No Prompt", "ConnectedMessage": "*** CONNECTED", "DisconnectedMessage": "*** DISCONNECTED", "MyCallCommand": "MY", "ConnectCommand": "C"}`, false},
		{"Test Empty Connected", `{"Name": "Test Empty Connected", "CommandPrompt": "cmd:", "ConnectedMessage": "", "DisconnectedMessage": "*** DISCONNECTED", "MyCallCommand": "MY", "ConnectCommand": "C"}`, false},
		{"Test Blank Disconnected", `{"Name": "Test Blank Disconnected", "CommandPrompt": "cmd:", "ConnectedMessage": "*** CONNECTED", "DisconnectedMessage": "  ", "MyCallCommand": "MY", "ConnectCommand": "C"}`, false},
		{"Test No MyCall", `{"Name": "Test No MyCall", "CommandPrompt": "cmd:", "ConnectedMessage": "*** CONNECTED", "DisconnectedMessage": "*** DISCONNECTED", "ConnectCommand": "C"}`, false},
		{"Test No Connect", `{"Name": "Test No Connect", "CommandPrompt": "cmd:", "ConnectedMessage": "*** CONNECTED", "DisconnectedMessage": "*** DISCONNECTED", "MyCallCommand": "MY", "ConnectCommand": ""}`, false},
		{"Test Negative Echo", `{"Name": "Test Negative Echo", ` + fields + `, "EchoTimeout": -1}`, false},
		{"Test Negative Reply", `{"Name": "Test Negative Reply", ` + fields + `, "ReplyTimeout": -100}`, false},
		{"Test Bad JSON", `{"Name": "Test Bad JSON", `, false},
	}
	dir := t.TempDir()
	for i, c := range cases {
		path := filepath.Join(dir, "tnc"+string(rune('a'+i))+".json")
		if err := os.WriteFile(path, []byte(c.json), 0o600); err != nil {
			t.Fatal(err)
		}
		_, existed := tncs[c.name]
		if existed {
			t.Fatalf("test TNC name %q collides with a builtin", c.name)
		}
		fetchOne(path)
		def, got := tncs[c.name]
		delete(tncs, c.name)
		if got != c.ok {
			t.Errorf("%s (%q): loaded=%v, want %v", path, c.name, got, c.ok)
		}
		if got && c.ok && (def.EchoTimeout != 300 || def.CommandPrompt != "cmd:") {
			t.Errorf("%q loaded incorrectly: %+v", c.name, def)
		}
	}
}
