// wppcron is the weekly packet practice cron job.  It receives, analyzes,
// responds to, and reports on SCCo weekly packet practice messages.  It is
// designed to be run as a cron job, once every 5 minutes.  config.yaml
// (configuration data) and wppsvr.db (message database) must exist in the
// current directory when the program is started.
package main

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/rothskeller/packet/v4/form/formdefs"
	"github.com/rothskeller/packet/v4/wppsvr/config"
	"github.com/rothskeller/packet/v4/wppsvr/retrieve"
	"github.com/rothskeller/packet/v4/wppsvr/store"
)

func main() {
	var (
		st  *store.Store
		err error
	)
	os.Setenv("TZ", "PST8PDT")
	// One of these chdirs should succeed; we don't care which one.  If
	// neither does, and there is no wppsvr.db in the current directory,
	// the store.Open will fail.
	os.Chdir("/Users/stever/src/packet/v2/wppdata")
	os.Chdir("/home/sccares/private/wpp")
	openLog()
	formdefs.UseInternalForms = true
	formdefs.RegisterForms()
	if st, err = store.Open(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	if err = config.Read(); err != nil {
		os.Exit(1)
	}
	// Capture and log any panics that occur during the run.
	defer func() {
		if panicked := recover(); panicked != nil {
			log.Printf("PANIC: %v", panicked)
			log.Print(string(debug.Stack()))
		}
	}()
	checkBBSes(st)    // retrieve and respond to check-in messages
	closeSessions(st) // close sessions that are ending and send reports
	openSessions(st)  // open sessions that should be running
}

// checkBBSes retrieves, analyzes, and responds to new messages in all running
// practice sessions.
func checkBBSes(st *store.Store) {
	retrieve.ForRunningSessions(st)
}
