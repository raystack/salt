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

type MuxServer struct {
	config Config
	GRPCServer
	HTTPServer
}

type MuxOption func(*muxOptions)

type muxOptions struct {
	grpcOptions
	httpOptions
}

func WithMuxGRPCServerOptions(opts ...grpc.ServerOption) MuxOption {
	return func(mos *muxOptions) {
		WithGRPCServerOptions(opts...)(&mos.grpcOptions)
	}
}

func WithMuxGRPCServer(grpcServer *grpc.Server) MuxOption {
	return func(mos *muxOptions) {
		WithGRPCServer(grpcServer)(&mos.grpcOptions)
	}
}

func WithMuxHTTPServer(httpServer *http.Server) MuxOption {
	return func(mos *muxOptions) {
		WithHTTPServer(httpServer)(&mos.httpOptions)
	}
}

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

func (s *MuxServer) Shutdown(ctx context.Context) {
	s.HTTPServer.Shutdown(ctx)
	s.GRPCServer.Shutdown(ctx)
}
