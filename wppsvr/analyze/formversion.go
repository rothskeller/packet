package analyze

import (
	"fmt"

	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/xscmsg"
)

// Problem codes
const (
	ProblemPIFOVersion = "PIFOVersion"
	ProblemFormVersion = "FormVersion"
)

func init() {
	ProblemLabel[ProblemPIFOVersion] = "PackItForms version out of date"
	ProblemLabel[ProblemFormVersion] = "form version out of date"
}

// checkFormVersion makes sure that the form embedded in the message (if any)
// used a current version of the form template.
func (a *Analysis) checkFormVersion() {
	var form *pktmsg.Form

	if xsc, ok := a.xsc.(interface{ Form() *pktmsg.Form }); ok {
		form = xsc.Form()
	} else {
		return
	}
	// Check the version of the PackItForms encoding.
	minimums := config.Get().MinimumVersions
	if xscmsg.OlderVersion(form.PIFOVersion, minimums["PackItForms"]) {
		a.problems = append(a.problems, &problem{
			code: ProblemPIFOVersion,
			response: fmt.Sprintf(`
This message used version %s of PackItForms to encode the form, but that
version is not current.  Please use PackItForms version %s or newer to encode
messages containing forms.
`, form.PIFOVersion, minimums[config.PackItForms]),
		})
	}
	// Check the version of the specific form type.
	if min := minimums[a.xsc.TypeTag()]; min != "" {
		if xscmsg.OlderVersion(form.FormVersion, min) {
			a.problems = append(a.problems, &problem{
				code: ProblemFormVersion,
				response: fmt.Sprintf(`
This message contains version %s of the %s, but that version is not current.
Please use version %s or newer of the form.  (You can get the newer form by
updating your PackItForms installation.)
`, form.FormVersion, a.xsc.TypeName(), min),
			})
		}
	}
}
