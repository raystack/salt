# Version

The `version` package provides utilities to fetch and compare software version information from GitHub releases. It helps check if a newer version is available and generates update notices.

## Features

- **Fetch Release Information**: Retrieve the latest release details from a GitHub repository.
- **Version Comparison**: Compare semantic versions to determine if an update is needed.
- **Update Notifications**: Generate user-friendly messages if a newer version is available.

## Installation

To include this package in your Go project, use:

```bash
go get github.com/raystack/salt/cli/release
```

## Usage

### 1. Fetching Release Information

You can use the `FetchInfo` function to fetch the latest release details from a GitHub repository.

```go
package main

import (
    "fmt"
    "github.com/raystack/salt/cli/release"
)

func main() {
    releaseURL := "https://api.github.com/repos/raystack/optimus/releases/latest"
    info, err := release.FetchInfo(releaseURL)
    if err != nil {
        fmt.Println("Error fetching release info:", err)
        return
    }
    fmt.Printf("Latest Version: %s\nDownload URL: %s\n", info.Version, info.TarURL)
}
```

### 2. Comparing Versions

Use `CompareVersions` to check if the current version is up-to-date with the latest release.

```go
current := "1.2.3"
latest := "1.2.4"
isLatest, err := version.CompareVersions(current, latest)
if err != nil {
    fmt.Println("Error comparing versions:", err)
} else if isLatest {
    fmt.Println("You are using the latest release!")
} else {
    fmt.Println("A newer release is available.")
}
```

### 3. Generating Update Notices

`UpdateNotice` generates a message prompting the user to update if a newer version is available.

```go
notice := version.CheckForUpdate("1.0.0", "raystack/optimus")
if notice != "" {
    fmt.Println(notice)
} else {
    fmt.Println("You are up-to-date!")
}
```