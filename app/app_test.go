package app_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/raystack/salt/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func nopLogger() *slog.Logger { return slog.New(slog.DiscardHandler) }

// freeAddr returns a "127.0.0.1:<port>" string using a port that is free at
// the time of the call. There is a small TOCTOU window, but it eliminates
// hardcoded-port flakes in CI.
func freeAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := ln.Addr().String()
	ln.Close()
	return addr
}

func TestNew(t *testing.T) {
	t.Run("creates app with defaults", func(t *testing.T) {
		a, err := app.New()
		require.NoError(t, err)
		assert.NotNil(t, a)
		assert.NotNil(t, a.Logger())
	})

	t.Run("sets logger", func(t *testing.T) {
		l := nopLogger()
		a, err := app.New(app.WithLogger(l))
		require.NoError(t, err)
		assert.Equal(t, l, a.Logger())
	})

	t.Run("returns error from option", func(t *testing.T) {
		badOpt := func(_ *app.App) error {
			return fmt.Errorf("bad option")
		}
		_, err := app.New(badOpt)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bad option")
	})
}

func TestAppStartAndShutdown(t *testing.T) {
	t.Run("starts with health check and h2c by default", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		addr := freeAddr(t)

		a, err := app.New(
			app.WithLogger(nopLogger()),
			app.WithAddr(addr),
		)
		require.NoError(t, err)

		errCh := make(chan error, 1)
		go func() { errCh <- a.Start(ctx) }()

		time.Sleep(100 * time.Millisecond)

		// Health check should be on by default at /ping
		resp, err := http.Get("http://" + addr + "/ping")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		cancel()

		select {
		case err := <-errCh:
			assert.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("shutdown timed out")
		}
	})

	t.Run("runs onStart hooks", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		addr := freeAddr(t)

		var hookRan bool
		a, err := app.New(
			app.WithAddr(addr),
			app.WithOnStart(func(_ context.Context) error {
				hookRan = true
				return nil
			}),
		)
		require.NoError(t, err)

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		a.Start(ctx)
		assert.True(t, hookRan)
	})

	t.Run("runs onStop hooks", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		addr := freeAddr(t)

		var hookRan bool
		a, err := app.New(
			app.WithAddr(addr),
			app.WithOnStop(func(_ context.Context) error {
				hookRan = true
				return nil
			}),
		)
		require.NoError(t, err)

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		a.Start(ctx)
		assert.True(t, hookRan)
	})

	t.Run("serves custom handler", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		addr := freeAddr(t)

		a, err := app.New(
			app.WithLogger(nopLogger()),
			app.WithAddr(addr),
			app.WithHandler("/hello", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, "world")
			})),
		)
		require.NoError(t, err)

		go a.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://" + addr + "/hello")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "world", string(body))

		cancel()
	})

	t.Run("applies explicit middleware", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		addr := freeAddr(t)

		addHeader := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Custom", "salt")
				next.ServeHTTP(w, r)
			})
		}

		a, err := app.New(
			app.WithLogger(nopLogger()),
			app.WithAddr(addr),
			app.WithHTTPMiddleware(addHeader),
			app.WithHandler("/test", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
		)
		require.NoError(t, err)

		go a.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://" + addr + "/test")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, "salt", resp.Header.Get("X-Custom"))

		cancel()
	})

	t.Run("onStart failure returns error and runs cleanup", func(t *testing.T) {
		addr := freeAddr(t)
		var cleanupRan bool
		a, err := app.New(
			app.WithAddr(addr),
			app.WithOnStart(func(_ context.Context) error {
				return fmt.Errorf("migration failed")
			}),
			app.WithOnStop(func(_ context.Context) error {
				cleanupRan = true
				return nil
			}),
		)
		require.NoError(t, err)

		err = a.Start(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "migration failed")
		// OnStop should NOT run — it runs only on graceful shutdown.
		// cleanup() (telemetry flush) runs, but not onStop hooks.
		assert.False(t, cleanupRan, "onStop hooks should not run on startup failure")
	})
}
