// Package server provides an HTTP server with h2c support.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/raystack/salt/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	defaultAddr        = ":8080"
	defaultGracePeriod = 10 * time.Second
	defaultHealthPath  = "/ping"
)

// Server is an HTTP server with h2c (HTTP/2 cleartext) support,
// health checks, HTTP middleware, and graceful shutdown.
//
// By default, h2c is enabled and a health check is served at /ping.
type Server struct {
	addr         string
	mux          *http.ServeMux
	h2c          bool
	healthPath   string
	gracePeriod  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
	logger       *slog.Logger
	httpMW       []func(http.Handler) http.Handler
	listenAddr   net.Addr // set after Start binds
}

// New creates a new Server with the given options.
// Defaults: h2c enabled, health check at /ping, grace period 10s.
func New(opts ...Option) *Server {
	s := &Server{
		addr:        defaultAddr,
		mux:         http.NewServeMux(),
		h2c:         true,
		healthPath:  defaultHealthPath,
		gracePeriod: defaultGracePeriod,
		logger:      slog.New(slog.DiscardHandler),
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.healthPath != "" {
		s.mux.HandleFunc(s.healthPath, healthHandler)
	}
	return s
}

// Start begins serving and blocks until the context is cancelled.
// It performs graceful shutdown when the context is done.
func (s *Server) Start(ctx context.Context) error {
	var handler http.Handler = s.mux

	// Apply HTTP middleware chain (outermost first).
	if len(s.httpMW) > 0 {
		handler = middleware.ChainHTTP(s.httpMW...)(handler)
	}

	if s.h2c {
		handler = h2c.NewHandler(handler, &http2.Server{})
	}

	srv := &http.Server{
		Handler:      handler,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
	}

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("server listen: %w", err)
	}

	s.listenAddr = ln.Addr()
	s.logger.Info("server started", "addr", s.listenAddr.String())

	errCh := make(chan error, 1)
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server serve: %w", err)
	case <-ctx.Done():
	}

	s.logger.Info("server shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.gracePeriod)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	s.logger.Info("server stopped")
	return nil
}

// ListenAddr returns the address the server is listening on.
// Only valid after Start has been called.
func (s *Server) ListenAddr() net.Addr {
	return s.listenAddr
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}
