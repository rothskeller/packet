package analyze

import (
	"fmt"
	"log"
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
		routing       = config.Get().FormRouting[a.msg.TypeCode()]
	)
	if routing == nil || (len(routing.ToICSPosition) == 0 && len(routing.ToLocation) == 0) {
		return
	}
	if len(routing.ToICSPosition) == 0 {
		foundPosition = true
	} else {
		log.Printf("%s %T\n", a.msg.TypeCode(), a.msg)
		actual := a.msg.SCCoForm().ToICSPosition
		for _, wanted := range routing.ToICSPosition {
			if actual == wanted {
				foundPosition = true
				break
			}
		}
	}
	if len(routing.ToLocation) == 0 {
		foundLocation = true
	} else {
		actual := a.msg.SCCoForm().ToLocation
		for _, wanted := range routing.ToLocation {
			if actual == wanted {
				foundLocation = true
				break
			}
		}
	}
	if !foundPosition && !foundLocation {
		var positions = make([]string, len(routing.ToICSPosition))
		for i, p := range routing.ToICSPosition {
			positions[i] = strconv.Quote(p)
		}
		var locations = make([]string, len(routing.ToLocation))
		for i, p := range routing.ToLocation {
			locations[i] = strconv.Quote(p)
		}
		a.problems = append(a.problems, &problem{
			code: ProblemFormDestination,
			response: fmt.Sprintf(`
This message form is addressed to ICS Position %q at Location %q.  %ss
should be addressed to %s at %s.
`, a.msg.SCCoForm().ToICSPosition, a.msg.SCCoForm().ToLocation, a.msg.TypeName(), english.Conjoin(positions, "or"), english.Conjoin(locations, "or")),
			references: refFormRouting,
		})
	}
	if !foundPosition && foundLocation {
		var positions = make([]string, len(routing.ToICSPosition))
		for i, p := range routing.ToICSPosition {
			positions[i] = strconv.Quote(p)
		}
		a.problems = append(a.problems, &problem{
			code: ProblemFormDestination,
			response: fmt.Sprintf(`
This message form is addressed to ICS Position %q.  %ss should be addressed to
ICS Position %s.
`, a.msg.SCCoForm().ToICSPosition, a.msg.TypeName(), english.Conjoin(positions, "or")),
			references: refFormRouting,
		})
	}
	if foundPosition && !foundLocation {
		var locations = make([]string, len(routing.ToLocation))
		for i, p := range routing.ToLocation {
			locations[i] = strconv.Quote(p)
		}
		a.problems = append(a.problems, &problem{
			code: ProblemFormDestination,
			response: fmt.Sprintf(`
This message form is addressed to Location %q.  %ss should be addressed to
Location %s.
`, a.msg.SCCoForm().ToLocation, a.msg.TypeName(), english.Conjoin(locations, "or")),
			references: refFormRouting,
		})
	}
}
