package all

import (
	"github.com/rothskeller/packet/typedmsg"
	"github.com/rothskeller/packet/xscmsg/checkin"
	"github.com/rothskeller/packet/xscmsg/checkout"
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/plaintext"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
)

// RegisterAll registers all of the message types defined in subpackages of
// xscmsg.
func RegisterAll() {
	typedmsg.Register(&delivrcpt.Type)
	typedmsg.Register(&readrcpt.Type)
	typedmsg.Register(&checkin.Type)
	typedmsg.Register(&checkout.Type)
	typedmsg.Register(&plaintext.Type)
}
