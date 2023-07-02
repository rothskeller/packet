// Package allmsg contains a function to register all message types defined in
// sibling packages.
package allmsg

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/ahfacstat"
	"github.com/rothskeller/packet/message/checkin"
	"github.com/rothskeller/packet/message/checkout"
	"github.com/rothskeller/packet/message/delivrcpt"
	"github.com/rothskeller/packet/message/eoc213rr"
	"github.com/rothskeller/packet/message/ics213"
	"github.com/rothskeller/packet/message/jurisstat"
	"github.com/rothskeller/packet/message/plaintext"
	"github.com/rothskeller/packet/message/racesmar"
	"github.com/rothskeller/packet/message/readrcpt"
	"github.com/rothskeller/packet/message/sheltstat"
	"github.com/rothskeller/packet/message/unkform"
)

// Register registers all message types defined in sibling packages.
func Register() {
	message.Register(&ahfacstat.Type)
	message.Register(&eoc213rr.Type)
	message.Register(&ics213.Type)
	message.Register(&jurisstat.Type)
	message.Register(&jurisstat.OldType)
	message.Register(&racesmar.Type)
	message.Register(&sheltstat.Type)
	message.Register(&unkform.Type)
	message.Register(&checkin.Type)
	message.Register(&checkout.Type)
	message.Register(&delivrcpt.Type)
	message.Register(&readrcpt.Type)
	message.Register(&plaintext.Type)
}
