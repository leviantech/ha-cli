package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Config struct {
	URL      string `json:"url"`
	Token    string `json:"token"`
	Interval int    `json:"interval"`
}

var appConfig Config

var rootCmd = &cobra.Command{
	Use:   "ha-cli",
	Short: "Home Assistant CLI wrapper",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "config" || cmd.Name() == "help" {
			return nil
		}

		home, err := os.UserHomeDir()
		var userConfig string
		if configDir := os.Getenv("HA_CONFIG_DIR"); configDir != "" {
			userConfig = filepath.Join(configDir, "config.json")
		} else if err == nil {
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

		if envURL := os.Getenv("HA_URL"); envURL != "" {
			appConfig.URL = envURL
		}
		if envToken := os.Getenv("HA_TOKEN"); envToken != "" {
			appConfig.Token = envToken
		}

		if appConfig.URL == "" || appConfig.Token == "" {
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
		return
	}
	json.Unmarshal(data, &appConfig)
	if appConfig.Interval <= 0 {
		appConfig.Interval = 300
	}
}
