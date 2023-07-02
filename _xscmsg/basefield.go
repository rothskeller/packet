package xscmsg

import (
	"strings"
)

// BaseField creates a new base field operating on the supplied string.  It must
// be wrapped by another type that provides Label, Size, and Help
// implementations.
func BaseField(vp *string) Field { return &baseField{vp} }

type baseField struct{ vp *string }

func (f baseField) Value() string          { return *f.vp }
func (f *baseField) SetValue(value string) { *f.vp = strings.TrimSpace(value) }
func (f baseField) Problem() string        { return "" }
func (f baseField) Choices() []string      { return nil }
func (f baseField) Hint() string           { return "" }

// These baseField methods panic.  They must be overridden by a surrounding
// concrete type.

func (f baseField) Label() string              { panic("baseField.Label must be overridden") }
func (f *baseField) Size() (width, height int) { panic("baseField.Size must be overridden") }
func (f baseField) Help() string               { panic("baseField.Help must be overridden") }
