package db

import (
	"time"

	"zombiezen.com/go/sqlite"
)

// St is a database statement.
type St struct {
	stmt   *sqlite.Stmt
	column int
}

// SQL calls the supplied function with a new statement.  Within the function,
// use the statement's BindXXX methods to bind values to any parameters in the
// SQL string, then call the statement's Step method to run the query and get
// the first result row.  Use the statement's ColumnXXX methods to fetch the
// values in the row.  Continue using the Step method to fetch subsequent rows.
// If the function needs to run the statement multiple times with different
// bound variables, call the statement's Reset method between runs.
//
// SQL and the statement methods panic on any error.
func SQL(conn *sqlite.Conn, sql string, fn func(st *St)) {
	var st St

	st.stmt = conn.Prep(sql)
	defer func() {
		if err := st.stmt.Reset(); err != nil {
			panic(err)
		}
	}()
	fn(&st)
}

// BindBlob binds the next statement parameter to the specified byte slice.
func (st *St) BindBlob(v []byte) {
	st.column++
	st.stmt.BindBytes(st.column, v)
}

// BindNullBlob binds the next statement parameter to the specified byte slice
// if it's non-nil, or to NULL if it's nil.  (Zero length is not nil.)
func (st *St) BindNullBlob(v []byte) {
	if v == nil {
		st.BindNull()
	} else {
		st.BindBlob(v)
	}
}

// BindBool binds the next statement parameter to the specified bool value.
func (st *St) BindBool(v bool) {
	st.column++
	st.stmt.BindBool(st.column, v)
}

// BindNullBool binds the next statement parameter to the specified bool value
// if it's true, or to NULL if it's false.
func (st *St) BindNullBool(v bool) {
	if !v {
		st.BindNull()
	} else {
		st.BindBool(v)
	}
}

// BindFloat binds the next statement parameter to the specified float value.
func (st *St) BindFloat(v float64) {
	st.column++
	st.stmt.BindFloat(st.column, v)
}

// BindNullFloat binds the next statement parameter to the specified float value
// if it's non-zero, or to NULL if it's zero.
func (st *St) BindNullFloat(v float64) {
	if v == 0 {
		st.BindNull()
	} else {
		st.BindFloat(v)
	}
}

// BindInt binds the next statement parameter to the specified int value.
func (st *St) BindInt(v int) {
	st.column++
	st.stmt.BindInt64(st.column, int64(v))
}

// BindNullInt binds the next statement parameter to the specified int value if
// it's non-zero, or to NULL if it's zero.
func (st *St) BindNullInt(v int) {
	if v == 0 {
		st.BindNull()
	} else {
		st.BindInt(v)
	}
}

// BindText binds the next statement parameter to the specified text value.
func (st *St) BindText(v string) {
	st.column++
	st.stmt.BindText(st.column, v)
}

// BindNullText binds the next statement parameter to the specified text value
// if it's non-empty, or to NULL if it's empty.
func (st *St) BindNullText(v string) {
	if v == "" {
		st.BindNull()
	} else {
		st.BindText(v)
	}
}

// BindTime binds the next statement parameter to the specified time value,
// using the supplied format.  If the time is zero, it binds an empty string.
func (st *St) BindTime(v time.Time, fmt string) {
	st.column++
	if v.IsZero() {
		st.stmt.BindText(st.column, "")
	} else {
		st.stmt.BindText(st.column, v.Format(fmt))
	}
}

// BindNullTime binds the next statement parameter to the specified time value
// if it's non-zero, or to NULL if it's zero.
func (st *St) BindNullTime(v time.Time, fmt string) {
	if v.IsZero() {
		st.BindNull()
	} else {
		st.BindTime(v, fmt)
	}
}

// BindNull binds the next statement parameter to NULL.
func (st *St) BindNull() {
	st.column++
	st.stmt.BindNull(st.column)
}

// Step reads the next row returned by the statement, if there is one, and
// returns whether there was one.
func (st *St) Step() (found bool) {
	var err error

	st.column = 0
	if found, err = st.stmt.Step(); err != nil {
		panic(err)
	}
	return found
}

// ColumnBlob reads the next statement column as a byte slice.  If the column is
// NULL, it returns nil.
func (st *St) ColumnBlob() (v []byte) {
	if !st.ColumnIsNull() {
		v = make([]byte, st.stmt.ColumnLen(st.column))
		st.stmt.ColumnBytes(st.column, v)
	}
	st.column++
	return v
}

// ColumnBool reads the next statement column as a bool value.  If the column is
// NULL, it returns false.
func (st *St) ColumnBool() (v bool) {
	if !st.ColumnIsNull() {
		v = st.stmt.ColumnBool(st.column)
	}
	st.column++
	return v
}

// ColumnFloat reads the next statement column as a float value.  If the column
// is NULL, it returns 0.
func (st *St) ColumnFloat() (v float64) {
	if !st.ColumnIsNull() {
		v = st.stmt.ColumnFloat(st.column)
	}
	st.column++
	return v
}

// ColumnInt reads the next statement column as an int value.  If the column is
// NULL, it returns 0.
func (st *St) ColumnInt() (v int) {
	if !st.ColumnIsNull() {
		v = st.stmt.ColumnInt(st.column)
	}
	st.column++
	return v
}

// ColumnText reads the next statement column as a string value.  If the column
// is NULL, it returns an empty string.
func (st *St) ColumnText() (v string) {
	if !st.ColumnIsNull() {
		v = st.stmt.ColumnText(st.column)
	}
	st.column++
	return v
}

// ColumnTime reads the next statement column as a time value with the supplied
// format.  If the column is NULL or contains an empty string, it returns the
// zero time.  If the column value can't be parsed with the supplied format,
// ColumnTime panics.
func (st *St) ColumnTime(fmt string) (v time.Time) {
	if !st.ColumnIsNull() {
		txt := st.stmt.ColumnText(st.column)
		if txt != "" {
			var err error
			if v, err = time.ParseInLocation(fmt, txt, time.Local); err != nil {
				panic("value in database column cannot be parsed as timestamp")
			}
		}
	}
	st.column++
	return v
}

// ColumnIsNull returns whether the next statement column has a NULL value.  It
// does not consume the column; the next ColumnXXX method will read it.
func (st *St) ColumnIsNull() bool {
	return st.stmt.ColumnType(st.column) == sqlite.TypeNull
}

// Reset resets the statement so that it can be run again.  It does not clear
// out the prior bound parameters, but it does reset the sequence, so if any
// BindXXX calls are going to be made after Reset, they all should be.
func (st *St) Reset() {
	st.column = 0
	if err := st.stmt.Reset(); err != nil {
		panic(err)
	}
}
