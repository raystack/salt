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

func TestServer(t *testing.T) {
	t.Run("health check enabled by default", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv := server.New(
			server.WithAddr("127.0.0.1:18923"),
		)

		go srv.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://127.0.0.1:18923/ping")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		body, _ := io.ReadAll(resp.Body)
		var result map[string]string
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err)
		assert.Equal(t, "ok", result["status"])

		cancel()
	})

	t.Run("h2c enabled by default", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv := server.New(
			server.WithAddr("127.0.0.1:18924"),
		)

		go srv.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		// HTTP/1.1 still works with h2c enabled
		resp, err := http.Get("http://127.0.0.1:18924/ping")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		cancel()
	})

	t.Run("serves custom handler", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "hello")
		})

		srv := server.New(
			server.WithAddr("127.0.0.1:18925"),
			server.WithHandler("/hello", handler),
		)

		go srv.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://127.0.0.1:18925/hello")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "hello", string(body))

		cancel()
	})

	t.Run("graceful shutdown completes", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		srv := server.New(
			server.WithAddr("127.0.0.1:0"),
			server.WithGracePeriod(1*time.Second),
		)

		errCh := make(chan error, 1)
		go func() { errCh <- srv.Start(ctx) }()

		time.Sleep(100 * time.Millisecond)
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
			server.WithAddr("127.0.0.1:18926"),
			server.WithHealthCheck("/healthz"),
		)

		go srv.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://127.0.0.1:18926/healthz")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		cancel()
	})

	t.Run("disable health check", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv := server.New(
			server.WithAddr("127.0.0.1:18927"),
			server.WithHealthCheck(""),
		)

		go srv.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://127.0.0.1:18927/ping")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		cancel()
	})
}
