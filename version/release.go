package version

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
	// ReleaseInfoTimeout sets the HTTP client timeout for fetching release info.
	ReleaseInfoTimeout = time.Second * 1

	// Release is the GitHub API URL template to fetch the latest release of a repository.
	Release = "https://api.github.com/repos/%s/releases/latest"
)

// Info holds information about a software release.
type Info struct {
	Version string // Version of the release
	TarURL  string // Tarball URL of the release
}

// ReleaseInfo fetches details related to the latest release from the provided URL.
//
// Parameters:
//   - releaseURL: The URL to fetch the latest release information from.
//     Example: "https://api.github.com/repos/raystack/optimus/releases/latest"
//
// Returns:
//   - An *Info struct containing the version and tarball URL.
//   - An error if the HTTP request or response parsing fails.
func ReleaseInfo(releaseURL string) (*Info, error) {
	httpClient := http.Client{
		Timeout: ReleaseInfoTimeout,
	}
	req, err := http.NewRequest(http.MethodGet, releaseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("User-Agent", "raystack/salt")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to reach releaseURL: %s", releaseURL)
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "failed to reach releaseURL: %s, status code: %d", releaseURL, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var releaseData struct {
		TagName string `json:"tag_name"`
		Tarball string `json:"tarball_url"`
	}
	if err = json.Unmarshal(body, &releaseData); err != nil {
		return nil, errors.Wrapf(err, "failed to parse JSON response: %s", string(body))
	}

	return &Info{
		Version: releaseData.TagName,
		TarURL:  releaseData.Tarball,
	}, nil
}

// IsCurrentLatest compares the current version with the latest version.
//
// Parameters:
//   - currVersion: The current version string.
//   - latestVersion: The latest version string.
//
// Returns:
//   - true if the current version is greater than or equal to the latest version.
//   - An error if version parsing fails.
func IsCurrentLatest(currVersion, latestVersion string) (bool, error) {
	currentV, err := version.NewVersion(currVersion)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse current version")
	}
	latestV, err := version.NewVersion(latestVersion)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse latest version")
	}
	return currentV.GreaterThanOrEqual(latestV), nil
}

// UpdateNotice generates a notice message if a newer version is available.
//
// Parameters:
//   - currentVersion: The current version string.
//   - githubRepo: The GitHub repository in the format "owner/repo".
//
// Returns:
//   - A string message prompting the user to update if a newer version is available.
//   - An empty string if there are no updates or if any errors occur.
func UpdateNotice(currentVersion, githubRepo string) string {
	info, err := ReleaseInfo(fmt.Sprintf(Release, githubRepo))
	if err != nil {
		return ""
	}
	latestVersion := info.Version
	isCurrentLatest, err := IsCurrentLatest(currentVersion, latestVersion)
	if err != nil || isCurrentLatest {
		return ""
	}
	return fmt.Sprintf("A new release (%s) is available, consider updating the client.", latestVersion)
}
