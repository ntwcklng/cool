package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ntwcklng/cool/utils"
)

// EnvironmentVariable represents an environment variable from the Coolify API
type EnvironmentVariable struct {
	ID               int    `json:"id"`
	UUID             string `json:"uuid"`
	ResourceableType string `json:"resourceable_type"`
	ResourceableID   int    `json:"resourceable_id"`
	IsLiteral        bool   `json:"is_literal"`
	IsMultiline      bool   `json:"is_multiline"`
	IsPreview        bool   `json:"is_preview"`
	IsRuntime        bool   `json:"is_runtime"`
	IsBuildtime      bool   `json:"is_buildtime"`
	IsShared         bool   `json:"is_shared"`
	IsShownOnce      bool   `json:"is_shown_once"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	RealValue        string `json:"real_value"`
	Version          string `json:"version"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// EnvCreateRequest represents the request body for creating an environment variable
type EnvCreateRequest struct {
	Key          string `json:"key"`
	Value        string `json:"value"`
	IsPreview    bool   `json:"is_preview"`
	IsLiteral    bool   `json:"is_literal"`
	IsMultiline  bool   `json:"is_multiline"`
	IsShownOnce  bool   `json:"is_shown_once"`
}

// EnvBulkUpdateRequest represents the request body for bulk updating environment variables
type EnvBulkUpdateRequest struct {
	Data []EnvCreateRequest `json:"data"`
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment variables for deployments",
	Long:  `Pull and push environment variables between Coolify deployments and local .env files.`,
}

var envPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull environment variables from deployment",
	Long:  `Fetch all environment variables from the linked Coolify deployment and save them to a local .env file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📥 Pulling environment variables from deployment...")

		// Check if cool.yaml exists in current directory
		deploymentUUID := getLocalDeploymentUUID()
		if deploymentUUID == "" {
			fmt.Println("❌ No deployment linked. Run 'cool link' first to link this directory to a deployment.")
			return
		}

		// Determine the filename to use
		envFile := ".env.local"
		fileExists := fileExists(envFile)

		if fileExists {
			fmt.Printf("⚠️  File %s already exists.\n", envFile)
			fmt.Print("Do you want to update existing variables? (y/n): ")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("❌ Operation cancelled.")
				return
			}
		} else {
			// Ask for filename if .env.local doesn't exist
			fmt.Print("Enter filename to save environment variables (default: .env.local): ")
			var input string
			fmt.Scanln(&input)
			if input != "" {
				envFile = input
			}
		}

		// Fetch environment variables from API
		envVars, err := fetchEnvironmentVariables(deploymentUUID)
		if err != nil {
			fmt.Printf("❌ Failed to fetch environment variables: %v\n", err)
			return
		}

		// Read existing env file if it exists
		existingVars := make(map[string]string)
		if fileExists {
			existingVars, err = readEnvFile(envFile)
			if err != nil {
				fmt.Printf("⚠️  Failed to read existing file: %v. Creating new file instead.\n", err)
			}
		}

		// Merge/update variables
		for _, envVar := range envVars {
			existingVars[envVar.Key] = envVar.Value
		}

		// Write to file
		err = writeEnvFile(envFile, existingVars)
		if err != nil {
			fmt.Printf("❌ Failed to write environment variables: %v\n", err)
			return
		}

		fmt.Printf("✅ Successfully pulled %d environment variables to %s\n", len(envVars), envFile)
	},
}

var envPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push environment variables to deployment",
	Long:  `Upload environment variables from a local .env file to the linked Coolify deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📤 Pushing environment variables to deployment...")

		// Check if cool.yaml exists in current directory
		deploymentUUID := getLocalDeploymentUUID()
		if deploymentUUID == "" {
			fmt.Println("❌ No deployment linked. Run 'cool link' first to link this directory to a deployment.")
			return
		}

		// Determine the filename to use
		envFile := ".env.local"
		if !fileExists(envFile) {
			fmt.Printf("⚠️  File %s does not exist.\n", envFile)
			fmt.Print("Enter the filename containing environment variables: ")
			var input string
			fmt.Scanln(&input)
			if input == "" {
				fmt.Println("❌ No filename provided. Operation cancelled.")
				return
			}
			envFile = input
			if !fileExists(envFile) {
				fmt.Printf("❌ File %s does not exist.\n", envFile)
				return
			}
		}

		// Read environment variables from file
		envVars, err := readEnvFile(envFile)
		if err != nil {
			fmt.Printf("❌ Failed to read environment file: %v\n", err)
			return
		}

		if len(envVars) == 0 {
			fmt.Println("⚠️  No environment variables found in file.")
			return
		}

		// Push environment variables using bulk update API
		err = bulkUpdateEnvironmentVariables(deploymentUUID, envVars)
		if err != nil {
			fmt.Printf("❌ Failed to push environment variables: %v\n", err)
			return
		}

		fmt.Printf("✅ Successfully pushed %d environment variables\n", len(envVars))
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envPullCmd)
	envCmd.AddCommand(envPushCmd)
}

// getLocalDeploymentUUID reads the DeploymentUUID from the local cool.yaml
func getLocalDeploymentUUID() string {
	localConfig := viper.New()
	localConfig.SetConfigName("cool")
	localConfig.SetConfigType("yaml")
	localConfig.AddConfigPath(".")

	if err := localConfig.ReadInConfig(); err != nil {
		return ""
	}

	return localConfig.GetString("DeploymentUUID")
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// fetchEnvironmentVariables fetches environment variables from the Coolify API
func fetchEnvironmentVariables(uuid string) ([]EnvironmentVariable, error) {
	apiURL := GetAPIURL()
	token := GetToken()

	url := fmt.Sprintf("%s/api/v1/applications/%s/envs", apiURL, uuid)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !utils.HandleHTTPResponse(resp, "fetch environment variables") {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var envVars []EnvironmentVariable
	err = json.Unmarshal(body, &envVars)
	if err != nil {
		return nil, err
	}

	return envVars, nil
}

// bulkUpdateEnvironmentVariables updates multiple environment variables via the Coolify API
func bulkUpdateEnvironmentVariables(uuid string, envVars map[string]string) error {
	apiURL := GetAPIURL()
	token := GetToken()

	url := fmt.Sprintf("%s/api/v1/applications/%s/envs/bulk", apiURL, uuid)

	// Convert map to array of EnvCreateRequest
	var envData []EnvCreateRequest
	for key, value := range envVars {
		envData = append(envData, EnvCreateRequest{
			Key:         key,
			Value:       value,
			IsPreview:   false,
			IsLiteral:   false,
			IsMultiline: strings.Contains(value, "\n"),
			IsShownOnce: false,
		})
	}

	reqBody := EnvBulkUpdateRequest{
		Data: envData,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !utils.HandleHTTPResponse(resp, "bulk update environment variables") {
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// readEnvFile reads a .env file and returns a map of key-value pairs
func readEnvFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first = sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		envVars[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envVars, nil
}

// writeEnvFile writes environment variables to a .env file
func writeEnvFile(filename string, envVars map[string]string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write header
	fmt.Fprintln(writer, "# Environment variables pulled from Coolify")
	fmt.Fprintln(writer, "# Generated by cool CLI")
	fmt.Fprintln(writer, "")

	// Write variables
	for key, value := range envVars {
		// Quote values that contain spaces or special characters
		if strings.ContainsAny(value, " \t\n\"'$\\") || value == "" {
			value = fmt.Sprintf(`"%s"`, strings.ReplaceAll(value, `"`, `\"`))
		}
		fmt.Fprintf(writer, "%s=%s\n", key, value)
	}

	return writer.Flush()
}
