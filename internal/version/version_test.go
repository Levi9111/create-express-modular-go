package version

import (
	"strings"
	"testing"
)

func TestVersionShort(t *testing.T) {
	s := Short()
	if s == "" {
		t.Fatal("expected non-empty Short() version")
	}
	if s != Version {
		t.Fatalf("expected %q, got %q", Version, s)
	}
}

func TestVersionInfo(t *testing.T) {
	info := Info()
	if !strings.Contains(info, "cem version") {
		t.Fatalf("expected info to contain 'cem version', got: %s", info)
	}
	if !strings.Contains(info, Version) {
		t.Fatalf("expected info to contain version %s, got: %s", Version, info)
	}
}
