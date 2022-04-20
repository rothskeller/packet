package config

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// An Interval represents a recurring moment in time, such as "every Tuesday at
// 20:00" or "every 5 minutes".
type Interval struct {
	year    []int
	month   []int
	day     []int
	weekday []int
	hour    []int
	minute  []int
}

// ParseInterval decodes the string form of an interval.  The string form is a
// whitespace-separated sequence of one or more terms.  Each term has the form
// of a keyword, an equals sign, and a list of values.  The keywords are YEAR,
// MONTH, DAY, WEEKDAY, HOUR, and MINUTE (not case sensitive).  The list of
// values is comma-separated, and may include ranges with a dash.  Values for
// weekday are strings (Monday, Tu, W, etc.); values for the other keywords are
// integers.
func ParseInterval(s string) (i *Interval, err error) {
	var (
		fields []string
		parts  []string
	)
	i = new(Interval)
	fields = strings.Fields(s)
	if len(fields) == 0 {
		return nil, errors.New("interval specification must have at least one term")
	}
	for _, field := range fields {
		parts = strings.Split(field, "=")
		if len(parts) == 1 {
			return nil, errors.New("interval specification has a term without an equals sign")
		} else if len(parts) > 2 {
			return nil, errors.New("interval specification has a term with multiple equals signs")
		}
		switch strings.ToLower(parts[0]) {
		case "year":
			err = addValues(parts[1], &i.year, 2020, 2099, strconv.Atoi)
		case "month":
			err = addValues(parts[1], &i.month, 1, 12, strconv.Atoi)
		case "day":
			err = addValues(parts[1], &i.day, 1, 31, strconv.Atoi)
		case "weekday":
			err = addValues(parts[1], &i.weekday, 0, 6, parseWeekday)
		case "hour":
			err = addValues(parts[1], &i.hour, 0, 23, strconv.Atoi)
		case "minute":
			err = addValues(parts[1], &i.minute, 0, 59, strconv.Atoi)
		default:
			err = fmt.Errorf("interval specification has an unknown keyword %q", parts[0])
		}
		if err != nil {
			return nil, err
		}
	}
	sortUniq(&i.year)
	sortUniq(&i.month)
	sortUniq(&i.day)
	sortUniq(&i.weekday)
	sortUniq(&i.hour)
	sortUniq(&i.minute)
	return i, nil
}

// addValues adds a set of values from the string s to the slice intlist.  Each
// value must be between min and max.  Each value is parsed from a string into
// an integer by calling parseFunc.  Note that the resulting intlist is
// unordered and may contain duplicates.
func addValues(s string, intlist *[]int, min, max int, parseFunc func(string) (int, error)) (err error) {
	var (
		ranges []string
		bounds []string
		start  int
		end    int
	)
	ranges = strings.Split(s, ",")
	for _, r := range ranges {
		bounds = strings.Split(r, "-")
		switch len(bounds) {
		case 1:
			if start, err = parseFunc(bounds[0]); err != nil {
				return errors.New("interval specification has an invalid value")
			}
			end = start
		case 2:
			if start, err = parseFunc(bounds[0]); err != nil {
				return errors.New("interval specification has an invalid value")
			}
			if end, err = parseFunc(bounds[1]); err != nil {
				return errors.New("interval specification has an invalid value")
			}
		default:
			return errors.New("interval specification has a list item with multiple dashes")
		}
		if start < min || end < start || end > max {
			return errors.New("interval specification has invalid values")
		}
		for i := start; i <= end; i++ {
			*intlist = append(*intlist, i)
		}
	}
	return nil
}

// parseWeekday is a parsing function used by addValues, which parses weekday
// names.  They are not case-sensitive, and any unambiguous abbreviation can be
// used.
func parseWeekday(s string) (weekday int, err error) {
	s = strings.ToLower(s)
	if len(s) >= 2 && len(s) <= 6 && s == "sunday"[:len(s)] {
		return int(time.Sunday), nil
	}
	if len(s) >= 1 && len(s) <= 6 && s == "monday"[:len(s)] {
		return int(time.Monday), nil
	}
	if len(s) >= 2 && len(s) <= 7 && s == "tuesday"[:len(s)] {
		return int(time.Tuesday), nil
	}
	if (len(s) >= 1 && len(s) <= 9 && s == "wednesday"[:len(s)]) || s == "weds" {
		return int(time.Wednesday), nil
	}
	if len(s) >= 2 && len(s) <= 8 && s == "thursday"[:len(s)] {
		return int(time.Thursday), nil
	}
	if len(s) >= 1 && len(s) <= 6 && s == "friday"[:len(s)] {
		return int(time.Friday), nil
	}
	if len(s) >= 2 && len(s) <= 8 && s == "saturday"[:len(s)] {
		return int(time.Saturday), nil
	}
	return 0, errors.New("interval specification contains an unknown weekday")
}

// sortUniq sorts a list of integers and removes duplicates.
func sortUniq(intlist *[]int) {
	sort.Ints(*intlist)
	j := 0
	for _, v := range *intlist {
		if j == 0 || v != (*intlist)[j-1] {
			(*intlist)[j] = v
			j++
		}
	}
	*intlist = (*intlist)[:j]
}

// Next returns the next instance of the interval occurring after the specified
// time.  If there is no next interval, it returns the zero time.
func (i *Interval) Next(from time.Time) (next time.Time) {
	next = time.Date(from.Year(), from.Month(), from.Day(), from.Hour(), from.Minute()+1, 0, 0, time.Local)
RESTART:
	// If we're not in a selected year, skip to one that is selected.
	if len(i.year) != 0 {
		var found bool
		for _, y := range i.year {
			if y > next.Year() {
				next = time.Date(y, 1, 1, 0, 0, 0, 0, time.Local)
			}
			if y == next.Year() {
				found = true
				break
			}
		}
		if !found {
			return time.Time{}
		}
	}
	// If we're not in a selected month, skip to one that is selected.
	if len(i.month) != 0 {
		var found bool
		for _, m := range i.month {
			if time.Month(m) > next.Month() {
				next = time.Date(next.Year(), time.Month(m), 1, 0, 0, 0, 0, time.Local)
			}
			if time.Month(m) == next.Month() {
				found = true
				break
			}
		}
		if !found {
			next = time.Date(next.Year()+1, 1, 1, 0, 0, 0, 0, time.Local)
			goto RESTART
		}
	}
	// If we're not on a selected day, skip to one that is selected.
	if len(i.day) != 0 {
		var found bool
		for _, d := range i.day {
			if d > next.Day() {
				maybe := time.Date(next.Year(), next.Month(), d, 0, 0, 0, 0, time.Local)
				if maybe.Month() == next.Month() {
					next = maybe
				} else {
					break
				}
			}
			if d == next.Day() {
				found = true
				break
			}
		}
		if !found {
			next = time.Date(next.Year(), next.Month()+1, 1, 0, 0, 0, 0, time.Local)
			goto RESTART
		}
	}
	// If we're not on a selected day of the week, try another day.
	if len(i.weekday) != 0 {
		var found bool
		for _, d := range i.weekday {
			if time.Weekday(d) == next.Weekday() {
				found = true
				break
			}
		}
		if !found {
			next = time.Date(next.Year(), next.Month(), next.Day()+1, 0, 0, 0, 0, time.Local)
			goto RESTART
		}
	}
	// If we're not at a selected hour, skip to one that is selected.
	if len(i.hour) != 0 {
		var found bool
		for _, h := range i.hour {
			if h > next.Hour() {
				next = time.Date(next.Year(), next.Month(), next.Day(), h, 0, 0, 0, time.Local)
			}
			if h == next.Hour() {
				found = true
				break
			}
		}
		if !found {
			next = time.Date(next.Year(), next.Month(), next.Day()+1, 0, 0, 0, 0, time.Local)
			goto RESTART
		}
	}
	// If we're not at a selected minute, skip to one that is selected.
	if len(i.minute) != 0 {
		var found bool
		for _, m := range i.minute {
			if m > next.Minute() {
				next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), m, 0, 0, time.Local)
			}
			if m == next.Minute() {
				found = true
				break
			}
		}
		if !found {
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour()+1, 0, 0, 0, time.Local)
			goto RESTART
		}
	}
	return next
}

// Prev returns the instance of the interval occurring before the specified
// time.  If there is no previous interval, it returns the zero time.
func (i *Interval) Prev(from time.Time) (prev time.Time) {
	if from.Second() == 0 && from.Nanosecond() == 0 {
		prev = time.Date(from.Year(), from.Month(), from.Day(), from.Hour(), from.Minute()-1, 0, 0, time.Local)
	} else {
		prev = time.Date(from.Year(), from.Month(), from.Day(), from.Hour(), from.Minute(), 0, 0, time.Local)
	}
RESTART:
	// If we're not in a selected year, skip to one that is selected.
	if len(i.year) != 0 {
		var found bool
		for idx := len(i.year) - 1; idx >= 0; idx-- {
			y := i.year[idx]
			if y < prev.Year() {
				prev = time.Date(y, 12, 31, 23, 59, 0, 0, time.Local)
			}
			if y == prev.Year() {
				found = true
				break
			}
		}
		if !found {
			return time.Time{}
		}
	}
	// If we're not in a selected month, skip to one that is selected.
	if len(i.month) != 0 {
		var found bool
		for idx := len(i.month) - 1; idx >= 0; idx-- {
			m := i.month[idx]
			if time.Month(m) < prev.Month() {
				prev = time.Date(prev.Year(), time.Month(m+1), 0, 23, 59, 0, 0, time.Local)
			}
			if time.Month(m) == prev.Month() {
				found = true
				break
			}
		}
		if !found {
			prev = time.Date(prev.Year()-1, 12, 31, 23, 59, 0, 0, time.Local)
			goto RESTART
		}
	}
	// If we're not on a selected day, skip to one that is selected.
	if len(i.day) != 0 {
		var found bool
		for idx := len(i.day) - 1; idx >= 0; idx-- {
			d := i.day[idx]
			if d < prev.Day() {
				maybe := time.Date(prev.Year(), prev.Month(), d, 23, 59, 0, 0, time.Local)
				if maybe.Month() == prev.Month() {
					prev = maybe
				} else {
					break
				}
			}
			if d == prev.Day() {
				found = true
				break
			}
		}
		if !found {
			prev = time.Date(prev.Year(), prev.Month(), 0, 23, 59, 0, 0, time.Local)
			goto RESTART
		}
	}
	// If we're not on a selected day of the week, try another day.
	if len(i.weekday) != 0 {
		var found bool
		for idx := len(i.weekday) - 1; idx >= 0; idx-- {
			d := i.weekday[idx]
			if time.Weekday(d) == prev.Weekday() {
				found = true
				break
			}
		}
		if !found {
			prev = time.Date(prev.Year(), prev.Month(), prev.Day()-1, 23, 59, 0, 0, time.Local)
			goto RESTART
		}
	}
	// If we're not at a selected hour, skip to one that is selected.
	if len(i.hour) != 0 {
		var found bool
		for idx := len(i.hour) - 1; idx >= 0; idx-- {
			h := i.hour[idx]
			if h < prev.Hour() {
				prev = time.Date(prev.Year(), prev.Month(), prev.Day(), h, 59, 0, 0, time.Local)
			}
			if h == prev.Hour() {
				found = true
				break
			}
		}
		if !found {
			prev = time.Date(prev.Year(), prev.Month(), prev.Day()-1, 23, 59, 0, 0, time.Local)
			goto RESTART
		}
	}
	// If we're not at a selected minute, skip to one that is selected.
	if len(i.minute) != 0 {
		var found bool
		for idx := len(i.minute) - 1; idx >= 0; idx-- {
			m := i.minute[idx]
			if m < prev.Minute() {
				prev = time.Date(prev.Year(), prev.Month(), prev.Day(), prev.Hour(), m, 0, 0, time.Local)
			}
			if m == prev.Minute() {
				found = true
				break
			}
		}
		if !found {
			prev = time.Date(prev.Year(), prev.Month(), prev.Day(), prev.Hour()-1, 59, 0, 0, time.Local)
			goto RESTART
		}
	}
	return prev
}

// Match returns whether the specified time matches the interval specification.
// (Seconds and nanoseconds are ignored.)
func (i *Interval) Match(t time.Time) bool {
	if len(i.year) != 0 {
		var found bool
		for _, y := range i.year {
			if y == t.Year() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if len(i.month) != 0 {
		var found bool
		for _, m := range i.month {
			if time.Month(m) == t.Month() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if len(i.day) != 0 {
		var found bool
		for _, d := range i.day {
			if d == t.Day() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if len(i.weekday) != 0 {
		var found bool
		for _, d := range i.weekday {
			if time.Weekday(d) == t.Weekday() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if len(i.hour) != 0 {
		var found bool
		for _, h := range i.hour {
			if h == t.Hour() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if len(i.minute) != 0 {
		var found bool
		for _, m := range i.minute {
			if m == t.Minute() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
