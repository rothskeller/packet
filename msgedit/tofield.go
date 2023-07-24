package msgedit

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/message"
)

var jnosMailboxRE = regexp.MustCompile(`(?i)^[A-Z][A-Z0-9]{0,5}$`)

func (me *messageEditor) applyToField() {
	if me.tofield == nil {
		me.tofield = &message.EditField{
			Label: "To",
			Value: strings.Join(me.env.To, ", "),
			Width: 80,
			Help:  "This is the list of addresses to which the message is sent.  Each address must be a JNOS mailbox name, a BBS network address, or an email address.  The addresses must be separated by commas.  At least one address is required.",
		}
	}
	addresses := strings.Split(me.tofield.Value, ",")
	j := 0
	me.tofield.Problem = ""
	for _, address := range addresses {
		if trim := strings.TrimSpace(address); trim != "" {
			addresses[j], j = trim, j+1
			if jnosMailboxRE.MatchString(address) {
				// do nothing
			} else if _, err := mail.ParseAddress(address); err == nil {
				// do nothing
			} else {
				me.tofield.Problem = fmt.Sprintf(`The "To" field contains %q, which is not a valid JNOS mailbox name, BBS network address, or email address.`, address)
			}
		}
	}
	me.tofield.Value = strings.Join(addresses[:j], ", ")
	me.env.To = addresses[:j]
	if me.tofield.Value == "" {
		me.tofield.Problem = `The "To" field is required.`
	}
}
