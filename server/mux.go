package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// MuxServer is an server to serve grpc requests and http requests on same host and port
//
// Deprecated: Prefer `mux` package instead of this.
type MuxServer struct {
	config Config
	GRPCServer
	HTTPServer
}

// MuxOption sets configs, properties or other parameters for the server.MuxServer
type MuxOption func(*muxOptions)

type muxOptions struct {
	grpcOptions
	httpOptions
}

// WithMuxGRPCServerOptions sets []grpc.ServerOption for the internal grpc server of server.MuxServer
func WithMuxGRPCServerOptions(opts ...grpc.ServerOption) MuxOption {
	return func(mos *muxOptions) {
		WithGRPCServerOptions(opts...)(&mos.grpcOptions)
	}
}

// WithMuxGRPCServer sets grpc.Server instance for the internal grpc server of server.MuxServer
func WithMuxGRPCServer(grpcServer *grpc.Server) MuxOption {
	return func(mos *muxOptions) {
		WithGRPCServer(grpcServer)(&mos.grpcOptions)
	}
}

// WithMuxHTTPServer sets http.Server instance for the internal http server of server.MuxServer
func WithMuxHTTPServer(httpServer *http.Server) MuxOption {
	return func(mos *muxOptions) {
		WithHTTPServer(httpServer)(&mos.httpOptions)
	}
}

// NewMux creates a new server.MuxServer instance with given config and server.MuxOption
//
// Deprecated: Prefer `mux` package instead of this.
func NewMux(config Config, options ...MuxOption) (*MuxServer, error) {
	mos := &muxOptions{}
	for _, opt := range options {
		opt(mos)
	}

	server := &MuxServer{config: config}
	if mos.grpcServer != nil {
		server.grpcServer = mos.grpcServer
	} else {
		server.grpcServer = grpc.NewServer(mos.grpcServerOptions...)
	}

	if mos.httpServer != nil {
		server.httpServer = mos.httpServer
	} else {
		server.httpServer = &http.Server{}
	}

	server.httpMux = http.NewServeMux()

	return server, nil
}

// Serve starts the configured grpc and http servers to serve requests
func (s *MuxServer) Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		return err
	}

	m := cmux.New(l)
	defer m.Close()

	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	s.httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.httpMux.ServeHTTP(w, r)
	})

	errorChannel := make(chan error)

	go func() {
		reflection.Register(s.grpcServer)
		err := s.grpcServer.Serve(grpcL)
		if err != nil {
			errorChannel <- err
		}
	}()

	go func() {
		err := s.httpServer.Serve(httpL)
		if err != nil {
			errorChannel <- err
		}
	}()

	go func() {
		err := m.Serve()
		if err != nil {
			errorChannel <- err
		}
	}()

	return <-errorChannel
}

// Shutdown gracefully stops the server, and kills the server when passed context is cancelled
func (s *MuxServer) Shutdown(ctx context.Context) {
	s.HTTPServer.Shutdown(ctx)
	s.GRPCServer.Shutdown(ctx)
}
