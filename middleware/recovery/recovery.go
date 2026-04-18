// Package recovery provides panic recovery middleware.
package recovery

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"connectrpc.com/connect"
)

// Option configures the recovery middleware.
type Option func(*config)

type config struct {
	logger  *slog.Logger
	handler func(ctx context.Context, p any) error
}

// WithLogger sets the logger for panic reporting.
func WithLogger(l *slog.Logger) Option {
	return func(c *config) { c.logger = l }
}

// WithHandler sets a custom panic handler. If it returns an error,
// that error is returned to the client.
func WithHandler(fn func(ctx context.Context, p any) error) Option {
	return func(c *config) { c.handler = fn }
}

func newConfig(opts []Option) *config {
	c := &config{logger: slog.New(slog.DiscardHandler)}
	for _, opt := range opts {
		opt(c)
	}
	if c.handler == nil {
		c.handler = func(_ context.Context, _ any) error {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("internal error"))
		}
	}
	return c
}

// NewInterceptor returns a Connect interceptor that recovers from panics.
func NewInterceptor(opts ...Option) connect.UnaryInterceptorFunc {
	cfg := newConfig(opts)
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
			defer func() {
				if p := recover(); p != nil {
					cfg.logger.Error("panic recovered", "panic", p, "stack", string(debug.Stack()))
					err = cfg.handler(ctx, p)
				}
			}()
			return next(ctx, req)
		}
	}
}

// HTTPMiddleware returns net/http middleware that recovers from panics.
func HTTPMiddleware(opts ...Option) func(http.Handler) http.Handler {
	cfg := newConfig(opts)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					cfg.logger.Error("panic recovered", "panic", p, "stack", string(debug.Stack()))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
