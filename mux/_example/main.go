package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/raystack/salt/common"
	"github.com/raystack/salt/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	grpcGateway := runtime.NewServeMux()

	httpMux := http.NewServeMux()
	httpMux.Handle("/api/", http.StripPrefix("/api", grpcGateway))

	log.Fatalf("server exited: %v", mux.Serve(
		ctx,
		mux.WithHTTPTarget(":8080", &http.Server{
			Handler:        httpMux,
			ReadTimeout:    120 * time.Second,
			WriteTimeout:   120 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}),
		mux.WithGRPCTarget(":8081", grpcServer),
		mux.WithGracePeriod(5*time.Second),
	))
}

type SlowCommonService struct {
	*common.CommonService
}

func (s SlowCommonService) GetVersion(ctx context.Context, req *commonv1.GetVersionRequest) (*commonv1.GetVersionResponse, error) {
	for i := 0; i < 5; i++ {
		log.Printf("dooing stuff")
		time.Sleep(1 * time.Second)
	}
	return s.CommonService.GetVersion(ctx, req)
}
