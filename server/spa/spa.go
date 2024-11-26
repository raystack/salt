package spa

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// Handler returns an HTTP handler for serving a Single Page Application (SPA).
//
// The handler serves static files from the specified directory in the embedded
// file system and falls back to serving the index file if a requested file is not found.
// This is useful for client-side routing in SPAs.
//
// Parameters:
//   - build: An embedded file system containing the build assets.
//   - dir: The directory within the embedded file system where the static files are located.
//   - index: The name of the index file (usually "index.html").
//   - gzip: If true, the response body will be compressed using gzip for clients that support it.
//
// Returns:
//   - An http.Handler that serves the SPA and optional gzip compression.
//   - An error if the file system or index file cannot be initialized.
func Handler(build embed.FS, dir string, index string, gzip bool) (http.Handler, error) {
	fsys, err := fs.Sub(build, dir)
	if err != nil {
		return nil, fmt.Errorf("couldn't create sub filesystem: %w", err)
	}

	if _, err = fsys.Open(index); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("ui is enabled but no index.html found: %w", err)
		} else {
			return nil, fmt.Errorf("ui assets error: %w", err)
		}
	}
	router := &router{index: index, fs: http.FS(fsys)}

	hlr := http.FileServer(router)

	if !gzip {
		return hlr, nil
	}
	return gziphandler.GzipHandler(hlr), nil
}
