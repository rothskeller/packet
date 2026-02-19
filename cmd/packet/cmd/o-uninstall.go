//go:build windows

package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/form/formdefs"
	"github.com/spf13/pflag"
	"golang.org/x/sys/windows/registry"
)

const (
	outpostUninstallSlug = `Uninstalls the packet software and disconnects from Outpost`
	outpostUninstallHelp = `
usage: packet outpost uninstall

The "packet outpost uninstall" command undoes the operations of the "packet outpost install" command, removing the packet software and its files from the system.  If the packet software was connected to Outpost, that connection is removed.
`
)

func cmdOutpostUninstall(args []string) (err error) {
	flags := pflag.NewFlagSet("o-uninstall", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdOutpostHelp([]string{"uninstall"})
	} else if err != nil {
		cio.Open().Error("%s", err.Error())
		return usage(outpostUninstallHelp)
	}
	if len(args) != 0 {
		return usage(outpostUninstallHelp)
	}
	if err = checkRunningAsAdmin(); err != nil {
		return err
	}
	if err = disconnectFromOutpost(); err != nil {
		return err
	}
	if err = removeFromStartMenu(); err != nil {
		return err
	}
	if err = removeFromRegistry(); err != nil {
		return err
	}
	removeFileTree()
	removeSelf()
	return nil
}

// disconnectFromOutpost removes any references to this Packet installation
// from any Outpost data directory/ies that we were installed into.
func disconnectFromOutpost() (err error) {
	var dir string

	// Open and read the outpost-data-dir.txt file.
	if data, err := os.ReadFile(outpostDataDirFile); os.IsNotExist(err) {
		slog.Debug("no outpost data dir to disconnect from")
		return nil // not connected to Outpost
	} else if err != nil {
		slog.Error("os.ReadFile", "f", outpostDataDirFile, "err", err)
		return fmt.Errorf("reading %s: %s", outpostDataDirFile, err)
	} else {
		dir = strings.TrimSpace(string(data))
	}
	if err = formdefs.UpdateLaunchFile(filepath.Join(dir, "Launch.ini"), nil, nil, nil, packetRoot); err != nil {
		return fmt.Errorf("updating %s\\Launch.ini: %s", dir, err)
	}
	if err = formdefs.UpdateLaunchFile(filepath.Join(dir, "Launch.local"), nil, nil, nil, packetRoot); err != nil {
		return fmt.Errorf("updating %s\\Launch.local: %s", dir, err)
	}
	return nil
}

// removeFromStartMenu removes the entries we added to the Windows Start Menu.
func removeFromStartMenu() (err error) {
	if err = os.Remove(uninstallLink); err != nil && !os.IsNotExist(err) {
		slog.Warn("os.Remove", "f", uninstallLink, "err", err)
	} else {
		slog.Info("removed uninstall link from Start Menu", "f", uninstallLink)
	}
	return nil
}

// removeFromRegistry removes the application entry from the registry.
func removeFromRegistry() (err error) {
	if err = registry.DeleteKey(registry.LOCAL_MACHINE, registryKey); err != nil {
		slog.Warn("registry.DeleteKey", "key", registryKey, "err", err)
	} else {
		slog.Info("removed registry key", "key", registryKey)
	}
	return nil

}

// removeFileTree removes as much of the C:\PackItForms tree as possible.
func removeFileTree() {
	// We're not going to be able to remove the whole thing, because our
	// own executable is in there and is in use.  And os.RemoveAll would
	// quit when it got to that and leave other stuff around that shouldn't
	// be.  So we'll remove each item in the directory separately.
	ents, _ := filepath.Glob(filepath.Join(packetRoot, "*"))
	for _, ent := range ents {
		if err := os.RemoveAll(ent); err != nil {
			slog.Warn("os.RemoveAll", "d", ent, "err", err)
		}
	}
	// And we will try to remove the directory too, just in case we're
	// running from somewhere else.
	os.RemoveAll(packetRoot)
	slog.Info("removed most/all of packet files", "d", packetRoot)
}

// removeSelf tries to remove our own executable and the residual directory.
func removeSelf() {
	var (
		fh  *os.File
		err error
	)
	// To do this, we'll create a small batch file in temp with
	// instructions to do the removals after a pause, and then let this
	// executable finish and exit.
	if fh, err = os.CreateTemp("", "uninstall*.cmd"); err != nil {
		slog.Warn("os.CreateTemp", "err", err)
		return
	}
	exe, _ := os.Executable()
	fmt.Fprintf(fh, "SLEEP 2\r\nDEL %s\r\nRMDIR %s\r\n", exe, packetRoot)
	fh.Close()
	if err = exec.Command(fh.Name()).Start(); err != nil {
		slog.Warn("cmd.Start", "err", err)
	}
}
