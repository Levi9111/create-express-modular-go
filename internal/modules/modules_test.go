package modules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRouteInjectionAndUnwiring(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cem-route-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	routesDir := filepath.Join(tempDir, "src/app/routes")
	_ = os.MkdirAll(routesDir, 0755)

	initialContent := `import { Router } from 'express';
// --- INJECT IMPORTS HERE ---

const router = Router();

const moduleRoutes: { path: string; route: Router }[] = [
  // --- INJECT ROUTES HERE ---
];

moduleRoutes.forEach((route) => router.use(route.path, route.route));

export default router;
`
	routesFile := filepath.Join(routesDir, "index.ts")
	if err := os.WriteFile(routesFile, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to write initial routes: %v", err)
	}

	importLine := "import { ProductRoutes } from '../modules/Product/product.route';"
	routeLine := "  { path: '/products', route: ProductRoutes },"

	// Test Injection
	if err := InjectRoute(tempDir, importLine, routeLine); err != nil {
		t.Fatalf("InjectRoute failed: %v", err)
	}

	contentBytes, _ := os.ReadFile(routesFile)
	content := string(contentBytes)
	if !strings.Contains(content, importLine) {
		t.Errorf("Import line was not injected properly")
	}
	if !strings.Contains(content, routeLine) {
		t.Errorf("Route line was not injected properly")
	}

	// Test Unwiring
	if err := RemoveRoute(tempDir, "Product"); err != nil {
		t.Fatalf("RemoveRoute failed: %v", err)
	}

	contentBytes, _ = os.ReadFile(routesFile)
	content = string(contentBytes)
	if strings.Contains(content, importLine) {
		t.Errorf("Import line was not removed")
	}
	if strings.Contains(content, routeLine) {
		t.Errorf("Route line was not removed")
	}
}

func TestEnvAddAndRemove(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cem-env-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configDir := filepath.Join(tempDir, "src/app/config")
	_ = os.MkdirAll(configDir, 0755)

	// Mock .env
	_ = os.WriteFile(filepath.Join(tempDir, ".env"), []byte("PORT=5000\nNODE_ENV=development\n"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, ".env.example"), []byte("PORT=5000\nNODE_ENV=development\n"), 0644)

	// Mock src/app/config/index.ts
	configContent := `import dotenv from 'dotenv';
import path from 'path';

dotenv.config({ path: path.join(process.cwd(), '.env') });

export default {
  NODE_ENV: process.env.NODE_ENV || 'development',
  port: process.env.PORT || 5000,
};
`
	_ = os.WriteFile(filepath.Join(configDir, "index.ts"), []byte(configContent), 0644)

	// Add Env
	if err := AddEnv(tempDir, []string{"STRIPE_SECRET_KEY"}); err != nil {
		t.Fatalf("AddEnv failed: %v", err)
	}

	envBytes, _ := os.ReadFile(filepath.Join(tempDir, ".env"))
	if !strings.Contains(string(envBytes), "STRIPE_SECRET_KEY=") {
		t.Errorf("STRIPE_SECRET_KEY not found in .env")
	}

	// Remove Env
	if err := RemoveEnv(tempDir, []string{"STRIPE_SECRET_KEY"}); err != nil {
		t.Fatalf("RemoveEnv failed: %v", err)
	}

	envBytes, _ = os.ReadFile(filepath.Join(tempDir, ".env"))
	if strings.Contains(string(envBytes), "STRIPE_SECRET_KEY") {
		t.Errorf("STRIPE_SECRET_KEY should have been removed from .env")
	}
}
