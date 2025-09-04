package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cool",
	Short: "A CLI tool for managing Coolify deployments",
	Long: `Cool is a command-line interface for managing your Coolify deployments.
It allows you to authenticate, view deployments, and trigger deployments
directly from your terminal.

Available commands:
  auth   - Set up authentication credentials
  deploy - View and trigger deployments`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Cool - Coolify CLI")
		fmt.Println()
		fmt.Println("Welcome to the Coolify command-line interface!")
		fmt.Println("Use 'cool --help' to see available commands.")
		fmt.Println()
		fmt.Println("Quick start:")
		fmt.Println("  cool auth   - Set up your API credentials")
		fmt.Println("  cool deploy - View and manage deployments")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
