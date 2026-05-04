package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	haURL   string
	haToken string
)

var rootCmd = &cobra.Command{
	Use:   "ha-cli",
	Short: "Home Assistant CLI wrapper",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "config" || cmd.Name() == "help" {
			return nil
		}

		home, err := os.UserHomeDir()
		var userConfig string
		if err == nil {
			userConfig = filepath.Join(home, ".ha-cli", "config.json")
		}

		if userConfig != "" && hasConfig(userConfig) {
			loadConfig(userConfig)
		} else {
			configFile := os.Getenv("HA_CONFIG")
			if configFile == "" && home != "" {
				configFile = filepath.Join(home, ".config", "home-assistant", "config.json")
			}

			if configFile != "" && hasConfig(configFile) {
				loadConfig(configFile)
			}
		}

		haURL = os.Getenv("HA_URL")
		haToken = os.Getenv("HA_TOKEN")

		if haURL == "" || haToken == "" {
			return fmt.Errorf("set HA_URL and HA_TOKEN environment variables, or run 'ha-cli config'")
		}

		return nil
	},
}

func hasConfig(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return // Ignore if config file doesn't exist or can't be read
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return
	}

	if url, ok := config["url"].(string); ok && os.Getenv("HA_URL") == "" {
		os.Setenv("HA_URL", url)
	}
	if token, ok := config["token"].(string); ok && os.Getenv("HA_TOKEN") == "" {
		os.Setenv("HA_TOKEN", token)
	}
}
