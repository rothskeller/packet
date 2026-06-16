// Package messageid has Encode and Decode methods to convert between standard
// Santa Clara County ARES/RACES message ID strings and their constituent parts
// (prefix, sequence number, and suffix).
package messageid

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/v4/errors"
)

var (
	msgIDPrefixRE      = regexp.MustCompile(`^[A-Z][A-Z0-9][A-Z0-9]|[0-9][A-Z][A-Z]$`)
	msgIDSuffixRE      = regexp.MustCompile(`^[A-Z]?$`)
	ErrInvalidPrefix   = errors.New("The message ID prefix must be a three-character code, one of XXX X## X#X XX# #XX, where X is any uppercase letter and # is any digit.")
	ErrInvalidSequence = errors.New("The message ID sequence number must be >= 0.")
	ErrInvalidSuffix   = errors.New("The message ID suffix must be a single uppercase letter.")
	ErrNoHyphen        = errors.New("The message ID must contain a hyphen.")
	ErrNoMessageID     = errors.New("The message ID is empty.")
	ErrNoSequence      = errors.New("The message ID must have digit(s) after the hyphen.")
	ErrNoSuffix        = errors.New("The message ID for a packet message must have an uppercase letter suffix (usually \"P\") after the sequence number.")
)

// Decode decodes a standard-format message ID into its constituent parts.  It
// returns error(s) (and best-effort results) if the message ID is not in
// standard format.  If permissive is true, lowercase is accepted.  If packet
// is true, the ID is expected to have a suffix.
func Decode(id string, permissive, packet bool) (prefix string, sequence int, suffix string, err error) {
	var (
		found bool
		seq   string
	)
	if id == "" {
		return "", 0, "", ErrNoMessageID
	}
	if permissive {
		id = strings.ToUpper(id)
	}
	if prefix, suffix, found = strings.Cut(id, "-"); !found {
		return "", 0, "", ErrNoHyphen
	}
	if !msgIDPrefixRE.MatchString(prefix) {
		err = errors.Join(err, ErrInvalidPrefix)
	}
	if idx := strings.IndexFunc(suffix, nondigit); idx > 0 {
		seq, suffix = suffix[:idx], suffix[idx:]
	} else if idx == 0 || suffix == "" {
		err = errors.Join(err, ErrNoSequence)
		return "", 0, "", err
	} else {
		seq, suffix = suffix, ""
	}
	sequence, _ = strconv.Atoi(seq)
	if !msgIDSuffixRE.MatchString(suffix) {
		err = errors.Join(err, ErrInvalidSuffix)
	} else if packet && suffix == "" {
		err = errors.Join(err, ErrNoSuffix)
	}
	return prefix, sequence, suffix, err
}
func nondigit(r rune) bool { return r < '0' || r > '9' }

// Encode encodes a standard-format message ID from its constituent parts.  It
// returns an error if the parts are invalid.
func Encode(prefix string, sequence int, suffix string) (id string, err error) {
	prefix = strings.ToUpper(prefix)
	suffix = strings.ToUpper(suffix)
	if !msgIDPrefixRE.MatchString(prefix) {
		err = errors.Join(err, ErrInvalidPrefix)
	}
	if sequence < 0 {
		err = errors.Join(err, ErrInvalidSequence)
	}
	if !msgIDSuffixRE.MatchString(suffix) {
		err = errors.Join(err, ErrInvalidSuffix)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%03d%s", prefix, sequence, suffix), nil
}

// Cleanup decodes a message ID permissively and re-encodes it strictly.  It
// returns any problems with the message ID.  If the input string can't be
// decoded as a message ID at all, it is returned unmodified with errors.  If
// packet is true, the ID is expected to have a suffix.
func Cleanup(id string, packet bool) (strict string, err error) {
	if prefix, sequence, suffix, err := Decode(id, true, packet); err != nil {
		return id, err
	} else {
		strict, _ = Encode(prefix, sequence, suffix)
		return strict, nil
	}
}

// Increment increments the sequence number in a message ID.  If the message
// ID can't be decoded (permissively), it returns an error.
func Increment(id string) (string, error) {
	if pfx, seq, sfx, err := Decode(id, true, false); err != nil {
		return "", err
	} else {
		return Encode(pfx, seq+1, sfx)
	}
}
