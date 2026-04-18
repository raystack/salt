// Package errorz provides error sanitization middleware for Connect services.
package errorz

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect"
)

// Option configures the error sanitization middleware.
type Option func(*config)

type config struct {
	verbose bool
	logger  *slog.Logger
}

// WithVerbose enables full error messages in responses.
// Useful for development/staging environments.
func WithVerbose(v bool) Option {
	return func(c *config) { c.verbose = v }
}

// WithLogger sets the logger for recording original errors before sanitization.
func WithLogger(l *slog.Logger) Option {
	return func(c *config) { c.logger = l }
}

func newConfig(opts []Option) *config {
	c := &config{logger: slog.New(slog.DiscardHandler)}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// NewInterceptor returns a Connect interceptor that sanitizes internal errors.
// Non-Connect errors are mapped to CodeInternal with a timestamp reference.
// Connect errors with known codes are passed through.
func NewInterceptor(opts ...Option) connect.UnaryInterceptorFunc {
	cfg := newConfig(opts)
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err == nil {
				return resp, nil
			}

			// If it's already a Connect error, preserve the code.
			var connectErr *connect.Error
			if errors.As(err, &connectErr) {
				if cfg.verbose {
					return resp, err
				}
				// Preserve code but sanitize message for client-facing codes.
				code := connectErr.Code()
				if code == connect.CodeInternal || code == connect.CodeUnknown {
					ref := time.Now().Unix()
					cfg.logger.Error("internal error",
						"error", err.Error(),
						"ref", ref,
					)
					return resp, connect.NewError(code, fmt.Errorf("internal error (ref: %d)", ref))
				}
				return resp, err
			}

			// Non-Connect error: sanitize completely.
			ref := time.Now().Unix()
			cfg.logger.Error("internal error",
				"error", err.Error(),
				"ref", ref,
			)
			if cfg.verbose {
				return resp, connect.NewError(connect.CodeInternal, err)
			}
			return resp, connect.NewError(connect.CodeInternal, fmt.Errorf("internal error (ref: %d)", ref))
		}
	}
}
