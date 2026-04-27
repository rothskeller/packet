//go:build !windows

package tnc

import (
	"os"
	"path/filepath"
)

// TNCDir returns the pathname of the directory that contains TNC definitions.
func TNCDir() string {
	var home = os.Getenv("HOME")

	if home == "" {
		return ""
	}
	return filepath.Join(home, ".local", "share", "packet", "TNCs")
}
