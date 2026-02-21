package cmd

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/pseudomsg"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/messageid"
)

const paramsSlug = `parameter formats and shortcuts`
const paramsHelp = `
Many "packet" commands take a «message-id» argument to identify a message.  It can have any of the following formats (none of which is case sensitive):
  - ⇥the complete local message ID of the message
  - ⇥the complete remote message ID of the message, if it is unique
  - ⇥just the numeric part of the local or remote message ID, if it is unique
  - ⇥a pound sign followed by the log entry number of any log entry describing the message (e.g., "#23")
Some commands can act on either a message or a log entry.  When using the log entry number form with those commands, add the command's --message (or -m) flag to tell it to act on the message rather than the log entry itself.

Many "packet" commands take a «log-entry» argument to identify a log entry.  It can have any of the above forms.  When one of the message ID forms is used, the command acts on all log entries describing that message.

Many "packet" commands take a «field» argument to identify one field of a message, a log entry, or the incident configuration.  It can have any of the following formats (not case sensitive):
  - ⇥The exact field name
  - ⇥The exact field name with all spaces removed (to avoid the need for quoting)
  - ⇥An at-sign followed by the PackItForms tag for the field, including its trailing period if any (e.g., "@5.")
  - ⇥A shortened version of the field name, such as "ocs" for "Operator Call Sign."
The last of these is allowed only when standard output is a terminal.
`

type matchMessageFlag uint8

const (
	// MMNoRemote indicates that the target message cannot be identified by
	// its remote ID, only by its local ID or log entry number.
	MMNoRemote matchMessageFlag = 1 << iota
	// MMNoAbbrev indicates that the target message cannot be identified by
	// an abbreviated message ID, only by a full message ID or log entry
	// number.
	MMNoAbbrev
	// MMMessageOnly indicates that the target must be an actual packet
	// message, not "configuration" or a hand-entered log entry.  If the
	// input is a log entry number, the result is the message described by
	// that log entry, not a pseudo-message for the log entry itself.
	MMMessageOnly
)

// matchMessage finds the message or pseudo-message identified by the input
// string.  Depending on the flags provided, the input can be:
//   - A local message ID, in which case the corresponding message is returned.
//   - An unambiguous remote message ID, in which case the corresponding message
//     is returned.
//   - An integer that is unambiguously the sequence number of a local or remote
//     message ID, in which case the corresponding message is returned.
//   - A pound sign followed by an integer, in which case a pseudo-message for
//     the log entry with that number is returned.
//   - The word "configuration" (or an abbreviation), in which case a
//     pseudo-message for the incident configuration is returned.
func matchMessage(i *incident.Incident, in string, flags matchMessageFlag) (msg message.Message, le *incident.LogEntry, err error) {
	if flags&MMMessageOnly == 0 && strings.HasPrefix("configuration", in) {
		return &pseudomsg.ConfigMessage{Config: i.Config.Clone()}, nil, nil
	}
	if strings.HasPrefix(in, "#") {
		var num int
		if num, err = strconv.Atoi(in[1:]); err != nil || num < 1 {
			return nil, nil, errors.NewF("%q is not a valid log entry number", in)
		}
		for _, e := range i.Log {
			if e.Ident == num {
				if e.Status == incident.StatusDeleted {
					return nil, nil, errors.NewF("Log entry #%d has been deleted and is not recoverable.", num)
				}
				if flags&MMMessageOnly != 0 {
					if e.Status == incident.StatusHandEntered {
						return nil, nil, errors.NewF("Log entry #%d is not associated with a packet message.", num)
					}
					if msg, err = i.GetMessageFromLogEntry(e); err != nil {
						return nil, nil, err
					}
					return msg, e, nil
				}
				return &pseudomsg.LogEntryMessage{LogEntry: e}, e, nil
			}
		}
		return nil, nil, errors.NewF("There is no log entry #%d.", num)
	}
	var (
		entry     *incident.LogEntry
		ambiguous bool
		num       int
	)
	if num, _ = strconv.Atoi(in); num != 0 && flags&MMNoAbbrev != 0 {
		return nil, nil, errors.New("This command does not accept abbreviated message IDs.  Please provide the complete message ID.")
	}
	for _, e := range i.Log {
		var local string

		if e.Status == incident.StatusDeleted || e.Status == incident.StatusHandEntered {
			continue
		}
		if local = e.LocalMsgID; e.Flags&incident.FIsReceipt != 0 {
			local = ""
		}
		if strings.EqualFold(local, in) {
			entry, ambiguous = e, false
			break
		}
		if flags&MMNoRemote == 0 && (strings.EqualFold(e.FromMsgID, in) || strings.EqualFold(e.ToMsgID, in)) {
			ambiguous = ambiguous || entry != nil
			entry = e
			continue
		}
		if num != 0 {
			if _, n, _, err := messageid.Decode(local, true, false); err == nil && n == num {
				ambiguous = ambiguous || (entry != nil && entry.LocalMsgID != local)
				entry = e
			} else if flags&MMNoRemote == 0 {
				if _, n, _, err = messageid.Decode(e.FromMsgID, true, false); err == nil && n == num {
					ambiguous = ambiguous || (entry != nil && entry.LocalMsgID != local)
					entry = e
				} else if _, n, _, err = messageid.Decode(e.ToMsgID, true, false); err == nil && n == num {
					ambiguous = ambiguous || (entry != nil && entry.LocalMsgID != local)
					entry = e
				}
			}
		}
	}
	if entry == nil {
		return nil, nil, errors.NewF("There is no message %q.", in)
	} else if ambiguous {
		return nil, nil, errors.NewF("The string %q is ambiguous: it identifies multiple messages.", in)
	}
	if msg, err = i.GetMessageFromLogEntry(entry); err != nil {
		return nil, nil, err
	} else if msg == nil {
		return nil, nil, errors.NewF("The file for message %s was not found.", in)
	} else {
		return msg, entry, nil
	}
}

var fieldNumberRE = regexp.MustCompile(`^\d+(?:[A-Za-z])?\.$`)

// matchField finds the message field that (best) matches the supplied field
// name.  If loose is true, it can be a partial, heuristic match.
func matchField(msg message.Message, in string, loose bool) (field.Field, error) {
	// First priority is a match on PIFO tag.
	if strings.HasPrefix(in, "@") {
		for f := range msg.Fields() {
			if f.Tag() == in[1:] {
				return f, nil
			}
		}
	}
	// Remaining comparisons are case-sensitive if the input contains any
	// uppercase letters.
	var caseSensitive bool
	var compare func(string, string) bool
	if strings.IndexFunc(in, func(r rune) bool { return r >= 'A' && r <= 'Z' }) >= 0 {
		compare, caseSensitive = func(a, b string) bool { return a == b }, true
	} else {
		compare, caseSensitive = strings.EqualFold, false
	}
	// Second priority is a case-insensitive match on full field name.
	for f := range msg.Fields() {
		if compare(f.Label(), in) {
			return f, nil
		}
	}
	// To avoid the need for quoting, we will also accept a case-insensitive
	// match on the full field name with spaces removed.
	if !strings.Contains(in, " ") {
		for f := range msg.Fields() {
			if compare(strings.ReplaceAll(f.Label(), " ", ""), in) {
				return f, nil
			}
		}
	}
	// If the input looks like a field number, we'll accept a prefix match with
	// it.
	if fieldNumberRE.MatchString(in) {
		insp := strings.ToLower(in) + " "
		for f := range msg.Fields() {
			if strings.HasPrefix(strings.ToLower(f.Label()), insp) {
				return f, nil
			}
		}
	}
	// Unless the loose flag is set, those are the only options.
	if !loose {
		return nil, errors.NewF("There is no field %q.", in)
	}
	// Now we look for fields whose name contains the same characters as the
	// input, in the same order, but also contains additional characters.
	// We return the field that has the smallest number of unmatched
	// capital letters in its name, and among those, the one with the
	// smallest number of unmatched characters.
	var match field.Field
	var missedUC, missedCH int
	for f := range msg.Fields() {
		if mUC, mCH, ok := matchFieldName(f.Label(), in, caseSensitive); ok {
			if match == nil || mUC < missedUC || (mUC == missedUC && mCH < missedCH) {
				match, missedUC, missedCH = f, mUC, mCH
			}
		}
	}
	if match == nil {
		return nil, errors.NewF("There is no field %q.", in)
	}
	return match, nil
}

// matchFieldName is a recursive function that determines whether the input is a
// valid shortening of the field name, and returns the heuristic scoring if so.
// It's a pretty expensive algorithm, but at human time scales it's negligible.
func matchFieldName(fname, in string, caseSensitive bool) (mUC, mCH int, ok bool) {
	if in == "" && fname == "" {
		// Nothing left of either string.  Perfect match.
		return 0, 0, true
	}
	if in == "" {
		// Nothing left of in, but we still have some fname.  Compute
		// the score.  Use a recursive call to get the score after
		// removing the first character of fname, then add the score
		// for that character.
		mUC, mCH, _ = matchFieldName(fname[1:], in, caseSensitive)
		if fname[0] >= 'A' && fname[0] <= 'Z' {
			mUC++
		}
		mCH++
		return mUC, mCH, true
	}
	if fname == "" {
		// Nothing left of fname, but we still have some in.  Not a
		// match at all.
		return 0, 0, false
	}
	var mUC1, mCH1, mUC2, mCH2 int
	var ok1, ok2 bool
	// If the lead characters of fname and in match, calculate the score
	// based on matching those two.
	if fname[0] == in[0] || (!caseSensitive && downcase(fname[0]) == downcase(in[0])) {
		mUC1, mCH1, ok1 = matchFieldName(fname[1:], in[1:], caseSensitive)
	}
	// Whether the lead characters match or not, also calculate the score
	// assuming they don't.
	mUC2, mCH2, ok2 = matchFieldName(fname[1:], in, caseSensitive)
	if fname[0] >= 'A' && fname[0] <= 'Z' {
		mUC2++
	}
	mCH2++
	// Return the better of the two scores.
	if !ok1 && !ok2 {
		return 0, 0, false
	}
	if !ok1 {
		return mUC2, mCH2, ok2
	}
	if !ok2 || mUC1 < mUC2 || (mUC1 == mUC2 && mCH1 < mCH2) {
		return mUC1, mCH1, ok1
	}
	return mUC2, mCH2, ok2
}

func downcase(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b + 'A' - 'a'
	}
	return b
}
