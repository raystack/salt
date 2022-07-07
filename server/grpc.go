package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// GRPCServer is an server to serve grpc requests
//
// Deprecated: Prefer `mux` package instead of this.
type GRPCServer struct {
	config     Config
	grpcServer *grpc.Server
}

// GRPCOption sets configs, properties or other parameters for the server.GRPCServer
type GRPCOption func(*grpcOptions)

type grpcOptions struct {
	grpcServerOptions []grpc.ServerOption
	grpcServer        *grpc.Server
}

// WithGRPCServerOptions sets []grpc.ServerOption for server.GRPCServer
func WithGRPCServerOptions(opts ...grpc.ServerOption) GRPCOption {
	return func(gos *grpcOptions) {
		gos.grpcServerOptions = opts
	}
}

// WithGRPCServer sets grpc.Server instance for server.GRPCServer
func WithGRPCServer(grpcServer *grpc.Server) GRPCOption {
	return func(gos *grpcOptions) {
		gos.grpcServer = grpcServer
	}
}

// NewGRPC creates a new server.GRPCServer instance with given config and server.GRPCOption
//
// Deprecated: Prefer `mux` package instead of this.
func NewGRPC(config Config, options ...GRPCOption) (*GRPCServer, error) {
	gos := &grpcOptions{}
	for _, opt := range options {
		opt(gos)
	}

	server := &GRPCServer{config: config}
	if gos.grpcServer != nil {
		server.grpcServer = gos.grpcServer
	} else {
		server.grpcServer = grpc.NewServer(gos.grpcServerOptions...)
	}
	return server, nil
}

// Serve starts the configured grpc server to serve requests
func (s *GRPCServer) Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		return err
	}

	reflection.Register(s.grpcServer)
	return s.grpcServer.Serve(l)
}

func (s *GRPCServer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.grpcServer.RegisterService(sd, ss)
}

// Shutdown gracefully stops the server, and kills the server when passed context is cancelled
func (s *GRPCServer) Shutdown(ctx context.Context) {
	s.grpcServer.GracefulStop()
	go func() {
		<-ctx.Done()
		s.grpcServer.Stop()
	}()
}

// RegisterHealth adds standard grpc health check service to grpc server
func (s *GRPCServer) RegisterHealth() *health.Server {
	hs := health.NewServer()
	s.RegisterService(&grpc_health_v1.Health_ServiceDesc, hs)
	return hs
}
