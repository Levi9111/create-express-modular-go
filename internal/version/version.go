package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the current semantic version of the CLI.
	// Can be overridden at build time via -ldflags "-X ...version.Version=vX.Y.Z"
	Version = "3.3.10-go.exp"

	// GitCommit is the git commit SHA injected at build time.
	GitCommit = "dev"

	// BuildDate is the ISO timestamp injected at build time.
	BuildDate = "unknown"
)

// Info returns the comprehensive build and environment information string.
func Info() string {
	return fmt.Sprintf("cem version %s (commit: %s, built: %s, %s/%s)",
		Version, GitCommit, BuildDate, runtime.GOOS, runtime.GOARCH)
}

// Short returns just the semantic version.
func Short() string {
	return Version
}
