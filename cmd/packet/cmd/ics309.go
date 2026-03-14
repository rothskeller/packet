package cmd

import (
	"log/slog"
	"os"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/spf13/pflag"
)

const (
	ics309Slug = `Create and show the ICS-309 log for the incident`
	ics309Help = `
usage: packet ics309
       packet 309

The "packet ics309" command shows the ICS-309 communications log for the incident in the system default PDF viewer.  If the ICS-309 PDF has not already been generated for the incident, it will generate it.
`
)

func cmdICS309(args []string) (err error) {
	var sig string

	flags := pflag.NewFlagSet("ics309", pflag.ContinueOnError)
	flags.StringVarP(&sig, "signature", "s", "", "Signature to add to the generated ICS-309")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"ics309"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(ics309Help)
	}
	if len(args) != 0 {
		return usage(ics309Help)
	}
	if sig != "" {
		// Remove any existing ICS-309 (and thus regenerate it)
		// because the signature could be different.
		os.Remove("ics309.pdf")
	}
	if _, err = os.Stat("ics309.pdf"); os.IsNotExist(err) {
		if sig == "" {
			if sig, err = askForSignature(); err != nil {
				return err
			}
		}
		if err = incRead(func(i *incident.Incident) error {
			return i.GenerateICS309(sig)
		}); err != nil {
			return err
		}
		slog.Debug("Generated missing ics309.pdf")
	} else if err != nil {
		return err
	}
	var showcmd = osdep.OpenFileCommand("ics309.pdf")
	if err := showcmd.Start(); err != nil {
		return errors.NewF("Unable to start PDF viewer: %s", err)
	}
	go func() { showcmd.Wait() }()
	return nil
}

func askForSignature() (string, error) {
	c := cio.Open()
	if c.InputIsTerm && c.OutputIsTerm {
		// If we're running interactively, prompt for a
		// signature.
		var valueWidth int
		for f := range incident.ICS309FormDef().AllFields() {
			if f.Tag == "Signature" {
				valueWidth = f.CharWidth(0)
			}
		}
		c.Confirm("NOTICE: ⇥By providing a signature, you are making a legal assertion that the ICS-309 log is accurate.  Do not provide a signature unless you are sure of that.")
		_, sig, err := c.EditField("Signature", 0, "", valueWidth, nil,
			`This is the signature to be added at the bottom of the ICS-309 form.`, "", false, false, nil)
		return sig, err
	}
	return "", nil
}
