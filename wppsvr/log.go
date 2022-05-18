package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// When development is true, log entries go to stderr instead of to a log file.
// It should be false in production deployments.
const development = false

var (
	logFH    *os.File // log file handle
	logMonth string   // month whose log file is open
)

// maybeReopenLog reopens the log file if we are now in a different month than
// we were when the current log file was opened.  In other words, this reopens
// the log just after midnight on the first of each month.
func maybeReopenLog() {
	if development {
		return
	}
	var currentMonth = time.Now().Format("2006-01")
	if currentMonth != logMonth {
		openLog()
	}
}

// openLog opens a log file named with the current year and month.
func openLog() {
	var (
		newFH    *os.File
		filename string
		err      error
	)
	if development {
		return
	}
	// Make sure the log directory exists.
	if err = os.MkdirAll("log", 0777); err != nil {
		log.Fatalf("ERROR: create log directory: %s", err)
	}
	// Open the log file and set it as the output of the standard logger.
	logMonth = time.Now().Format("2006-01")
	filename = filepath.Join("log", logMonth)
	if newFH, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		log.Fatalf("ERROR: create log file: %s", err)
	}
	log.SetOutput(newFH)
	// Close the previous month's log file if any.
	if logFH != nil {
		logFH.Close()
	}
	logFH = newFH
}
