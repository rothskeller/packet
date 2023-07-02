package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/rothskeller/packet/jnos"
	"github.com/rothskeller/packet/jnos/simulator"
	"github.com/rothskeller/packet/wppsvr/analyze"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/interval"
	"github.com/rothskeller/packet/wppsvr/report"
	"github.com/rothskeller/packet/wppsvr/retrieve"
	"github.com/rothskeller/packet/wppsvr/store"
	_ "github.com/rothskeller/packet/xscmsg/all"
)

func main() {
	var (
		st      *store.Store
		session store.Session
		fh      *os.File
		conn    *jnos.Conn
		sim     *simulator.Simulator
		err     error
	)
	if fh, err = os.Open("messages.2022-05-17.txt"); err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer fh.Close()
	if sim, err = simulator.Start(map[string]io.Reader{"": fh}, ""); err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	if st, err = store.Open(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	session = store.Session{
		CallSign:     "PKTTUE",
		Name:         "SVECS Net",
		Prefix:       "TUE",
		Start:        time.Date(2022, 5, 11, 0, 0, 0, 0, time.Local),
		End:          time.Date(2022, 5, 17, 20, 0, 0, 0, time.Local),
		ReportToText: []string{"packet@scc-ares-races.groups.io"},
		ToBBSes:      []string{"W2XSC"},
		DownBBSes:    []string{"W3XSC"},
		Retrieve: []*store.Retrieval{{
			When:    "MINUTE=0",
			BBS:     "W2XSC",
			Mailbox: "PKTTUE",
		}},
		MessageTypes: []string{"plain"},
	}
	session.Retrieve[0].Interval = interval.Parse(session.Retrieve[0].When)
	st.CreateSession(&session)
	if err = config.Read(analyze.ProblemLabels); err != nil {
		os.Exit(1)
	}
	retrieve.ForSession(st, &session)
	conn = retrieve.ConnectToBBS("W2XSC", "x")
	if conn == nil {
		panic("nil")
	}
	report.Send(st, conn, &session)
	sim.Stop()
	for _, sent := range sim.Sent() {
		fmt.Print("\n------------------------------------------------------------------------------\n\n")
		fmt.Print(sent)
	}
}
