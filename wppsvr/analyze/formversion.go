package analyze

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbPIFOVersion.Code] = ProbPIFOVersion
	Problems[ProbFormVersion.Code] = ProbFormVersion
}

type form interface {
	Form() *pktmsg.Form
}

// ProbPIFOVersion is raised when the message contains a PackItForms form whose
// PackItForms version is too old.
var ProbPIFOVersion = &Problem{
	Code:  "PIFOVersion",
	Label: "PackItForms version out of date",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		var f *pktmsg.Form
		if xsc, ok := a.xsc.(form); ok {
			f = xsc.Form()
			return xscmsg.OlderVersion(f.PIFOVersion, config.Get().MinimumVersions["PackItForms"]), ""
		}
		return false, ""
	},
	Variables: variableMap{
		"ACTUALVER": func(a *Analysis) string {
			return a.xsc.(form).Form().PIFOVersion
		},
		"EXPECTVER": func(a *Analysis) string {
			return config.Get().MinimumVersions["PackItForms"]
		},
	},
}

// ProbFormVersion is raised when the message contains a form whose version is
// too old.
var ProbFormVersion = &Problem{
	Code:  "FormVersion",
	Label: "form version out of date",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		var f *pktmsg.Form
		if xsc, ok := a.xsc.(form); ok {
			f = xsc.Form()
			if min := config.Get().MinimumVersions[a.xsc.TypeTag()]; min != "" {
				return xscmsg.OlderVersion(f.FormVersion, min), ""
			}
		}
		return false, ""
	},
	Variables: variableMap{
		"ACTUALVER": func(a *Analysis) string {
			return a.xsc.(form).Form().FormVersion
		},
		"EXPECTVER": func(a *Analysis) string {
			return config.Get().MinimumVersions[a.xsc.TypeTag()]
		},
	},
}
