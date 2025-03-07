package main

import (
	"testing"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
)

// TestReplaceTemplateValues checks if placeholders are correctly replaced
func TestReplaceTemplateValues(t *testing.T) {
	// Initialize Koanf with test values
	k = koanf.New(".")
	k.Load(confmap.Provider(map[string]interface{}{
		"username": "Alice",
		"email":    "alice@example.com",
		"role":     "Admin",
	}, "."), nil)

	// Define test cases
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello {{ username }}!", "Hello Alice!"},
		{"Your email is {{ email }}.", "Your email is alice@example.com."},
		{"Your role is {{ role }}.", "Your role is Admin."},
		{"Hello {{ username }}, your role is {{ role }}.", "Hello Alice, your role is Admin."},
		{"No replacements here.", "No replacements here."},
		{"Missing key: {{ missing_key }}.", "Missing key: ."},              // Missing key returns empty string
		{"Mixed: {{ username }} - {{ missing_key }}.", "Mixed: Alice - ."}, // One valid, one missing
	}

	// Run test cases
	for _, tc := range tests {
		output := template(tc.input, map[string]string{})
		if output != tc.expected {
			t.Errorf("For input: %q\nExpected: %q\nGot: %q", tc.input, tc.expected, output)
		}
	}
}

// TestGetTemplateKeyParam checks if it fetches the correct param of a key
func TestGetTemplateKeyParam(t *testing.T) {
	// Initialize Koanf with test values
	k = koanf.New(".")
	k.Load(confmap.Provider(map[string]interface{}{
		"i":     "Some text {{res[hello]}}",
		"j":     "Some other text {{dummy}}, this is the one {{i}}",
		"dummy": "dummy",
	}, "."), nil)

	// Define test cases
	tests := []struct {
		input    string
		expected string
	}{
		{"Test: {{ j }}!", "hello"},
	}

	// Run test cases
	for _, tc := range tests {
		_, output := templateGetKeyParams(tc.input, "res")
		if output != tc.expected {
			t.Errorf("For input: %q\nExpected: %q\nGot: %q", tc.input, tc.expected, output)
		}
	}
}
