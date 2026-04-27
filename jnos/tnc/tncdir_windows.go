//go:build windows

package tnc

// TNCDir returns the pathname of the directory that contains TNC definitions.
func TNCDir() string {
	return `C:\PackItForms\TNCs`
}
