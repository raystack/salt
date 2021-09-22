package server

import (
	"context"
	"fmt"
	"net/http"
)

type HTTPServer struct {
	config     Config
	httpServer *http.Server
	// httpMux is used for allowing addition of custom handlers to the http server
	httpMux *http.ServeMux
}

type HTTPOption func(*httpOptions)

type httpOptions struct {
	httpServer *http.Server
}

func WithHTTPServer(httpServer *http.Server) HTTPOption {
	return func(hos *httpOptions) {
		hos.httpServer = httpServer
	}
}

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

func (s *HTTPServer) Serve() error {
	s.httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.httpMux.ServeHTTP(w, r)
	})
	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) RegisterHandler(pattern string, handler http.Handler) {
	s.httpMux.Handle(pattern, handler)
}

func (s *HTTPServer) SetGateway(gw *GRPCGateway) {
	s.httpMux.Handle("/api/", http.StripPrefix("/api", gw.gwmux))
}

func (s *HTTPServer) Shutdown(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
	go func() {
		<-ctx.Done()
		s.httpServer.Close()
	}()
}
