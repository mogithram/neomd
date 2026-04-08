package ui

import (
	"strings"
	"testing"
)

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user@example.com", "u***@example.com"},
		{"Name <user@example.com>", "Name <u***@example.com>"},
		{"a@b.com", "a***@b.com"},
		{"", ""},
		{"no-at-sign", "no-at-sign"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := maskEmail(tt.input)
			if got != tt.want {
				t.Errorf("maskEmail(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// isURLSchemeAllowed replicates the inline URL scheme check from model.go Update().
func isURLSchemeAllowed(url string) bool {
	lower := strings.ToLower(url)
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

func TestURLSchemeValidation(t *testing.T) {
	tests := []struct {
		url     string
		allowed bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"HTTP://EXAMPLE.COM", true},
		{"https://secure.example.com/path?q=1", true},
		{"javascript:alert(1)", false},
		{"ftp://files.example.com", false},
		{"data:text/html,<h1>hi</h1>", false},
		{"", false},
		{"file:///etc/passwd", false},
		{"mailto:user@example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := isURLSchemeAllowed(tt.url)
			if got != tt.allowed {
				t.Errorf("isURLSchemeAllowed(%q) = %v, want %v", tt.url, got, tt.allowed)
			}
		})
	}
}
