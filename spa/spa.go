package spa

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// Init creates a single page application manager.
func Init(build embed.FS, dir string, index string) (*SPA, error) {
	fsys, err := fs.Sub(build, dir)
	if err != nil {
		panic(fmt.Errorf("couldn't create sub filesystem: %w", err))
	}

	_, err = fsys.Open(index)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("ui is enabled but no index.html found: %w", err)
		} else {
			return nil, fmt.Errorf("ui assets error: %w", err)
		}
	}

	return &SPA{&RouterFS{index: index, FS: http.FS(fsys)}}, nil
}

type SPA struct {
	FS http.FileSystem
}

// Handler return a file server http handler for single page
// application. This handler can be mounted on your mux server.
func (spa *SPA) Handler() http.Handler {
	return http.FileServer(spa.FS)
}

// GZipHandler gzip the response body, for clients which support it
// Usually it also can be left to proxies like Nginx, this method
// is useful when that's undesirable.
func (spa *SPA) GZipHandler() http.Handler {
	handler := spa.Handler()
	return gziphandler.GzipHandler(handler)
}
