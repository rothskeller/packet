package analyze

import (
	"strconv"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbFormDestination.Code] = ProbFormDestination
	Problems[ProbFormToICSPosition.Code] = ProbFormToICSPosition
	Problems[ProbFormToLocation.Code] = ProbFormToLocation
}

// ProbFormDestination is raised when neither the form's To ICS Position field
// nor its To Location field match the recommended routing.
var ProbFormDestination = &Problem{
	Code: "FormDestination",
	detect: func(a *Analysis) bool {
		// This check does not apply unless the message is of a type
		// that has recommended To ICS Position and To Location values.
		want := config.Get().FormRouting[a.xsc.Type.Tag]
		if want == nil || len(want.ToICSPosition) == 0 && len(want.ToLocation) == 0 {
			return false
		}
		// The check.
		return !inList(want.ToICSPosition, a.xsc.KeyField(xscmsg.FToICSPosition).Value) &&
			!inList(want.ToLocation, a.xsc.KeyField(xscmsg.FToLocation).Value)
	},
	Variables: variableMap{
		"ACTUALLOC":   varActualLoc,
		"ACTUALPOSN":  varActualPosn,
		"EXPECTLOCS":  varExpectLocs,
		"EXPECTPOSNS": varExpectPosns,
	},
}

// ProbFormToICSPosition is raised when a form's To ICS Position field doesn't
// match the recommended routing (but its To Location field does).
var ProbFormToICSPosition = &Problem{
	Code: "FormToICSPosition",
	ifnot: []*Problem{
		// This check does not apply if we already found that both the
		// To ICS Position and To Location fields are off.
		ProbFormDestination,
	},
	detect: func(a *Analysis) bool {
		// This check does not apply unless the message is of a type
		// that has recommended To ICS Position values.
		want := config.Get().FormRouting[a.xsc.Type.Tag]
		if want == nil || len(want.ToICSPosition) == 0 {
			return false
		}
		// The check.
		return !inList(want.ToICSPosition, a.xsc.KeyField(xscmsg.FToICSPosition).Value)
	},
	Variables: variableMap{
		"ACTUALPOSN":  varActualPosn,
		"EXPECTPOSNS": varExpectPosns,
	},
}

// ProbFormToLocation is raised when a form's To Location field doesn't match
// the recommended routing (but its To ICS Position field does).
var ProbFormToLocation = &Problem{
	Code: "FormToLocation",
	ifnot: []*Problem{
		// This check does not apply if we already found that both the
		// To ICS Position and To Location fields are off.
		ProbFormDestination,
	},
	detect: func(a *Analysis) bool {
		// This check does not apply unless the message is of a type
		// that has recommended To Location values.
		want := config.Get().FormRouting[a.xsc.Type.Tag]
		if want == nil || len(want.ToLocation) == 0 {
			return false
		}
		// The check.
		return !inList(want.ToLocation, a.xsc.KeyField(xscmsg.FToLocation).Value)
	},
	Variables: variableMap{
		"ACTUALLOC":  varActualLoc,
		"EXPECTLOCS": varExpectLocs,
	},
}

func varActualLoc(a *Analysis) string {
	return a.xsc.KeyField(xscmsg.FToLocation).Value
}
func varActualPosn(a *Analysis) string {
	return a.xsc.KeyField(xscmsg.FToICSPosition).Value
}
func varExpectLocs(a *Analysis) string {
	var locations []string
	for _, loc := range config.Get().FormRouting[a.xsc.Type.Tag].ToLocation {
		locations = append(locations, strconv.Quote(loc))
	}
	return english.Conjoin(locations, "or")
}
func varExpectPosns(a *Analysis) string {
	var positions []string
	for _, pos := range config.Get().FormRouting[a.xsc.Type.Tag].ToICSPosition {
		positions = append(positions, strconv.Quote(pos))
	}
	return english.Conjoin(positions, "or")
}
