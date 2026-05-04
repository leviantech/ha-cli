package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create a test config file
	configData := []byte(`{
		"url": "http://test-ha:8123",
		"token": "secret-test-token"
	}`)
	err := os.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Ensure environment variables are clear before testing
	os.Unsetenv("HA_URL")
	os.Unsetenv("HA_TOKEN")

	// Call the function
	loadConfig(configPath)

	// Check if env vars were set correctly
	if os.Getenv("HA_URL") != "http://test-ha:8123" {
		t.Errorf("expected HA_URL 'http://test-ha:8123', got '%s'", os.Getenv("HA_URL"))
	}
	if os.Getenv("HA_TOKEN") != "secret-test-token" {
		t.Errorf("expected HA_TOKEN 'secret-test-token', got '%s'", os.Getenv("HA_TOKEN"))
	}

	// Clean up env vars after test
	os.Unsetenv("HA_URL")
	os.Unsetenv("HA_TOKEN")
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
}
