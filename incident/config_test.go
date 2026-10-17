package incident

import (
	"encoding/json/v2"
	"testing"
	"time"
)

// configJSON is the JSON written for testConfig by the json/v2 `format` tags
// these fields used before Go 1.27 removed support for them.
const configJSON = `{"opStart":"2026-10-16T08:05","opEnd":"2026-10-17T20:00","bulletinChecks":{"a":{"d":3600},"b":{"d":90},"c":{"d":1.5},"d":{"d":0},"e":{"d":-1}},"noSendReceipts":false,"allowVoice":false}`

var testConfig = Config{
	OpStart: OpTime{time.Date(2026, 10, 16, 8, 5, 0, 0, time.Local)},
	OpEnd:   OpTime{time.Date(2026, 10, 17, 20, 0, 59, 9, time.UTC)},
	BulletinChecks: map[string]CheckFrequency{
		"a": {time.Hour}, "b": {90 * time.Second}, "c": {1500 * time.Millisecond}, "d": {0}, "e": {-time.Second},
	},
}

func TestConfigMarshal(t *testing.T) {
	if b, err := json.Marshal(testConfig, json.Deterministic(true)); err != nil {
		t.Fatal(err)
	} else if string(b) != configJSON {
		t.Errorf("got  %s\nwant %s", b, configJSON)
	}
	if b, err := json.Marshal(Config{}); err != nil {
		t.Fatal(err)
	} else if want := `{"noSendReceipts":false,"allowVoice":false}`; string(b) != want {
		t.Errorf("got  %s\nwant %s", b, want)
	}
}

func TestConfigUnmarshal(t *testing.T) {
	var c Config

	if err := json.Unmarshal([]byte(configJSON), &c, json.RejectUnknownMembers(true)); err != nil {
		t.Fatal(err)
	}
	if want := time.Date(2026, 10, 16, 8, 5, 0, 0, time.UTC); !c.OpStart.Equal(want) {
		t.Errorf("OpStart = %v, want %v", c.OpStart, want)
	}
	if want := time.Date(2026, 10, 17, 20, 0, 0, 0, time.UTC); !c.OpEnd.Equal(want) {
		t.Errorf("OpEnd = %v, want %v", c.OpEnd, want)
	}
	for area, cf := range testConfig.BulletinChecks {
		if c.BulletinChecks[area] != cf {
			t.Errorf("BulletinChecks[%q] = %v, want %v", area, c.BulletinChecks[area], cf)
		}
	}
	for _, in := range []string{`{}`, `{"opStart":null}`, `{"bulletinChecks":{"x":{}}}`} {
		c = Config{}
		if err := json.Unmarshal([]byte(in), &c, json.RejectUnknownMembers(true)); err != nil {
			t.Errorf("%s: %v", in, err)
		} else if !c.OpStart.IsZero() || c.BulletinChecks["x"].Duration != 0 {
			t.Errorf("%s: got non-zero value", in)
		}
	}
	for _, in := range []string{`{"opStart":"bad"}`, `{"opStart":5}`, `{"bulletinChecks":{"x":{"d":"1"}}}`, `{"bulletinChecks":{"x":{"d":1,"q":2}}}`} {
		if err := json.Unmarshal([]byte(in), &c, json.RejectUnknownMembers(true)); err == nil {
			t.Errorf("%s: expected error", in)
		}
	}
}
