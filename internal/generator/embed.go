package generator

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed templates/base/*
var baseTemplates embed.FS

// ExtractBaseTemplates writes the embedded base configuration files into projectPath.
func ExtractBaseTemplates(projectPath string) error {
	entries, err := fs.ReadDir(baseTemplates, "templates/base")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := baseTemplates.ReadFile("templates/base/" + entry.Name())
		if err != nil {
			return err
		}

		targetName := entry.Name()
		if targetName == "gitignore" {
			targetName = ".gitignore"
		}

		destPath := filepath.Join(projectPath, targetName)
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return err
		}
	}
	return nil
}
