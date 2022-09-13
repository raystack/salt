package mux

import (
	"net/http"
	"time"

	"google.golang.org/grpc"
)

// Option values can be used with Serve() for customisation.
type Option func(m *muxServer) error

// WithHTTP registers the http-server for use in Serve().
func WithHTTP(h http.Handler) Option {
	return func(m *muxServer) error {
		m.httpHandler = h
		return nil
	}
}

// WithGRPC registers the gRPC-server for use in Serve().
func WithGRPC(server *grpc.Server) Option {
	return func(m *muxServer) error {
		if server == nil {
			server = grpc.NewServer()
		}
		m.grpcServer = server
		return nil
	}
}

// WithGracePeriod sets the wait duration for graceful shutdown.
func WithGracePeriod(d time.Duration) Option {
	return func(m *muxServer) error {
		if d <= 0 {
			d = defaultGracePeriod
		}
		m.gracePeriod = d
		return nil
	}
}
