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

// openLog opens a log file named with the current year and month.
func openLog() {
	var (
		logMonth string
		logFH    *os.File
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
	if logFH, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		log.Fatalf("ERROR: create log file: %s", err)
	}
	log.SetOutput(logFH)
}
