package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ntwcklng/cool/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(deployCmd)
}

type Deployment struct {
	ID              int    `json:"id"`
	ApplicationName string `json:"name"`
	DeploymentUUID  string `json:"uuid"`
	FQDN            string `json:"fqdn"`
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "View and trigger deployments",
	Long:  `List all available deployments and trigger a deployment for the selected application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Fetching deployments...")
		fmt.Println()
		
		viper.SetConfigFile(ConfigFilePath)
		if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
			fmt.Println("⚙️  Configuration not found. Setting up authentication first...")
			fmt.Println()
			authCmd.Run(authCmd, args)
			return
		}

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("❌ Error reading config file: %v\n", err)
			fmt.Println("💡 Try running 'cool auth' to reconfigure")
			return
		}

		apiURL := GetAPIURL()
		token := GetToken()

		if apiURL == "" || token == "" {
			fmt.Println("❌ Missing API URL or token in configuration")
			fmt.Println("💡 Running authentication setup...")
			fmt.Println()
			authCmd.Run(authCmd, args)
			return
		}

		client := &http.Client{}
		req, _ := http.NewRequest("GET", GetAllDeploymentsURL(), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Error fetching deployments: %v\n", err)
			fmt.Println("💡 Check your internet connection and API URL")
			return
		}
		defer resp.Body.Close()

		if !utils.HandleHTTPResponse(resp, "fetching deployments") {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("❌ Error reading response: %v\n", err)
			return
		}

		if len(body) == 0 {
			fmt.Println("⚠️  No deployments found or API returned empty response")
			return
		}

		var deployments []Deployment
		if err := json.Unmarshal(body, &deployments); err != nil {
			fmt.Printf("❌ Error parsing deployments data: %v\n", err)
			return
		}

		if len(deployments) == 0 {
			fmt.Println("📭 No deployments available")
			return
		}

		fmt.Printf("📋 Found %d deployment(s):\n", len(deployments))
		fmt.Println()
		for i, d := range deployments {
			fmt.Printf("  %d) %s\n", i+1, d.ApplicationName)
			fmt.Printf("     🌐 %s\n", d.FQDN)
			if i < len(deployments)-1 {
				fmt.Println()
			}
		}
		fmt.Println()
		var choice int
		fmt.Print("🎯 Select deployment (1-" + fmt.Sprintf("%d", len(deployments)) + "): ")
		_, err = fmt.Scanln(&choice)
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
		fmt.Printf("🚀 Triggering deployment for: %s\n", selectedDeployment.ApplicationName)
		fmt.Printf("🌐 URL: %s\n", selectedDeployment.FQDN)
		fmt.Println()
		
		deploy(GetDeploymentURL(selectedDeployment.DeploymentUUID), token)
	},
}

func deploy(URL, token string) {
	fmt.Println("⏳ Sending deployment request...")
	
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Printf("❌ Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Deployment request failed: %v\n", err)
		fmt.Println("💡 Check your internet connection and try again")
		return
	}
	defer resp.Body.Close()

	if utils.HandleHTTPResponse(resp, "deployment") {
		fmt.Println("✅ Deployment triggered successfully!")
		fmt.Println("💡 Check your Coolify dashboard for deployment progress")
	}
}
