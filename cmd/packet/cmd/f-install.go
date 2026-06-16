package cmd

import (
	"io"
	"os"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/form/formdefs"
	"github.com/spf13/pflag"
)

const (
	formsInstallSlug = `Install a forms bundle`
	formsInstallHelp = `
usage: packet forms install { «filename» | «url» }

The "packet forms install" command installs a forms bundle from the specified location (either a local filename or an https:// URL).  The forms bundle must be signed with the digital signature for this version of the packet software.  It will replace any previously installed version of the same bundle, even if that version is newer.
`
)

func cmdFormsInstall(args []string) (err error) {
	var bundle, readme string

	flags := pflag.NewFlagSet("f-install", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdFormsHelp([]string{"install"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(formsInstallHelp)
	}
	if len(args) != 1 {
		return usage(formsInstallHelp)
	}
	if bundle, readme, err = formdefs.InstallBundle(args[0]); err != nil {
		return err
	}
	if readme != "" {
		io.WriteString(os.Stdout, readme)
	} else {
		cio.Open().Confirm("Installed forms bundle %q.\n", bundle)
	}
	return nil
}
