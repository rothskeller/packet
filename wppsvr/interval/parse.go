// Package interval defines the Interval data type, which represents an
// algorithmically-defined series of moments in time.
package interval

import (
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Parse parses a string representing an Interval and returns that Interval.  If
// the string cannot be parsed as a valid Interval, ParseInterval returns nil.
// The string form of an Interval has the following syntax (entirely
// case-insensitive):
//
//     interval = real-interval
//              | "NEVER"
//
//     real-interval = simple-interval "AND" simple-interval
//                   | simple-interval "OR" simple-interval
//                   | simple-interval
//
//     simple-interval = "NOT" positive-interval
//                     | positive-interval
//
//     positive-interval = "(" real-interval ")"
//                       | interval-term { interval-term }
//
//     interval-term = keyword "=" values
//
//     keyword = "YEAR" | "MONTH" | "WEEKDAY" | "DAY" | "HOUR" | "MINUTE"
//
//     values = value
//            | value "-" value
//            | value "," values
//
//     value = integer   /* for keywords other than WEEKDAY */
//           | month     /* for keyword MONTH only */
//           | weekday   /* for keyword WEEKDAY only */
//
//     month = /* any unambiguous prefix of a month name */
//
//     weekday = /* any unambiguous prefix of a weekday name */
//             | "WEDS"
//
// Semantically, an interval comprises all moments in time (to one-minute
// granularity) that satisfy all of the interval terms in the interval, with
// ANDs, ORs, NOTs, and parentheses interpreted as one would expect.  The
// interval "NEVER" has no moments in time.
func Parse(s string) Interval {
	var (
		tokens []string
		token  string
	)
	// Tokenize the string.
	for _, c := range s {
		if unicode.IsSpace(c) {
			if token != "" {
				tokens = append(tokens, token)
				token = ""
			}
			continue
		}
		if c == '(' || c == ')' || c == '=' || c == ',' || c == '-' {
			if token != "" {
				tokens = append(tokens, token)
				token = ""
			}
			tokens = append(tokens, string(c))
			continue
		}
		token += string(unicode.ToUpper(c))
	}
	if token != "" {
		tokens = append(tokens, token)
	}
	return parseInterval(tokens)
}

func parseInterval(tokens []string) (i Interval) {
	//     interval = real-interval
	//              | "NEVER"
	if len(tokens) == 1 && tokens[0] == "NEVER" {
		return orInterval{nil}
	}
	if tokens, i := parseRealInterval(tokens); i != nil && len(tokens) == 0 {
		return i
	}
	return nil
}

func parseRealInterval(tokens []string) ([]string, Interval) {
	//     real-interval = simple-interval "AND" simple-interval
	//                   | simple-interval "OR" simple-interval
	//                   | simple-interval
	var s1, s2 Interval
	if tokens, s1 = parseSimpleInterval(tokens); s1 == nil {
		return nil, nil
	}
	if len(tokens) != 0 && (tokens[0] == "AND" || tokens[0] == "OR") {
		conj := tokens[0]
		if tokens, s2 = parseSimpleInterval(tokens[1:]); s2 == nil {
			return nil, nil
		}
		if conj == "AND" {
			return tokens, andInterval{[]Interval{s1, s2}}
		}
		return tokens, orInterval{[]Interval{s1, s2}}
	}
	return tokens, s1
}

func parseSimpleInterval(tokens []string) ([]string, Interval) {
	//     simple-interval = "NOT" positive-interval
	//                     | positive-interval
	if len(tokens) != 0 && tokens[0] == "NOT" {
		if tokens, i := parsePositiveInterval(tokens[1:]); i != nil {
			return tokens, notInterval{i}
		}
		return nil, nil
	}
	return parsePositiveInterval(tokens)
}

func parsePositiveInterval(tokens []string) ([]string, Interval) {
	//     positive-interval = "(" real-interval ")"
	//                       | interval-term { interval-term }
	var (
		terms []Interval
		term  Interval
	)
	if len(tokens) != 0 && tokens[0] == "(" {
		if tokens, i := parseRealInterval(tokens[1:]); i != nil {
			if len(tokens) != 0 && tokens[0] == ")" {
				return tokens[1:], i
			}
		}
		return nil, nil
	}
	if tokens, term = parseIntervalTerm(tokens); term == nil {
		return nil, nil
	}
	terms = append(terms, term)
	for term != nil {
		var nt []string
		nt, term = parseIntervalTerm(tokens)
		if term != nil {
			terms = append(terms, term)
			tokens = nt
		}
	}
	if len(terms) == 1 {
		return tokens, terms[0]
	}
	return tokens, andInterval{terms}
}

func parseIntervalTerm(tokens []string) ([]string, Interval) {
	//     interval-term = keyword "=" values
	//     keyword = "YEAR" | "MONTH" | "WEEKDAY" | "DAY" | "HOUR" | "MINUTE"
	var (
		keyword string
		list    []int
	)
	if len(tokens) < 3 || tokens[1] != "=" {
		return nil, nil
	}
	switch keyword = strings.ToUpper(tokens[0]); keyword {
	case "YEAR":
		tokens, list = parseValues(tokens[2:], parseYear)
	case "MONTH":
		tokens, list = parseValues(tokens[2:], parseMonth)
	case "WEEKDAY":
		tokens, list = parseValues(tokens[2:], parseWeekday)
	case "DAY":
		tokens, list = parseValues(tokens[2:], parseDay)
	case "HOUR":
		tokens, list = parseValues(tokens[2:], parseHour)
	case "MINUTE":
		tokens, list = parseValues(tokens[2:], parseMinute)
	default:
		return nil, nil
	}
	if list == nil {
		return nil, nil
	}
	return tokens, termInterval{keyword, list}
}

func parseValues(tokens []string, parseValue func(string) int) (_ []string, values []int) {
	//     values = value
	//            | value "-" value
	//            | value "," values
	for {
		if len(tokens) == 0 {
			return nil, nil
		}
		value := parseValue(tokens[0])
		if value < 0 {
			return nil, nil
		}
		tokens = tokens[1:]
		if len(tokens) == 0 {
			values = append(values, value)
			return tokens, uniqueValues(values)
		}
		if tokens[0] == "-" {
			if len(tokens) == 1 {
				return nil, nil
			}
			value2 := parseValue(tokens[1])
			if value2 < value {
				return nil, nil
			}
			tokens = tokens[2:]
			for i := value; i <= value2; i++ {
				values = append(values, i)
			}
		} else {
			values = append(values, value)
		}
		if len(tokens) == 0 || tokens[0] != "," {
			return tokens, uniqueValues(values)
		}
		tokens = tokens[1:]
	}
}
func uniqueValues(values []int) []int {
	sort.Ints(values)
	j := 0
	for _, value := range values {
		if j == 0 || value != values[j-1] {
			values[j] = value
			j++
		}
	}
	return values[:j]
}

var monthNames = []string{"", "JANUARY", "FEBRUARY", "MARCH", "APRIL", "MAY", "JUNE", "JULY", "AUGUST", "SEPTEMBER", "OCTOBER", "NOVEMBER", "DECEMBER"}
var weekdayNames = []string{"SUNDAY", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY"}

func parseYear(s string) int { return parseInt(s, 2020, 2050) }
func parseMonth(s string) int {
	if s != "" && s[0] >= '0' && s[0] <= '9' {
		return parseInt(s, 1, 12)
	}
	return parseAbbrev(s, monthNames)
}
func parseWeekday(s string) int {
	if v := parseAbbrev(s, weekdayNames); v != -1 {
		return v
	}
	if strings.EqualFold(s, "WEDS") {
		return 3
	}
	return -1
}
func parseDay(s string) int    { return parseInt(s, 1, 31) }
func parseHour(s string) int   { return parseInt(s, 0, 23) }
func parseMinute(s string) int { return parseInt(s, 0, 59) }

func parseInt(s string, min, max int) (value int) {
	var err error

	if value, err = strconv.Atoi(s); err != nil || value < min || value > max {
		return -1
	}
	return value
}
func parseAbbrev(s string, values []string) int {
	var found = -1
	var multiple bool

	s = strings.ToUpper(s)
	for i, v := range values {
		if s == v {
			return i
		}
		if len(s) < len(v) && s == v[:len(s)] {
			if found != -1 {
				multiple = true
			} else {
				found = i
			}
		}
	}
	if !multiple {
		return found
	}
	return -1
}
