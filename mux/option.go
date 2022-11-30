package mux

import (
	"net/http"
	"time"

	"google.golang.org/grpc"
)

// Option values can be used with Serve() for customisation.
type Option func(m *muxServer) error

func WithHTTPTarget(addr string, srv *http.Server) Option {
	srv.Addr = addr
	return func(m *muxServer) error {
		m.targets = append(m.targets, httpServeTarget{Server: srv})
		return nil
	}
}

func WithGRPCTarget(addr string, srv *grpc.Server) Option {
	return func(m *muxServer) error {
		m.targets = append(m.targets, gRPCServeTarget{
			Addr:   addr,
			Server: srv,
		})
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
