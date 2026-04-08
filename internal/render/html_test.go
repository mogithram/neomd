package render

import (
	"strings"
	"testing"
)

func TestToHTML_Bold(t *testing.T) {
	out, err := ToHTML("**bold** text")
	if err != nil {
		t.Fatalf("ToHTML returned error: %v", err)
	}
	if !strings.Contains(out, "<strong>bold</strong>") {
		t.Errorf("expected <strong>bold</strong> in output, got:\n%s", out)
	}
}

func TestToHTML_HTMLWrapper(t *testing.T) {
	out, err := ToHTML("hello")
	if err != nil {
		t.Fatalf("ToHTML returned error: %v", err)
	}
	if !strings.HasPrefix(out, "<!DOCTYPE html>") {
		t.Errorf("expected output to start with <!DOCTYPE html>, got:\n%.80s...", out)
	}
	if !strings.Contains(out, "<body>") {
		t.Errorf("expected <body> in output")
	}
}

func TestToHTML_GFMTable(t *testing.T) {
	md := "| A | B |\n|---|---|\n| 1 | 2 |\n"
	out, err := ToHTML(md)
	if err != nil {
		t.Fatalf("ToHTML returned error: %v", err)
	}
	if !strings.Contains(out, "<table>") {
		t.Errorf("expected <table> in output, got:\n%s", out)
	}
}

func TestToHTML_CodeBlock(t *testing.T) {
	md := "```go\nfmt.Println(\"hi\")\n```\n"
	out, err := ToHTML(md)
	if err != nil {
		t.Fatalf("ToHTML returned error: %v", err)
	}
	if !strings.Contains(out, "<pre>") {
		t.Errorf("expected <pre> in output, got:\n%s", out)
	}
}

func TestToHTML_Empty(t *testing.T) {
	out, err := ToHTML("")
	if err != nil {
		t.Fatalf("ToHTML returned error for empty input: %v", err)
	}
	if !strings.HasPrefix(out, "<!DOCTYPE html>") {
		t.Errorf("expected DOCTYPE even for empty input, got:\n%.80s...", out)
	}
}

func TestToANSI_Smoke(t *testing.T) {
	_, err := ToANSI("# Hello\n\nSome **bold** text.", "dark", 80)
	if err != nil {
		t.Fatalf("ToANSI returned error: %v", err)
	}
}
