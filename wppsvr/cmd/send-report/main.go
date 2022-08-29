// send-report sends the report for the specified session to the specified
// email addresses.
//
// usage: send-report session-date email-address...
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rothskeller/packet/wppsvr/analyze"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/report"
	"github.com/rothskeller/packet/wppsvr/store"
	_ "github.com/rothskeller/packet/xscmsg/all" // register message types
)

func main() {
	var (
		date     time.Time
		st       *store.Store
		sessions []*store.Session
		session  *store.Session
		rep      *report.Report
		err      error
	)
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: send-report session-date email-address...\n")
		os.Exit(2)
	}
	if date, err = time.Parse("2006-01-02", os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %q is not a date\n", os.Args[1])
		fmt.Fprintf(os.Stderr, "usage: send-report session-date email-address...\n")
		os.Exit(2)
	}
	if err = config.Read(analyze.KnownProblems()); err != nil {
		log.Fatal(err)
	}
	if st, err = store.Open(); err != nil {
		log.Fatal(err)
	}
	sessions = st.GetSessions(
		time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local),
		time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, time.Local),
	)
	switch len(sessions) {
	case 0:
		fmt.Fprintf(os.Stderr, "ERROR: no sessions on %s\n", os.Args[1])
		os.Exit(1)
	case 1:
		session = sessions[0]
	default:
		fmt.Fprintf(os.Stderr, "ERROR: multiple sessions on %s\n", os.Args[1])
		os.Exit(1)
	}
	if session == nil {
		fmt.Fprintf(os.Stderr, "usage: reanalyze session-end-date\n")
		os.Exit(2)
	}
	rep = report.Generate(st, session)
	if err := rep.SendHTML(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
