package mux

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

const defaultGracePeriod = 5 * time.Second

// Serve starts a TCP listener and serves the registered protocol
// servers (i.e., http, grpc, etc.) and blocks until server exits.
// Context can be cancelled to perform graceful shutdown.
// Use WithHTTP() and WithGRPC() to register procol servers.
func Serve(ctx context.Context, listenAddr string, opts ...Option) error {
	var mux cmuxWrapper
	for _, opt := range opts {
		if err := opt(&mux); err != nil {
			return err
		}
	}
	if mux.httpServer == nil && mux.grpcServer == nil {
		return errors.New("at-least one of http & grpc server must be set")
	}

	return mux.Serve(ctx, listenAddr)
}

type cmuxWrapper struct {
	httpServer  *http.Server
	grpcServer  *grpc.Server
	gracePeriod time.Duration
}

func (cmw *cmuxWrapper) Serve(baseCtx context.Context, addr string) error {
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	m := cmux.New(l)
	defer m.Close()

	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	errorCh := make(chan error, 3)
	serveOnListener(httpL, cmw.httpServer, errorCh)
	serveOnListener(grpcL, cmw.grpcServer, errorCh)

	go func() {
		if err := m.Serve(); err != nil {
			errorCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return cmw.shutDown()

	case e := <-errorCh:
		_ = cmw.shutDown()
		return e
	}
}

func (cmw *cmuxWrapper) shutDown() error {
	if cmw.grpcServer != nil {
		cmw.grpcServer.GracefulStop()
	}

	if cmw.httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cmw.gracePeriod)
		defer cancel()

		return cmw.httpServer.Shutdown(shutdownCtx)
	}

	return nil
}

func serveOnListener(l net.Listener, server interface{ Serve(l net.Listener) error }, errCh chan<- error) {
	if server == nil {
		return
	}

	go func() {
		err := server.Serve(l)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) || errors.Is(err, grpc.ErrServerStopped) {
				return
			}
			errCh <- err
			return
		}
	}()
}
