package plaintext

import (
	"github.com/rothskeller/packet/message"
)

// Compare compares two messages.  It returns a score indicating how closely
// they match, and the detailed comparisons of each field in the message.  The
// comparison is not symmetric:  the receiver of the call is the "expected"
// message and the argument is the "actual" message.
func (exp *PlainText) Compare(actual message.Message) (score, outOf int, fields []*message.CompareField) {
	panic("not implemented")
}
