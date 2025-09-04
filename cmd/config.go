package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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

	viper.SetDefault("apiurl", "")
	viper.SetDefault("token", "")

	if _, err := os.Stat(ConfigFilePath); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("⚠️  Warning: Could not read config file: %v\n", err)
		}
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
