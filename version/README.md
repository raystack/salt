# Version

The `version` package provides utilities to fetch and compare software version information from GitHub releases. It helps check if a newer version is available and generates update notices.

## Features

- **Fetch Release Information**: Retrieve the latest release details from a GitHub repository.
- **Version Comparison**: Compare semantic versions to determine if an update is needed.
- **Update Notifications**: Generate user-friendly messages if a newer version is available.

## Installation

To include this package in your Go project, use:

```bash
go get github.com/raystack/salt/version
```

## Usage

### 1. Fetching Release Information

You can use the `ReleaseInfo` function to fetch the latest release details from a GitHub repository.

```go
package main

import (
    "fmt"
    "github.com/raystack/salt/version"
)

func main() {
    releaseURL := "https://api.github.com/repos/raystack/optimus/releases/latest"
    info, err := version.ReleaseInfo(releaseURL)
    if err != nil {
        fmt.Println("Error fetching release info:", err)
        return
    }
    fmt.Printf("Latest Version: %s\nDownload URL: %s\n", info.Version, info.TarURL)
}
```

### 2. Comparing Versions

Use `IsCurrentLatest` to check if the current version is up-to-date with the latest release.

```go
currVersion := "1.2.3"
latestVersion := "1.2.4"
isLatest, err := version.IsCurrentLatest(currVersion, latestVersion)
if err != nil {
    fmt.Println("Error comparing versions:", err)
} else if isLatest {
    fmt.Println("You are using the latest version!")
} else {
    fmt.Println("A newer version is available.")
}
```

### 3. Generating Update Notices

`UpdateNotice` generates a message prompting the user to update if a newer version is available.

```go
notice := version.UpdateNotice("1.0.0", "raystack/optimus")
if notice != "" {
    fmt.Println(notice)
} else {
    fmt.Println("You are up-to-date!")
}
```

## API Reference

### Functions

- `ReleaseInfo(releaseURL string) (*Info, error)`: Fetches the latest release information from the given GitHub API URL.
- `IsCurrentLatest(currVersion, latestVersion string) (bool, error)`: Compares the current version with the latest version using semantic versioning.
- `UpdateNotice(currentVersion, githubRepo string) string`: Returns an update notice if a newer version is available, or an empty string if up-to-date.

### Structs

- `type Info`: Contains details about a release.
    - `Version`: The version string (e.g., "v1.2.3").
    - `TarURL`: The tarball URL for downloading the release.

## Environment Variables

- The `User-Agent` header in HTTP requests is set to `raystack/salt` to comply with GitHub's API requirements.

## Error Handling

- Uses `github.com/pkg/errors` to wrap errors for better error context.
- Returns errors when HTTP requests fail, or when JSON parsing or version comparison fails.

## Dependencies

- `github.com/hashicorp/go-version`: For semantic version comparison.
- `github.com/pkg/errors`: For enhanced error wrapping.
