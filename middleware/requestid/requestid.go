// Package requestid provides request ID propagation middleware.
package requestid

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
)

// Header is the HTTP header used to propagate request IDs.
const Header = "X-Request-ID"

type ctxKey struct{}

// FromContext returns the request ID from the context, or empty string if not set.
func FromContext(ctx context.Context) string {
	if id, ok := ctx.Value(ctxKey{}).(string); ok {
		return id
	}
	return ""
}

// NewContext returns a new context with the given request ID.
func NewContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

func extractOrGenerate(headers http.Header) string {
	if id := headers.Get(Header); id != "" {
		return id
	}
	return uuid.New().String()
}

// NewInterceptor returns a Connect interceptor that propagates or generates request IDs.
func NewInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			id := extractOrGenerate(req.Header())
			ctx = NewContext(ctx, id)
			resp, err := next(ctx, req)
			if resp != nil {
				resp.Header().Set(Header, id)
			}
			return resp, err
		}
	}
}

// HTTPMiddleware returns net/http middleware that propagates or generates request IDs.
func HTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := extractOrGenerate(r.Header)
			ctx := NewContext(r.Context(), id)
			w.Header().Set(Header, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
