package mux

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

const (
	defaultGracePeriod = 10 * time.Second
	grpcRequestMarker  = "application/grpc"
)

// Serve starts a TCP listener and serves the registered protocol servers and blocks
// until server exits. Context can be cancelled to perform graceful shutdown.
// Use WithHTTP() and WithGRPC() to register protocol servers.
func Serve(ctx context.Context, listenAddr string, opts ...Option) error {
	var mux muxServer
	for _, opt := range opts {
		if err := opt(&mux); err != nil {
			return err
		}
	}
	if mux.httpHandler == nil && mux.grpcServer == nil {
		return errors.New("at-least one of http & grpc server must be set")
	}
	return mux.Serve(ctx, listenAddr)
}

type muxServer struct {
	httpHandler http.Handler
	grpcServer  *grpc.Server
	gracePeriod time.Duration
}

func (pmux *muxServer) Serve(baseCtx context.Context, addr string) error {
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()

	httpServer := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(pmux.muxedHandler(), &http2.Server{}),
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[ERROR] server exited with error: %v", err)

			// force-cancel the context so that graceful shutdown sequence exits as well.
			cancel()
		}
	}()

	<-ctx.Done()

	// context has been cancelled (due to cancelled base-context or due to
	// server exit).
	shutdownCtx, cancel := context.WithTimeout(context.Background(), pmux.gracePeriod)
	defer cancel()

	err := httpServer.Shutdown(shutdownCtx)

	if pmux.grpcServer != nil {
		pmux.grpcServer.GracefulStop()
	}
	return err
}

func (pmux *muxServer) muxedHandler() http.Handler {
	// if only one of gRPC and HTTP are set, no need to multiplex.
	if pmux.grpcServer == nil {
		return pmux.httpHandler
	} else if pmux.httpHandler == nil {
		return pmux.grpcServer
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), grpcRequestMarker) {
			pmux.grpcServer.ServeHTTP(w, r)
		} else {
			pmux.httpHandler.ServeHTTP(w, r)
		}
	})
}
