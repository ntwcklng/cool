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
	Short: "auth",
	Long:  `auth`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("initConfig")
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
			fmt.Println("Found Config File:", viper.ConfigFileUsed())
			fmt.Println("Use this Config (y) or change the API-Url + Token (n)", ConfigFilePath)
			fmt.Scanln(&answer)

			if strings.ToLower(answer) == "y" {
				fmt.Println("Using Config File:", viper.ConfigFileUsed())
				return
			}
		}
	}
	var apiURL, token string
	fmt.Print("Enter your coolify URL: ")
	fmt.Scanln(&apiURL)
	fmt.Print("Token: ")
	fmt.Scanln(&token)
	viper.Set("apiURL", apiURL)
	viper.Set("token", token)
	if err := viper.WriteConfigAs(ConfigFilePath); err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}
	fmt.Println("Created config file:", viper.ConfigFileUsed())
}
