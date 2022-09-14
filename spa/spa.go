package spa

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// Handler return a file server http handler for single page
// application. This handler can be mounted on your mux server.
//
// If gzip is set true, handler gzip the response body, for clients
// which support it. Usually it also can be left to proxies like Nginx,
// this method is useful when that's undesirable.
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
