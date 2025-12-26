package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dnsbuster",
	Short: "Fast DNS enumeration tool written in Go",
	Long:  "dnsbuster is a modern Go rewrite of dnsenum with concurrency and clean architecture.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
