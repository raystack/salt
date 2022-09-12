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
func Handler(build embed.FS, dir string, index string) (http.Handler, error) {
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
	router := &router{index: index, fs: http.FS(fsys)}

	return http.FileServer(router), nil
}

// GZipHandler gzip the response body, for clients which support it
// Usually it also can be left to proxies like Nginx, this method
// is useful when that's undesirable.
func GZipHandler(build embed.FS, dir string, index string) (http.Handler, error) {
	handler, err := Handler(build, dir, index)
	if err != nil {
		return nil, err
	}
	return gziphandler.GzipHandler(handler), nil
}
