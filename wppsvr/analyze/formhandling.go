package analyze

import (
	"fmt"

	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/wppsvr/config"
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
		ho      pktmsg.HandlingOrder
		routing = config.Get().FormRouting[a.msg.TypeCode()]
	)
	if routing == nil {
		return
	}
	switch handling := routing.HandlingOrder; handling {
	case "":
		return
	case "computed":
		ho = a.msg.(interface{ RecommendedHandlingOrder() pktmsg.HandlingOrder }).RecommendedHandlingOrder()
	default:
		ho, _ = pktmsg.ParseHandlingOrder(handling)
	}
	if f := a.msg.Form(); f.HandlingOrder != ho {
		a.problems = append(a.problems, &problem{
			code: ProblemFormHandlingOrder,
			response: fmt.Sprintf(`
This message has handling order %s.  According to the SCCo ARES/RACES
Recommended Form Routing document, it should have handling order %s.
`, f.HandlingOrder, ho.String()),
			references: refFormRouting,
		})
	}
}
