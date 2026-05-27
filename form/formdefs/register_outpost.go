//go:build windows

package formdefs

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// launchINI is the contents of the file written to the *.ini sidecar of every
// *.launch file in the forms directory.  The $ADDON$ variable is replaced with
// the addon name.
const launchINI = `[$ADDON$]
bang-id=!$ADDON$!
fmt-date=mm/dd/yyyy
fmt-time=hh:mm
launch-submitted=never
launch-presend=always
launch-sent=always
launch-retrieved=always
launch-unread=always
launch-read=always
launch-convert=always
msg-dir=C:\PackItForms\Spool
cmd-new=C:\PackItForms\pifo.exe outpost new "$ADDON$" "{{ADDON_MSG_TYPE}}"
  "{{SETUP_NEW_MSG_NUMBER}}" "{{SETUP_ID_LEGAL_CALL}}" "{{SETUP_ID_LEGAL_NAME}}"
  "{{SETUP_ID_TAC_CALL}}" "{{SETUP_ID_TAC_NAME}}" "{{SETUP_ID_ACTIVE_CALL}}"
cmd-draft=C:\PackItForms\pifo.exe outpost edit "{{MSG_FILENAME}}" "{{MSG_INDEX}}"
cmd-ready=C:\PackItForms\pifo.exe outpost edit "{{MSG_FILENAME}}" "{{MSG_INDEX}}"
cmd-sent=C:\PackItForms\pifo.exe outpost view "{{MSG_FILENAME}}"
cmd-unread=C:\PackItForms\pifo.exe outpost view "{{MSG_FILENAME}}"
  "{{MSG_LOCAL_ID}}" "{{SETUP_ID_LEGAL_CALL}}" "{{SETUP_ID_LEGAL_NAME}}"
  "{{MSG_DATETIME_OP_RCVD}}"
cmd-retrieved=C:\PackItForms\pifo.exe outpost view "{{MSG_FILENAME}}"
  "{{MSG_LOCAL_ID}}" "{{SETUP_ID_LEGAL_CALL}}" "{{SETUP_ID_LEGAL_NAME}}"
  "{{MSG_DATETIME_OP_RCVD}}"
cmd-read=C:\PackItForms\pifo.exe outpost view "{{MSG_FILENAME}}"
  "{{MSG_LOCAL_ID}}" "{{SETUP_ID_LEGAL_CALL}}" "{{SETUP_ID_LEGAL_NAME}}"
  "{{MSG_DATETIME_OP_RCVD}}"
cmd-convert=C:\PackItForms\pifo.exe outpost convert "{{MSG_FILENAME}}"
  "{{MSG_STATE}}" "{{MSG_LOCAL_ID}}" "{{SETUP_ID_LEGAL_CALL}}"
  "{{SETUP_ID_LEGAL_NAME}}" "{{MSG_DATETIME_OP_RCVD}}" "{{SPOOL_DIR}}"
  "{{COPY_NAMES}}"
`

// oldPIFOPath is the path prefix of references to pre-4.0 PackItForms.  All
// such references are unconditionally removed.
const oldPIFOPath = `C:\PackItForms\Outpost\`

var includeLineRE = regexp.MustCompile(`(?i)^INCLUDE\s*(.*)$`)

// UpdateOutpostConfiguration updates the Outpost configuration files to point
// to the current set of "addons" we support.  If a remove map is supplied, all
// addons that are keys in that map are removed from the configuration.  This
// function is a no-op if packet is not connected to Outpost.
func UpdateOutpostConfiguration(datadir string, remove map[string]bool) (err error) {
	// In a fresh install of 164B, the Outpost Launch.ini file contains a
	// reference to the old SCCoPIFO.launch file, a LINE, and the generic
	// ICS-213 form, and there is no Outpost Launch.local file.  Earlier
	// installations could have any random set of entries in either or both
	// files.  The Launch.ini file gets replaced when Outpost is updated;
	// the Launch.local file (if any) does not.
	//
	// Our algorithm here is, first, to remove all non-commented entries
	// from the Launch.ini file.  We don't want any of them.  And then, to
	// update or add entries in Launch.local, creating it if necessary.
	var (
		iniFName    string
		launchFiles []string
		localFName  string
		addons      = make(map[string]string)
	)
	iniFName = filepath.Join(datadir, "Launch.ini")
	if _, err = os.Stat(iniFName); err != nil {
		slog.Error("os.Stat", "f", iniFName, "err", err)
		return fmt.Errorf("can't locate Outpost data directory: %s", err)
	}
	// Make sure the spool directory exists.
	const spoolDir = `C:\PackItForms\Spool`
	if err = os.MkdirAll(spoolDir, 0777); err != nil {
		slog.Error("os.MkdirAll", "d", spoolDir, "err", err)
		return fmt.Errorf("can't create spool directory: %s", err)
	}
	// Make sure "remove" is a real map and all its keys are downcased.
	if remove == nil {
		remove = make(map[string]bool)
	}
	for k := range remove {
		if k != strings.ToLower(k) {
			remove[strings.ToLower(k)] = true
			delete(remove, k)
		}
	}
	// What addons do we have in our forms?
	launchFiles, _ = filepath.Glob(filepath.Join(FormsDir(), "*", "*.launch"))
	for _, lf := range launchFiles {
		addonName := strings.TrimSuffix(strings.ToLower(filepath.Base(lf)), ".launch")
		if addons[addonName] != "" {
			slog.Error("addon multiply defined", "addon", addonName, "1", addons[addonName], "2", lf)
			return fmt.Errorf("addon %q defined in both %s and %s", addonName, addons[addonName], lf)
		}
		addons[addonName] = lf
		delete(remove, addonName)
		if err = writeAddonINI(lf); err != nil {
			return err
		}
	}
	// Update Launch.ini to remove all non-commented lines from it.
	if err = emptyLaunchINI(iniFName); err != nil {
		return fmt.Errorf("can't update %s: %s", iniFName, err)
	}
	// Update Launch.local to update any remaining paths and/or add
	// remaining items.
	localFName = filepath.Join(datadir, "Launch.local")
	if err = UpdateLaunchLocal(localFName, addons, remove, oldPIFOPath); err != nil {
		return fmt.Errorf("can't update %s: %s", localFName, err)
	}
	return nil
}

// writeAddonINI (re)writes the «addon».ini sidecar to a «addon».launch file.
func writeAddonINI(launchFile string) (err error) {
	var base = launchFile[:len(launchFile)-7] // remove .launch, ignoring case
	var addonName = filepath.Base(base)
	var ini = strings.ReplaceAll(launchINI, "$ADDON$", addonName)
	ini = strings.ReplaceAll(ini, "\n", "\r\n")
	var inifn = base + ".ini"
	if err = os.WriteFile(inifn, []byte(ini), 0666); err != nil {
		slog.Error("os.WriteFile", "f", inifn, "err", err)
		return fmt.Errorf("can't create %s.ini: %s", base, err)
	}
	slog.Info("wrote addon.ini file", "f", inifn)
	return nil
}

// emptyLaunchINI removes all non-commented lines from the Outpost Launch.ini
// file.  (Comments start with a '/'.)
func emptyLaunchINI(filename string) (err error) {
	var (
		fh    *os.File
		scan  *bufio.Scanner
		lines []string
	)
	if fh, err = os.Open(filename); err != nil && !os.IsNotExist(err) {
		slog.Error("os.Open", "f", filename, "err", err)
		return err
	}
	if err == nil {
		scan = bufio.NewScanner(fh)
		for scan.Scan() {
			line := scan.Text()
			if strings.HasPrefix(line, "/") || strings.TrimSpace(line) == "" {
				lines = append(lines, line)
			}
		}
		fh.Close()
	}
	if err = os.WriteFile(filename, []byte(strings.Join(lines, "\r\n")), 0666); err != nil {
		slog.Error("os.WriteFile", "f", filename, "err", err)
		return err
	}
	return nil
}

// UpdateLaunchLocal modifies or creates the Launch.local file.  The entries in
// "remove" or matching "removePath" are removed.  Leading and trailing LINE
// lines are removed.  Any entries in "add" that already exist in the file are
// updated.  If there are remaining entries in "add", they are added at the end
// of the file, with a LINE separator before each unless the file was otherwise
// empty.  An error is returned only if the file cannot be read or written.
func UpdateLaunchLocal(filename string, add map[string]string, remove map[string]bool, removePath string) (err error) {
	var (
		fh       *os.File
		scan     *bufio.Scanner
		lines    []string
		modified bool
	)
	if fh, err = os.Open(filename); err != nil && !os.IsNotExist(err) {
		slog.Error("os.Open", "f", filename, "err", err)
		return err
	}
	if err == nil {
		scan = bufio.NewScanner(fh)
		for scan.Scan() {
			line := scan.Text()
			match := includeLineRE.FindStringSubmatch(line)
			if match == nil {
				lines = append(lines, line)
				continue
			}
			launchfile := strings.TrimSpace(match[1])
			if !strings.HasSuffix(strings.ToLower(launchfile), ".launch") {
				lines = append(lines, line)
				continue
			}
			addonName := strings.TrimSuffix(strings.ToLower(filepath.Base(launchfile)), ".launch")
			if sb := add[addonName]; sb != "" {
				if !strings.EqualFold(sb, launchfile) {
					slog.Info("modified include in launch file", "f", filename, "from", launchfile, "to", sb)
					line = "INCLUDE " + sb
					modified = true
				}
				lines = append(lines, line)
				delete(add, addonName)
				remove[addonName] = true
				continue
			} else if remove[addonName] {
				slog.Info("removed include from launch file", "f", filename, "include", launchfile)
				modified = true
				continue
			} else if strings.HasPrefix(launchfile, removePath) {
				slog.Info("removed include from launch file", "f", filename, "include", launchfile)
				modified = true
				continue
			}
			lines = append(lines, line)
		}
		fh.Close()
	}
	// Ensure the file starts with a single blank line and no LINE entries.
	// However, don't set the modified flag; we won't write a new file based
	// only on that change.
	for len(lines) != 0 && (lines[0] == "" || lines[0] == "LINE") {
		lines = lines[1:]
	}
	lines = append([]string{""}, lines...)
	// Ensure the file has no trailing LINE entries or blank lines (other
	// than the initial blank line).
	for len(lines) > 1 && (lines[len(lines)-1] == "" || lines[len(lines)-1] == "LINE") {
		lines = lines[:len(lines)-1]
	}
	// Add any remaining "add" entries to the file.  Put a LINE entry before
	// each, unless the file has only a blank line.
	for addonName, launchfile := range add {
		slog.Info("added include to launch file", "f", filename, "include", launchfile)
		if len(lines) != 1 {
			lines = append(lines, "LINE")
		}
		lines = append(lines, "INCLUDE "+launchfile)
		modified = true
		delete(add, addonName)
	}
	if !modified {
		slog.Debug("no changes to launch file", "f", filename)
		return nil
	}
	if len(lines) != 0 && lines[len(lines)-1] != "" {
		lines = append(lines, "") // file should end with blank line
	}
	if err = os.WriteFile(filename, []byte(strings.Join(lines, "\r\n")), 0666); err != nil {
		slog.Error("os.WriteFile", "f", filename, "err", err)
		return err
	}
	return nil
}
