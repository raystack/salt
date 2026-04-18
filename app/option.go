package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/raystack/salt/config"
	"github.com/raystack/salt/server"
	"github.com/raystack/salt/telemetry"
)

// Option configures an App.
type Option func(*App) error

// WithConfig loads configuration into the target struct.
// The target must be a pointer to a struct. Config is loaded eagerly
// so that subsequent options can reference fields from it.
func WithConfig(target interface{}, loaderOpts ...config.Option) Option {
	return func(_ *App) error {
		loader := config.NewLoader(loaderOpts...)
		return loader.Load(target)
	}
}

// WithLogger sets the logger for the app and all components.
// The logger is propagated to the server.
func WithLogger(l *slog.Logger) Option {
	return func(a *App) error {
		if l != nil {
			a.logger = l
		}
		return nil
	}
}

// WithTelemetry configures OpenTelemetry.
// Telemetry is initialized when Start() is called.
func WithTelemetry(cfg telemetry.Config) Option {
	return func(a *App) error {
		a.telCfg = &cfg
		return nil
	}
}

// WithAddr sets the server listen address (default ":8080").
func WithAddr(addr string) Option {
	return func(a *App) error {
		a.serverOps = append(a.serverOps, server.WithAddr(addr))
		return nil
	}
}

// WithHandler registers an HTTP handler at the given pattern.
// Use this for ConnectRPC handlers, REST endpoints, SPA handlers, etc.
func WithHandler(pattern string, handler http.Handler) Option {
	return func(a *App) error {
		a.serverOps = append(a.serverOps, server.WithHandler(pattern, handler))
		return nil
	}
}

// WithHTTPMiddleware adds HTTP middleware to the server.
// Use middleware.DefaultHTTP(logger) for the standard chain (recovery,
// request ID, request logging, CORS), or compose your own.
func WithHTTPMiddleware(mw ...func(http.Handler) http.Handler) Option {
	return func(a *App) error {
		a.serverOps = append(a.serverOps, server.WithHTTPMiddleware(mw...))
		return nil
	}
}

// WithGracePeriod sets the shutdown grace period (default 10s).
func WithGracePeriod(d time.Duration) Option {
	return func(a *App) error {
		a.serverOps = append(a.serverOps, server.WithGracePeriod(d))
		return nil
	}
}

// WithServer passes options directly to the underlying server.
// Use this for server options that don't have an app-level wrapper,
// e.g. timeouts:
//
//	app.WithServer(
//	    server.WithReadTimeout(60 * time.Second),
//	    server.WithIdleTimeout(120 * time.Second),
//	)
func WithServer(opts ...server.Option) Option {
	return func(a *App) error {
		a.serverOps = append(a.serverOps, opts...)
		return nil
	}
}

// WithOnStart registers a function to run after infrastructure is ready
// but before the server starts. Use for migrations, seed data, etc.
func WithOnStart(fn func(context.Context) error) Option {
	return func(a *App) error {
		a.onStart = append(a.onStart, fn)
		return nil
	}
}

// WithOnStop registers a function to run during graceful shutdown,
// after the server stops but before infrastructure cleanup.
func WithOnStop(fn func(context.Context) error) Option {
	return func(a *App) error {
		a.onStop = append(a.onStop, fn)
		return nil
	}
}
