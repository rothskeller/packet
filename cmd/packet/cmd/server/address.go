package server

import (
	"fmt"

	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/spf13/cobra"
)

var addressCmd = &cobra.Command{
	Use:          "address",
	Short:        "Prints the URL of the internal web server",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if addr, err := server.GetAddress(false); err != nil {
			return err
		} else {
			fmt.Println(addr)
		}
		return nil
	},
}

func init() {
	Command.AddCommand(addressCmd)
}
