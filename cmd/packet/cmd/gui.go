package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/cmd/packet/osdep"
	"github.com/rothskeller/packet/v4/cmd/packet/server"
	"github.com/rothskeller/packet/v4/form/formdefs"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/spf13/pflag"
)

const (
	guiSlug = `Start the graphical (web) interface`
	guiHelp = `
usage: packet gui [«incident-dir»]
       packet web [«incident-dir»]

The "packet gui" command starts the web-based graphical user interface to the packet software.

If an incident directory is given, the graphical interface will start in that incident.  Otherwise, if the current directory is an incident directory, it will start there.  Otherwise it will ask the user to select or create an incident directory.
`
)

func cmdGUI(args []string) (err error) {
	var (
		dir     string
		created bool
		address string
		cmd     *exec.Cmd
		wg      *sync.WaitGroup
	)
	flags := pflag.NewFlagSet("gui", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"gui"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(guiHelp)
	}
	if len(args) > 1 {
		return usage(guiHelp)
	}
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
	if address, err = server.GetAddress(false); err != nil {
		return fmt.Errorf("finding server: %s", err)
	}
	if address == "" {
		// There's no server running; we need to start one.  But first,
		// before starting the server is the proper time to check for
		// forms updates.  Errors are logged but not returned.
		_ = formdefs.CheckForUpdates(false, true)
		if address, wg, err = guiStartServer(); err != nil {
			return fmt.Errorf("starting server: %s", err)
		}
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
	if wg != nil {
		wg.Wait()
	}
	return nil
}

// guiStartServerGoroutine starts the packet server in a goroutine of the
// current process, with a WaitGroup that waits for it to exit.
func guiStartServerGoroutine() (address string, wg *sync.WaitGroup, err error) {
	if err = formdefs.RegisterForms(); err != nil {
		return "", nil, err
	}
	wg = new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		server.Start(0, false)
		wg.Done()
	}()
	if address, err = server.GetAddress(true); err != nil {
		return "", wg, err
	} else if address == "" {
		return "", wg, errors.New("The packet server did not start up within the expected time.")
	}
	return address, wg, nil
}
