package cmd

import (
	"log/slog"
	"os"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/cmd/cmdutil"
	"github.com/rothskeller/packet/incident"
	"github.com/spf13/cobra"
)

var ics309Cmd = &cobra.Command{
	Use:     "ics309 [-s signature]",
	Aliases: []string{"309"},
	Short:   "Show the ICS-309 log for the incident",
	Long: `Creates ics309.pdf in the incident directory, if it does not already exist, containing the PDF-rendered ICS-309 communications log for the incident. Then, opens that file in the system-default PDF viewer if any.

The signature for the generated log can be provided with the --signature (or -s) flag.  In interactive mode, the software will prompt for it if not provided.  Note that by providing a signature, you are making a legal assertion that the log is accurate.  Do not sign it unless you are sure of that. `,
	Args:                  cobra.NoArgs,
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			dir string
			sig string
		)
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		if sig, _ = cmd.Flags().GetString("signature"); sig != "" {
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
			if err = incident.Read(dir, func(i *incident.Incident) error {
				return i.GenerateICS309(sig)
			}); err != nil {
				return err
			}
			slog.Debug("Generated missing ics309.pdf")
		} else if err != nil {
			return err
		}
		return cmdutil.ShowPDF("ics309.pdf")
	},
}

func init() {
	ics309Cmd.Flags().StringP("signature", "s", "", "Signature to add to the generated ICS-309")
	RootCmd.AddCommand(ics309Cmd)
}

func askForSignature() (string, error) {
	c := cio.Open()
	if c.InputIsTerm && c.OutputIsTerm {
		// If we're running interactively, prompt for a
		// signature.
		var valueWidth int
		for f := range incident.ICS309FormDef().AllFields() {
			if f.Tag == "Signature" {
				valueWidth = f.CharWidth()
			}
		}
		c.Confirm("NOTICE: ⇥By providing a signature, you are making a legal assertion that the ICS-309 log is accurate.  Do not provide a signature unless you are sure of that.")
		_, sig, err := c.EditField("Signature", 0, "", valueWidth, nil,
			`This is the signature to be added at the bottom of the ICS-309 form.`, "", false, false, nil)
		return sig, err
	}
	return "", nil
}
