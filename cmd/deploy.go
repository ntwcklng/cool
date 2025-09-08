package cmd

import (
	"fmt"
	"os"

	"github.com/ntwcklng/cool/pkg/types"
	"github.com/ntwcklng/cool/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(deployCmd)
}

var (
	selectedDeployment types.Deployment
	deployCmd          = &cobra.Command{
		Use:   "deploy",
		Short: "View and trigger deployments",
		Long:  `List all available deployments and trigger a deployment for the selected application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println()
			v := viper.New()
			v.SetConfigName("cool")
			v.SetConfigType("yaml")
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("❌ Could not get current directory: %v\n", err)
				return
			}
			v.AddConfigPath(pwd)
			if err := v.ReadInConfig(); err == nil {
				if v.GetString("DeploymentUUID") != "" {
					fmt.Printf("🕵🏻‍♂️ A cool.yaml found with DeploymentUUID: %s\n", v.GetString("DeploymentUUID"))
					selectedDeployment.DeploymentUUID = v.GetString("DeploymentUUID")
					selectedDeployment.ApplicationName = v.GetString("ApplicationName")
					selectedDeployment.FQDN = v.GetString("FQDN")
				}
			} else {
				fmt.Println("💡 No cool.yaml found in current directory. Listing all deployments...")
				deployments := ListAllApplications()
				selectedDeployment = utils.Select(deployments, "Select an application to deploy:")
				if (selectedDeployment == types.Deployment{}) {
					fmt.Println("❌ No deployment selected. Exiting.")
					return
				}

				fmt.Printf("✅ Selected deployment: %s (UUID: %s)\n", selectedDeployment.ApplicationName, selectedDeployment.DeploymentUUID)
			}
			fmt.Println()
			fmt.Printf("🚀 Triggering deployment for: %s\n", selectedDeployment.ApplicationName)
			fmt.Printf("🌐 URL: %s\n", selectedDeployment.FQDN)
			fmt.Println()

			Deploy(GetDeploymentURL(selectedDeployment.DeploymentUUID))
		},
	}
)
