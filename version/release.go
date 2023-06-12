package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

var (
	ReleaseInfoTimeout = time.Second * 1
	Release            = "https://api.github.com/repos/%s/releases/latest"
)

type Info struct {
	Version string
	TarURL  string
}

// ReleaseInfo fetches details related to provided release URL
// releaseURL should point to a specific version
// for example: https://api.github.com/repos/raystack/optimus/releases/latest
func ReleaseInfo(releaseURL string) (*Info, error) {
	httpClient := http.Client{
		Timeout: ReleaseInfoTimeout,
	}
	req, err := http.NewRequest(http.MethodGet, releaseURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request")
	}
	req.Header.Set("User-Agent", "raystack/salt")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to reach releaseURL: %s", releaseURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "failed to reach releaseURL: %s, returned: %d", releaseURL, resp.StatusCode)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read response body")
	}

	var releaseBody struct {
		TagName string `json:"tag_name"`
		Tarball string `json:"tarball_url"`
	}
	if err = json.Unmarshal(body, &releaseBody); err != nil {
		return nil, errors.Wrapf(err, "failed to parse: %s", string(body))
	}

	return &Info{
		Version: releaseBody.TagName,
		TarURL:  releaseBody.Tarball,
	}, nil
}

// IsCurrentLatest returns true if the current version string is greater than
// or equal to latestVersion as per semantic versioning
func IsCurrentLatest(currVersion, latestVersion string) (bool, error) {
	currentV, err := version.NewVersion(currVersion)
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse current version")
	}
	latestV, err := version.NewVersion(latestVersion)
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse latest version")
	}
	if currentV.GreaterThanOrEqual(latestV) {
		return true, nil
	}
	return false, nil
}

// UpdateNotice returns a notice message if there is a newer version available
// Note:  all errors are ignored
func UpdateNotice(currentVersion, githubRepo string) string {
	info, err := ReleaseInfo(fmt.Sprintf(Release, githubRepo))
	if err != nil {
		return ""
	}
	latestVersion := info.Version
	isCurrentLatest, err := IsCurrentLatest(currentVersion, latestVersion)
	if err != nil {
		return ""
	}
	if isCurrentLatest {
		return ""
	}
	return fmt.Sprintf("A new release (%s) is available, consider updating the client.", info.Version)
}
