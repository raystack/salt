package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/raystack/salt/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var GRPCMiddlewaresInterceptor = grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	grpc_recovery.UnaryServerInterceptor(),
	grpc_ctxtags.UnaryServerInterceptor(),
	grpc_prometheus.UnaryServerInterceptor,
	grpc_zap.UnaryServerInterceptor(zap.NewExample()),
))

func main() {
	grpcPort := 8000
	httpPort := 8080
	muxPort := 9000
	gatewayClientPort := grpcPort
	go grpcS(grpcPort)
	go httpS(httpPort, gatewayClientPort)
	muxS(muxPort)
}

func httpS(httpPort, gatewayClientPort int) {
	ctx, cancelFunc := context.WithCancel(server.HandleSignals(context.Background()))
	defer cancelFunc()

	s, err := server.NewHTTP(server.Config{
		Port: httpPort,
	})
	if err != nil {
		panic(err)
	}

	gw, err := server.NewGateway("", gatewayClientPort)
	if err != nil {
		panic(err)
	}

	s.SetGateway("/api", gw)
	s.RegisterHandler("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}))

	go s.Serve()
	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()

	s.Shutdown(shutdownCtx)
}

func grpcS(grpcPort int) {
	ctx, cancelFunc := context.WithCancel(server.HandleSignals(context.Background()))
	defer cancelFunc()

	s, err := server.NewGRPC(server.Config{
		Port: grpcPort,
	}, server.WithGRPCServerOptions(GRPCMiddlewaresInterceptor))
	if err != nil {
		panic(err)
	}

	go s.Serve()
	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()

	s.Shutdown(shutdownCtx)
}

func muxS(muxPort int) {
	ctx, cancelFunc := context.WithCancel(server.HandleSignals(context.Background()))
	defer cancelFunc()

	s, err := server.NewMux(server.Config{
		Port: muxPort,
	}, server.WithMuxGRPCServerOptions(GRPCMiddlewaresInterceptor))
	if err != nil {
		panic(err)
	}

	gatewayClientPort := muxPort
	gw, err := server.NewGateway("", gatewayClientPort)
	if err != nil {
		panic(err)
	}

	s.SetGateway("/api", gw)

	s.RegisterHandler("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}))

	go s.Serve()
	<-ctx.Done()
	// clean anything that needs to be closed etc like common server implementation etc
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()

	s.Shutdown(shutdownCtx)
}
