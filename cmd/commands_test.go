package cmd

import (
	"net/http"
	"net/http/httptest"
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
}
