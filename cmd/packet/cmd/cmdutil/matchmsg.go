package cmdutil

import (
	"strconv"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/messageid"
)

type MatchMessageFlag uint8

const (
	// MMNoRemote indicates that the target message cannot be identified by
	// its remote ID, only by its local ID or log entry number.
	MMNoRemote MatchMessageFlag = 1 << iota
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

// MatchMessage finds the message or pseudo-message identified by the input
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
func MatchMessage(i *incident.Incident, in string, flags MatchMessageFlag) (msg message.Message, le *incident.LogEntry, err error) {
	if flags&MMMessageOnly == 0 && strings.HasPrefix("configuration", in) {
		return &ConfigMessage{Config: i.Config}, nil, nil
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
				return &LogEntryMessage{LogEntry: e}, e, nil
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
