// Package all imports (and, therefore, registers) all known XSC message types.
package all

import (
	_ "github.com/rothskeller/packet/xscmsg/ahfacstat" // .
	_ "github.com/rothskeller/packet/xscmsg/checkin"   // .
	_ "github.com/rothskeller/packet/xscmsg/checkout"  // .
	_ "github.com/rothskeller/packet/xscmsg/delivrcpt" // .
	_ "github.com/rothskeller/packet/xscmsg/eoc213rr"  // .
	_ "github.com/rothskeller/packet/xscmsg/ics213"    // .
	_ "github.com/rothskeller/packet/xscmsg/jurisstat" // .
	_ "github.com/rothskeller/packet/xscmsg/racesmar"  // .
	_ "github.com/rothskeller/packet/xscmsg/readrcpt"  // .
	_ "github.com/rothskeller/packet/xscmsg/sheltstat" // .
)
