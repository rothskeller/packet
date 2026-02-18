package forms

import (
	"io"
	"os"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/form/formdefs"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:          "install filename|URL",
	Short:        "Installs a forms bundle",
	Long:         `Installs the specified forms bundle, which must be either a URL starting with https:// or a local filename.  This will replace any existing bundle with the same name, even if it is newer than the one being installed.  Any README text associated with the forms bundle is displayed to standard output.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var bundle, readme string

		if bundle, readme, err = formdefs.InstallBundle(args[0]); err != nil {
			return err
		}
		if readme != "" {
			io.WriteString(os.Stdout, readme)
		} else {
			cio.Open().Confirm("Installed forms bundle %q.\n", bundle)
		}
		return nil
	},
}

func init() {
	Command.AddCommand(installCmd)
}
