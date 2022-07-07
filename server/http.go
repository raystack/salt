package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// HTTPServer is an server to serve http requests
//
// Deprecated: Prefer `mux` package instead of this.
type HTTPServer struct {
	config     Config
	httpServer *http.Server
	// httpMux is used for allowing addition of custom handlers to the http server
	httpMux *http.ServeMux
}

// HTTPOption sets configs, properties or other parameters for the server.HTTPServer
type HTTPOption func(*httpOptions)

type httpOptions struct {
	httpServer *http.Server
}

// WithHTTPServer sets http.Server instance for server.HTTPServer
func WithHTTPServer(httpServer *http.Server) HTTPOption {
	return func(hos *httpOptions) {
		hos.httpServer = httpServer
	}
}

// NewHTTP creates a new server.HTTPServer instance with given config and server.HTTPOption
//
// Deprecated: Prefer `mux` package instead of this.
func NewHTTP(config Config, options ...HTTPOption) (*HTTPServer, error) {
	hos := &httpOptions{}
	for _, opt := range options {
		opt(hos)
	}

	server := &HTTPServer{config: config}
	if hos.httpServer != nil {
		server.httpServer = hos.httpServer
	} else {
		server.httpServer = &http.Server{}
	}
	server.httpServer.Addr = fmt.Sprintf("%s:%d", config.Host, config.Port)
	server.httpMux = http.NewServeMux()

	return server, nil
}

// Serve starts the configured http server to serve requests
func (s *HTTPServer) Serve() error {
	s.httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.httpMux.ServeHTTP(w, r)
	})
	return s.httpServer.ListenAndServe()
}

// RegisterHandler registers provided pattern and handler on the http server
func (s *HTTPServer) RegisterHandler(pattern string, handler http.Handler) {
	s.httpMux.Handle(pattern, handler)
}

// SetGateway sets a server.GRPCGateway instance on the http server to be proxy requests to a grpc service
func (s *HTTPServer) SetGateway(patternPrefix string, gw *GRPCGateway) {
	prefix := strings.TrimSuffix(patternPrefix, "/")
	pattern := prefix + "/"
	s.httpMux.Handle(pattern, http.StripPrefix(prefix, gw.gwmux))
}

// Shutdown gracefully stops the server, and kills the server when passed context is cancelled
func (s *HTTPServer) Shutdown(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
	go func() {
		<-ctx.Done()
		s.httpServer.Close()
	}()
}
