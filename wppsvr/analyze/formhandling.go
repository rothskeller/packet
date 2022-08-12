package analyze

import (
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbFormHandlingOrder.Code] = ProbFormHandlingOrder
}

type formWithHandlingOrder interface {
	HandlingOrder() (string, xscmsg.HandlingOrder)
}

// ProbFormHandlingOrder is raised when the handling order in a form does not
// conform to the recommended form routing.
var ProbFormHandlingOrder = &Problem{
	Code:  "FormHandlingOrder",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		var (
			want    xscmsg.HandlingOrder
			have    xscmsg.HandlingOrder
			routing *config.FormRouting
		)
		// What handling order do we have?
		if f, ok := a.xsc.(formWithHandlingOrder); ok {
			_, have = f.HandlingOrder()
		} else {
			return false, ""
		}
		// What handling order are we supposed to have?
		if routing = config.Get().FormRouting[a.xsc.TypeTag()]; routing == nil {
			return false, ""
		}
		if routing.HandlingOrder == "computed" {
			want = config.ComputedRecommendedHandlingOrder[a.xsc.TypeTag()](a.xsc)
		} else {
			want, _ = xscmsg.ParseHandlingOrder(routing.HandlingOrder)
		}
		if want == 0 {
			return false, ""
		}
		// Return whether they match.
		return have != want, ""
	},
	Variables: variableMap{
		"ACTUALHO": func(a *Analysis) string {
			actual, _ := a.xsc.(formWithHandlingOrder).HandlingOrder()
			return actual
		},
		"EXPECTHO": func(a *Analysis) string {
			expectstr := config.Get().FormRouting[a.xsc.TypeTag()].HandlingOrder
			if expectstr == "computed" {
				return config.ComputedRecommendedHandlingOrder[a.xsc.TypeTag()](a.xsc).String()
			}
			return expectstr
		},
	},
}
