// Package all imports (and, therefore, registers) all known XSC message types.
package all

import (
	_ "steve.rothskeller.net/packet/xscmsg/ahfacstat" // .
	_ "steve.rothskeller.net/packet/xscmsg/checkin"   // .
	_ "steve.rothskeller.net/packet/xscmsg/checkout"  // .
	_ "steve.rothskeller.net/packet/xscmsg/delivrcpt" // .
	_ "steve.rothskeller.net/packet/xscmsg/eoc213rr"  // .
	_ "steve.rothskeller.net/packet/xscmsg/ics213"    // .
	_ "steve.rothskeller.net/packet/xscmsg/jurisstat" // .
	_ "steve.rothskeller.net/packet/xscmsg/racesmar"  // .
	_ "steve.rothskeller.net/packet/xscmsg/readrcpt"  // .
	_ "steve.rothskeller.net/packet/xscmsg/sheltstat" // .
)
