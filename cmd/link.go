package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(linkCmd)
}

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link a local project to a deployment",
	Long:  `Link your local project directory to a specific deployment for easier management.`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("❌ Could not get current directory: %v\n", err)
			return
		}
		fmt.Printf("🔗 Linking current directory (%s) to deployment...\n", pwd)
		v := viper.New()
		v.SetConfigName("cool")
		v.SetConfigType("yaml")
		v.AddConfigPath(pwd)
		if err := v.ReadInConfig(); err == nil {
			if v.GetString("deploymentUUID") == "" {
				fmt.Println("❌ cool.yaml found but no deploymentUUID set.")
				fmt.Println("💡 Please run 'cool link' in a directory without an existing cool.yaml to create a new one.")
				return
			}
			fmt.Printf("🕵🏻‍♂️ A cool.yaml found with DeploymentUUID: %s\n", v.GetString("deploymentUUID"))
		} else {
			deployments := ListAllApplications()
			var choice int
			fmt.Print("🎯 Select deployment (1-" + fmt.Sprintf("%d", len(deployments)) + "): ")
			_, err := fmt.Scanln(&choice)
			if err != nil {
				fmt.Printf("❌ Invalid input: %v\n", err)
				return
			}

			if choice < 1 || choice > len(deployments) {
				fmt.Printf("❌ Invalid choice. Please select a number between 1 and %d\n", len(deployments))
				return
			}

			selectedDeployment := deployments[choice-1]
			fmt.Println()
			v.Set("DeploymentUUID", selectedDeployment.DeploymentUUID)
			v.Set("ApplicationName", selectedDeployment.ApplicationName)
			v.Set("FQDN", selectedDeployment.FQDN)
			v.WriteConfigAs(pwd + "/cool.yaml")
			fmt.Printf("Created a new cool.yaml\n")
			fmt.Println("✅ Successfully linked current directory to deployment.")
		}
	},
}
