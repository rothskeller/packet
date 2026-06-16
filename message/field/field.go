package field

import (
	"github.com/rothskeller/packet/v4/message/msgifc"
)

// Field is the interface satisfied by all message fields.  Note that this is
// more than just PackItForms fields; there are fields relating to metadata as
// well, and non-forms messages have fields also.
type Field = msgifc.Field

type ChoicePair = msgifc.ChoicePair

type ComparedField = msgifc.ComparedField
