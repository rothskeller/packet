//go:build windows

package osdep

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"golang.org/x/sys/windows"
)

const (
	// ServerStopFile is the file that, when touched, causes the running
	// server to stop.
	ServerStopFile = `C:\PackItForms\stop`

	allBytes = ^uint32(0)
)

var (
	// AddressFile is the pathname of the file containing the address of the
	// running server for the current user.
	AddressFile string
	// DefaultsFile is the pathname of the file containing the incident
	// configuration defaults.
	DefaultsFile string
	// LogsDir is the pathname of the directory containing the packet server
	// log files for the current user.
	LogsDir string
	// HomeDir is the user's home directory (or the root, if for some reason
	// the user's home directory isn't discernable).
	HomeDir string
)

func init() {
	if HomeDir, _ = os.UserHomeDir(); HomeDir == "" {
		HomeDir = "C:\\"
	}
	if ad := os.Getenv("APPDATA"); ad != "" {
		DefaultsFile = filepath.Join(ad, "packet.json")
	}
	AddressFile = `C:\PackItForms\server.url`
	LogsDir = `C:\PackItForms\Log`
}

// ReadLock locks the file for reading.
func ReadLock(fh *os.File) (err error) {
	var ol windows.Overlapped

	return windows.LockFileEx(windows.Handle(fh.Fd()), 0, 0, allBytes, allBytes, &ol)
}

// WriteLock locks the file for writing.
func WriteLock(fh *os.File) (err error) {
	var ol windows.Overlapped

	return windows.LockFileEx(windows.Handle(fh.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK, 0, allBytes, allBytes, &ol)
}

// Unlock unlocks the file.
func Unlock(fh *os.File) (err error) {
	var ol windows.Overlapped

	return windows.UnlockFileEx(windows.Handle(fh.Fd()), 0, allBytes, allBytes, &ol)
}

// DetachChild is the SysProcAttr to use in a *exec.Cmd when creating a process
// that should run independent of its parent (i.e., the server).
var DetachChild = &syscall.SysProcAttr{
	HideWindow:    true,
	CreationFlags: windows.CREATE_NEW_PROCESS_GROUP | windows.DETACHED_PROCESS,
}

// OpenURLCommand is the command to cause the system default browser to open a
// URL.
func OpenURLCommand(url string) *exec.Cmd {
	// This command is used instead of the more common "cmd.exe /c start"
	// because it does not flash a command terminal window.
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
}

// OpenFileCommand is the command to open a file using the system default
// application for its file type.
func OpenFileCommand(file string) *exec.Cmd {
	// This command is used instead of the more common "cmd.exe /c start"
	// because it does not flash a command terminal window.
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", file)
}

// IsAdmin returns whether the user is an Administrator.
func IsAdmin() bool {
	if fh, err := os.Open(`\\.\PHYSICALDRIVE0`); err == nil {
		fh.Close()
		return true
	}
	return false
}

// PrintPDFCommand is the command to send a PDF file to the system default printer.  It may return nil if no such command is available.
func PrintPDFCommand(file string) *exec.Cmd {
	serverPrintOnce.Do(setServerPrintCmd)
	if serverPrintCmd != "" {
		return exec.Command(serverPrintCmd, file)
	}
	return nil
}

var serverPrintCmd string
var serverPrintOnce sync.Once

func setServerPrintCmd() {
	const path = `C:\Program Files (x86)\SCCo Packet\PdfToPrinter.exe`

	if _, err := os.Stat(path); err == nil {
		serverPrintCmd = path
	}
}
