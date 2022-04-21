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
		form          = a.msg.SCCoForm()
	)
	if form == nil || routing == nil || (len(routing.ToICSPosition) == 0 && len(routing.ToLocation) == 0) {
		return
	}
	if len(routing.ToICSPosition) == 0 {
		foundPosition = true
	} else {
		log.Printf("%s %T\n", a.msg.TypeCode(), a.msg)
		for _, wanted := range routing.ToICSPosition {
			if form.ToICSPosition == wanted {
				foundPosition = true
				break
			}
		}
	}
	if len(routing.ToLocation) == 0 {
		foundLocation = true
	} else {
		for _, wanted := range routing.ToLocation {
			if form.ToLocation == wanted {
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
`, form.ToICSPosition, form.ToLocation, a.msg.TypeName(), english.Conjoin(positions, "or"), english.Conjoin(locations, "or")),
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
`, form.ToICSPosition, a.msg.TypeName(), english.Conjoin(positions, "or")),
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
`, form.ToLocation, a.msg.TypeName(), english.Conjoin(locations, "or")),
			references: refFormRouting,
		})
	}
}
