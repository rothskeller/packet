// Package allmsg contains a function to register all message types defined in
// sibling packages.
package allmsg

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/ahfacstat"
	"github.com/rothskeller/packet/xscmsg/checkin"
	"github.com/rothskeller/packet/xscmsg/checkout"
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/eoc213rr"
	"github.com/rothskeller/packet/xscmsg/ics213"
	"github.com/rothskeller/packet/xscmsg/jurisstat"
	"github.com/rothskeller/packet/xscmsg/plaintext"
	"github.com/rothskeller/packet/xscmsg/racesmar"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
	"github.com/rothskeller/packet/xscmsg/sheltstat"
	"github.com/rothskeller/packet/xscmsg/unkform"
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
