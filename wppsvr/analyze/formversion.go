package analyze

import (
	"fmt"
	"strconv"
	"strings"

	"steve.rothskeller.net/packet/wppsvr/config"
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
	// This check only applies to messages with encoded forms.  It does not
	// apply to corrupted forms.
	var form = a.msg.Form()
	if form == nil || form.CorruptForm {
		return
	}
	minimums := config.Get().MinimumVersions
	if older, valid := older(form.PIFOVersion, minimums["PackItForms"]); older || !valid {
		a.problems = append(a.problems, &problem{
			code:    ProblemPIFOVersion,
			subject: "PackItForms version out of date",
			response: fmt.Sprintf(`
This message used version %s of PackItForms to encode the form, but that
version is not current.  Please use PackItForms version %s or newer to encode
messages containing forms.
`, form.PIFOVersion, minimums[config.PackItForms]),
		})
	}
	if minimums[form.FormName] != "" {
		if older, valid := older(form.FormVersion, minimums[form.FormName]); older || !valid {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormVersion,
				subject: "form version out of date",
				response: fmt.Sprintf(`
This message contains version %s of the %s, but that version is not current.
Please use version %s or newer of the form.  (You can get the newer form by
updating your PackItForms installation.)
`, form.FormVersion, form.TypeName(), minimums[form.FormName]),
			})
		}
	}
}

// older compares two version numbers.  Version numbers are a dot-separated
// sequence of parts, each of which is compared independently.  The parts are
// compared numerically if they parse as integers, and as strings otherwise.
// The valid flag is set if any component of the first version number is a valid
// integer.
func older(a, b string) (older, valid bool) {
	aparts := strings.Split(a, ".")
	bparts := strings.Split(b, ".")
	for len(aparts) != 0 && len(bparts) != 0 {
		anum, aerr := strconv.Atoi(aparts[0])
		if aerr == nil {
			valid = true
		}
		bnum, berr := strconv.Atoi(bparts[0])
		if aerr == nil && berr == nil {
			if anum < bnum {
				return true, valid
			}
			if anum > bnum {
				return false, valid
			}
		} else {
			if aparts[0] < bparts[0] {
				return true, valid
			}
			if aparts[0] > bparts[0] {
				return false, valid
			}
		}
		aparts = aparts[1:]
		bparts = bparts[1:]
	}
	if len(bparts) != 0 {
		return true, valid
	}
	return false, valid
}
