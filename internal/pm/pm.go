package pm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PackageManager string

const (
	NPM  PackageManager = "npm"
	Yarn PackageManager = "yarn"
	PNPM PackageManager = "pnpm"
	Bun  PackageManager = "bun"
)

// DetectPM detects the active package manager.
func DetectPM(cwd string) PackageManager {
	if _, err := os.Stat(filepath.Join(cwd, "bun.lockb")); err == nil {
		return Bun
	}
	if _, err := os.Stat(filepath.Join(cwd, "bun.lock")); err == nil {
		return Bun
	}
	if _, err := os.Stat(filepath.Join(cwd, "pnpm-lock.yaml")); err == nil {
		return PNPM
	}
	if _, err := os.Stat(filepath.Join(cwd, "yarn.lock")); err == nil {
		return Yarn
	}
	if _, err := os.Stat(filepath.Join(cwd, "package-lock.json")); err == nil {
		return NPM
	}

	ua := os.Getenv("npm_config_user_agent")
	if strings.HasPrefix(ua, "bun") {
		return Bun
	}
	if strings.HasPrefix(ua, "pnpm") {
		return PNPM
	}
	if strings.HasPrefix(ua, "yarn") {
		return Yarn
	}

	return NPM
}

// InitialInstallCmd returns the optimized command and args for first-time install.
func InitialInstallCmd(pm PackageManager) (string, []string) {
	switch pm {
	case Bun:
		return "bun", []string{"install"}
	case PNPM:
		return "pnpm", []string{"install", "--prefer-offline"}
	case Yarn:
		return "yarn", []string{"install"}
	default:
		return "npm", []string{"install", "--no-audit", "--no-fund", "--no-progress", "--loglevel=error", "--prefer-offline"}
	}
}

// AddPackagesCmd returns command and args to add dependencies.
func AddPackagesCmd(pm PackageManager, packages []string, dev bool) (string, []string) {
	var base []string
	switch pm {
	case Bun:
		if dev {
			base = []string{"add", "-d"}
		} else {
			base = []string{"add"}
		}
		return "bun", append(base, packages...)
	case PNPM:
		if dev {
			base = []string{"add", "-D"}
		} else {
			base = []string{"add"}
		}
		return "pnpm", append(base, packages...)
	case Yarn:
		if dev {
			base = []string{"add", "-D"}
		} else {
			base = []string{"add"}
		}
		return "yarn", append(base, packages...)
	default:
		if dev {
			base = []string{"install", "-D"}
		} else {
			base = []string{"install"}
		}
		return "npm", append(base, packages...)
	}
}

// RunScriptCmd returns the executable and arguments to run an npm script.
func RunScriptCmd(pm PackageManager, script string) (string, []string) {
	switch pm {
	case NPM:
		return "npm", []string{"run", script}
	default:
		return string(pm), []string{script}
	}
}

// RunCommand executes a command in the given directory and returns output or error.
func RunCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// RunSilentCommand executes silently and returns combined output on error.
func RunSilentCommand(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%w: %s", err, string(out))
	}
	return string(out), nil
}
