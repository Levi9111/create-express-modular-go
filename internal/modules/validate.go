package modules

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	identifierRegex  = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
	envKeyRegex      = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	projectNameRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
)

// ValidateIdentifier ensures module/middleware names are safe single identifiers without path traversal.
func ValidateIdentifier(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fmt.Errorf("identifier cannot be empty")
	}
	if !identifierRegex.MatchString(trimmed) {
		return fmt.Errorf("invalid identifier %q: must start with a letter and contain only letters, digits, or underscores", name)
	}
	return nil
}

// ValidateEnvKey ensures environment variable names are valid shell/JS identifiers.
func ValidateEnvKey(key string) error {
	trimmed := strings.TrimSpace(key)
	if trimmed == "" {
		return fmt.Errorf("environment variable key cannot be empty")
	}
	if !envKeyRegex.MatchString(trimmed) {
		return fmt.Errorf("invalid environment variable key %q: must start with a letter or underscore and contain only letters, digits, or underscores", key)
	}
	return nil
}

// ValidateProjectName ensures project names do not contain path separators or traversal sequences.
func ValidateProjectName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if trimmed == "." {
		return nil
	}
	if strings.Contains(trimmed, "..") || !projectNameRegex.MatchString(trimmed) {
		return fmt.Errorf("invalid project name %q: must be '.' or start with alphanumeric and contain only letters, digits, '.', '_', or '-'", name)
	}
	return nil
}
