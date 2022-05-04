package analyze

import (
	"fmt"

	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/xscmsg"
)

// Problem codes
const (
	ProblemFormHandlingOrder = "FormHandlingOrder"
)

func init() {
	ProblemLabel[ProblemFormHandlingOrder] = "incorrect handling order for form"
}

// checkFormHandlingOrder determines whether the message has the correct
// handling order based on the form contents.
func (a *Analysis) checkFormHandlingOrder() {
	var (
		want    xscmsg.HandlingOrder
		routing = config.Get().FormRouting[a.xsc.TypeTag()]
	)
	if routing == nil {
		return
	}
	switch handling := routing.HandlingOrder; handling {
	case "":
		return
	case "computed":
		want = config.ComputedRecommendedHandlingOrder[a.xsc.TypeTag()](a.xsc)
		if want == 0 {
			return
		}
	default:
		want, _ = xscmsg.ParseHandlingOrder(handling)
	}
	if f, ok := a.xsc.(interface {
		HandlingOrder() (string, xscmsg.HandlingOrder)
	}); ok {
		if _, have := f.HandlingOrder(); have != want {
			a.problems = append(a.problems, &problem{
				code: ProblemFormHandlingOrder,
				response: fmt.Sprintf(`
This message has handling order %s.  According to the SCCo ARES/RACES
Recommended Form Routing document, it should have handling order %s.
`, f.HandlingOrder, want.String()),
				references: refFormRouting,
			})
		}
	}
}
