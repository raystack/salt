// Package requestlog provides request logging middleware.
package requestlog

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/raystack/salt/middleware/requestid"
)

// Option configures the request logging middleware.
type Option func(*config)

type config struct {
	logger *slog.Logger
	filter func(procedure string) bool
}

// WithLogger sets the logger.
func WithLogger(l *slog.Logger) Option {
	return func(c *config) { c.logger = l }
}

// WithFilter sets a filter function. If it returns true for a procedure,
// that procedure will not be logged. Useful for skipping health checks.
func WithFilter(fn func(procedure string) bool) Option {
	return func(c *config) { c.filter = fn }
}

func newConfig(opts []Option) *config {
	c := &config{
		logger: slog.New(slog.DiscardHandler),
		filter: func(string) bool { return false },
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// NewInterceptor returns a Connect interceptor that logs requests.
func NewInterceptor(opts ...Option) connect.UnaryInterceptorFunc {
	cfg := newConfig(opts)
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			procedure := req.Spec().Procedure
			if cfg.filter(procedure) {
				return next(ctx, req)
			}

			start := time.Now()
			resp, err := next(ctx, req)
			duration := time.Since(start)

			rid := requestid.FromContext(ctx)
			if err != nil {
				cfg.logger.Error("request completed",
					"procedure", procedure,
					"duration", duration.String(),
					"error", err.Error(),
					"request_id", rid,
				)
			} else {
				cfg.logger.Info("request completed",
					"procedure", procedure,
					"duration", duration.String(),
					"request_id", rid,
				)
			}
			return resp, err
		}
	}
}

// HTTPMiddleware returns net/http middleware that logs requests.
func HTTPMiddleware(opts ...Option) func(http.Handler) http.Handler {
	cfg := newConfig(opts)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if cfg.filter(path) {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sw, r)
			duration := time.Since(start)

			rid := requestid.FromContext(r.Context())
			cfg.logger.Info("request completed",
				"method", r.Method,
				"path", path,
				"status", sw.status,
				"duration", duration.String(),
				"request_id", rid,
			)
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
