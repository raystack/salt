package spa

import (
	"errors"
	"io/fs"
	"net/http"
)

// router is the http filesystem which only serves files
// and prevent the directory traversal.
type router struct {
	index string
	fs    http.FileSystem
}

// Open inspects the URL path to locate a file within the static dir.
// If a file is found, it will be served. If not, the file located at
// the index path on the SPA handler will be served.
func (r *router) Open(name string) (http.File, error) {
	file, err := r.fs.Open(name)

	if err == nil {
		return file, nil
	}
	// Serve index if file does not exist.
	if errors.Is(err, fs.ErrNotExist) {
		file, err := r.fs.Open(r.index)
		return file, err
	}

	return nil, err
}
