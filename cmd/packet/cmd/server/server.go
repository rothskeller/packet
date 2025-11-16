package server

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "server",
	Short: "Commands to control the internal web server",
}
