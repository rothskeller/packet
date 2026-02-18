package outpost

import (
	"fmt"
	"os"

	"github.com/rothskeller/packet/form/formdefs"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "outpost",
	Short: "Commands used by Outpost",
	Long:  `These are commands invoked by Outpost to work with non-plain-text messages. These commands are not normally used by humans directly.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := formdefs.RegisterForms(); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: %s\n", err)
		}
	},
}
