package pktmsg

// IsForm returns whether the Message is a PackItForms form.
func IsForm(m Message) bool {
	return m.PIFOVersion() != nil
}

// Finalized returns whether the Message has been finalized, i.e., has been
// transmitted or received over the air.  Such messages may still be modified in
// small ways, e.g. to add a destination message number when a delivery receipt
// is received, but from a human perspective they are no longer changeable.
func Finalized(m Message) bool {
	return m.SentDate().Value() != ""
}
