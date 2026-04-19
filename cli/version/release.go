package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-version"
)

var (
	timeout   = time.Second * 1
	apiFormat = "https://api.github.com/repos/%s/releases/latest"
)

type releaseInfo struct {
	version string
	tarURL  string
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
func CheckForUpdate(currentVersion, repo string) string {
	releaseURL := fmt.Sprintf(apiFormat, repo)
	info, err := fetchInfo(releaseURL)
	if err != nil {
		return ""
	}

	isLatest, err := compareVersions(currentVersion, info.version)
	if err != nil || isLatest {
		return ""
	}

	return fmt.Sprintf("A new release (%s) is available. consider updating to latest version.", info.version)
}
