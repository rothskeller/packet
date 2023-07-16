// fakereceive reads a message from stdin and hands it to wppsvr as if wppsvr
// had read it from a BBS.  Any response messages are emitted to stdout.
//
// usage: fakereceive frommailbox frombbs < messagetext
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/rothskeller/packet/message/allmsg"
	"github.com/rothskeller/packet/wppsvr/analyze"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
)

func main() {
	var (
		mailbox string
		bbs     string
		st      *store.Store
		session *store.Session
		rawmsg  []byte
		err     error
	)
	allmsg.Register()
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: fakereceive frommailbox frombbs < messagetext\n")
		os.Exit(2)
	}
	mailbox, bbs = os.Args[1], os.Args[2]
	if err = config.Read(); err != nil {
		log.Fatal(err)
	}
	if st, err = store.Open(); err != nil {
		log.Fatal(err)
	}
	if _, ok := config.Get().BBSes[bbs]; !ok {
		fmt.Fprintf(os.Stderr, "ERROR: bbs %q is not defined in configuration\n", bbs)
		os.Exit(1)
	}
	for _, s := range st.GetRunningSessions() {
		if s.CallSign == mailbox {
			session = s
			break
		}
	}
	if session == nil {
		fmt.Fprintf(os.Stderr, "ERROR: no running session for mailbox %q\n", mailbox)
		os.Exit(1)
	}
	rawmsg, _ = io.ReadAll(os.Stdin)
	analysis := analyze.Analyze(st, session, bbs, string(rawmsg))
	responses := analysis.Responses(st)
	for _, response := range responses {
		response.SendTime = time.Now()
	}
	analysis.Commit(st)
	for _, response := range responses {
		st.SaveResponse(response)
		fmt.Printf("To: %s\nSubject: %s\n\n%s", response.To, response.Subject, response.Body)
	}
}
