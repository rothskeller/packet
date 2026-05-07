// Package interval defines the Interval data type, which represents an
// algorithmically-defined series of moments in time.
package interval

import (
	"time"
)

// An Interval represents an algorithmically-defined series of moments in time,
// to one-minute granularity.
type Interval interface {
	// Match returns true if the Interval contains the specified time.
	// Seconds and Nanoseconds are ignored.
	Match(time.Time) bool
}

// neverInterval never matches any time.
type neverInterval struct{}

func (neverInterval) Match(time.Time) bool { return false }

// andInterval matches the logical AND of its child intervals.
type andInterval struct{ is []Interval }

func (ai andInterval) Match(t time.Time) bool {
	for _, i := range ai.is {
		if !i.Match(t) {
			return false
		}
	}
	return true
}

// orInterval matches the logical OR of its child intervals.
type orInterval struct{ is []Interval }

func (oi orInterval) Match(t time.Time) bool {
	for _, i := range oi.is {
		if i.Match(t) {
			return true
		}
	}
	return false
}

// notInterval matches the logical NOT of its child intervals.
type notInterval struct{ i Interval }

func (ni notInterval) Match(t time.Time) bool {
	return !ni.i.Match(t)
}

// termInterval matches one component of the time against a list of allowed
// values.
type termInterval struct {
	k string
	v []int
}

func (ti termInterval) Match(t time.Time) bool {
	var value int

	switch ti.k {
	case "YEAR":
		value = t.Year()
	case "MONTH":
		value = int(t.Month())
	case "WEEKDAY":
		value = int(t.Weekday())
	case "DAY":
		value = t.Day()
	case "HOUR":
		value = t.Hour()
	case "MINUTE":
		value = t.Minute()
	}
	for _, v := range ti.v {
		if v == value {
			return true
		}
	}
	return false
}
