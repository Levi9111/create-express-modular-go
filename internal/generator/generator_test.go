package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Levi9111/create-express-modular-go/internal/pm"
)

func TestGenerateProject(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cem-gen-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	targetDir := filepath.Join(tempDir, "test-api")

	opts := ProjectOptions{
		ProjectPath:   targetDir,
		ProjectName:   "test-api",
		DB:            DbMongoose,
		Validator:     ValidatorZod,
		UseAuth:       true,
		TokenDelivery: TokenCookie,
		PM:            pm.NPM,
		UseDocker:     true,
		UseSwagger:    true,
		SkipInstall:   true,
	}

	if err := GenerateProject(opts); err != nil {
		t.Fatalf("GenerateProject failed: %v", err)
	}

	// Verify key files exist
	expectedFiles := []string{
		"package.json",
		"tsconfig.json",
		"src/server.ts",
		"src/app.ts",
		"src/app/routes/index.ts",
		"src/app/config/index.ts",
		"src/app/config/swagger.ts",
		"src/app/modules/Auth/auth.controller.ts",
		"src/app/modules/Auth/auth.service.ts",
		"src/app/modules/Auth/auth.route.ts",
		"src/app/modules/Auth/auth.validation.ts",
		"src/app/utils/QueryBuilder.ts",
		"Dockerfile",
		"docker-compose.yml",
		"README.md",
		"AGENTS.md",
		"CLAUDE.md",
		"cem-cli.json",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(targetDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file missing: %s", file)
		}
	}
}

func TestGenerateProjectPrisma(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cem-prisma-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	targetDir := filepath.Join(tempDir, "prisma-api")

	opts := ProjectOptions{
		ProjectPath:   targetDir,
		ProjectName:   "prisma-api",
		DB:            DbPrisma,
		Validator:     ValidatorZod,
		UseAuth:       true,
		TokenDelivery: TokenHeader,
		PM:            pm.Bun,
		UseDocker:     false,
		UseSwagger:    false,
		SkipInstall:   true,
	}

	if err := GenerateProject(opts); err != nil {
		t.Fatalf("GenerateProject failed: %v", err)
	}

	expectedFiles := []string{
		"package.json",
		"prisma/schema.prisma",
		"src/app/utils/prisma.ts",
		"src/app/modules/Auth/auth.controller.ts",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(targetDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected Prisma file missing: %s", file)
		}
	}
}
