package plaintext

// Validate checks the contents of the message for compliance with rules
// enforced by standard Santa Clara County packet software (Outpost and
// PackItForms).  It returns a list of strings describing problems that those
// programs would flag or block.
func (m *PlainText) Validate() (problems []string) {
	if m.Subject == "" {
		problems = append(problems, `A message subject is required.`)
	}
	if m.Body == "" {
		problems = append(problems, `The message body cannot be empty.`)
	}
	return problems
}
