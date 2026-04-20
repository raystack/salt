package server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/raystack/salt/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// startServer starts a server on a random port, waits for it to be ready,
// and returns the base URL. The server shuts down when ctx is cancelled.
func startServer(t *testing.T, ctx context.Context, srv *server.Server) string {
	t.Helper()
	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start(ctx) }()

	// Wait for the server to bind.
	require.Eventually(t, func() bool {
		return srv.ListenAddr() != nil
	}, 2*time.Second, 10*time.Millisecond, "server did not start")

	return "http://" + srv.ListenAddr().String()
}

func TestServer(t *testing.T) {
	t.Run("health check enabled by default", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv := server.New(server.WithAddr("127.0.0.1:0"))
		base := startServer(t, ctx, srv)

		resp, err := http.Get(base + "/ping")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		body, _ := io.ReadAll(resp.Body)
		var result map[string]string
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err)
		assert.Equal(t, "ok", result["status"])
	})

	t.Run("h2c enabled by default", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv := server.New(server.WithAddr("127.0.0.1:0"))
		base := startServer(t, ctx, srv)

		// HTTP/1.1 still works with h2c enabled
		resp, err := http.Get(base + "/ping")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("serves custom handler", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "hello")
		})

		srv := server.New(
			server.WithAddr("127.0.0.1:0"),
			server.WithHandler("/hello", handler),
		)
		base := startServer(t, ctx, srv)

		resp, err := http.Get(base + "/hello")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "hello", string(body))
	})

	t.Run("graceful shutdown completes", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		srv := server.New(
			server.WithAddr("127.0.0.1:0"),
			server.WithGracePeriod(1*time.Second),
		)

		errCh := make(chan error, 1)
		go func() { errCh <- srv.Start(ctx) }()

		require.Eventually(t, func() bool {
			return srv.ListenAddr() != nil
		}, 2*time.Second, 10*time.Millisecond)

		cancel()

		select {
		case err := <-errCh:
			assert.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("shutdown did not complete in time")
		}
	})

	t.Run("custom health check path", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv := server.New(
			server.WithAddr("127.0.0.1:0"),
			server.WithHealthCheck("/healthz"),
		)
		base := startServer(t, ctx, srv)

		resp, err := http.Get(base + "/healthz")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("disable health check", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv := server.New(
			server.WithAddr("127.0.0.1:0"),
			server.WithHealthCheck(""),
		)
		base := startServer(t, ctx, srv)

		resp, err := http.Get(base + "/ping")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
