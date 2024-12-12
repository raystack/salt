package releaser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

var (
	// Timeout sets the HTTP client timeout for fetching release info.
	Timeout = time.Second * 1

	// APIFormat is the GitHub API URL template to fetch the latest release of a repository.
	APIFormat = "https://api.github.com/repos/%s/releases/latest"
)

// Info holds information about a software release.
type Info struct {
	Version string // Version of the release
	TarURL  string // Tarball URL of the release
}

// FetchInfo fetches details related to the latest release from the provided URL.
//
// Parameters:
//   - releaseURL: The URL to fetch the latest release information from.
//     Example: "https://api.github.com/repos/raystack/optimus/releases/latest"
//
// Returns:
//   - An *Info struct containing the release and tarball URL.
//   - An error if the HTTP request or response parsing fails.
func FetchInfo(url string) (*Info, error) {
	httpClient := http.Client{Timeout: Timeout}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}
	req.Header.Set("User-Agent", "raystack/salt")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch release information from URL: %s", url)
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
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var data struct {
		TagName string `json:"tag_name"`
		Tarball string `json:"tarball_url"`
	}
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, errors.Wrapf(err, "failed to parse JSON response")
	}

	return &Info{
		Version: data.TagName,
		TarURL:  data.Tarball,
	}, nil
}

// CompareVersions compares the current release with the latest release.
//
// Parameters:
//   - currVersion: The current release string.
//   - latestVersion: The latest release string.
//
// Returns:
//   - true if the current release is greater than or equal to the latest release.
//   - An error if release parsing fails.
func CompareVersions(current, latest string) (bool, error) {
	currentVersion, err := version.NewVersion(current)
	if err != nil {
		return false, errors.Wrap(err, "invalid current version format")
	}

	latestVersion, err := version.NewVersion(latest)
	if err != nil {
		return false, errors.Wrap(err, "invalid latest version format")
	}

	return currentVersion.GreaterThanOrEqual(latestVersion), nil
}

// CheckForUpdate generates a message indicating if an update is available.
//
// Parameters:
//   - currentVersion: The current version string (e.g., "v1.0.0").
//   - repo: The GitHub repository in the format "owner/repo".
//
// Returns:
//   - A string containing the update message if a newer version is available.
//   - An empty string if the current version is up-to-date or if an error occurs.
func CheckForUpdate(currentVersion, repo string) string {
	releaseURL := fmt.Sprintf(APIFormat, repo)
	info, err := FetchInfo(releaseURL)
	if err != nil {
		return ""
	}

	isLatest, err := CompareVersions(currentVersion, info.Version)
	if err != nil || isLatest {
		return ""
	}

	return fmt.Sprintf("A new release (%s) is available. consider updating to latest version.", info.Version)
}
