package message

import (
	"strings"
)

// OlderVersion compares two version numbers, and returns true if the first one
// is older than the second one.  Each version number is split into a sequence
// of dot-separated parts.  Any such parts that start with a digit and contain
// non-digits are further split into the part with the leading digits and the
// remainder.  The parts in the first version number are then compared one by
// one with the corresponding parts of the second number until a difference is
// found.  When both parts contain only digits, they are compared numerically
// (e.g. 2 < 10).  When both parts contain non-digits, they are compared as
// case-insensitive strings.  When the parts are of mixed type, a part
// containing non-digits is "older" than the absence of a part, which is "older"
// than a part containing only digits.
//
// This algorithm ensures that "undefined" is older than any numeric version,
// and that "3.15alpha" is older than "3.15".
func OlderVersion(a, b string) bool {
	for {
		var apart, bpart string
		var astr, bstr bool

		if a == "" && b == "" {
			return false
		}
		apart, a, astr = versionPart(a)
		bpart, b, bstr = versionPart(b)
		switch {
		case apart == "":
			switch {
			case bpart == "":
				break
			case bstr:
				return false
			default: // !bstr
				return true
			}
		case astr:
			switch {
			case bpart == "":
				return true
			case bstr:
				apart = strings.ToLower(apart)
				bpart = strings.ToLower(bpart)
				if apart != bpart {
					return apart < bpart
				}
			default: // !bstr
				return true
			}
		default: // !astr
			switch {
			case bpart == "":
				return false
			case bstr:
				return false
			default: // !bstr
				if apart == bpart {
					// nothing
				} else if len(apart) == len(bpart) {
					return apart < bpart
				} else {
					return len(apart) < len(bpart)
				}
			}
		}
	}
}
func versionPart(s string) (part, rest string, nonDigits bool) {
LOOP:
	for s != "" {
		switch s[0] {
		case '.':
			s = s[1:]
			break LOOP
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			part += string(s[0])
			s = s[1:]
		default:
			if part == "" || nonDigits {
				nonDigits = true
				part += string(s[0])
				s = s[1:]
			} else {
				break LOOP
			}
		}
	}
	for !nonDigits && len(part) > 1 && part[0] == '0' {
		part = part[1:]
	}
	return part, s, nonDigits
}

// SmartJoin joins the two provided strings with the provided separator, but
// only when both are non-empty.  If one is empty, it returns the other.  If
// both are empty, it returns an empty string.
func SmartJoin(a, b, sep string) string {
	if a != "" && b != "" {
		return a + sep + b
	}
	if a == "" {
		return b
	}
	return a
}
