package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput captures stdout during function execution.
func captureOutput(f func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = orig

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestColorFunctions(t *testing.T) {
	if Cyan("test") == "" {
		t.Fatal("expected non-empty Cyan string")
	}
	if Green("test") == "" {
		t.Fatal("expected non-empty Green string")
	}
	if Red("test") == "" {
		t.Fatal("expected non-empty Red string")
	}
}

func TestPrintBanner(t *testing.T) {
	out := captureOutput(func() {
		PrintBanner("3.3.10-test")
	})
	if !strings.Contains(out, "CEM") || !strings.Contains(out, "3.3.10-test") {
		t.Fatalf("expected banner with CEM and version, got: %s", out)
	}
}

func TestSectionHeaderAndBullets(t *testing.T) {
	out := captureOutput(func() {
		SectionHeader("Testing Section")
		Bullet("Key", "Value")
		Substep("Substep item")
		Warn("Warning message")
		Success("Success message")
		Err("Error message")
	})

	if !strings.Contains(out, "Testing Section") {
		t.Error("expected section title in output")
	}
	if !strings.Contains(out, "Key") || !strings.Contains(out, "Value") {
		t.Error("expected bullet key/val in output")
	}
}

func TestNewSpinner(t *testing.T) {
	s := NewSpinner("Loading test...")
	if s == nil {
		t.Fatal("expected non-nil spinner")
	}
	if !strings.Contains(s.Suffix, "Loading test...") {
		t.Errorf("expected suffix to contain loading message, got: %s", s.Suffix)
	}
}

func TestPrintSummary(t *testing.T) {
	out := captureOutput(func() {
		PrintSummary("my-test-app", "mongoose", "zod", true, true, true)
	})
	if !strings.Contains(out, "my-test-app") || !strings.Contains(out, "mongoose") {
		t.Fatalf("expected summary to contain project name and db, got: %s", out)
	}
}
