// Package store contains code that interacts directly with the database.
package db

import (
	"time"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

// Transaction executes the supplied function in a transaction.  It may call the
// function multiple times, in the face of retriable errors, so the function
// must be idempotent.  If a non-retriable error occurs, it panics.  If the
// function calls the DoNotCommit method, the transaction rolls back rather than
// committing, but does not panic.  The write parameter indicates whether the
// transaction should start in write-exclusive mode.
func Transaction(conn *sqlite.Conn, write bool, fn func() error) error {
	delay := 8 * time.Millisecond
	start := time.Now()
	const maxDelay = 1024 * time.Millisecond
	const maxWait = 30 * time.Second
	for {
		if ran, err := attemptTransaction(conn, write, fn); ran {
			return err // success
		}
		if time.Since(start) > maxWait {
			panic("too many retries")
		}
		time.Sleep(delay)
		if delay < maxDelay {
			delay *= 2
		}
	}
}

func attemptTransaction(conn *sqlite.Conn, write bool, fn func() error) (ran bool, err error) {
	var endfn func(*error)

	if write {
		endfn, err = sqlitex.ExclusiveTransaction(conn)
	} else {
		endfn = sqlitex.Transaction(conn)
	}
	if code := sqlite.ErrCode(err).ToPrimary(); code == sqlite.ResultBusy || code == sqlite.ResultLocked {
		return false, nil
	} else if err != nil {
		return false, err
	}
	defer endfn(&err)
	err = fn()
	return true, err
}

// attemptTransaction is a single transaction attempt.  It returns true if the
// transaction succeeded, false if it failed with a retriable error.  It panics
// if the transaction fails with a non-retriable error.  It returns a non-nil
// error only if one is returned by the function (in which case the transaction
// is rolled back).
func xattemptTransaction(conn *sqlite.Conn, write bool, fn func() error) (ran bool, err error) {
	// Set up the transaction cleanup handler.
	defer func() {
		ran = cleanupTransaction(conn, recover(), err)
	}()
	// Run the transaction.
	beginTransaction(conn, write)
	err = fn()
	println(0, err)
	return // return value is set in defer above.
}

// cleanup finalizes a transaction started with attemptTransaction.
func cleanupTransaction(conn *sqlite.Conn, panicked interface{}, err error) bool {
	// Roll back the transaction if it failed or its no-commit flag is set.
	if panicked != nil || err != nil {
		println(1)
		goto ROLLBACK
	}
	if e2 := commitTransaction(conn); e2 != nil {
		println(2)
		panicked = e2
		goto ROLLBACK
	}
	println(3)
	return true

ROLLBACK:
	// Roll back the transaction.  If we aren't already panicking and the rollback
	// fails, we'll panic with that cause.
	if e2 := rollbackTransaction(conn); e2 != nil && panicked == nil {
		panicked = e2
	}
	// If we rolled back due to the function returning an error, and had no other
	// issue, return true.
	if panicked == nil {
		println(4)
		return true
	}
	// If we rolled back due to a retriable error, return false.
	if e2, ok := panicked.(error); ok {
		if code := sqlite.ErrCode(e2).ToPrimary(); code == sqlite.ResultBusy || code == sqlite.ResultLocked {
			println(5)
			return false
		}
	}
	// In all other cases, re-raise the panic.
	println(6)
	panic(panicked)
}

// beginTransaction begins a transaction in the database.
func beginTransaction(conn *sqlite.Conn, write bool) (err error) {
	var (
		stmt *sqlite.Stmt
		cmd  = "BEGIN"
	)
	if write {
		cmd += " IMMEDIATE"
	}
	stmt = conn.Prep(cmd)
	if _, err = stmt.Step(); err != nil {
		return err
	}
	return stmt.Reset()
}

// commitTransaction commits a transaction in the database.
func commitTransaction(conn *sqlite.Conn) (err error) {
	println(">commit")
	defer println("<commit")
	stmt := conn.Prep("COMMIT")
	if _, err = stmt.Step(); err != nil {
		return err
	}
	return stmt.Reset()
}

// rollbackTransaction commits a transaction in the database.
func rollbackTransaction(conn *sqlite.Conn) (err error) {
	stmt := conn.Prep("ROLLBACK")
	if _, err = stmt.Step(); err != nil {
		return err
	}
	return stmt.Reset()
}
