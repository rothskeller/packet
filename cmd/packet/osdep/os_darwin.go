//go:build darwin

package osdep

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

const (
	// ServerStopFile is the file that, when touched, causes the running
	// server to stop.
	ServerStopFile = "/tmp/packet-stop"
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
	if HomeDir = os.Getenv("HOME"); HomeDir == "" {
		fmt.Fprintln(os.Stderr, "ERROR: can't locate data files: $HOME not set")
		os.Exit(1)
	}
	AddressFile = HomeDir + "/.local/state/packet/server.url"
	LogsDir = HomeDir + "/.local/state/packet/log"
	DefaultsFile = HomeDir + "/.config/packet/packet.json"
}

// ReadLock locks the file for reading.
func ReadLock(fh *os.File) (err error) {
	return syscall.Flock(int(fh.Fd()), syscall.LOCK_SH)
}

// WriteLock locks the file for writing.
func WriteLock(fh *os.File) (err error) {
	return syscall.Flock(int(fh.Fd()), syscall.LOCK_EX)
}

// Unlock unlocks the file.
func Unlock(fh *os.File) (err error) {
	return syscall.Flock(int(fh.Fd()), syscall.LOCK_UN)
}

// DetachChild is the SysProcAttr to use in a *exec.Cmd when creating a process
// that should run independent of its parent (i.e., the server).
var DetachChild = &syscall.SysProcAttr{
	// Setsid:  true,
	// Setpgid: true,
	// Noctty: true,
}

// OpenURLCommand is the command to cause the system default browser to open a
// URL.
func OpenURLCommand(url string) *exec.Cmd {
	return exec.Command("open", url)
}

// OpenFileCommand is the command to open a file using the system default
// application for its file type.
func OpenFileCommand(file string) *exec.Cmd {
	return exec.Command("open", file)
}

// IsAdmin returns whether the user is an Administrator.
func IsAdmin() bool { return false } // only valid on Windows

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
	var err error

	if serverPrintCmd, err = exec.LookPath("lpr"); err != nil || serverPrintCmd == "" {
		if serverPrintCmd, err = exec.LookPath("lp"); err != nil {
			serverPrintCmd = ""
		}
	}
}
