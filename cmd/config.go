package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ntwcklng/cool/pkg/types"
	"github.com/ntwcklng/cool/utils"
	"github.com/spf13/viper"
)

var (
	home, _        = os.UserHomeDir()
	ConfigFileName = ".cool.yaml"
	ConfigFilePath = filepath.Join(home, ConfigFileName)
)

func init() {
	viper.SetConfigFile(ConfigFilePath)
	viper.SetConfigType("yaml")

	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		fmt.Println("⚙️  Configuration not found. Setting up authentication first...")
		fmt.Println()
		authCmd.Run(authCmd, []string{})
		return
	}
	viper.SetConfigFile(ConfigFilePath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("❌ Error reading config file: %v\n", err)
		fmt.Println("💡 Try running 'cool auth' to reconfigure")
		return
	}
}

func ListAllApplications() []types.Deployment {
	// Diese Funktion wird in deploy.go definiert

	fmt.Println("🚀 Fetching deployments...")
	fmt.Println()

	apiURL := GetAPIURL()
	token := GetToken()

	if apiURL == "" || token == "" {
		fmt.Println("❌ Missing API URL or token in configuration")
		fmt.Println("💡 Running authentication setup...")
		fmt.Println()
		authCmd.Run(authCmd, []string{})
		return []types.Deployment{}
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", GetAllDeploymentsURL(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error fetching deployments: %v\n", err)
		fmt.Println("💡 Check your internet connection and API URL")
		return []types.Deployment{}
	}
	defer resp.Body.Close()

	if !utils.HandleHTTPResponse(resp, "fetching deployments") {
		return []types.Deployment{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Error reading response: %v\n", err)
		return []types.Deployment{}
	}

	if len(body) == 0 {
		fmt.Println("⚠️  No deployments found or API returned empty response")
		return []types.Deployment{}
	}

	var deployments []types.Deployment
	if err := json.Unmarshal(body, &deployments); err != nil {
		fmt.Printf("❌ Error parsing deployments data: %v\n", err)
		return []types.Deployment{}
	}

	if len(deployments) == 0 {
		fmt.Println("📭 No deployments available")
		return []types.Deployment{}
	}

	return deployments
}

func Deploy(URL string) {
	token := GetToken()
	if token == "" {
		fmt.Println("❌ No valid token found. Please authenticate first.")
		fmt.Println("💡 Try running 'cool auth' to reconfigure")
		return
	}
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

func GetAPIURL() string {
	raw := strings.TrimSpace(viper.GetString("apiurl"))
	if raw == "" {
		return ""
	}

	// Falls Scheme fehlt, http als Default anfügen
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "https://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return raw // fallback, falls ungültig
	}

	// Nur Scheme + Host zurückgeben, kein Pfad, keine Query, kein Fragment
	return u.Scheme + "://" + u.Host
}

func GetToken() string {
	return viper.GetString("token")
}

func GetAllDeploymentsURL() string {
	return GetAPIURL() + "/api/v1/applications"
}

func GetDeploymentURL(uuid string) string {
	return GetAPIURL() + "/api/v1/deploy?uuid=" + uuid
}
