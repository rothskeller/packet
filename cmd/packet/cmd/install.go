//go:build windows

package cmd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gotopkg/mslnk/pkg/mslnk"
	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/form/formdefs"
	"github.com/rothskeller/packet/form/pifover"
	"github.com/rothskeller/packet/forms"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	packetRoot         = `C:\PackItForms`
	packetExe          = packetRoot + `\packet.exe`
	pifoExe            = packetRoot + `\pifo.exe`
	outpostDataDirFile = packetRoot + `\outpost-data-dir.txt`
	startMenuDir       = `C:\ProgramData\Microsoft\Windows\Start Menu\Programs\SCCo Packet`
	uninstallLink      = startMenuDir + `\Uninstall PackItForms.lnk`
	registryKey        = `Software\Microsoft\Windows\CurrentVersion\Uninstall\SCCoPackItForms`
	oldPackItForms1    = `C:\PackItForms\Outpost\SCCo\bin\SCCoPIFO.exe`
	oldPackItForms2    = `C:\PackItForms\Outpost\SCCo\bin\Outpost_Forms.exe`
)

var installCmd = &cobra.Command{
	Use:   "install [outpost-data-dir]",
	Short: "Installs the packet software and connects it to Outpost",
	Long: `
The "packet install" command connects the packet software to Windows and
Outpost.  The packet software does not need to be "installed;" it can be run
from the command line without any installation process.  However, installing it
makes it easier to invoke and to maintain:
  - It adds entries for the packet software in the Windows Start menu.
  - It adds entries for the packet software in the Windows Registry, so that
    the packet software appears in the list of installed applications and can
    be uninstalled from there.
  - It updates the PackItForms 3.x C:\PackItForms\Outpost\SCCo\manual.cmd
    script, if it exists, to start manual mode operations.

If an Outpost data directory is given on the command line, "packet install"
will also modify the Outpost configuration so that Outpost will use the packet
software for forms messages.`,
	Args:                  cobra.RangeArgs(0, 1),
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var odd string

		if err = checkRunningAsAdmin(); err != nil {
			return err
		}
		stopServers()
		if err = installExecutables(); err != nil {
			return err
		}
		if err = addToStartMenu(); err != nil {
			return err
		}
		if err = addToRegistry(); err != nil {
			return err
		}
		if err = formdefs.RegisterForms(); err != nil && formdefs.FormsFS == forms.EmbeddedForms {
			return err
		}
		if len(args) != 0 {
			odd = args[0]
		}
		if err = connectToOutpost(odd); err != nil {
			return err
		}
		cleanupOldPackItForms()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}

// checkRunningAsAdmin ensures that the program is running with Administrator
// privileges.  If not, it relaurches itself with Administrator privileges, and
// the non-privileged copy exits.  The only possible error return is if we're
// not running as Admin and we can't relaunch.
func checkRunningAsAdmin() (err error) {
	if osdep.IsAdmin() {
		return nil
	}
	// We're not running as admin.  Relaunch with admin privileges.
	verb := "runas"
	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exe, _ := os.Executable()
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwd, _ := os.Getwd()
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	args := strings.Join(os.Args[1:], " ")
	argPtr, _ := syscall.UTF16PtrFromString(args)
	showCmd := int32(1)
	if err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd); err != nil {
		slog.Error("relaunch as Administrator", "err", err)
		return fmt.Errorf("relaunching as Administrator: %s", err)
	}
	slog.Debug("relaunched as Administrator")
	os.Exit(0)
	return nil // not reachable
}

// stopServers stops any running servers so that we can replace or remove their
// files.
func stopServers() {
	if stopOldServer(pifoExe, "server", "stop") ||
		stopOldServer(oldPackItForms1, "stop") ||
		stopOldServer(oldPackItForms2, "stop") {
		time.Sleep(2 * time.Second)
	}
}

// stopOldServer stops the old PackItForms server at the specified path, if
// it's running.
func stopOldServer(exeFile string, args ...string) bool {
	if _, err := os.Stat(exeFile); err != nil {
		return false
	}
	cmd := exec.Command(exeFile, args...)
	if err := cmd.Run(); err != nil {
		slog.Warn("stop old server", "exe", exeFile, "err", err)
		return false
	} else {
		slog.Info("stopped old server", "exe", exeFile)
		return true
	}
}

// installExecutables makes sure that the packet executable is installed at
// C:\PackItForms\packet.exe (console mode) and C:\PackItForms\pifo.exe (GUI
// mode).
func installExecutables() (err error) {
	var selfFile string

	// Make sure the packetRoot directory exists.
	if err = os.MkdirAll(packetRoot, 0777); err != nil {
		slog.Error("os.MkdirAll "+packetRoot, "err", err)
		return fmt.Errorf("can't create %s: %s", packetRoot, err)
	}
	// If that isn't the location we're running from, copy our own
	// executable there.
	if selfFile, err = os.Executable(); err != nil {
		slog.Error("os.Executable", "err", err)
		return fmt.Errorf("can't locate installer executable: %s", err)
	}
	if !strings.EqualFold(packetExe, selfFile) {
		if err = copyFile(selfFile, packetExe); err != nil {
			return fmt.Errorf("can't install %s: %s", packetExe, err)
		}
	}
	// Unconditionally copy our own executable to pifo.exe and mark it to
	// be a GUI app.  (This is a bit kludgey, but it saves us having to
	// deliver two large executables that differ by a single byte.)
	if err = copyFile(selfFile, pifoExe); err != nil {
		return fmt.Errorf("can't install %s: %s", pifoExe, err)
	}
	if err = makeGUI(pifoExe); err != nil {
		os.Remove(pifoExe)
		return fmt.Errorf("can't mark %s as GUI app: %s", pifoExe, err)
	}
	return nil
}

// copyFile copies a file.
func copyFile(src, dest string) (err error) {
	var (
		sf *os.File
		df *os.File
	)
	if sf, err = os.Open(src); err != nil {
		slog.Error("os.Open", "f", src, "err", err)
		return fmt.Errorf("%s: %s", src, err)
	}
	defer sf.Close()
	if df, err = os.Create(dest); err != nil {
		slog.Error("os.Create", "f", dest, "err", err)
		return fmt.Errorf("%s: %s", dest, err)
	}
	if _, err = io.Copy(df, sf); err != nil {
		slog.Error("io.Copy", "src", src, "dest", dest, "err", err)
		df.Close()
		os.Remove(dest)
		return fmt.Errorf("%s: %s", dest, err)
	}
	if err = df.Close(); err != nil {
		slog.Error("os.Close", "f", dest, "err", err)
		os.Remove(dest)
		return fmt.Errorf("%s: %s", dest, err)
	}
	slog.Info("copied", "src", src, "dest", dest)
	return nil
}

// makeGUI marks a Windows executable as using the GUI rather than CUI
// subsystem, i.e., it doesn't open a console window when invoked.
func makeGUI(exeFile string) (err error) {
	var (
		ef     *os.File
		buf    [4]byte
		offset uint32
	)
	if ef, err = os.OpenFile(exeFile, os.O_RDWR, 0666); err != nil {
		slog.Error("os.Open", "f", exeFile, "err", err)
		return err
	}
	defer ef.Close()
	if _, err = ef.ReadAt(buf[:], 0x3c); err != nil {
		slog.Error("fh.ReadAt", "f", exeFile, "off", 0x3c, "err", err)
		return err
	}
	offset = binary.LittleEndian.Uint32(buf[:])
	offset += 0x5c
	if _, err = ef.ReadAt(buf[:1], int64(offset)); err != nil {
		slog.Error("fh.ReadAt", "f", exeFile, "off", offset, "err", err)
		return err
	}
	if buf[0] != 3 {
		slog.Error("wrong exe subsystem", "f", exeFile, "subsys", buf[0])
		return fmt.Errorf("subsystem is %d, expected 3", buf[0])
	}
	buf[0] = 2
	if _, err = ef.WriteAt(buf[:1], int64(offset)); err != nil {
		slog.Error("fh.WriteAt", "f", exeFile, "off", offset, "err", err)
		return err
	}
	if err = ef.Close(); err != nil {
		slog.Error("fh.Close", "f", exeFile, "err", err)
		return err
	}
	slog.Info("marked as GUI application", "f", exeFile)
	return nil
}

// addToStartMenu creates the desired entries in the Start Menu for the packet
// software.
func addToStartMenu() (err error) {
	if err = os.MkdirAll(startMenuDir, 0777); err != nil {
		slog.Error("os.MkdirAll", "d", startMenuDir, "err", err)
		return fmt.Errorf("create Start Menu folder: %s", err)
	}
	link := &mslnk.ShellLink{
		ShellLinkHeader: mslnk.Header(),
		LinkTargetIDList: mslnk.LinkTargetIDList{
			ItemIDList: []mslnk.ItemID{
				mslnk.ItemIDCLSID(mslnk.ItemIDMagic["MY_COMPUTER"]),
				mslnk.ItemIDDrive(pifoExe[:3]),
				mslnk.ItemIDFile(pifoExe[3:]),
			},
		},
		StringData: mslnk.StringData{
			"CommandLineArguments": mslnk.StringDataStruct("uninstall"),
		},
	}
	link.LinkTargetIDList.Size()
	link.ShellLinkHeader.LinkFlags["HasLinkTargetIDList"] = true
	link.ShellLinkHeader.LinkFlags["HasArguments"] = true
	link.ShellLinkHeader.LinkFlags["RunAsUser"] = true // as Administrator
	link.ShellLinkHeader.FileAttributes["FILE_ATTRIBUTE_NORMAL"] = true
	link.ShellLinkHeader.Update()
	if err = link.Save(uninstallLink); err != nil {
		slog.Error("link.Save", "f", uninstallLink, "err", err)
		return fmt.Errorf("create Uninstall entry in Start Menu folder: %s", err)
	}
	slog.Info("added to Start Menu", "f", uninstallLink)
	return nil
}

// addToRegistry creates the desired entries in the registry to allow the
// packet software to be listed and uninstalled through the standard Windows
// mechanism.
func addToRegistry() (err error) {
	var (
		key  registry.Key
		size uint32 = 20000
	)
	if stat, err := os.Stat(packetExe); err == nil {
		size = uint32((stat.Size() + 1023) / 1024)
	}
	if key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, registryKey, registry.SET_VALUE); err != nil {
		slog.Error("registry.CreateKey", "err", err)
		return fmt.Errorf("create registry key: %s", err)
	}
	err = errors.Join(
		key.SetStringValue("DisplayName", "SCCo PackItForms"),
		key.SetStringValue("UninstallString", `C:\PackItForms\pifo.exe uninstall`),
		key.SetStringValue("Publisher", "Santa Clara County ARES/RACES"),
		key.SetStringValue("URLInfoAbout", "https://www.scc-ares-races.org"),
		key.SetStringValue("DisplayVersion", pifover.PIFOVersion),
		key.SetDWordValue("VersionMajor", pifover.PIFOVersionMajor),
		key.SetDWordValue("VersionMinor", pifover.PIFOVersionMinor),
		key.SetDWordValue("NoModify", 1),
		key.SetDWordValue("NoRepair", 1),
		key.SetDWordValue("EstimatedSize", size),
	)
	if err != nil {
		slog.Error("SetValues", "err", err)
		return fmt.Errorf("write registry values: %s", err)
	}
	slog.Info("added to registry", "key", registryKey)
	return nil
}

// connectToOutpost updates the Outpost configuration in the specified
// directory to point to this packet software as the handler for all of the
// Outpost addon types it supports.
func connectToOutpost(outpostDataDir string) (err error) {
	// If we're invoked from Git bash or similar, the slashes might be
	// wrong.
	outpostDataDir = strings.ReplaceAll(outpostDataDir, "/", `\`)
	// If they didn't give us a directory, maybe we already know it from a
	// previous install?
	if outpostDataDir == "" {
		if data, err := os.ReadFile(outpostDataDirFile); err == nil {
			dir := strings.TrimSpace(string(data))
			if _, err = os.Stat(filepath.Join(dir, "Launch.ini")); err == nil {
				outpostDataDir = dir
			}
		}
	}
	// If we still don't have an Outpost directory, try the common location.
	if outpostDataDir == "" {
		if _, err = os.Stat(`C:\SCCo Packet\Launch.ini`); err == nil {
			outpostDataDir = `C:\SCCo Packet`
		} else {
			slog.Debug("not connecting to Outpost: no data directory")
			return nil // no Outpost to connect to.
		}
	}
	// If the location we're given doesn't have a launch.ini, it's not
	// correct.
	if _, err = os.Stat(filepath.Join(outpostDataDir, "Launch.ini")); err != nil {
		slog.Error("invalid data directory", "d", outpostDataDir, "err", "no Launch.ini")
		return fmt.Errorf("%s: not a valid Outpost data directory (no Launch.ini file)", outpostDataDir)
	}
	// Record the location.
	if err = os.WriteFile(outpostDataDirFile, []byte(outpostDataDir), 0666); err != nil {
		slog.Error("os.WriteFile", "err", err)
		return fmt.Errorf("recording location of Outpost: %s: %s", outpostDataDirFile, err)
	}
	slog.Info("Outpost data directory", "d", outpostDataDir)
	// Update the Outpost configuration files.  If RegisterForms found any
	// updates, this may already have been done, but it's idempotent so
	// doing it again won't hurt.
	if err = formdefs.UpdateOutpostConfiguration(outpostDataDir, nil); err != nil {
		return fmt.Errorf("registering forms with Outpost: %s", err)
	}
	return nil
}

// cleanupOldPackItForms removes everything related to an SCCoPIFO 3.x
// installation.  (It does not touch non-SCCoPIFO PackItForms installations.)
// Any errors are ignored.
func cleanupOldPackItForms() {
	// We don't need to remove references from the Outpost configuration
	// because formdefs.UpdateOutpostConfiguration already did that.

	// Remove old registry entry.
	const oldRegistryKey = `Software\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall\SCCoPackItForms`
	if err := registry.DeleteKey(registry.LOCAL_MACHINE, oldRegistryKey); err != nil && !os.IsNotExist(err) {
		slog.Warn("registry.DeleteKey", "key", oldRegistryKey, "err", err)
	} else if err == nil {
		slog.Info("deleted registry key", "key", oldRegistryKey)
	}
	// Remove Start Menu entries.
	ents, _ := filepath.Glob(filepath.Join(startMenuDir, "Uninstall SCCo PackItForms for Outpost*.lnk"))
	for _, ent := range ents {
		if err := os.Remove(ent); err != nil {
			slog.Warn("os.Remove", "f", ent, "err", err)
		} else {
			slog.Info("removed old Start Menu entry", "f", ent)
		}
	}
	// Remove old PackItForms files.
	const oldDir = `C:\PackItForms\Outpost`
	if err := os.RemoveAll(oldDir); err != nil {
		slog.Warn("os.RemoveAll", "d", oldDir, "err", err)
	} else {
		slog.Info("removed old PackItForms files", "d", oldDir)
	}
}
