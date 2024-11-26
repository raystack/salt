package oidc

import (
	"context"
	"crypto/rand"
	_ "embed" // for embedded html
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

const (
	codeParam  = "code"
	stateParam = "state"

	errParam     = "error"
	errDescParam = "error_description"
)

//go:embed redirect.html
var callbackResponsePage string

func encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)
	return encoded
}

// https://tools.ietf.org/html/rfc7636#section-4.1)
func randomBytes(length int) ([]byte, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const csLen = byte(len(charset))
	output := make([]byte, 0, length)
	for {
		buf := make([]byte, length)
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			return nil, fmt.Errorf("failed to read random bytes: %v", err)
		}
		for _, b := range buf {
			// Avoid bias by using a value range that's a multiple of 62
			if b < (csLen * 4) {
				output = append(output, charset[b%csLen])

				if len(output) == length {
					return output, nil
				}
			}
		}
	}
}

func browserAuthzHandler(ctx context.Context, redirectURL, authCodeURL string) (code string, state string, err error) {
	if err := openURL(authCodeURL); err != nil {
		return "", "", err
	}

	u, err := url.Parse(redirectURL)
	if err != nil {
		return "", "", err
	}

	code, state, err = waitForCallback(ctx, fmt.Sprintf(":%s", u.Port()))
	if err != nil {
		return "", "", err
	}
	return code, state, nil
}

func waitForCallback(ctx context.Context, addr string) (code, state string, err error) {
	var cb struct {
		code  string
		state string
		err   error
	}

	stopCh := make(chan struct{})
	srv := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cb.code, cb.state, cb.err = parseCallbackRequest(r)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "text/html")
			_, _ = w.Write([]byte(callbackResponsePage))

			// try to flush to ensure the page is shown to user before we close
			// the server.
			if fl, ok := w.(http.Flusher); ok {
				fl.Flush()
			}

			close(stopCh)
		}),
	}

	go func() {
		select {
		case <-stopCh:
			_ = srv.Close()

		case <-ctx.Done():
			cb.err = ctx.Err()
			_ = srv.Close()
		}
	}()

	if serveErr := srv.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		return "", "", serveErr
	}
	return cb.code, cb.state, cb.err
}

func parseCallbackRequest(r *http.Request) (code string, state string, err error) {
	if err = r.ParseForm(); err != nil {
		return "", "", err
	}

	state = r.Form.Get(stateParam)
	if state == "" {
		return "", "", errors.New("missing state parameter")
	}

	if errorCode := r.Form.Get(errParam); errorCode != "" {
		// Got error from provider. Passing through.
		return "", "", fmt.Errorf("%s: %s", errorCode, r.Form.Get(errDescParam))
	}

	code = r.Form.Get(codeParam)
	if code == "" {
		return "", "", errors.New("missing code parameter")
	}

	return code, state, nil
}

// openURL opens the specified URL in the default application registered for
// the URL scheme.
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
		// If we don't escape &, cmd will ignore everything after the first &.
		url = strings.Replace(url, "&", "^&", -1)

	case "darwin":
		cmd = "open"

	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
