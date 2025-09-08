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
			if v.GetString("DeploymentUUID") == "" {
				fmt.Println("❌ cool.yaml found but no DeploymentUUID set.")
				fmt.Println("💡 Please run 'cool link' in a directory without an existing cool.yaml to create a new one.")
				return
			}
			fmt.Printf("🕵🏻‍♂️ A cool.yaml found with DeploymentUUID: %s\n", v.GetString("DeploymentUUID"))
		} else {
			deployments := ListAllApplications()
			selectedDeployment = utils.Select(deployments, "Select deployment to link to:")
			if (selectedDeployment == types.Deployment{}) {
				fmt.Println("❌ No deployment selected. Exiting.")
				return
			}

			fmt.Printf("✅ Selected deployment: %s (UUID: %s)\n", selectedDeployment.ApplicationName, selectedDeployment.DeploymentUUID)
			fmt.Println()
			v.Set("DeploymentUUID", selectedDeployment.DeploymentUUID)
			v.Set("ApplicationName", selectedDeployment.ApplicationName)
			v.Set("FQDN", selectedDeployment.FQDN)

			configPath := pwd + "/cool.yaml"
			if err := v.WriteConfigAs(configPath); err != nil {
				fmt.Printf("❌ Error creating cool.yaml: %v\n", err)
				return
			}

			fmt.Printf("✅ Created cool.yaml with deployment: %s\n", selectedDeployment.ApplicationName)
			fmt.Printf("🌐 FQDN: %s\n", selectedDeployment.FQDN)
			fmt.Println("✅ Successfully linked current directory to deployment.")
		}
	},
}
