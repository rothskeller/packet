package analyze

import (
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbFormHandlingOrder.Code] = ProbFormHandlingOrder
}

// ProbFormHandlingOrder is raised when the handling order in a form does not
// conform to the recommended form routing.
var ProbFormHandlingOrder = &Problem{
	Code: "FormHandlingOrder",
	detect: func(a *Analysis) bool {
		// This check does not apply unless the message is of a type
		// that has a recommended Handling value.
		var want xscmsg.HandlingOrder
		if routing := config.Get().FormRouting[a.xsc.Type.Tag]; routing != nil {
			if routing.HandlingOrder == "computed" {
				want = config.ComputedRecommendedHandlingOrder[a.xsc.Type.Tag](a.xsc)
			} else {
				want, _ = xscmsg.ParseHandlingOrder(routing.HandlingOrder)
			}
		}
		if want == 0 {
			return false
		}
		// The check.
		have, _ := xscmsg.ParseHandlingOrder(a.xsc.KeyField(xscmsg.FHandling).Value)
		return have != want
	},
	Variables: variableMap{
		"ACTUALHO": func(a *Analysis) string {
			return a.xsc.KeyField(xscmsg.FHandling).Value
		},
		"EXPECTHO": func(a *Analysis) string {
			expectstr := config.Get().FormRouting[a.xsc.Type.Tag].HandlingOrder
			if expectstr == "computed" {
				return config.ComputedRecommendedHandlingOrder[a.xsc.Type.Tag](a.xsc).String()
			}
			return expectstr
		},
	},
}
