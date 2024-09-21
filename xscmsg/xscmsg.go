// Package xscmsg contains a function to register all public Santa Clara County
// message types (i.e., those defined in subpackages of xscmsg).
package xscmsg

import (
	_ "github.com/rothskeller/packet/xscmsg/ahfacstat"
	_ "github.com/rothskeller/packet/xscmsg/checkin"
	_ "github.com/rothskeller/packet/xscmsg/checkout"
	_ "github.com/rothskeller/packet/xscmsg/delivrcpt"
	_ "github.com/rothskeller/packet/xscmsg/eoc213rr"
	_ "github.com/rothskeller/packet/xscmsg/ics213"
	_ "github.com/rothskeller/packet/xscmsg/jurisstat"
	_ "github.com/rothskeller/packet/xscmsg/plaintext"
	_ "github.com/rothskeller/packet/xscmsg/racesmar"
	_ "github.com/rothskeller/packet/xscmsg/readrcpt"
	_ "github.com/rothskeller/packet/xscmsg/sheltstat"
	_ "github.com/rothskeller/packet/xscmsg/unkform"
)

// Register registers all message types defined in sibling packages.
func Register() {
	// This is a no-op.  Just importing this package is sufficient.
}
