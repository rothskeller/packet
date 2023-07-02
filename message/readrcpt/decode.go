package readrcpt

import (
	"regexp"
	"strings"
)

// readReceiptRE matches the first lines of a read receipt message.  Its
// substrings are the read time and the To address.
var readReceiptRE = regexp.MustCompile(`^!RR!(.+)\n.*\n\nTo: (.+)`)

func decode(subject, body string) *ReadReceipt {
	if !strings.HasPrefix(subject, "READ: ") {
		return nil
	}
	if match := readReceiptRE.FindStringSubmatch(body); match != nil {
		return &ReadReceipt{
			MessageSubject: subject[6:],
			MessageTo:      match[2],
			ReadTime:       match[1],
		}
	}
	return nil
}
