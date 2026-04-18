// Package app provides a service lifecycle manager for raystack services.
package app

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/raystack/salt/server"
	"github.com/raystack/salt/telemetry"
)

// App is a service lifecycle manager that wires together configuration,
// logging, database, telemetry, and HTTP serving with graceful shutdown.
//
// Defaults: h2c enabled, health check at /ping.
type App struct {
	logger    *slog.Logger
	telCfg    *telemetry.Config
	telClean  func()
	serverOps []server.Option
	onStart   []func(context.Context) error
	onStop    []func(context.Context) error
}

// New creates a new App by applying the given options.
func New(opts ...Option) (*App, error) {
	a := &App{
		logger: slog.New(slog.DiscardHandler),
	}
	for _, opt := range opts {
		if err := opt(a); err != nil {
			return nil, fmt.Errorf("app option: %w", err)
		}
	}
	return a, nil
}

// Run is the simplest entry point: creates an App, starts it with signal
// handling (SIGINT, SIGTERM), and blocks until shutdown completes.
func Run(opts ...Option) error {
	a, err := New(opts...)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return a.Start(ctx)
}

// Start initializes all components and starts the server.
// It blocks until the context is cancelled, then performs graceful shutdown.
func (a *App) Start(ctx context.Context) error {
	// Initialize telemetry if configured.
	if a.telCfg != nil {
		cleanup, err := telemetry.Init(ctx, *a.telCfg, a.logger)
		if err != nil {
			return fmt.Errorf("app telemetry: %w", err)
		}
		a.telClean = cleanup
	}

	// Run onStart hooks.
	for _, fn := range a.onStart {
		if err := fn(ctx); err != nil {
			a.cleanup()
			return fmt.Errorf("app on_start: %w", err)
		}
	}

	// Build server with logger.
	opts := make([]server.Option, len(a.serverOps), len(a.serverOps)+1)
	copy(opts, a.serverOps)
	opts = append(opts, server.WithLogger(a.logger))
	srv := server.New(opts...)

	err := srv.Start(ctx)

	// Shutdown sequence.
	a.stop(context.Background())
	return err
}

// Logger returns the app's logger.
func (a *App) Logger() *slog.Logger {
	return a.logger
}

func (a *App) stop(ctx context.Context) {
	for _, fn := range a.onStop {
		if err := fn(ctx); err != nil {
			a.logger.Error("app on_stop hook error", "error", err)
		}
	}
	a.cleanup()
}

func (a *App) cleanup() {
	if a.telClean != nil {
		a.telClean()
	}
}
