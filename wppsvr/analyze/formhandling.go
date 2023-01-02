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
		var (
			want    xscmsg.HandlingOrder
			have    xscmsg.HandlingOrder
			routing *config.FormRouting
		)
		// What handling order do we have?
		if f := a.xsc.KeyField(xscmsg.FHandling); f != nil {
			have, _ = xscmsg.ParseHandlingOrder(f.Value)
		} else {
			return false
		}
		// What handling order are we supposed to have?
		if routing = config.Get().FormRouting[a.xsc.Type.Tag]; routing == nil {
			return false
		}
		if routing.HandlingOrder == "computed" {
			want = config.ComputedRecommendedHandlingOrder[a.xsc.Type.Tag](a.xsc)
		} else {
			want, _ = xscmsg.ParseHandlingOrder(routing.HandlingOrder)
		}
		if want == 0 {
			return false
		}
		// Return whether they match.
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
