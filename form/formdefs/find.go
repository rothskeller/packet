package formdefs

import (
	"github.com/rothskeller/packet/form"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/message"
)

// Find finds the defined message type whose form definition matches the
// supplied predicate.  It returns nil if no such message type was found.
func Find(pred func(fd *formdef.FormDef) bool) message.MType {
	for m := range message.AllTypes() {
		if ft, ok := m.(form.FormType); ok {
			if pred(ft.FormDef) {
				return m
			}
		}
	}
	return nil
}
