package xscmsg

import (
	"regexp"
	"time"
)

var callSignRE = regexp.MustCompile(`^[A-Z][A-Z0-9]{2,5}$`)

// ValidCallSign verifies that a string contains a call sign, which may be
// either FCC or tactical.
func ValidCallSign(value string) bool {
	return callSignRE.MatchString(value)
}

// ValidDate verifies that a string contains a valid date.
func ValidDate(value string) bool {
	if t, err := time.ParseInLocation("01/02/2006", value, time.Local); err != nil || value != t.Format("01/02/2006") {
		return false
	}
	return true
}

var fccCallSignRE = regexp.MustCompile(`^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})$`)

// ValidFCCCallSign verifies that a string contains an FCC call sign.
func ValidFCCCallSign(value string) bool {
	return fccCallSignRE.MatchString(value)
}

var messageIDRE = regexp.MustCompile(`^(?:[0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-(?:[1-9][0-9]{3,}|[0-9]{3})[A-Z]?$`)

// ValidMessageID verifies that a string contains a message ID.  (It need not
// be a packet message ID.)
func ValidMessageID(value string) bool {
	return messageIDRE.MatchString(value)
}

var packetMessageIDRE = regexp.MustCompile(`^(?:[0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-(?:[1-9][0-9]{3,}|[0-9]{3})[PMR]$`)

// ValidPacketMessageID verifies that a string contains a packet message ID.
func ValidPacketMessageID(value string) bool {
	return packetMessageIDRE.MatchString(value)
}

// ValidRestrictedValue verifies that a string contains one of the specified
// allowed values.
func ValidRestrictedValue(value string, allowed []string) bool {
	for _, a := range allowed {
		if value == a {
			return true
		}
	}
	return false
}

var timeRE = regexp.MustCompile(`^(?:[01][0-9]|2[0-3]):[0-5][0-9]$`)

// ValidTime verifies that a string contains a valid time.
func ValidTime(value string) bool {
	return value == "24:00" || timeRE.MatchString(value)
}
