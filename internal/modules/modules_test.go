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
	if _, err := AddEnv(tempDir, []string{"STRIPE_SECRET_KEY"}); err != nil {
		t.Fatalf("AddEnv failed: %v", err)
	}

	envBytes, _ := os.ReadFile(filepath.Join(tempDir, ".env"))
	if !strings.Contains(string(envBytes), "STRIPE_SECRET_KEY=") {
		t.Errorf("STRIPE_SECRET_KEY not found in .env")
	}

	// Remove Env
	if _, err := RemoveEnv(tempDir, []string{"STRIPE_SECRET_KEY"}); err != nil {
		t.Fatalf("RemoveEnv failed: %v", err)
	}

	envBytes, _ = os.ReadFile(filepath.Join(tempDir, ".env"))
	if strings.Contains(string(envBytes), "STRIPE_SECRET_KEY") {
		t.Errorf("STRIPE_SECRET_KEY should have been removed from .env")
	}
}

func TestValidationAndConfigSync(t *testing.T) {
	if err := ValidateProjectName("../escape"); err == nil {
		t.Errorf("expected error for traversal project name")
	}
	if err := ValidateProjectName("my-valid-app"); err != nil {
		t.Errorf("unexpected error for valid project name: %v", err)
	}
	if err := ValidateIdentifier("../../evil"); err == nil {
		t.Errorf("expected error for traversal identifier")
	}
	if err := ValidateIdentifier("OrderItem_1"); err != nil {
		t.Errorf("unexpected error for valid identifier: %v", err)
	}
	if err := ValidateEnvKey("BAD=KEY"); err == nil {
		t.Errorf("expected error for invalid env key")
	}
	if err := ValidateEnvKey("STRIPE_SECRET_KEY"); err != nil {
		t.Errorf("unexpected error for valid env key: %v", err)
	}

	tempDir, err := os.MkdirTemp("", "cem-env-cfg-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_ = os.WriteFile(filepath.Join(tempDir, ".env"), []byte("PORT=5000\n"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, ".env.example"), []byte("PORT=5000\n"), 0644)
	cfgDir := filepath.Join(tempDir, "src/app/config")
	_ = os.MkdirAll(cfgDir, 0755)
	_ = os.WriteFile(filepath.Join(cfgDir, "index.ts"), []byte("export default {\n  port: process.env.PORT || 5000,\n};\n"), 0644)

	if _, err := AddEnv(tempDir, []string{"REDIS_URL"}); err != nil {
		t.Fatalf("AddEnv failed: %v", err)
	}
	cfgBytes, _ := os.ReadFile(filepath.Join(cfgDir, "index.ts"))
	if !strings.Contains(string(cfgBytes), "redis_url: process.env.REDIS_URL,") {
		t.Errorf("expected redis_url in src/app/config/index.ts, got:\n%s", string(cfgBytes))
	}

	if _, err := RemoveEnv(tempDir, []string{"REDIS_URL"}); err != nil {
		t.Fatalf("RemoveEnv failed: %v", err)
	}
	cfgBytes, _ = os.ReadFile(filepath.Join(cfgDir, "index.ts"))
	if strings.Contains(string(cfgBytes), "REDIS_URL") {
		t.Errorf("expected REDIS_URL to be removed from src/app/config/index.ts")
	}
}

func TestAddModuleWithFieldsAndCRUD(t *testing.T) {
	if _, err := ParseFields("bad_field"); err == nil {
		t.Errorf("expected error for field missing type")
	}
	if _, err := ParseFields("price:unknown"); err == nil {
		t.Errorf("expected error for unsupported field type")
	}
	if _, err := ParseFields("role:enum(a'|evil)"); err == nil {
		t.Errorf("expected error for unsafe enum value")
	}
	if camelCase("OrderItem") != "orderItem" {
		t.Errorf("expected camelCase('OrderItem') = 'orderItem', got %q", camelCase("OrderItem"))
	}

	tempDir, err := os.MkdirTemp("", "cem-fields-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_ = os.MkdirAll(filepath.Join(tempDir, "src/app/routes"), 0755)
	_ = os.WriteFile(filepath.Join(tempDir, "src/app/routes/index.ts"), []byte("// --- INJECT IMPORTS HERE ---\nconst moduleRoutes = [\n  // --- INJECT ROUTES HERE ---\n];\n"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "cem-cli.json"), []byte(`{"version":"1.0.0","project":{"name":"t","db":"mongoose","validator":"zod","auth":false,"docker":false,"packageManager":"npm"}}`), 0644)

	added, err := AddModuleWithFields(tempDir, []string{"Catalog"}, "title:string,price:number,inStock:boolean?,category:enum(book|tech),owner:ref(User),tags:string[]?")
	if err != nil {
		t.Fatalf("AddModuleWithFields failed: %v", err)
	}
	if len(added) != 1 || added[0] != "Catalog" {
		t.Fatalf("expected [Catalog], got %v", added)
	}

	modDir := filepath.Join(tempDir, "src/app/modules/Catalog")
	ifaceBytes, _ := os.ReadFile(filepath.Join(modDir, "catalog.interface.ts"))
	if !strings.Contains(string(ifaceBytes), "Types.ObjectId") || !strings.Contains(string(ifaceBytes), "category: 'book' | 'tech';") {
		t.Errorf("unexpected interface content:\n%s", string(ifaceBytes))
	}

	modelBytes, _ := os.ReadFile(filepath.Join(modDir, "catalog.model.ts"))
	if !strings.Contains(string(modelBytes), "enum: ['book', 'tech']") || !strings.Contains(string(modelBytes), "ref: 'User'") {
		t.Errorf("unexpected model content:\n%s", string(modelBytes))
	}

	serviceBytes, _ := os.ReadFile(filepath.Join(modDir, "catalog.service.ts"))
	if !strings.Contains(string(serviceBytes), "new QueryBuilder(CatalogModel.find(), query)") {
		t.Errorf("expected QueryBuilder CRUD in service, got:\n%s", string(serviceBytes))
	}

	valBytes, _ := os.ReadFile(filepath.Join(modDir, "catalog.validation.ts"))
	if !strings.Contains(string(valBytes), "z.enum(['book', 'tech'])") {
		t.Errorf("unexpected validation content:\n%s", string(valBytes))
	}
}
