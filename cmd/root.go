package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cool",
	Short: "cool",
	Long:  `cool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cool")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
