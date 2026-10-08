package pm

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDetectPMFromLockfiles(t *testing.T) {
	tests := []struct {
		lockfile string
		expected PackageManager
	}{
		{"bun.lockb", Bun},
		{"bun.lock", Bun},
		{"pnpm-lock.yaml", PNPM},
		{"yarn.lock", Yarn},
		{"package-lock.json", NPM},
	}

	for _, tt := range tests {
		t.Run(tt.lockfile, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "pm-test-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			filePath := filepath.Join(tmpDir, tt.lockfile)
			if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
				t.Fatalf("failed to create lockfile: %v", err)
			}

			detected := DetectPM(tmpDir)
			if detected != tt.expected {
				t.Errorf("expected %s for lockfile %s, got %s", tt.expected, tt.lockfile, detected)
			}
		})
	}
}

func TestDetectPMFromUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		ua       string
		expected PackageManager
	}{
		{"bun", "bun/1.0.0 npm/? node/v18.0.0 linux x64", Bun},
		{"pnpm", "pnpm/8.0.0 npm/? node/v18.0.0 linux x64", PNPM},
		{"yarn", "yarn/1.22.19 npm/? node/v18.0.0 linux x64", Yarn},
		{"npm", "npm/9.0.0 node/v18.0.0 linux x64", NPM},
		{"empty", "", NPM},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "pm-ua-test-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			t.Setenv("npm_config_user_agent", tt.ua)
			detected := DetectPM(tmpDir)
			if detected != tt.expected {
				t.Errorf("expected %s for user-agent %q, got %s", tt.expected, tt.ua, detected)
			}
		})
	}
}

func TestInitialInstallCmd(t *testing.T) {
	tests := []struct {
		pm       PackageManager
		expected string
	}{
		{Bun, "bun"},
		{PNPM, "pnpm"},
		{Yarn, "yarn"},
		{NPM, "npm"},
	}

	for _, tt := range tests {
		name, args := InitialInstallCmd(tt.pm)
		if name != tt.expected {
			t.Errorf("expected cmd name %s, got %s", tt.expected, name)
		}
		if len(args) == 0 {
			t.Errorf("expected non-empty args for %s", tt.pm)
		}
	}
}

func TestAddPackagesCmd(t *testing.T) {
	pkgs := []string{"express", "cors"}

	// Production dependencies
	testsProd := []struct {
		pm           PackageManager
		expectedName string
		expectedArgs []string
	}{
		{Bun, "bun", []string{"add", "express", "cors"}},
		{PNPM, "pnpm", []string{"add", "express", "cors"}},
		{Yarn, "yarn", []string{"add", "express", "cors"}},
		{NPM, "npm", []string{"install", "express", "cors"}},
	}

	for _, tt := range testsProd {
		name, args := AddPackagesCmd(tt.pm, pkgs, false)
		if name != tt.expectedName {
			t.Errorf("expected %s, got %s", tt.expectedName, name)
		}
		if !reflect.DeepEqual(args, tt.expectedArgs) {
			t.Errorf("expected args %v, got %v", tt.expectedArgs, args)
		}
	}

	// Dev dependencies
	testsDev := []struct {
		pm           PackageManager
		expectedName string
		expectedArgs []string
	}{
		{Bun, "bun", []string{"add", "-d", "express", "cors"}},
		{PNPM, "pnpm", []string{"add", "-D", "express", "cors"}},
		{Yarn, "yarn", []string{"add", "-D", "express", "cors"}},
		{NPM, "npm", []string{"install", "-D", "express", "cors"}},
	}

	for _, tt := range testsDev {
		name, args := AddPackagesCmd(tt.pm, pkgs, true)
		if name != tt.expectedName {
			t.Errorf("expected %s, got %s", tt.expectedName, name)
		}
		if !reflect.DeepEqual(args, tt.expectedArgs) {
			t.Errorf("expected dev args %v, got %v", tt.expectedArgs, args)
		}
	}
}

func TestRunScriptCmd(t *testing.T) {
	// NPM uses "npm run dev"
	name, args := RunScriptCmd(NPM, "dev")
	if name != "npm" || !reflect.DeepEqual(args, []string{"run", "dev"}) {
		t.Errorf("unexpected NPM script cmd: %s %v", name, args)
	}

	// Bun uses "bun dev"
	name, args = RunScriptCmd(Bun, "dev")
	if name != "bun" || !reflect.DeepEqual(args, []string{"dev"}) {
		t.Errorf("unexpected Bun script cmd: %s %v", name, args)
	}
}

func TestRunSilentCommand(t *testing.T) {
	tmpDir := os.TempDir()
	out, err := RunSilentCommand(tmpDir, "echo", "hello-cem")
	if err != nil {
		t.Fatalf("unexpected error running echo: %v", err)
	}
	if !reflect.DeepEqual(out, "hello-cem\n") {
		t.Errorf("expected 'hello-cem\n', got %q", out)
	}
}
