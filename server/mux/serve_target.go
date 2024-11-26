package mux

import (
	"context"
	"errors"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

type serveTarget interface {
	Address() string
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
}

type httpServeTarget struct {
	*http.Server
}

func (h httpServeTarget) Address() string { return h.Addr }

func (h httpServeTarget) Serve(l net.Listener) error {
	if err := h.Server.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

type gRPCServeTarget struct {
	Addr string
	*grpc.Server
}

func (g gRPCServeTarget) Address() string { return g.Addr }

func (g gRPCServeTarget) Shutdown(ctx context.Context) error {
	signal := make(chan struct{})
	go func() {
		defer close(signal)

		g.GracefulStop()
	}()

	select {
	case <-ctx.Done():
		g.Stop()
		return errors.New("graceful stop failed")

	case <-signal:
	}

	return nil
}
