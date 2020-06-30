package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cectl",
		Short: "CLI for CloudEvents",
		Long:  "Experimental CLI for CloudEvents",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
