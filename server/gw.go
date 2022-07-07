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

// GatewayOption sets configs, properties or other parameters for the server.GRPCGateway
type GatewayOption func(*gatewayOptions)

type gatewayOptions struct {
	gwMuxOptions []runtime.ServeMuxOption
	gwmux        *runtime.ServeMux
}

// WithGatewayMuxOptions sets []runtime.ServeMuxOption for server.GRPCGateway
func WithGatewayMuxOptions(opts ...runtime.ServeMuxOption) GatewayOption {
	return func(gwo *gatewayOptions) {
		gwo.gwMuxOptions = opts
	}
}

// WithGRPCGateway sets runtime.ServeMux instance for server.GRPCGateway
func WithGRPCGateway(gwmux *runtime.ServeMux) GatewayOption {
	return func(gwo *gatewayOptions) {
		gwo.gwmux = gwmux
	}
}

// NewGateway creates a new server.GRPCGateway to proxy grpc requests to specified host and port.
func NewGateway(host string, port int, opts ...GatewayOption) (*GRPCGateway, error) {
	gwo := &gatewayOptions{}
	for _, opt := range opts {
		opt(gwo)
	}

	gateway := &GRPCGateway{address: fmt.Sprintf("%s:%d", host, port)}
	if gwo.gwmux != nil {
		gateway.gwmux = gwo.gwmux
	} else {
		gateway.gwmux = runtime.NewServeMux(gwo.gwMuxOptions...)
	}

	return gateway, nil
}

// RegisterHandler helps in adding routes and handlers to be used for proxying requests to grpc service given the grpc-gateway generated Register*ServiceHandlerFromEndpoint function
func (s *GRPCGateway) RegisterHandler(ctx context.Context, f func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) error {
	if err := f(ctx, s.gwmux, s.address, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		return errors.Wrap(err, "RegisterHandler")
	}
	return nil
}
