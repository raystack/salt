// Package cors provides CORS middleware with Connect-specific defaults.
package cors

import (
	"net/http"
	"strconv"
	"strings"
)

// Option configures the CORS middleware.
type Option func(*config)

type config struct {
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
	maxAge         int
}

// WithAllowedOrigins sets the allowed origins. Use "*" to allow all.
func WithAllowedOrigins(origins ...string) Option {
	return func(c *config) { c.allowedOrigins = origins }
}

// WithAllowedMethods sets the allowed HTTP methods.
func WithAllowedMethods(methods ...string) Option {
	return func(c *config) { c.allowedMethods = methods }
}

// WithAllowedHeaders sets the allowed request headers.
func WithAllowedHeaders(headers ...string) Option {
	return func(c *config) { c.allowedHeaders = headers }
}

// WithMaxAge sets the max age (in seconds) for preflight cache.
func WithMaxAge(seconds int) Option {
	return func(c *config) { c.maxAge = seconds }
}

// Defaults returns sensible CORS defaults for ConnectRPC services.
// Includes Connect-specific headers.
func Defaults() []Option {
	return []Option{
		WithAllowedOrigins("*"),
		WithAllowedMethods("GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"),
		WithAllowedHeaders(
			"Content-Type",
			"Connect-Protocol-Version",
			"Connect-Timeout-Ms",
			"Grpc-Timeout",
			"X-Grpc-Web",
			"X-User-Agent",
			"X-Request-ID",
			"Authorization",
		),
		WithMaxAge(7200),
	}
}

func newConfig(opts []Option) *config {
	c := &config{}
	// Apply defaults first, then user overrides.
	for _, opt := range Defaults() {
		opt(c)
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Middleware returns net/http CORS middleware.
func Middleware(opts ...Option) func(http.Handler) http.Handler {
	cfg := newConfig(opts)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Vary", "Origin")

			if !isOriginAllowed(cfg.allowedOrigins, origin) {
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.allowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.allowedHeaders, ", "))

			if cfg.maxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.maxAge))
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(allowed []string, origin string) bool {
	for _, a := range allowed {
		if a == "*" || a == origin {
			return true
		}
	}
	return false
}
