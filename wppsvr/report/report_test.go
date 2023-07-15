package report

import (
	"testing"
	"time"

	"github.com/rothskeller/packet/message/allmsg"
	"github.com/rothskeller/packet/wppsvr/store"
)

var fakeSession1 = store.Session{
	ID:       1,
	CallSign: "PKTMON",
	Name:     "SPECS Net",
	Start:    time.Date(2022, 4, 12, 0, 0, 0, 0, time.Local),
	End:      time.Date(2022, 4, 18, 20, 0, 0, 0, time.Local),
}
var fakeSession2 = store.Session{
	ID:       2,
	CallSign: "PKTEST",
	Name:     "Test Check-Ins",
	Start:    time.Date(2022, 4, 18, 0, 0, 0, 0, time.Local),
	End:      time.Date(2022, 4, 18, 23, 59, 0, 0, time.Local),
	Flags:    store.ExcludeFromWeek,
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
	Retrieve:     []*store.Retrieval{{LastRun: time.Date(2022, 4, 19, 20, 0, 0, 0, time.Local)}},
}

type fakeStore struct{}

func (fakeStore) GetSessionMessages(sessionID int) []*store.Message {
	switch sessionID {
	case 1:
		return []*store.Message{
			{
				LocalID:      "TST-001P",
				FromAddress:  "k6sny@w1xsc.ampr.org",
				FromCallSign: "K6SNY",
				Score:        100,
			},
			{
				LocalID:      "TST-002P",
				FromAddress:  "kc6rsc@w1xsc.ampr.org",
				FromCallSign: "KC6RSC",
				Score:        100,
			},
		}
	case 3:
		return []*store.Message{
			{
				LocalID:      "TST-003P",
				FromAddress:  "kc6rsc@w1xsc.ampr.org",
				FromCallSign: "KC6RSC",
				FromBBS:      "W1XSC",
				Jurisdiction: "SNY",
				MessageType:  "plain",
				Score:        100,
				Summary:      "OK",
			},
			{
				LocalID:      "TST-004P",
				FromAddress:  "kc6rsc@w1xsc.ampr.org",
				FromCallSign: "KC6RSC",
				FromBBS:      "W1XSC",
				Jurisdiction: "SNY",
				MessageType:  "plain",
				Score:        100,
				Summary:      "OK",
			},
			{
				LocalID:      "TST-005P",
				FromAddress:  "aa6bt@w3xsc.ampr.org",
				FromCallSign: "AA6BT",
				FromBBS:      "W3XSC",
				Jurisdiction: "Unknown",
				MessageType:  "plain",
				Score:        77,
				Summary:      "multiple issues",
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

const expected = `==== SCCo ARES/RACES Packet Practice Report
==== for SVECS Net on Tuesday, April 19, 2022

2 unique call signs (3 for the week)

EXPECTATIONS:  OA municipal status form or plain text message sent to PKTTUE
at W2XSC between Wed 2022-04-13 00:00 and Tue 2022-04-19 20:00; not sent from
W3XSC.

---- RESULTS
88% Average Score
 2  Counted
 1  Duplicate

---- MESSAGES
AA6BT   @W3XSC   (???)   77%  multiple issues
KC6RSC  @W1XSC*  (SNY)  100%  OK
* multiple messages from this address; only the last one counts

---- SENT FROM
1  W1XSC
1  W3XSC (simulated outage)

---- JURISDICTION
1  SNY
1  ???

---- MESSAGE TYPE
2  plain

This report was generated on Tuesday, April 19, 2022 at 20:00 by wppsvr.
`

func TestReport(t *testing.T) {
	allmsg.Register()
	now = func() time.Time { return time.Date(2022, 4, 19, 20, 0, 1, 0, time.Local) }
	actual := Generate(fakeStore{}, &fakeSession3).RenderPlainText()
	if actual != expected {
		t.Errorf("incorrect report output:\n%s", actual)
	}
}
