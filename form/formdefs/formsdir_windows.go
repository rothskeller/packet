//go:build windows

package formdefs

import (
	"path/filepath"

	"github.com/rothskeller/packet/v4/form/pifover"
)

// FormsDir returns the pathname of the directory that should contain the local
// forms cache.
func FormsDir() string {
	return filepath.Join(`C:\PackItForms\Forms`, pifover.FormsDirVersion)
}
