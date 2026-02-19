//go:build windows

package cmd

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/spf13/pflag"
)

const (
	outpostEditSlug = `Start an editor on an existing unsent message`
	outpostEditHelp = `
usage: packet outpost edit msgfile msgindex

Command called by Outpost to start editing an existing (unsent) message. The parameters are:

msgfile   ⇥Filename of the file containing the existing message
msgindex  ⇥Index of the message in the Outpost database
`
)

func cmdOutpostEdit(args []string) (err error) {
	var (
		msgfile string
		values  = make(url.Values)
	)
	flags := pflag.NewFlagSet("o-edit", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdOutpostHelp([]string{"edit"})
	} else if err != nil {
		cio.Open().Error("%s", err.Error())
		return usage(outpostEditHelp)
	}
	if len(args) != 2 {
		return usage(outpostEditHelp)
	}
	registerForms()
	// Verify that the message file is reachable.
	msgfile = decodeOutpostArg(args[0])
	if _, err := os.Stat(msgfile); err != nil {
		slog.Error("os.Stat", "f", msgfile, "err", err)
		return fmt.Errorf("%s: %s", msgfile, err)
	}
	values.Set("msgfile", msgfile)
	values.Set("index", decodeOutpostArg(args[1]))
	if err := outpostNECommon(values, "/outpost-edit"); err != nil {
		return err
	}
	slog.Info("opened browser to edit message", "f", msgfile)
	return nil
}
