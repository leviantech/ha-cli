package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestGetDomain(t *testing.T) {
	tests := []struct {
		entity   string
		expected string
	}{
		{"light.living_room", "light"},
		{"switch.fan", "switch"},
		{"scene.movie_night", "scene"},
		{"script.goodnight", "script"},
		{"automation.motion_lights", "automation"},
		{"climate.thermostat", "climate"},
		{"sensor", "sensor"}, // No dot, should return the whole string
		{"", ""},
	}

	for _, test := range tests {
		t.Run(test.entity, func(t *testing.T) {
			result := getDomain(test.entity)
			if result != test.expected {
				t.Errorf("getDomain(%q) = %q; want %q", test.entity, result, test.expected)
			}
		})
	}
}

func setupMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case "GET":
			switch r.URL.Path {
			case "/api/states/light.living_room":
				w.Write([]byte(`{"state": "on", "entity_id": "light.living_room"}`))
			case "/api/states":
				w.Write([]byte(`[
					{"entity_id": "light.living_room", "state": "on"},
					{"entity_id": "switch.fan", "state": "off"}
				]`))
			case "/api/":
				w.Write([]byte(`{"message": "API running"}`))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		case "POST":
			w.Write([]byte(`[]`))
		}
	}))
}

func TestCommandExecution(t *testing.T) {
	server := setupMockServer()
	defer server.Close()

	haURL = server.URL
	haToken = "test-token"

	tests := []struct {
		name string
		cmd  *cobra.Command
		args []string
	}{
		{"state", stateCmd, []string{"light.living_room"}},
		{"states", statesCmd, []string{"light.living_room"}},
		{"on", onCmd, []string{"light.living_room"}},
		{"on with brightness", onCmd, []string{"light.living_room", "200"}},
		{"off", offCmd, []string{"light.living_room"}},
		{"toggle", toggleCmd, []string{"light.living_room"}},
		{"scene", sceneCmd, []string{"movie_night"}},
		{"script", scriptCmd, []string{"goodnight"}},
		{"automation", automationCmd, []string{"motion_lights"}},
		{"climate", climateCmd, []string{"climate.thermostat", "22"}},
		{"list all", listCmd, []string{"all"}},
		{"list specific", listCmd, []string{"lights"}},
		{"search", searchCmd, []string{"living"}},
		{"call", callCmd, []string{"light", "turn_on", `{"entity_id":"light.room"}`}},
		{"call no data", callCmd, []string{"light", "turn_on"}},
		{"info", infoCmd, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.RunE(tt.cmd, tt.args)
			if err != nil {
				t.Errorf("Command %s failed: %v", tt.name, err)
			}
		})
	}

	// Test command error branches (e.g., bad API url)
	haURL = "http://\x00" // Invalid URL to trigger API errors
	for _, tt := range tests {
		if tt.name != "info" && tt.name != "list all" && tt.name != "search" {
			// just pick one to test API error
			err := tt.cmd.RunE(tt.cmd, tt.args)
			if err == nil {
				t.Errorf("Expected error for command %s with bad URL", tt.name)
			}
		}
	}
}

func TestConfigCmd(t *testing.T) {
	// Backup and restore stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Provide mock input: url \n token \n
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdin = r

	go func() {
		w.WriteString("http://test-interactive:8123\n")
		w.WriteString("interactive-token\n")
		w.Close()
	}()

	// Change home dir so we don't overwrite real user config
	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome) // for os.UserHomeDir to pick up, though go might not always respect it
	// Actually os.UserHomeDir on mac relies on HOME env var mostly

	err = configCmd.RunE(configCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootPersistentPreRunE(t *testing.T) {
	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)

	os.Unsetenv("HA_URL")
	os.Unsetenv("HA_TOKEN")
	os.Unsetenv("HA_CONFIG")

	// Should fail if missing url/token and not config/help
	err := rootCmd.PersistentPreRunE(stateCmd, []string{})
	if err == nil {
		t.Errorf("expected error when HA_URL/HA_TOKEN are missing")
	}

	// Should succeed if command is config
	err = rootCmd.PersistentPreRunE(configCmd, []string{})
	if err != nil {
		t.Errorf("expected no error for config command: %v", err)
	}

	// Should succeed if HA_URL and HA_TOKEN are set
	os.Setenv("HA_URL", "http://test")
	os.Setenv("HA_TOKEN", "token")
	err = rootCmd.PersistentPreRunE(stateCmd, []string{})
	if err != nil {
		t.Errorf("expected no error when env vars are set: %v", err)
	}
}
