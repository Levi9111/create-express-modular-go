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
		Version: "1.0.0",
		Project: CemProject{
			Name:              "test-api",
			DB:                "mongoose",
			Validator:         "zod",
			Auth:              true,
			AuthTokenDelivery: "cookie",
			Docker:            true,
			PackageManager:    "npm",
		},
		Structure: CemStructure{
			SrcDir:         "src",
			AppDir:         "src/app",
			ModulesDir:     "src/app/modules",
			MiddlewaresDir: "src/app/middlewares",
			UtilsDir:       "src/app/utils",
			ConfigDir:      "src/app/config",
		},
		Features: CemFeatures{
			Testing:   false,
			Websocket: false,
			Swagger:   true,
		},
		State: &CemState{
			Modules:     []string{},
			Middlewares: []string{},
			Env:         []string{},
		},
	}

	if err := SaveConfig(tempDir, cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	loaded, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loaded.Project.Name != "test-api" {
		t.Errorf("Expected Project.Name 'test-api', got '%s'", loaded.Project.Name)
	}
	if loaded.Project.DB != "mongoose" {
		t.Errorf("Expected Project.DB 'mongoose', got '%s'", loaded.Project.DB)
	}
	if !loaded.Project.Auth {
		t.Errorf("Expected Project.Auth true, got false")
	}

	// Verify file exists
	if _, err := os.Stat(filepath.Join(tempDir, "cem-cli.json")); err != nil {
		t.Errorf("cem-cli.json does not exist: %v", err)
	}
}
