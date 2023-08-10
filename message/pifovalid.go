package message

// PIFOValid checks the contents of the message for compliance with rules
// enforced by standard Santa Clara County packet software (Outpost and
// PackItForms).  It returns a list of strings describing problems that
// those programs would flag or block.
func (bm *BaseMessage) PIFOValid() (problems []string) {
	for _, f := range bm.Fields {
		if p := f.PIFOValid(f); p != "" {
			problems = append(problems, p)
		}
	}
	return problems
}
