// wppsvr is a server that retrieves, analyzes, responds to, and reports on SCCo
// weekly packet practice messages.  config.yaml (configuration data) and
// wppsvr.db (message database) must exist in the current directory when the
// server is started.  The server runs forever, operating periodically as
// dictated by the configuration.  The configuration is re-read periodically and
// can be changed while the server is running.
//
// If multiple instances of the server are started concurrently, all but the
// first of them will exit immediately and silently.  As a result, the server
// program can be invoked frequently (e.g., via cron job) as a sort of
// poor-man's HA strategy: if a previous copy is already running, it will do
// nothing, and if not, the new copy will take over.
package main

import (
	"log"
	"os"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/rothskeller/packet/message/allmsg"
	"github.com/rothskeller/packet/wppsvr/analyze"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/retrieve"
	"github.com/rothskeller/packet/wppsvr/store"
	"github.com/rothskeller/packet/wppsvr/webserver"
)

func main() {
	var (
		st  *store.Store
		err error
	)
	openLog()
	ensureSingleton()
	allmsg.Register()
	if st, err = store.Open(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	if err = config.Read(analyze.ProblemLabels); err != nil {
		os.Exit(1)
	}
	if err = webserver.Run(st); err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	for { // repeat forever
		step(st)
	}
}

// step gets executed every 5 minutes, and handles all of the scheduled
// activities of the wppsvr system.
func step(st *store.Store) {
	// Capture any panics that occur during the step.  Log them but don't
	// let them kill the whole program.
	defer func() {
		if panicked := recover(); panicked != nil {
			log.Printf("PANIC: %v", panicked)
			log.Print(string(debug.Stack()))
		}
		sleep5min()
	}()
	maybeReopenLog()                   // at midnight on the first of each month
	config.Read(analyze.ProblemLabels) // re-read config in case it has changed
	checkBBSes(st)                     // retrieve and respond to check-in messages
	closeSessions(st)                  // close sessions that are ending and send reports
	openSessions(st)                   // open sessions that should be running
}

// lockFH is the singleton lock file used in ensureSingleton.  It is declared at
// global scope so that it never gets garbage collected.
var lockFH *os.File

// ensureSingleton makes sure there is only one instance of wppsvr running at a
// time.  Redundant instances exit immediately and silently.
func ensureSingleton() {
	var err error

	// Open (or create) the run.lock file.
	if lockFH, err = os.OpenFile("run.lock", os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		log.Fatalf("ERROR: open run.lock: %s", err)
	}
	// Acquire an exclusive lock on the run.lock file.
	switch err = syscall.Flock(int(lockFH.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err {
	case nil:
		// Lock successfully acquired, so we are the only running
		// instance.  We will hold the lock until our process exits.
		return
	case syscall.EWOULDBLOCK:
		// Another process has the lock, so there is already another
		// running instance.  Exit immediately and silently.
		os.Exit(0)
	default:
		// Unable to acquire the lock, for some reason other than
		// another process holding it.  Report the error and exit.
		log.Fatalf("ERROR: lock run.lock: %s", err)
	}
}

// checkBBSes retrieves, analyzes, and responds to new messages in all running
// practice sessions.
func checkBBSes(st *store.Store) {
	retrieve.ForRunningSessions(st)
}

// sleep5min sleeps until the clock next reaches a multiple of 5 minutes.
func sleep5min() {
	var now = time.Now()
	var next = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(),
		(now.Minute()/5+1)*5, 0, 0, time.Local)
	var delta = time.Until(next)
	time.Sleep(delta)
}
