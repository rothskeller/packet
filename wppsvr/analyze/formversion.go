package analyze

import (
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
		if f := a.xsc.RawForm; f != nil {
			return xscmsg.OlderVersion(f.PIFOVersion, config.Get().MinimumVersions["PackItForms"])
		}
		return false
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
		if f := a.xsc.RawForm; f != nil {
			if min := config.Get().MinimumVersions[a.xsc.Type.Tag]; min != "" {
				return xscmsg.OlderVersion(f.FormVersion, min)
			}
		}
		return false
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
