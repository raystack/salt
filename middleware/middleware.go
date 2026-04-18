// Package middleware provides Connect interceptors and HTTP middleware.
package middleware

import (
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/raystack/salt/middleware/cors"
	"github.com/raystack/salt/middleware/errorz"
	"github.com/raystack/salt/middleware/recovery"
	"github.com/raystack/salt/middleware/requestid"
	"github.com/raystack/salt/middleware/requestlog"
)

// Default returns the standard raystack Connect interceptor chain:
// recovery → requestid → requestlog → errorz
func Default(l *slog.Logger) []connect.Interceptor {
	return []connect.Interceptor{
		recovery.NewInterceptor(recovery.WithLogger(l)),
		requestid.NewInterceptor(),
		requestlog.NewInterceptor(requestlog.WithLogger(l)),
		errorz.NewInterceptor(errorz.WithLogger(l)),
	}
}

// DefaultHTTP returns the standard raystack HTTP middleware chain:
// recovery → requestid → requestlog → cors
func DefaultHTTP(l *slog.Logger, corsOpts ...cors.Option) func(http.Handler) http.Handler {
	return ChainHTTP(
		recovery.HTTPMiddleware(recovery.WithLogger(l)),
		requestid.HTTPMiddleware(),
		requestlog.HTTPMiddleware(requestlog.WithLogger(l)),
		cors.Middleware(corsOpts...),
	)
}

// ChainHTTP chains net/http middleware in order.
// The first middleware wraps outermost (processes request first).
func ChainHTTP(mws ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}
