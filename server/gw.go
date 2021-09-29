package server

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// GRPCGateway helps in registering grpc-gateway proxy handlers for a grpc service on server.HTTPServer
type GRPCGateway struct {
	// gwmux is the grpc-gateway proxy multiplexer
	gwmux   *runtime.ServeMux
	address string
}

// NewGateway creates a new server.GRPCGateway to proxy grpc requests to specified host and port
func NewGateway(host string, port int) (*GRPCGateway, error) {
	return &GRPCGateway{
		gwmux:   runtime.NewServeMux(),
		address: fmt.Sprintf("%s:%d", host, port),
	}, nil
}

// RegisterHandler helps in adding routes and handlers to be used for proxying requests to grpc service given the grpc-gateway generated Register*ServiceHandlerFromEndpoint function
func (s *GRPCGateway) RegisterHandler(ctx context.Context, f func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) error {
	if err := f(ctx, s.gwmux, s.address, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		return errors.Wrap(err, "RegisterHandler")
	}
	return nil
}
