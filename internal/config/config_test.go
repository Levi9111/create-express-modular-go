package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigSaveAndLoad(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cem-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := &CemConfig{
		Name:           "test-api",
		DB:             "mongoose",
		Validator:      "zod",
		Auth:           true,
		TokenDelivery:  "cookie",
		PackageManager: "npm",
		Docker:         true,
		Swagger:        true,
	}

	if err := SaveConfig(tempDir, cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	loaded, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loaded.Name != "test-api" {
		t.Errorf("Expected Name 'test-api', got '%s'", loaded.Name)
	}
	if loaded.DB != "mongoose" {
		t.Errorf("Expected DB 'mongoose', got '%s'", loaded.DB)
	}
	if !loaded.Auth {
		t.Errorf("Expected Auth true, got false")
	}

	// Verify file exists
	if _, err := os.Stat(filepath.Join(tempDir, "cem-cli.json")); err != nil {
		t.Errorf("cem-cli.json does not exist: %v", err)
	}
}
