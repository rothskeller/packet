package report

import (
	"testing"
	"time"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
	_ "github.com/rothskeller/packet/xscmsg/all"
)

var fakeConfig = config.Config{
	Problems: map[string]*config.ProblemConfig{
		"ToBBSDown":                   {ActionFlags: config.ActionDontCount | config.ActionReport},
		"SubjectFormat":               {ActionFlags: config.ActionError | config.ActionReport},
		"MultipleMessagesFromAddress": {ActionFlags: config.ActionReport},
	},
}

var fakeSession1 = store.Session{
	ID:       1,
	CallSign: "PKTMON",
	Name:     "SPECS Net",
	Start:    time.Date(2022, 4, 12, 0, 0, 0, 0, time.Local),
	End:      time.Date(2022, 4, 18, 20, 0, 0, 0, time.Local),
}
var fakeSession2 = store.Session{
	ID:              2,
	CallSign:        "PKTEST",
	Name:            "Test Check-Ins",
	Start:           time.Date(2022, 4, 18, 0, 0, 0, 0, time.Local),
	End:             time.Date(2022, 4, 18, 23, 59, 0, 0, time.Local),
	ExcludeFromWeek: true,
}
var fakeSession3 = store.Session{
	ID:           3,
	CallSign:     "PKTTUE",
	Name:         "SVECS Net",
	Start:        time.Date(2022, 4, 13, 0, 0, 0, 0, time.Local),
	End:          time.Date(2022, 4, 19, 20, 0, 0, 0, time.Local),
	ToBBSes:      []string{"W2XSC"},
	DownBBSes:    []string{"W3XSC"},
	MessageTypes: []string{"MuniStat", "plain"},
}

type fakeStore struct{}

func (fakeStore) GetSessionMessages(sessionID int) []*store.Message {
	switch sessionID {
	case 1:
		return []*store.Message{
			{
				FromAddress:  "k6sny@w1xsc.ampr.org",
				FromCallSign: "K6SNY",
			},
			{
				FromAddress:  "kc6rsc@w1xsc.ampr.org",
				FromCallSign: "KC6RSC",
			},
		}
	case 3:
		return []*store.Message{
			{
				FromAddress:  "kc6rsc@w1xsc.ampr.org",
				FromCallSign: "KC6RSC",
				FromBBS:      "W1XSC",
				Subject:      "STR-100P_I_MuniStat_Sunnyvale",
				Jurisdiction: "SNY",
				MessageType:  "plain",
				Problems:     nil,
			},
			{
				FromAddress:  "kc6rsc@w1xsc.ampr.org",
				FromCallSign: "KC6RSC",
				FromBBS:      "W1XSC",
				Subject:      "STR-100P_I_MuniStat_Sunnyvale",
				Jurisdiction: "SNY",
				MessageType:  "plain",
				Problems:     nil,
			},
			{
				FromAddress:  "aa6bt@w3xsc.ampr.org",
				FromCallSign: "AA6BT",
				FromBBS:      "W3XSC",
				Subject:      "BLAH",
				Jurisdiction: "SJC",
				MessageType:  "plain",
				Problems:     []string{"ToBBSDown", "SubjectFormat"},
				Actions:      config.ActionReport | config.ActionError,
			},
		}
	default:
		panic("unknown sessionID")
	}
}

func (fakeStore) GetSessions(start, end time.Time) []*store.Session {
	if !start.Equal(time.Date(2022, 4, 17, 0, 0, 0, 0, time.Local)) ||
		!end.Equal(time.Date(2022, 4, 24, 0, 0, 0, 0, time.Local)) {
		panic("wrong dates for GetSessionsEnding")
	}
	return []*store.Session{&fakeSession1, &fakeSession2, &fakeSession3}
}

func (fakeStore) UpdateSession(*store.Session) { panic("not implemented") }
func (fakeStore) NextMessageID(string) string  { panic("not implemented") }

const expected = `--------------- SCCo ARES/RACES Packet Practice Report ----------------
                for SVECS Net on Tuesday, April 19, 2022

                  2 unique call signs (3 for the week)


EXPECTATIONS------------------------------      RESULTS-----
Message type:  OA municipal status form or      OK         1
               plain text message               ERROR      1
Sent to:       PKTTUE at W2XSC                  Duplicate  1
Sent between:  Wed 2022-04-13 00:00 and
               Tue 2022-04-19 20:00
Not sent from: W3XSC

SOURCES-----------      JURISDICTIONS      TYPES---
W1XSC   1               SJC  1             plain  2
W3XSC*  1               SNY  1
* simulated outage

MESSAGES----------------------------------------------------------------
AA6BT   W3XSC   SJC  ERROR: multiple issues
KC6RSC  W1XSC*  SNY  OK
* multiple messages from this address; only the last one counts


This report was generated on Tuesday, April 19, 2022 at 20:00 by wppsvr.
`

func TestReport(t *testing.T) {
	now = func() time.Time { return time.Date(2022, 4, 19, 20, 0, 1, 0, time.Local) }
	config.SetConfig(&fakeConfig)
	actual := Generate(fakeStore{}, &fakeSession3).RenderPlainText()
	if actual != expected {
		t.Errorf("incorrect report output:\n%s", actual)
	}
}
