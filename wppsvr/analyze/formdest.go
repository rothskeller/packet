package analyze

import (
	"fmt"
	"strconv"

	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/wppsvr/english"
)

// Problem codes
const (
	ProblemFormDestination = "FormDestination"
)

func init() {
	ProblemLabel[ProblemFormDestination] = "incorrect destination for form"
}

// checkFormDestination determines whether the message has the correct
// destination.
func (a *Analysis) checkFormDestination() {
	var (
		foundPosition bool
		foundLocation bool
		havePos       string
		haveLoc       string
		want          = config.Get().FormRouting[a.xsc.TypeTag()]
	)
	if f, ok := a.xsc.(interface{ Routing() (string, string) }); ok {
		havePos, haveLoc = f.Routing()
	} else {
		return
	}
	if len(want.ToICSPosition) == 0 {
		foundPosition = true
	} else {
		for _, wantPos := range want.ToICSPosition {
			if havePos == wantPos {
				foundPosition = true
				break
			}
		}
	}
	if len(want.ToLocation) == 0 {
		foundLocation = true
	} else {
		for _, wantLoc := range want.ToLocation {
			if haveLoc == wantLoc {
				foundLocation = true
				break
			}
		}
	}
	if !foundPosition && !foundLocation {
		var positions = make([]string, len(want.ToICSPosition))
		for i, p := range want.ToICSPosition {
			positions[i] = strconv.Quote(p)
		}
		var locations = make([]string, len(want.ToLocation))
		for i, p := range want.ToLocation {
			locations[i] = strconv.Quote(p)
		}
		a.problems = append(a.problems, &problem{
			code: ProblemFormDestination,
			response: fmt.Sprintf(`
This message form is addressed to ICS Position %q at Location %q.  %ss
should be addressed to %s at %s.
`, havePos, haveLoc, a.xsc.TypeName(), english.Conjoin(positions, "or"), english.Conjoin(locations, "or")),
			references: refFormRouting,
		})
	}
	if !foundPosition && foundLocation {
		var positions = make([]string, len(want.ToICSPosition))
		for i, p := range want.ToICSPosition {
			positions[i] = strconv.Quote(p)
		}
		a.problems = append(a.problems, &problem{
			code: ProblemFormDestination,
			response: fmt.Sprintf(`
This message form is addressed to ICS Position %q.  %ss should be addressed to
ICS Position %s.
`, havePos, a.xsc.TypeName(), english.Conjoin(positions, "or")),
			references: refFormRouting,
		})
	}
	if foundPosition && !foundLocation {
		var locations = make([]string, len(want.ToLocation))
		for i, p := range want.ToLocation {
			locations[i] = strconv.Quote(p)
		}
		a.problems = append(a.problems, &problem{
			code: ProblemFormDestination,
			response: fmt.Sprintf(`
This message form is addressed to Location %q.  %ss should be addressed to
Location %s.
`, haveLoc, a.xsc.TypeName(), english.Conjoin(locations, "or")),
			references: refFormRouting,
		})
	}
}
