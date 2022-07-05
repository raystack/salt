package main

import (
	"context"
	"log"
	"net/http"
	"time"

	commonv1 "go.buf.build/odpf/gw/odpf/proton/odpf/common/v1"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/salt/common"
	"github.com/odpf/salt/mux"
)

func main() {
	ctx := context.Background()

	grpcServer := grpc.NewServer()
	grpcGateway := runtime.NewServeMux()

	commonSvc := common.New(nil)
	grpcServer.RegisterService(&commonv1.CommonService_ServiceDesc, commonSvc)
	if err := commonv1.RegisterCommonServiceHandlerServer(ctx, grpcGateway, commonSvc); err != nil {
		panic(err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/api/", http.StripPrefix("/api", grpcGateway))

	log.Fatalf("server exited: %v", mux.Serve(ctx, "localhost:8080",
		mux.WithHTTP(&http.Server{Handler: httpMux}),
		mux.WithGRPC(grpcServer),
		mux.WithGracePeriod(5*time.Second),
	))
}
