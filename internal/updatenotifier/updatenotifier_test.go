package updatenotifier

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestIsUpdateAvailable(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		latest   string
		expected bool
	}{
		{"newer patch", "3.3.10-go.exp", "v3.3.11", true},
		{"newer minor", "3.3.10-go.exp", "3.4.0", true},
		{"newer major", "3.3.10-go.exp", "v4.0.0-go.1", true},
		{"same version", "3.3.10-go.exp", "v3.3.10", false},
		{"older version", "3.3.10-go.exp", "v3.3.9", false},
		{"empty latest", "3.3.10-go.exp", "", false},
		{"empty current", "", "3.4.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUpdateAvailable(tt.current, tt.latest)
			if got != tt.expected {
				t.Errorf("IsUpdateAvailable(%q, %q) = %v, want %v", tt.current, tt.latest, got, tt.expected)
			}
		})
	}
}

func TestCheckForUpdatesFromServerAndCache(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tag_name": "v3.4.0"}`))
	}))
	defer srv.Close()

	cachePath := filepath.Join(t.TempDir(), "update-cache.json")

	got1 := checkForUpdatesFrom(srv.URL, cachePath)
	if got1 != "v3.4.0" {
		t.Fatalf("expected v3.4.0, got %q", got1)
	}
	if calls != 1 {
		t.Fatalf("expected 1 HTTP call, got %d", calls)
	}

	// Second call should read from cache
	got2 := checkForUpdatesFrom(srv.URL, cachePath)
	if got2 != "v3.4.0" {
		t.Fatalf("expected cached v3.4.0, got %q", got2)
	}
	if calls != 1 {
		t.Fatalf("expected cached result without additional HTTP call, got %d calls", calls)
	}
}

func TestCheckForUpdatesHandlesNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	cachePath := filepath.Join(t.TempDir(), "update-cache-404.json")
	got := checkForUpdatesFrom(srv.URL, cachePath)
	if got != "" {
		t.Fatalf("expected empty string on 404, got %q", got)
	}
}

