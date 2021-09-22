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

type GRPCServer struct {
	config     Config
	grpcServer *grpc.Server
}

type GRPCOption func(*grpcOptions)

type grpcOptions struct {
	grpcServerOptions []grpc.ServerOption
	grpcServer        *grpc.Server
}

func WithGRPCServerOptions(opts ...grpc.ServerOption) GRPCOption {
	return func(gos *grpcOptions) {
		gos.grpcServerOptions = opts
	}
}

func WithGRPCServer(grpcServer *grpc.Server) GRPCOption {
	return func(gos *grpcOptions) {
		gos.grpcServer = grpcServer
	}
}

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

func (s *GRPCServer) Shutdown(ctx context.Context) {
	s.grpcServer.GracefulStop()
	go func() {
		<-ctx.Done()
		s.grpcServer.Stop()
	}()
}

func (s *GRPCServer) RegisterHealth() *health.Server {
	hs := health.NewServer()
	s.RegisterService(&grpc_health_v1.Health_ServiceDesc, hs)
	return hs
}
