package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update cool CLI to the latest version",
	Long: `Check for and install the latest version of cool CLI from GitHub.
This will download, build, and replace the current binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Checking for latest version...")

		resp, err := http.Get("https://api.github.com/repos/ntwcklng/cool/releases/latest")
		if err != nil {
			fmt.Printf("❌ Error checking GitHub: %v\n", err)
			return
		}
		defer resp.Body.Close()

		var release GitHubRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			fmt.Printf("❌ Error decoding response: %v\n", err)
			return
		}

		currentVersion := "v1.0.2"
		if release.TagName == currentVersion {
			fmt.Printf("✅ You are already on the latest version: %s\n", currentVersion)
			return
		}

		fmt.Printf("🔄 Updating from %s to %s\n", currentVersion, release.TagName)

		tmpDir := filepath.Join(os.TempDir(), "cool_update")
		os.RemoveAll(tmpDir)

		fmt.Println("📥 Cloning latest source...")
		if err := exec.Command("git", "clone", "https://github.com/ntwcklng/cool.git", tmpDir).Run(); err != nil {
			fmt.Printf("❌ Error cloning repo: %v\n", err)
			return
		}

		fmt.Println("🔨 Building new binary...")
		binaryPath := filepath.Join(tmpDir, "cool")
		buildCmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", binaryPath, "./cmd")
		buildCmd.Dir = tmpDir
		if err := buildCmd.Run(); err != nil {
			fmt.Printf("❌ Error building: %v\n", err)
			return
		}

		execPath, err := os.Executable()
		if err != nil {
			fmt.Printf("❌ Error getting executable path: %v\n", err)
			return
		}

		fmt.Printf("📦 Replacing current binary at %s\n", execPath)
		if err := os.Rename(binaryPath, execPath); err != nil {
			fmt.Printf("❌ Error replacing binary: %v\n", err)
			return
		}

		fmt.Println("🎉 Update completed successfully!")
		os.RemoveAll(tmpDir)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
