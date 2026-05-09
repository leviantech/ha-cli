package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecute(t *testing.T) {
	// Backup original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Provide a safe command for Execute so it doesn't call os.Exit
	os.Args = []string{"ha-cli", "help"}
	Execute()
}

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	loadConfig(filepath.Join(tempDir, "does-not-exist.json"))

	err := os.WriteFile(configPath, []byte(`{bad-json`), 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}
	loadConfig(configPath)

	configData := []byte(`{
		"url": "http://test-ha:8123",
		"token": "secret-test-token",
		"interval": 120
	}`)
	err = os.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	appConfig = Config{}
	loadConfig(configPath)

	if appConfig.URL != "http://test-ha:8123" {
		t.Errorf("expected URL 'http://test-ha:8123', got '%s'", appConfig.URL)
	}
	if appConfig.Token != "secret-test-token" {
		t.Errorf("expected Token 'secret-test-token', got '%s'", appConfig.Token)
	}
	if appConfig.Interval != 120 {
		t.Errorf("expected Interval 120, got %d", appConfig.Interval)
	}
}

func TestHasConfig(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// File doesn't exist yet
	if hasConfig(configPath) {
		t.Errorf("expected hasConfig to return false for non-existent file")
	}

	// Create file
	err := os.WriteFile(configPath, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// File exists now
	if !hasConfig(configPath) {
		t.Errorf("expected hasConfig to return true for existing file")
	}
	
	// Test directory
	if hasConfig(tempDir) {
		t.Errorf("expected hasConfig to return false for directory")
	}
}
