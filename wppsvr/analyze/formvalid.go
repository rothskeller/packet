package analyze

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
)

// Problem codes
const (
	ProblemFormCorrupt = "FormCorrupt"
	ProblemFormInvalid = "FormInvalid"
)

func init() {
	ProblemLabel[ProblemFormCorrupt] = "incorrectly encoded form"
	ProblemLabel[ProblemFormInvalid] = "invalid form contents"
}

// checkValidForm makes sure that the form embedded in the message (if any) was
// properly encoded, includes all required fields, and has valid values for all
// fields.
func (a *Analysis) checkValidForm() {
	// If the message contains a corrupt form encoding, pktmsg.IsForm will
	// return true, but our recognizer will have classified it as plain
	// text.  Detect that and report a corrupt form.
	if pktmsg.IsForm(a.msg.Body) {
		if _, ok := a.xsc.(*config.PlainTextMessage); ok {
			a.problems = append(a.problems, &problem{
				code: ProblemFormCorrupt,
				response: `
This message appears to contain an encoded form, but the encoding is
incorrect.  It appears to have been created or edited by software other than
the current PackItForms software.  Please use current PackItForms software to
encode messages containing forms.
`,
			})
			return
		}
	}
	// If the form doesn't validate correctly, report that.
	if xsc, ok := a.xsc.(interface{ Validate(bool) []string }); ok {
		if problems := xsc.Validate(true); len(problems) != 0 {
			response := "\nThis message contains a form with invalid contents:\n"
			for _, p := range problems {
				response += "    " + p + "\n"
			}
			response += "Please verify the correctness of the form before sending.\n"
			a.problems = append(a.problems, &problem{
				code:     ProblemFormInvalid,
				response: response,
			})
		}
	}
}
