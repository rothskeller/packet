package xscmsg

import (
	"fmt"
	"regexp"
)

// XSCSubject contains all of the information encoded into an SCCo-standard
// subject line.
type XSCSubject struct {
	MessageNumber     string
	Severity          MessageSeverity
	SeverityCode      string
	HandlingOrder     HandlingOrder
	HandlingOrderCode string
	FormTag           string
	Subject           string
}

// subjectLineRE extracts the message number, severity code, handling order
// code, form name, and subject from the Subject line, assuming it's properly
// formatted.
var subjectLineRE = regexp.MustCompile(`(?i)^([A-Z0-9]+-?[0-9]+[A-Z]?)_(?:([A-Z])/)?([A-Z])_(?:([^_\s]+)_)?([^_\s]+(?:\s.*|$))`)

// ParseSubject parses the message Subject to see if it is an SCCo-standard
// subject line, and if so returns the data extracted from it.  (If not, it
// returns nil.)
func ParseSubject(subjline string) (s *XSCSubject) {
	if match := subjectLineRE.FindStringSubmatch(subjline); match != nil {
		s = new(XSCSubject)
		s.MessageNumber = match[1]
		s.SeverityCode = match[2]
		s.Severity, _ = ParseSeverity(match[2])
		s.HandlingOrderCode = match[3]
		s.HandlingOrder, _ = ParseHandlingOrder(match[3])
		s.FormTag = match[4]
		s.Subject = match[5]
		return s
	}
	return nil
}

// EncodeSubject encodes an SCCo-standard subject line.
func EncodeSubject(msgnum string, handling HandlingOrder, formtag, subject string) string {
	if formtag != "" {
		return fmt.Sprintf("%s_%s_%s_%s", msgnum, handling.Code(), formtag, subject)
	}
	return fmt.Sprintf("%s_%s_%s", msgnum, handling.Code(), subject)
}
