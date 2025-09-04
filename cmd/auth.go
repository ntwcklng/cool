package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(authCmd)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Set up authentication credentials",
	Long:  `Configure your Coolify API URL and authentication token for the CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔐 Setting up authentication...")
		fmt.Println()
		initConfig()
	},
}

func initConfig() {
	viper.SetConfigFile(ConfigFilePath)

	viper.SetConfigType("yaml")

	if _, err := os.Stat(ConfigFilePath); err == nil {
		err := viper.ReadInConfig()
		if err == nil {
			var answer string
			fmt.Printf("✅ Found existing configuration at: %s\n", viper.ConfigFileUsed())
			fmt.Printf("   API URL: %s\n", GetAPIURL())
			fmt.Printf("   Token: %s\n", maskToken(GetToken()))
			fmt.Println()
			fmt.Print("Use existing configuration? (y/n): ")
			fmt.Scanln(&answer)

			if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
				fmt.Println("✅ Using existing configuration")
				return
			}
			fmt.Println("📝 Updating configuration...")
			fmt.Println()
		}
	}
	var apiURL, token string
	fmt.Print("🌐 Enter your Coolify URL (e.g., https://coolify.example.com): ")
	fmt.Scanln(&apiURL)
	fmt.Print("🔑 Enter your API token: ")
	fmt.Scanln(&token)
	
	fmt.Println()
	fmt.Println("💾 Saving configuration...")
	
	viper.Set("apiURL", apiURL)
	viper.Set("token", token)
	if err := viper.WriteConfigAs(ConfigFilePath); err != nil {
		fmt.Printf("❌ Error writing config file: %v\n", err)
		return
	}
	fmt.Printf("✅ Configuration saved to: %s\n", ConfigFilePath)
	fmt.Printf("   API URL: %s\n", GetAPIURL())
	fmt.Printf("   Token: %s\n", maskToken(token))
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return strings.Repeat("*", len(token))
	}
	return token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
}
