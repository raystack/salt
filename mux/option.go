package mux

import (
	"net/http"
	"time"

	"google.golang.org/grpc"
)

// Option values can be used with Serve() for customisation.
type Option func(m *cmuxWrapper) error

// WithHTTP registers the http-server for use in Serve().
func WithHTTP(server *http.Server) Option {
	return func(m *cmuxWrapper) error {
		if server == nil {
			server = &http.Server{}
		}
		m.httpServer = server
		return nil
	}
}

// WithGRPC registers the gRPC-server for use in Serve().
func WithGRPC(server *grpc.Server) Option {
	return func(m *cmuxWrapper) error {
		if server == nil {
			server = grpc.NewServer()
		}
		m.grpcServer = server
		return nil
	}
}

// WithGracePeriod sets the wait duration for graceful shutdown.
func WithGracePeriod(d time.Duration) Option {
	return func(m *cmuxWrapper) error {
		if d == 0 {
			d = defaultGracePeriod
		}
		m.gracePeriod = d
		return nil
	}
}
