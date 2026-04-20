package server

import (
	"log/slog"
	"net/http"
	"time"
)

// Option configures a Server.
type Option func(*Server)

// WithAddr sets the listen address (default ":8080").
func WithAddr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithoutH2C disables HTTP/2 cleartext support.
// H2C is enabled by default for ConnectRPC compatibility.
func WithoutH2C() Option {
	return func(s *Server) {
		s.h2c = false
	}
}

// WithHandler registers an HTTP handler at the given pattern on the server's mux.
func WithHandler(pattern string, handler http.Handler) Option {
	return func(s *Server) {
		s.mux.Handle(pattern, handler)
	}
}

// WithHealthCheck sets the health check endpoint path.
// Default is "/ping". Pass an empty string to disable.
func WithHealthCheck(path string) Option {
	return func(s *Server) {
		s.healthPath = path
	}
}

// WithGracePeriod sets the maximum duration to wait for in-flight
// requests to complete during shutdown (default 10s).
func WithGracePeriod(d time.Duration) Option {
	return func(s *Server) {
		if d > 0 {
			s.gracePeriod = d
		}
	}
}

// WithLogger sets the logger for server lifecycle events.
func WithLogger(l *slog.Logger) Option {
	return func(s *Server) {
		if l != nil {
			s.logger = l
		}
	}
}

// WithHTTPMiddleware adds HTTP middleware to the server.
// Middleware is applied in order (first wraps outermost).
func WithHTTPMiddleware(mw ...func(http.Handler) http.Handler) Option {
	return func(s *Server) {
		s.httpMW = append(s.httpMW, mw...)
	}
}

// WithReadTimeout sets the maximum duration for reading the entire request.
// Zero means no timeout.
func WithReadTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = d
	}
}

// WithWriteTimeout sets the maximum duration for writing the response.
// Zero means no timeout.
func WithWriteTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = d
	}
}

// WithIdleTimeout sets the maximum duration to wait for the next request
// on a keep-alive connection. Zero means no timeout.
func WithIdleTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.idleTimeout = d
	}
}
