# SPA Server Package

The `spa` package provides a simple HTTP handler to serve Single Page Applications (SPAs) from an embedded file system, with optional gzip compression. This is particularly useful for applications that need to serve static assets and handle client-side routing, where all paths should fall back to an `index.html` file.

## Features

- **Serve Embedded Static Files**: Serve files directly from an embedded filesystem (`embed.FS`), making deployments easier.
- **SPA Support with Client-Side Routing**: Automatically serves `index.html` when a requested file is not found, allowing for client-side routing.
- **Optional Gzip Compression**: Optionally compresses responses with gzip for clients that support it.

## Installation

Add the package to your Go project by running:

```bash
go get github.com/raystack/spa
```

## Usage

Here’s an example of using `spa` to serve a Single Page Application from an embedded file system.

### Embed Your Static Files

Embed your static files (like `index.html`, JavaScript, CSS, etc.) using Go’s `embed` package:

```go
//go:embed all:build
var content embed.FS
```

### Setting Up the Server

Use the `Handler` function to create an HTTP handler that serves your SPA with optional gzip compression.

```go
package main

import (
    "embed"
    "log"
    "net/http"

    "github.com/raystack/spa"
)

//go:embed all:build
var content embed.FS

func main() {
    handler, err := spa.Handler(content, "build", "index.html", true)
    if err != nil {
        log.Fatalf("failed to initialize SPA handler: %v", err)
    }

    http.Handle("/", handler)
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}
```

In this example:
- `content`: Embedded filesystem containing the build directory.
- `"build"`: The directory within the embedded filesystem where the static files are located.
- `"index.html"`: The fallback file to serve when a requested file isn’t found, typically used for client-side routing.
- `true`: Enables gzip compression for supported clients.

## API Reference

### `Handler`

```go
func Handler(build embed.FS, dir string, index string, gzip bool) (http.Handler, error)
```

Creates an HTTP handler to serve an SPA with optional gzip compression.

- **Parameters**:
    - `build`: The embedded file system containing the static files.
    - `dir`: The subdirectory within `build` where static files are located.
    - `index`: The fallback file (usually "index.html") to serve when a requested file isn’t found.
    - `gzip`: If `true`, responses will be compressed with gzip for clients that support it.

- **Returns**: An `http.Handler` for serving the SPA, or an error if initialization fails.

### `router`

The `router` struct is an HTTP file system wrapper that prevents directory traversal and supports client-side routing by serving `index.html` for unknown paths.

## Example Scenarios

- **Deploying a Go-Based SPA**: Use `spa` to embed and serve your frontend from within your Go binary.
- **Supporting Client-Side Routing**: Serve a fallback `index.html` page for any route that doesn't match an existing file, supporting SPAs with dynamic URLs.
- **Optional Compression**: Enable gzip for production deployments to reduce bandwidth usage.
