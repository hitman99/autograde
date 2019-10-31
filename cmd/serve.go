package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var serveHttpCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {

		os.Exit(0)
	},
}
