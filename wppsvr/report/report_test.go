package report

import (
	"testing"
	"time"

	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/wppsvr/store"
	_ "steve.rothskeller.net/packet/xscmsg/all"
)

var fakeConfig = config.Config{
	ProblemActionFlags: map[string]config.Action{
		"ToBBSDown":                   config.ActionDontCount | config.ActionReport,
		"SubjectFormat":               config.ActionError | config.ActionReport,
		"MultipleMessagesFromAddress": config.ActionReport,
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
	ID:                     2,
	CallSign:               "PKTEST",
	Name:                   "Test Check-Ins",
	Start:                  time.Date(2022, 4, 18, 0, 0, 0, 0, time.Local),
	End:                    time.Date(2022, 4, 18, 23, 59, 0, 0, time.Local),
	ExcludeFromWeekSummary: true,
}
var fakeSession3 = store.Session{
	ID:                  3,
	CallSign:            "PKTTUE",
	Name:                "SVECS Net",
	Start:               time.Date(2022, 4, 13, 0, 0, 0, 0, time.Local),
	End:                 time.Date(2022, 4, 19, 20, 0, 0, 0, time.Local),
	GenerateWeekSummary: true,
	ToBBSes:             []string{"W2XSC"},
	DownBBSes:           []string{"W3XSC"},
	MessageTypes:        []string{"MuniStat", "plain"},
}

type fakeStore struct{}

func (fakeStore) GetSessionMessages(sessionID int) []*store.Message {
	switch sessionID {
	case 1:
		return []*store.Message{
			{
				Valid:        true,
				FromAddress:  "k6sny@w1xsc.ampr.org",
				FromCallSign: "K6SNY",
			},
			{
				Valid:        true,
				FromAddress:  "kc6rsc@w1xsc.ampr.org",
				FromCallSign: "KC6RSC",
			},
		}
	case 3:
		return []*store.Message{
			{
				Valid:        true,
				Correct:      true,
				FromAddress:  "kc6rsc@w1xsc.ampr.org",
				FromCallSign: "KC6RSC",
				FromBBS:      "W1XSC",
				Subject:      "STR-100P_I_MuniStat_Sunnyvale",
				Problems:     nil,
			},
			{
				Valid:        true,
				Correct:      true,
				FromAddress:  "kc6rsc@w1xsc.ampr.org",
				FromCallSign: "KC6RSC",
				FromBBS:      "W1XSC",
				Subject:      "STR-100P_I_MuniStat_Sunnyvale",
				Problems:     nil,
			},
			{
				Valid:        true,
				Correct:      false,
				FromAddress:  "aa6bt@w3xsc.ampr.org",
				FromCallSign: "AA6BT",
				FromBBS:      "W3XSC",
				Subject:      "BLAH",
				Problems:     []string{"ToBBSDown", "SubjectFormat"},
			},
		}
	default:
		panic("unknown sessionID")
	}
}

func (fakeStore) GetSessionsEnding(start, end time.Time) []*store.Session {
	if !start.Equal(time.Date(2022, 4, 17, 0, 0, 0, 0, time.Local)) ||
		!end.Equal(time.Date(2022, 4, 19, 20, 0, 0, 0, time.Local)) {
		panic("wrong dates for GetSessionsEnding")
	}
	return []*store.Session{&fakeSession1, &fakeSession2}
}

func (fakeStore) UpdateSession(*store.Session) { panic("not implemented") }
func (fakeStore) NextMessageID(string) string  { panic("not implemented") }

const expected = `====  SCCo ARES/RACES Packet Practice Report  ====
==== for SVECS Net on Tuesday, April 19, 2022 ====

This practice session expected an OA jurisdiction status form or plain text
message sent to PKTTUE at W2XSC, between 00:00 on Wednesday, April 13 and
20:00 on Tuesday, April 19, 2022.  W3XSC was simulated "down" for this
practice session.

Total messages:     3
Unique addresses:   2
Correct messages:   1  (50%)
Unique call signs:  2  [report this count to the net]
  for the week:     3
Messages from:
  W1XSC:            1
  W3XSC:            1  (simulated down)

---- The following messages were counted in this report: ----
aa6bt@w3xsc.ampr.org           BLAH
  ^ message to incorrect BBS (simulated down)
  ^ incorrect subject line format
kc6rsc@w1xsc.ampr.org          STR-100P_I_MuniStat_Sunnyvale
  ^ multiple messages from this address

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
