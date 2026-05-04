package cmd

import (
	"testing"
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
