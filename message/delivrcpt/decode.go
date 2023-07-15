package delivrcpt

import (
	"regexp"
	"strings"
)

// deliveryReceiptRE matches the first lines of a delivery receipt message.  Its
// substrings are the local message ID, the delivery time, and the To address.
var deliveryReceiptRE = regexp.MustCompile(`^!LMI!([^!]+)!DR!(.+)\n.*\nTo: (.+)\nSubject:.*\nwas delivered on.*\nRecipient's Local.*\n`)

func decode(subject, body string) *DeliveryReceipt {
	if !strings.HasPrefix(subject, "DELIVERED: ") {
		return nil
	}
	if match := deliveryReceiptRE.FindStringSubmatch(body); match != nil {
		return &DeliveryReceipt{
			MessageSubject: subject[11:],
			MessageTo:      match[3],
			LocalMessageID: match[1],
			DeliveredTime:  match[2],
			ExtraText:      strings.TrimSpace(body[len(match[0]):]),
		}
	}
	return nil
}
