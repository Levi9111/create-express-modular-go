package updatenotifier

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultReleaseURL = "https://api.github.com/repos/Levi9111/create-express-modular-go/releases/latest"
	cacheFileName     = "create-express-modular-go-update.json"
	cacheTTL          = 24 * time.Hour
	requestTimeout    = 1500 * time.Millisecond
)

type updateCache struct {
	Timestamp     int64  `json:"timestamp"`
	LatestVersion string `json:"latestVersion"`
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func cacheFilePath() string {
	return filepath.Join(os.TempDir(), cacheFileName)
}

func readCache(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}

	var c updateCache
	if err := json.Unmarshal(data, &c); err != nil {
		return "", false
	}

	if time.Since(time.UnixMilli(c.Timestamp)) < cacheTTL {
		return c.LatestVersion, true
	}
	return "", false
}

func writeCache(path, latestVersion string) {
	c := updateCache{
		Timestamp:     time.Now().UnixMilli(),
		LatestVersion: latestVersion,
	}
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o600)
}

// CheckForUpdates queries GitHub Releases for the latest Go CLI version with 24-hour local caching.
// Returns the latest version tag or an empty string if unavailable.
func CheckForUpdates() string {
	return checkForUpdatesFrom(defaultReleaseURL, cacheFilePath())
}

func checkForUpdatesFrom(releaseURL, cachePath string) string {
	if cached, ok := readCache(cachePath); ok {
		return cached
	}

	client := &http.Client{Timeout: requestTimeout}
	req, err := http.NewRequest(http.MethodGet, releaseURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "create-express-modular-go")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		writeCache(cachePath, "")
		return ""
	}

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		writeCache(cachePath, "")
		return ""
	}

	ver := strings.TrimSpace(rel.TagName)
	writeCache(cachePath, ver)
	return ver
}

// IsUpdateAvailable returns true if latest is strictly greater than current (semver X.Y.Z).
func IsUpdateAvailable(current, latest string) bool {
	if current == "" || latest == "" || current == latest {
		return false
	}

	cParts := parseSemver(current)
	lParts := parseSemver(latest)

	for i := 0; i < 3; i++ {
		if lParts[i] > cParts[i] {
			return true
		}
		if lParts[i] < cParts[i] {
			return false
		}
	}
	return false
}

func parseSemver(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	if idx := strings.IndexByte(v, '-'); idx != -1 {
		v = v[:idx]
	}
	parts := strings.Split(v, ".")
	var out [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		n, err := strconv.Atoi(parts[i])
		if err == nil {
			out[i] = n
		}
	}
	return out
}

