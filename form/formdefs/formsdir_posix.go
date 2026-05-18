//go:build !windows

package formdefs

import (
	"os"
	"path/filepath"

	"github.com/rothskeller/packet/form/pifover"
)

// FormsDir returns the pathname of the directory that should contain the local
// forms cache.  The directory does not have to exist.  It returns an empty
// string if the pathname cannot be determined.
func FormsDir() string {
	var home = os.Getenv("HOME")

	if home == "" {
		return ""
	}
	return filepath.Join(home, ".local", "share", "packet", pifover.FormsDirVersion)
}

// only applies to Windows
func UpdateOutpostConfiguration() (err error) { return nil }

func writeAddonINI(_ string) error { return nil }
