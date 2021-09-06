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
// for example: https://api.github.com/repos/odpf/optimus/releases/latest
func ReleaseInfo(releaseURL string) (*Info, error) {
	httpClient := http.Client{
		Timeout: ReleaseInfoTimeout,
	}
	req, err := http.NewRequest(http.MethodGet, releaseURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request")
	}
	req.Header.Set("User-Agent", "odpf/salt")
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

// UpdateMsg returns notification message if there is a greater version available
// then current in github release channel
// Note: all errors are silently ignored
func UpdateMsg(currentVersion, githubRepo string) string {
	info, err := ReleaseInfo(fmt.Sprintf(Release, githubRepo))
	if err != nil {
		return ""
	}
	isLatest, err := IsCurrentLatest(currentVersion, info.Version)
	if err != nil {
		return ""
	}
	if isLatest {
		return fmt.Sprintf("a newer version is available: %s, consider updating the client", info.Version)
	}
	return ""
}
