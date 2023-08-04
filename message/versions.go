package message

import (
	"strconv"
	"strings"
)

// OlderVersion compares two version numbers, and returns true if the first one
// is older than the second one.  Each version number is a dot-separated
// sequence of parts, each of which is compared independently with the
// corresponding part in the other version number.  The parts are compared
// numerically if they parse as integers, and as strings otherwise.  However, a
// part that starts with a digit is always "newer" than a part that doesn't.
// (This results in "undefined" being treated as infinitely old, which is what
// we want.)
func OlderVersion(a, b string) bool {
	aparts := strings.Split(a, ".")
	bparts := strings.Split(b, ".")
	for len(aparts) != 0 && len(bparts) != 0 {
		anum, aerr := strconv.Atoi(aparts[0])
		bnum, berr := strconv.Atoi(bparts[0])
		if aerr == nil && berr == nil {
			if anum < bnum {
				return true
			}
			if anum > bnum {
				return false
			}
		} else if startsWithDigit(aparts[0]) && !startsWithDigit(bparts[0]) {
			return false
		} else if !startsWithDigit(aparts[0]) && startsWithDigit(bparts[0]) {
			return true
		} else {
			if aparts[0] < bparts[0] {
				return true
			}
			if aparts[0] > bparts[0] {
				return false
			}
		}
		aparts = aparts[1:]
		bparts = bparts[1:]
	}
	return len(bparts) != 0
}
func startsWithDigit(s string) bool {
	return s != "" && s[0] >= '0' && s[0] <= '9'
}
