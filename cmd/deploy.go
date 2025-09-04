package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
	Short: "deploy",
	Long:  `deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigFile(ConfigFilePath)
		if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
			fmt.Println("Config File not found, running auth")
			authCmd.Run(authCmd, args)
		}

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Error reading config file:", err)
			return
		}

		apiURL := GetAPIURL()
		token := GetToken()

		if apiURL == "" || token == "" {
			fmt.Println("API URL or Token not found in config file, running auth")
			authCmd.Run(authCmd, args)
			return

		}

		client := &http.Client{}
		req, _ := http.NewRequest("GET", GetAllDeploymentsURL(), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error fetching deployments:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body:", err)
			return
		}

		if len(body) == 0 {
			fmt.Println("⚠️ API returned empty body")
			return
		}

		var deployments []Deployment
		if err := json.Unmarshal(body, &deployments); err != nil {
			fmt.Println("Error unmarshalling deployments:", err)
			return
		}

		if len(deployments) == 0 {
			fmt.Println("⚠️ No deployments found")
			return
		}
		fmt.Println(string(body))
		if err := json.Unmarshal(body, &deployments); err != nil {
			fmt.Println("Error unmarshalling deployments:", err)
			return
		}

		fmt.Println("Deployments:")
		for i, d := range deployments {
			fmt.Printf("%d) %s - (%s)\n", i+1, d.ApplicationName, d.FQDN)
		}
		var choice int
		fmt.Print("Choose a deployment to deploy: ")
		fmt.Scanln(&choice)
		if choice < 1 || choice > len(deployments) {
			fmt.Println("Invalid choice")
			return
		}
		selectedDeployment := deployments[choice-1]
		fmt.Println("Deploying", selectedDeployment.ApplicationName, "with UUID", GetDeploymentURL(selectedDeployment.DeploymentUUID))
		deploy(GetDeploymentURL(selectedDeployment.DeploymentUUID), token)
	},
}

func deploy(URL, token string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Deploy Error", err)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Deploy Error", err)
	}
	fmt.Println("Deploy Success")
	defer resp.Body.Close()
}
