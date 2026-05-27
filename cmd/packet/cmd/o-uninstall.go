//go:build windows

package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/server"
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
		cio.Open().Error(err)
		return usage(outpostUninstallHelp)
	}
	if len(args) != 0 {
		return usage(outpostUninstallHelp)
	}
	if err = stopServer(); err != nil {
		return err
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
	removeSelf()
	return nil
}

// stopServer stops the running server if any.
func stopServer() (err error) {
	if addr, err := server.GetAddress(false); err != nil {
		return err
	} else if addr == "" {
		return nil
	} else if resp, err := http.Post(addr+"/stop", "text/plain", nil); err != nil {
		slog.Error("POST /stop", "url", addr, "err", err)
		return err
	} else if resp.StatusCode >= 400 {
		slog.Error("POST /stop", "url", addr, "code", resp.StatusCode, "status", resp.Status)
		resp.Body.Close()
		return fmt.Errorf("stop request: %s", resp.Status)
	} else {
		slog.Info("sent /stop to server")
		resp.Body.Close()
		return nil
	}
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
	if err = formdefs.UpdateLaunchLocal(filepath.Join(dir, "Launch.local"), nil, nil, packetRoot); err != nil {
		return fmt.Errorf("updating %s\\Launch.local: %s", dir, err)
	}
	os.Remove(outpostDataDirFile)
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
	fmt.Fprintf(fh, "SLEEP 2\r\nDEL %s\r\nDEL %s\r\n", packetExe, pifoExe)
	fh.Close()
	if err = exec.Command(fh.Name()).Start(); err != nil {
		slog.Warn("cmd.Start", "err", err)
	}
}
