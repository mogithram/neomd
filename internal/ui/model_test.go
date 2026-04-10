package ui

import (
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sspaeti/neomd/internal/config"
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

func TestMergeAutoBCC(t *testing.T) {
	tests := []struct {
		name    string
		bcc     string
		autoBCC string
		want    string
	}{
		{
			name:    "append when empty",
			bcc:     "",
			autoBCC: "archive@example.com",
			want:    "archive@example.com",
		},
		{
			name:    "append when distinct",
			bcc:     "team@example.com",
			autoBCC: "archive@example.com",
			want:    "team@example.com, archive@example.com",
		},
		{
			name:    "dedupe bare and named address",
			bcc:     "Archive <archive@example.com>",
			autoBCC: "archive@example.com",
			want:    "Archive <archive@example.com>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeAutoBCC(tt.bcc, tt.autoBCC); got != tt.want {
				t.Fatalf("mergeAutoBCC(%q, %q) = %q, want %q", tt.bcc, tt.autoBCC, got, tt.want)
			}
		})
	}
}

func TestCollectRcptTo(t *testing.T) {
	got := collectRcptTo(
		"Alice <alice@example.com>, bob@example.com",
		"bob@example.com, Carol <carol@example.com>",
		"alice@example.com, dave@example.com",
	)
	want := []string{"alice@example.com", "bob@example.com", "carol@example.com", "dave@example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("collectRcptTo() = %#v, want %#v", got, want)
	}
}

func TestPresendSMTPAccount(t *testing.T) {
	cfg := &config.Config{
		Accounts: []config.AccountConfig{
			{Name: "Personal", From: "me@example.com"},
			{Name: "Work", From: "me@work.example"},
		},
		Senders: []config.SenderConfig{
			{Name: "Support", From: "support@example.com", Account: "Work"},
		},
	}
	m := Model{
		cfg:      cfg,
		accounts: cfg.ActiveAccounts(),
		accountI: 0,
	}

	t.Run("selected account uses its own SMTP account", func(t *testing.T) {
		m.presendFromI = 1
		if got := m.presendSMTPAccount().Name; got != "Work" {
			t.Fatalf("presendSMTPAccount() = %q, want %q", got, "Work")
		}
	})

	t.Run("sender alias resolves to referenced account", func(t *testing.T) {
		m.presendFromI = 2
		if got := m.presendSMTPAccount().Name; got != "Work" {
			t.Fatalf("presendSMTPAccount() = %q, want %q", got, "Work")
		}
	})
}

func TestActiveFolderUsesOffTabFolder(t *testing.T) {
	m := Model{
		cfg: &config.Config{
			Folders: config.FoldersConfig{
				Inbox:  "INBOX",
				Drafts: "Drafts",
				Spam:   "Spam",
			},
		},
		folders:       []string{"Inbox"},
		activeFolderI: 0,
	}

	m.offTabFolder = "Drafts"
	if got := m.activeFolder(); got != "Drafts" {
		t.Fatalf("activeFolder() with Drafts off-tab = %q, want %q", got, "Drafts")
	}

	m.offTabFolder = "Spam"
	if got := m.activeFolder(); got != "Spam" {
		t.Fatalf("activeFolder() with Spam off-tab = %q, want %q", got, "Spam")
	}
}

func TestUpdateInboxEscClearsCommittedFilter(t *testing.T) {
	m := Model{
		filterText: "invoice",
		inbox:      newInboxList(80, 20, "", ""),
		folders:    []string{"Inbox"},
		cfg: &config.Config{
			Folders: config.FoldersConfig{Inbox: "INBOX"},
		},
	}

	next, _ := m.updateInbox(tea.KeyMsg{Type: tea.KeyEsc})
	got := next.(Model)
	if got.filterText != "" {
		t.Fatalf("filterText = %q, want empty", got.filterText)
	}
	if got.filterActive {
		t.Fatal("filterActive should be false after esc")
	}
}
