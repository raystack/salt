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
	"github.com/odpf/salt/common"
	"github.com/odpf/salt/server"
	commonv1 "go.buf.build/odpf/gw/odpf/proton/odpf/common/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var Server = &commonv1.Version{
	Version: "v1.0.2",
}

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
		Host: "",
	})
	if err != nil {
		panic(err)
	}

	gw, err := server.NewGateway("", gatewayClientPort)
	if err != nil {
		panic(err)
	}
	gw.RegisterHandler(ctx, commonv1.RegisterCommonServiceHandlerFromEndpoint)

	s.SetGateway(gw)
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
		Host: "",
	}, server.WithGRPCServerOptions(GRPCMiddlewaresInterceptor))
	if err != nil {
		panic(err)
	}

	s.RegisterService(&commonv1.CommonService_ServiceDesc,
		common.New(Server),
	)

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
		Host: "",
	}, server.WithMuxGRPCServerOptions(GRPCMiddlewaresInterceptor))
	if err != nil {
		panic(err)
	}

	gatewayClientPort := muxPort
	gw, err := server.NewGateway("", gatewayClientPort)
	if err != nil {
		panic(err)
	}
	gw.RegisterHandler(ctx, commonv1.RegisterCommonServiceHandlerFromEndpoint)

	s.SetGateway(gw)

	s.RegisterHandler("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}))

	s.RegisterService(&commonv1.CommonService_ServiceDesc,
		common.New(Server),
	)

	go s.Serve()
	<-ctx.Done()
	// clean anything that needs to be closed etc like common server implementation etc
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()

	s.Shutdown(shutdownCtx)
}
