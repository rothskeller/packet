package analyze

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
	// This check only applies to messages with encoded forms.
	var form = a.msg.Form()
	if form == nil {
		return
	}
	if form.CorruptForm {
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
	if vform, ok := a.msg.(interface{ Valid() bool }); ok && !vform.Valid() {
		a.problems = append(a.problems, &problem{
			code: ProblemFormInvalid,
			response: `
This message contains a form with invalid contents.  One or more fields of the
form have invalid values, or required form fields are not filled in.  Please
verify the correctness of the form before sending.
`,
		})
	}
}
