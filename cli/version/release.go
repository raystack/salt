package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/hashicorp/go-version"
)

var (
	timeout   = time.Second * 1
	apiFormat = "https://api.github.com/repos/%s/releases/latest"
	cacheTTL  = 24 * time.Hour
)

type releaseInfo struct {
	version string
	tarURL  string
}

type cacheEntry struct {
	CheckedAt     time.Time `json:"checked_at"`
	LatestVersion string    `json:"latest_version"`
}

func fetchInfo(url string) (*releaseInfo, error) {
	httpClient := http.Client{Timeout: timeout}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("User-Agent", "raystack/salt")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release information from URL: %s: %w", url, err)
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from URL: %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data struct {
		TagName string `json:"tag_name"`
		Tarball string `json:"tarball_url"`
	}
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &releaseInfo{
		version: data.TagName,
		tarURL:  data.Tarball,
	}, nil
}

func compareVersions(current, latest string) (bool, error) {
	currentVersion, err := version.NewVersion(current)
	if err != nil {
		return false, fmt.Errorf("invalid current version format: %w", err)
	}

	latestVersion, err := version.NewVersion(latest)
	if err != nil {
		return false, fmt.Errorf("invalid latest version format: %w", err)
	}

	return currentVersion.GreaterThanOrEqual(latestVersion), nil
}

// CheckForUpdate checks GitHub for a newer release and returns an update
// message if one is available. Returns an empty string if up-to-date or
// if the check fails.
//
// Results are cached for 24 hours to avoid hitting GitHub on every invocation.
// The cache is stored at ~/.config/raystack/<repo>/state.json.
func CheckForUpdate(currentVersion, repo string) string {
	// Check cache first.
	if latest, ok := readCache(repo); ok {
		return buildMessage(currentVersion, latest)
	}

	// Fetch from GitHub.
	releaseURL := fmt.Sprintf(apiFormat, repo)
	info, err := fetchInfo(releaseURL)
	if err != nil {
		return ""
	}

	// Cache the result.
	writeCache(repo, info.version)

	return buildMessage(currentVersion, info.version)
}

func buildMessage(current, latest string) string {
	isLatest, err := compareVersions(current, latest)
	if err != nil || isLatest {
		return ""
	}
	return fmt.Sprintf("A new release (%s) is available. consider updating to latest version.", latest)
}

func cachePath(repo string) string {
	dir := configDir()
	return filepath.Join(dir, "raystack", repo, "state.json")
}

func readCache(repo string) (string, bool) {
	data, err := os.ReadFile(cachePath(repo))
	if err != nil {
		return "", false
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", false
	}

	if time.Since(entry.CheckedAt) > cacheTTL {
		return "", false
	}

	return entry.LatestVersion, true
}

func writeCache(repo, latestVersion string) {
	path := cachePath(repo)
	os.MkdirAll(filepath.Dir(path), 0755)

	entry := cacheEntry{
		CheckedAt:     time.Now(),
		LatestVersion: latestVersion,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	os.WriteFile(path, data, 0644)
}

func configDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir
	}
	if runtime.GOOS == "windows" {
		if dir := os.Getenv("APPDATA"); dir != "" {
			return dir
		}
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}
