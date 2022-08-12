package analyze

import (
	"strconv"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
)

func init() {
	Problems[ProbFormDestination.Code] = ProbFormDestination
}

type formWithRouting interface{ Routing() (string, string) }

// ProbFormDestination is raised when a form's destination doesn't match the
// recommended routing.
var ProbFormDestination = &Problem{
	Code:  "FormDestination",
	Label: "incorrect destination for form",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		want := config.Get().FormRouting[a.xsc.TypeTag()]
		if want == nil || len(want.ToICSPosition) == 0 && len(want.ToLocation) == 0 {
			return false, ""
		}
		actpos, actloc := a.xsc.(formWithRouting).Routing()
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
			_, actloc := a.xsc.(formWithRouting).Routing()
			return actloc
		},
		"ACTUALPOSN": func(a *Analysis) string {
			actpos, _ := a.xsc.(formWithRouting).Routing()
			return actpos
		},
		"EXPECTLOCS": func(a *Analysis) string {
			var locations []string
			for _, loc := range config.Get().FormRouting[a.xsc.TypeTag()].ToLocation {
				locations = append(locations, strconv.Quote(loc))
			}
			return english.Conjoin(locations, "or")
		},
		"EXPECTPOSNS": func(a *Analysis) string {
			var positions []string
			for _, pos := range config.Get().FormRouting[a.xsc.TypeTag()].ToICSPosition {
				positions = append(positions, strconv.Quote(pos))
			}
			return english.Conjoin(positions, "or")
		},
	},
	references: refFormRouting,
}
