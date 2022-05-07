package analyze

// This file defines a framework for testing the analysis package.  Every *.yaml
// file in the testdata tree, except config.yaml, describes a test case,
// including configuration and session setup, the message to be analyzed, the
// expected results of the analysis, and the expected response messages that
// should be generated.  The code in this file runs the analysis and tests the
// result.

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
	"gopkg.in/yaml.v3"

	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/wppsvr/store"
	_ "steve.rothskeller.net/packet/xscmsg/all"
)

type testdata struct {
	Now       time.Time        `yaml:"now"`
	Config    *config.Config   `yaml:"config"`
	SeenHash  string           `yaml:"seenHash"`
	Session   *store.Session   `yaml:"session"`
	ToBBS     string           `yaml:"toBBS"`
	Message   string           `yaml:"message"`
	Stored    *store.Message   `yaml:"stored"`
	Responses []*responseCheck `yaml:"responses"`
}
type responseCheck struct {
	store.Response `yaml:",inline"`
	BodyREs        []string `yaml:"bodyREs"`
}

func TestAnalyze(t *testing.T) {
	var testfiles []string
	log.SetOutput(io.Discard)
	filepath.WalkDir("testdata", func(path string, info fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".yaml") && path != "testdata/config.yaml" {
			testfiles = append(testfiles, path)
		}
		return nil
	})
	for _, testfile := range testfiles {
		t.Run(testfile[:len(testfile)-5], func(t *testing.T) {
			testAnalyze(t, testfile)
		})
	}
}

func testAnalyze(t *testing.T, testfile string) {
	var testdata testdata

	// First, we need a partial configuration.  We'll read that from
	// testdata/config.yaml, and then allow it to be modified by the test's
	// yaml file.
	os.Chdir("testdata")
	config.Read()
	os.Chdir("..")
	testdata.Config = config.Get()
	// We also need a session definition, which again, the test's yaml file
	// can modify.
	testdata.Session = &store.Session{
		ID:           42,
		CallSign:     "PKTTUE",
		Name:         "SVECS Net",
		Prefix:       "TUE",
		Start:        time.Date(2022, 1, 5, 0, 0, 0, 0, time.Local),
		End:          time.Date(2022, 1, 11, 20, 0, 0, 0, time.Local),
		ToBBSes:      []string{"W4XSC"},
		DownBBSes:    []string{"W2XSC"},
		MessageTypes: []string{"plain"},
	}
	// By default, we'll assume that "now" is a moment after the end of the
	// session.
	testdata.Now = time.Date(2022, 1, 11, 20, 0, 1, 0, time.Local)
	now = func() time.Time { return testdata.Now }
	// By default, we'll assume the test message is sent to the correct BBS.
	testdata.ToBBS = "W4XSC"
	// Now, read the test data file, allowing it to override any of the
	// above as well as defining the test message and expected output.
	fh, _ := os.Open(testfile)
	defer fh.Close()
	dec := yaml.NewDecoder(fh)
	dec.KnownFields(true)
	if err := dec.Decode(&testdata); err != nil {
		t.Fatal(err)
	}
	testdata.Config.Validate()
	// We'll need a fake store for the analyzer to use.
	store := &fakeStore{seenHash: testdata.SeenHash, nextID: 100}
	// Run the analysis.
	a := Analyze(store, testdata.Session, testdata.ToBBS, testdata.Message)
	responses := a.Responses(store)
	a.Commit(store)
	// First of all, did the analysis store the expected number of analyzed
	// messages (zero or one)?
	if testdata.Stored != nil && len(store.saved) == 0 {
		t.Error("no analysis saved to store")
	}
	if testdata.Stored == nil && len(store.saved) != 0 {
		t.Errorf("unexpected analysis saved to store: %s", spew.Sdump(store.saved[0]))
	}
	if len(store.saved) > 1 {
		t.Error("multiple analyses saved to store")
	}
	// If we both expected and got an analysis, is it correct?
	if testdata.Stored != nil && len(store.saved) != 0 {
		// The raw message in the analysis should be the same as the
		// input message.  We'll assign that here so it doesn't have to
		// be redundantly provided in every test file.
		testdata.Stored.Message = testdata.Message
		// We'll simply assume that the hash is correct, rather than
		// need to compute hashes for each test.
		testdata.Stored.Hash = store.saved[0].Hash
		// Fill in some defaults for the expected analysis.
		if testdata.Stored.LocalID == "" {
			testdata.Stored.LocalID = "TUE-100P"
		}
		if testdata.Stored.Session == 0 {
			testdata.Stored.Session = 42
		}
		if testdata.Stored.ToBBS == "" {
			testdata.Stored.ToBBS = "W4XSC"
		}
		// With those changes made, the expected and actual analysis
		// should compare identically.
		for _, diff := range deep.Equal(testdata.Stored, store.saved[0]) {
			t.Errorf("analysis mismatch: %s", diff)
		}
	}
	for i := 0; i < len(testdata.Responses) || i < len(responses); i++ {
		if i >= len(testdata.Responses) {
			t.Errorf("unexpected response: %s", spew.Sdump(responses[i]))
		} else if i >= len(responses) {
			t.Errorf("missing expected response: %s", spew.Sdump(testdata.Responses[i]))
		} else {
			checkResponse(t, testdata.Responses[i], responses[i])
		}
	}
}

func checkResponse(t *testing.T, want *responseCheck, have *store.Response) {
	// Fill in a few defaults so the test file doesn't have to specify
	// common stuff.
	if want.ResponseTo == "" {
		want.ResponseTo = "TUE-100P"
	}
	if want.SenderBBS == "" {
		want.SenderBBS = "W4XSC"
	}
	if want.SenderCall == "" {
		want.SenderCall = "PKTTUE"
	}
	// If the test file didn't specify a body for the response, copy over
	// the one we got so that they'll compare identically.
	if want.Body == "" {
		want.Body = have.Body
	}
	// Now compare the responses.
	for _, diff := range deep.Equal(&want.Response, have) {
		t.Errorf("response mismatch: %s", diff)
	}
	// The test file can also specify regular expressions that should be
	// matched by the body (often an easier way to write the test).  Check
	// those.
	for _, restr := range want.BodyREs {
		re, err := regexp.Compile(restr)
		if err != nil {
			t.Errorf("invalid RE in test: %s", err)
			return
		}
		if !re.MatchString(have.Body) {
			t.Errorf("response body does not match RE %q: %s", restr, spew.Sdump(have))
		}
	}
}

type fakeStore struct {
	seenHash string
	nextID   int
	saved    []*store.Message
}

func (f *fakeStore) HasMessageHash(hash string) string {
	if f.seenHash == hash {
		return "XXX-000P"
	}
	return ""
}

func (f *fakeStore) NextMessageID(prefix string) string {
	f.nextID++
	return fmt.Sprintf("%s-%03dP", prefix, f.nextID-1)
}

func (f *fakeStore) SaveMessage(m *store.Message) {
	f.saved = append(f.saved, m)
}
