package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoAPIRequest(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check headers
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Handle specific endpoints
		switch r.URL.Path {
		case "/api/states/light.living_room":
			if r.Method == "GET" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"state": "on", "entity_id": "light.living_room"}`))
			}
		case "/api/services/light/turn_on":
			if r.Method == "POST" {
				var payload map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if payload["entity_id"] == "light.living_room" {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`[{"entity_id": "light.living_room", "state": "on"}]`))
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Set global variables for testing
	haURL = server.URL
	haToken = "test-token"

	t.Run("GET Request", func(t *testing.T) {
		resp, err := doAPIRequest("GET", "/api/states/light.living_room", nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		var data map[string]string
		if err := json.Unmarshal(resp, &data); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}
		if data["state"] != "on" {
			t.Errorf("expected state 'on', got '%s'", data["state"])
		}
	})

	t.Run("POST Request", func(t *testing.T) {
		payload := map[string]string{"entity_id": "light.living_room"}
		resp, err := doAPIRequest("POST", "/api/services/light/turn_on", payload)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		var data []map[string]string
		if err := json.Unmarshal(resp, &data); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}
		if len(data) == 0 || data[0]["state"] != "on" {
			t.Errorf("expected state 'on', got response: %v", data)
		}
	})

	t.Run("Unauthorized Request", func(t *testing.T) {
		haToken = "invalid-token"
		_, err := doAPIRequest("GET", "/api/states/light.living_room", nil)
		if err == nil {
			t.Fatal("expected error due to unauthorized token, got none")
		}
		haToken = "test-token" // Restore token
	})
}
