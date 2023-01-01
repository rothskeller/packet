package analyze

import (
	"strconv"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbFormDestination.Code] = ProbFormDestination
}

// ProbFormDestination is raised when a form's destination doesn't match the
// recommended routing.
var ProbFormDestination = &Problem{
	Code:  "FormDestination",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		want := config.Get().FormRouting[a.xsc.Type.Tag]
		if want == nil || len(want.ToICSPosition) == 0 && len(want.ToLocation) == 0 {
			return false, ""
		}
		actpos := a.xsc.KeyField(xscmsg.FToICSPosition).Value
		actloc := a.xsc.KeyField(xscmsg.FToLocation).Value
		if len(want.ToICSPosition) != 0 && !inList(want.ToICSPosition, actpos) {
			if len(want.ToLocation) != 0 && !inList(want.ToLocation, actloc) {
				return true, "both"
			}
			return true, "pos"
		}
		if len(want.ToLocation) != 0 && !inList(want.ToLocation, actloc) {
			return true, "loc"
		}
		return false, ""
	},
	Variables: variableMap{
		"ACTUALLOC": func(a *Analysis) string {
			return a.xsc.KeyField(xscmsg.FToLocation).Value
		},
		"ACTUALPOSN": func(a *Analysis) string {
			return a.xsc.KeyField(xscmsg.FToICSPosition).Value
		},
		"EXPECTLOCS": func(a *Analysis) string {
			var locations []string
			for _, loc := range config.Get().FormRouting[a.xsc.Type.Tag].ToLocation {
				locations = append(locations, strconv.Quote(loc))
			}
			return english.Conjoin(locations, "or")
		},
		"EXPECTPOSNS": func(a *Analysis) string {
			var positions []string
			for _, pos := range config.Get().FormRouting[a.xsc.Type.Tag].ToICSPosition {
				positions = append(positions, strconv.Quote(pos))
			}
			return english.Conjoin(positions, "or")
		},
	},
}
