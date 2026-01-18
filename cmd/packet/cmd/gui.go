package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/rothskeller/packet/incident"
)

var guiCmd = &cobra.Command{
	Use:   "gui [incident-directory]",
	Short: "Starts a packet GUI",
	Long: `Opens a browser window for management of packet messages in an incident
directory.  The incident directory is chosen as follows:
  - If a directory is given on the command line, it is used.  It will be
    created if it does not already exist, and initialized as a new incident if
    it does not already contain incident data.
  - If the current working directory is an incident directory, it is used.
  - Otherwise, if the most-recently-used incident directory still exists and
    contains incident data, it is used.
  - Otherwise, the user is prompted to choose or create an incident directory.

Generally speaking, a new directory should be used for each incident.  In other
words, any time you would start a new ICS-309 communications log, you should
work in a new incident directory.`,
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	Args:                  cobra.RangeArgs(0, 1),
	RunE: func(_ *cobra.Command, args []string) (err error) {
		var (
			dir     string
			created bool
			address string
			cmd     *exec.Cmd
		)
		if len(args) > 0 {
			// We were given an incident directory on the command
			// line.  Validate it.
			dir = args[0]
			if abs, err := filepath.Abs(dir); err != nil {
				return err
			} else {
				dir = abs
			}
			if !incident.IsIncident(dir) {
				// It's not currently an incident directory.
				// Is it a legal one?
				if incident.IsUnsafeIncidentDir(dir) {
					return errors.New("incident data should not be in root, home, Desktop, or Documents; make an incident-specific directory instead")
				}
				// Does it exist, or can we create it?
				if err = os.MkdirAll(dir, 0777); err != nil {
					return err
				}
				if err = incident.Create(dir, func(*incident.Incident) error { return nil }); err != nil {
					return err
				}
				created = true
			}
			// It's good to use.
		}
		if dir == "" {
			// Nothing on the command line.  Is the current working
			// dir an incident?
			if cwd, err := os.Getwd(); err == nil && incident.IsIncident(cwd) {
				dir = cwd
			}
		}
		if dir == "" {
			// Still no incident dir chosen.  Get the last one we
			// used.  If it's still an incident, use it.
			if dirs := incident.GetIncDefaults().IncidentDirs; len(dirs) != 0 {
				if last := dirs[len(dirs)-1]; incident.IsIncident(last) {
					dir = last
				}
			}
		}
		// Get the server address.  This also starts the server if not already
		// running.
		if address, err = server.GetAddress(true); err != nil {
			return fmt.Errorf("starting server: %s", err)
		}
		// Build the request URL.
		if dir == "" {
			cwd, _ := os.Getwd()
			address += "/incident-open?dir=" + url.QueryEscape(cwd)
		} else if created {
			address += "/incident-config?dir=" + url.QueryEscape(dir)
		} else {
			address += "/incident?dir=" + url.QueryEscape(dir)
		}
		// Open that URL in a browser.
		cmd = osdep.OpenURLCommand(address)
		if err = cmd.Run(); err != nil {
			slog.Error("cmd.Run", "url", address, "err", err)
			return fmt.Errorf("open GUI mode in browser: %s", err)
		}
		slog.Info("opened browser for GUI mode", "inc", dir)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(guiCmd)
}
