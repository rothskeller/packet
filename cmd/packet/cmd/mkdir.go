package cmd

import (
	"os"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/errors"
	"github.com/spf13/pflag"
)

const (
	mkdirSlug = `Create a directory`
	mkdirHelp = `
usage: mkdir «directory»

The "mkdir" command creates the named directory, which must not already exist.  It does not switch to the newly created directory; use the "cd" command for that.
`
)

func cmdMkdir(args []string) (err error) {
	flags := pflag.NewFlagSet("mkdir", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"mkdir"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(mkdirHelp)
	}
	if len(args) != 1 {
		return usage(mkdirHelp)
	}
	if err = os.Mkdir(args[0], 0777); err != nil {
		return errors.NewF("Unable to create directory: %s", err)
	}
	return nil
}
