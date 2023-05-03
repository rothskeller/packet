// Package all imports all defined XSC message types, and therefore registers
// them with xscmsg.
package all

import (
	_ "github.com/rothskeller/packet/xscmsg/checkin"   // .
	_ "github.com/rothskeller/packet/xscmsg/checkout"  // .
	_ "github.com/rothskeller/packet/xscmsg/delivrcpt" // .
	_ "github.com/rothskeller/packet/xscmsg/readrcpt"  // .
)
