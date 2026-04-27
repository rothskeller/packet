package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/form/formdefs"
	"github.com/spf13/pflag"
)

const (
	formsUpdateSlug = `Check for forms updates`
	formsUpdateHelp = `
usage: packet forms update

The "packet forms update" command checks online for updates to any installed forms bundles.  If any updates are found, they are installed.
`
)

func cmdFormsUpdate(args []string) (err error) {
	var (
		rmfile string
		data   []byte
	)
	flags := pflag.NewFlagSet("f-update", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdFormsHelp([]string{"update"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(formsUpdateHelp)
	}
	if len(args) != 0 {
		return usage(formsUpdateHelp)
	}
	if err = formdefs.CheckForUpdates(true, true); err != nil {
		return err
	}
	rmfile = filepath.Join(formdefs.FormsDir(), "README.txt")
	if data, err = os.ReadFile(rmfile); err == nil {
		os.Remove(rmfile)
		os.Stdout.Write(data)
	} else if !os.IsNotExist(err) {
		slog.Error("os.ReadFile", "f", rmfile, "err", err)
	}
	return nil
}
