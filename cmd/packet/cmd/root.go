package cmd

import (
	"errors"
	"log/slog"

	"github.com/rothskeller/packet/cmd/packet/cmd/forms"
	"github.com/rothskeller/packet/cmd/packet/cmd/outpost"
	"github.com/rothskeller/packet/cmd/packet/cmd/server"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "packet",
	Short: "Open a shell for packet commands",
	Long:  `TBD`, // TODO: add help text
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented") // TODO: implement shell
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	defer func() {
		if p := recover(); p != nil {
			slog.Error("PANIC", "p", p)
		}
	}()
	return RootCmd.Execute()
}

func init() {
	cobra.MousetrapHelpText = "" // running without a console is OK
	RootCmd.AddCommand(forms.Command)
	RootCmd.AddCommand(outpost.Command)
	RootCmd.AddCommand(server.Command)
}
