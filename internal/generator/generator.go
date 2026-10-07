package generator

import (
	"fmt"
	"os"

	"github.com/Levi9111/create-express-modular-go/internal/config"
	"github.com/Levi9111/create-express-modular-go/internal/ui"
)

// GenerateProject runs the complete scaffolding pipeline for a new project.
func GenerateProject(opts ProjectOptions) error {
	// 1. Target directory creation
	if err := os.MkdirAll(opts.ProjectPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", opts.ProjectName, err)
	}

	ui.SectionHeader("Scaffolding")
	ui.Bullet("Project", opts.ProjectName)
	ui.Bullet("Database", string(opts.DB))
	ui.Bullet("Validator", string(opts.Validator))
	if opts.UseAuth {
		ui.Bullet("Auth", "yes")
	} else {
		ui.Bullet("Auth", "no")
	}
	if opts.UseDocker {
		ui.Bullet("Docker", "yes")
	} else {
		ui.Bullet("Docker", "no")
	}
	if opts.UseSwagger {
		ui.Bullet("Swagger", "yes")
	} else {
		ui.Bullet("Swagger", "no")
	}
	fmt.Println()

	// 2. Base template extraction
	spin := ui.NewSpinner("Writing project files...")
	spin.Start()

	if err := ExtractBaseTemplates(opts.ProjectPath); err != nil {
		spin.Stop()
		return fmt.Errorf("failed to extract base templates: %w", err)
	}

	// 3. Core architecture
	if err := scaffoldCore(opts); err != nil {
		spin.Stop()
		return fmt.Errorf("failed to scaffold core files: %w", err)
	}

	// 4. DB files
	dbBlock, err := scaffoldDB(opts)
	if err != nil {
		spin.Stop()
		return fmt.Errorf("failed to scaffold database files: %w", err)
	}

	// 5. Validator files
	valBlock, err := scaffoldValidator(opts)
	if err != nil {
		spin.Stop()
		return fmt.Errorf("failed to scaffold validator files: %w", err)
	}

	// 6. Global error handler
	if err := scaffoldGlobalErrorHandler(opts.ProjectPath, dbBlock, valBlock); err != nil {
		spin.Stop()
		return fmt.Errorf("failed to scaffold error handler: %w", err)
	}

	spin.Stop()
	ui.Success("Base architecture scaffolded")
	ui.Substep("cem-cli.json   ·  src/app.ts  ·  src/server.ts  ·  .env")
	ui.Substep("src/app/config/index.ts")
	ui.Substep("src/app/errors/    (AppError + db-specific handlers)")
	ui.Substep("src/app/utils/     (catchAsync · sendResponse · logger · welcomePage · validateRequest)")
	ui.Substep("src/app/middlewares/  (globalErrorHandler.middleware · notFound.middleware)")
	ui.Substep("src/app/routes/index.ts")

	// 7. Auth scaffold
	if opts.UseAuth {
		authSpin := ui.NewSpinner("Scaffolding Auth module...")
		authSpin.Start()
		if err := scaffoldAuth(opts); err != nil {
			authSpin.Stop()
			ui.Err("Auth scaffolding failed: " + err.Error())
		} else {
			authSpin.Stop()
			ui.Success("Auth module scaffolded")
			ui.Substep("src/app/modules/Auth/  (controller · service · route · validation · interface)")
			ui.Substep("src/app/utils/jwt.utils.ts  ·  src/app/middlewares/auth.middleware.ts")
			ui.Substep("Token delivery: " + string(opts.TokenDelivery))
			fmt.Println()
			ui.Warn("Replace stub credentials in auth.service.ts before going to production.")
			fmt.Println()
		}
	}

	// 8. Docker scaffold
	if opts.UseDocker {
		dockerSpin := ui.NewSpinner("Generating Docker files...")
		dockerSpin.Start()
		if err := scaffoldDocker(opts); err != nil {
			dockerSpin.Stop()
			ui.Err("Docker scaffold failed: " + err.Error())
		} else {
			dockerSpin.Stop()
			ui.Success("Docker files generated")
			ui.Substep("Dockerfile  ·  .dockerignore  ·  docker-compose.yml")
			fmt.Println()
		}
	}

	// 9. Documentation
	docSpin := ui.NewSpinner("Generating documentation...")
	docSpin.Start()
	if err := scaffoldDocs(opts); err != nil {
		docSpin.Stop()
		ui.Err("Documentation generation failed: " + err.Error())
	} else {
		docSpin.Stop()
		ui.Success("README.md, AGENTS.md & CLAUDE.md generated")
	}

	// 10. Package.json configuration
	if err := configurePackageJson(opts); err != nil {
		return fmt.Errorf("failed to configure package.json: %w", err)
	}

	// 11. Save cem-cli.json
	cfg := &config.CemConfig{
		Name:           opts.ProjectName,
		DB:             string(opts.DB),
		Validator:      string(opts.Validator),
		Auth:           opts.UseAuth,
		TokenDelivery:  string(opts.TokenDelivery),
		Docker:         opts.UseDocker,
		Swagger:        opts.UseSwagger,
		PackageManager: string(opts.PM),
	}
	return config.SaveConfig(opts.ProjectPath, cfg)
}
