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

// ProbPIFOVersion is raised when the message contains a PackItForms form whose
// PackItForms version is too old.
var ProbPIFOVersion = &Problem{
	Code: "PIFOVersion",
	detect: func(a *Analysis) bool {
		// This check only applies to forms.
		var form *pktmsg.Form
		if form = a.xsc.RawForm; form == nil {
			return false
		}
		// The check.
		return xscmsg.OlderVersion(form.PIFOVersion, config.Get().MinimumVersions["PackItForms"])
	},
	Variables: variableMap{
		"ACTUALVER": func(a *Analysis) string {
			return a.xsc.RawForm.PIFOVersion
		},
		"EXPECTVER": func(a *Analysis) string {
			return config.Get().MinimumVersions["PackItForms"]
		},
	},
}

// ProbFormVersion is raised when the message contains a form whose version is
// too old.
var ProbFormVersion = &Problem{
	Code: "FormVersion",
	detect: func(a *Analysis) bool {
		// This check only applies to forms.
		var form *pktmsg.Form
		if form = a.xsc.RawForm; form == nil {
			return false
		}
		// This check only applies to forms for which we have a minimum
		// version.  The config enforces that we have a minimum version
		// for all known form types, but we won't have one for an
		// unknown form type.
		var min string
		if min = config.Get().MinimumVersions[a.xsc.Type.Tag]; min == "" {
			return false
		}
		// The check.
		return xscmsg.OlderVersion(form.FormVersion, min)
	},
	Variables: variableMap{
		"ACTUALVER": func(a *Analysis) string {
			return a.xsc.RawForm.FormVersion
		},
		"EXPECTVER": func(a *Analysis) string {
			return config.Get().MinimumVersions[a.xsc.Type.Tag]
		},
	},
}
