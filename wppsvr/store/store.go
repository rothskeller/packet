// Package store manages the message database.
package store

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3" // database driver
)

// Store gives access to the database store.
type Store struct {
	dbh   *sql.DB
	mutex sync.RWMutex
}

// Open opens the database store.
func Open() (store *Store, err error) {
	store = new(Store)
	store.dbh, err = sql.Open("sqlite3",
		"file:wppsvr.db?mode=rw&cache=shared&_busy_timeout=1000&_txlock=immediate&_foreign_keys=1")
	// mode=rw overrides the default rwc and causes failure if the database doesn't exist.
	// cache=shared is a performance optimization when multiple threads use the same database.
	// _busy_timeout=1000 sets a 1-second wait time before lock contention failure.
	// _txlock=immediate causes transactions to acquire locks immediately rather than lazily.
	// _foreign_keys=1 turns on foreign key reference enforcement.
	if err != nil {
		return nil, fmt.Errorf("open database: %s", err)
	}
	return store, nil
}

// split splits a string on a semicolon.  Unlike strings.Split, it returns an
// empty list if the input is an empty string.
func split(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ";")
}
