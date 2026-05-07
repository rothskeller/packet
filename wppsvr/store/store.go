// Package store manages the message database.
package store

import (
	"fmt"
	"strings"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

// Store gives access to the database store.
type Store struct {
	conn *sqlite.Conn
}

// Open opens the database store.
func Open() (store *Store, err error) {
	store = new(Store)
	store.conn, err = sqlite.OpenConn("wppsvr.db", sqlite.OpenReadWrite)
	if err != nil {
		return nil, fmt.Errorf("open database: %s", err)
	}
	if err = sqlitex.ExecuteTransient(store.conn, "PRAGMA journal_mode = TRUNCATE;", nil); err != nil {
		return nil, err
	}
	if err = sqlitex.ExecuteTransient(store.conn, "PRAGMA foreign_keys = ON;", nil); err != nil {
		return nil, err
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
