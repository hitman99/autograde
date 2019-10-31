package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var debug bool

var rootCmd = &cobra.Command{
	Use: "autograde",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "set debug to true or false")
	rootCmd.AddCommand(serveHttpCmd)
}
