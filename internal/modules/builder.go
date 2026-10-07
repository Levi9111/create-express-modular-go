package modules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Levi9111/create-express-modular-go/internal/pm"
	"github.com/Levi9111/create-express-modular-go/internal/ui"
)

// RunDev starts the development server with live reload.
func RunDev(projectRoot string) error {
	ui.Success("Starting development server...")
	return pm.RunCommand(projectRoot, "npx", "tsx", "watch", "src/server.ts")
}

// RunBuild compiles TypeScript to dist/.
func RunBuild(projectRoot string) error {
	ui.Success("Compiling TypeScript to dist/...")
	return pm.RunCommand(projectRoot, "npx", "tsc")
}

// RunStart runs the production server from dist/server.js.
func RunStart(projectRoot string) error {
	distServer := filepath.Join(projectRoot, "dist/server.js")
	if _, err := os.Stat(distServer); os.IsNotExist(err) {
		ui.Warn("Build output not found. Running 'cem build' first...")
		if err := RunBuild(projectRoot); err != nil {
			return err
		}
	}
	ui.Success("Starting production server...")
	return pm.RunCommand(projectRoot, "node", "dist/server.js")
}

// RunCheck runs type-check, lint, and format checks.
func RunCheck(projectRoot string, fix bool) error {
	ui.SectionHeader("Running Preflight Checks")

	// 1. Typecheck
	spin := ui.NewSpinner("Type checking (tsc)...")
	spin.Start()
	_, err := pm.RunSilentCommand(projectRoot, "npx", "tsc", "--noEmit")
	spin.Stop()
	if err != nil {
		ui.Err("TypeScript type-check failed:")
		fmt.Println(err.Error())
		return err
	}
	ui.Success("TypeScript type-check passed")

	// 2. ESLint
	if fix {
		spin = ui.NewSpinner("Linting & fixing (eslint)...")
		spin.Start()
		_, _ = pm.RunSilentCommand(projectRoot, "npx", "eslint", "--fix", "src")
		spin.Stop()
		ui.Success("ESLint issues auto-fixed")
	} else {
		spin = ui.NewSpinner("Linting (eslint)...")
		spin.Start()
		_, err = pm.RunSilentCommand(projectRoot, "npx", "eslint", "src")
		spin.Stop()
		if err != nil {
			ui.Warn("ESLint found issues. Run 'cem fix' or 'cem check --fix' to resolve.")
		} else {
			ui.Success("ESLint check passed")
		}
	}

	// 3. Prettier
	if fix {
		spin = ui.NewSpinner("Formatting (prettier)...")
		spin.Start()
		_, _ = pm.RunSilentCommand(projectRoot, "npx", "prettier", "--write", "src")
		spin.Stop()
		ui.Success("Prettier formatting applied")
	} else {
		spin = ui.NewSpinner("Format checking (prettier)...")
		spin.Start()
		_, err = pm.RunSilentCommand(projectRoot, "npx", "prettier", "--check", "src")
		spin.Stop()
		if err != nil {
			ui.Warn("Prettier formatting issues detected. Run 'cem fix' to format.")
		} else {
			ui.Success("Prettier formatting check passed")
		}
	}

	return nil
}

// RunFix runs eslint --fix and prettier --write.
func RunFix(projectRoot string) error {
	return RunCheck(projectRoot, true)
}

// RunEject detaches the project from create-express-modular CLI dependency.
func RunEject(projectRoot string) error {
	ui.Warn("Ejecting from create-express-modular CLI dependency...")

	pkgPath := filepath.Join(projectRoot, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return err
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}

	// Replace scripts
	scripts, ok := pkg["scripts"].(map[string]interface{})
	if !ok {
		scripts = make(map[string]interface{})
	}
	scripts["dev"] = "tsx watch src/server.ts"
	scripts["build"] = "tsc"
	scripts["start"] = "node dist/server.js"
	scripts["check"] = "tsc --noEmit && eslint src && prettier --check src"
	scripts["fix"] = "eslint --fix src && prettier --write src"
	pkg["scripts"] = scripts

	// Remove create-express-modular dependency
	if devDeps, ok := pkg["devDependencies"].(map[string]interface{}); ok {
		delete(devDeps, "create-express-modular")
	}

	out, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(pkgPath, out, 0644); err != nil {
		return err
	}

	ui.Success("Project ejected successfully! Standalone scripts added to package.json.")
	return nil
}
