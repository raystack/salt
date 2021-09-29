# server

server package helps in setting up grpc server, http server or a mux server that runs both http and grpc on same port on a host
It exposes multiple options to configure each of the servers.

## Usage

### HTTP server

```go
    // context to be Done when SIGINT or SIGTERM is received
	ctx, cancelFunc := context.WithCancel(server.HandleSignals(context.Background()))
	defer cancelFunc()

    // create server.HTTPServer
	s, err := server.NewHTTP(server.Config{
		Port: httpPort,
	})
	if err != nil {
		panic(err)
	}

    // add a handler
	s.RegisterHandler("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}))

    // start serving
	go s.Serve()
    // wait for ctx context to be set done, by SIGINT or SIGTERM
	<-ctx.Done()

    // set a timer for graceful shutdown, if expired will force kill the server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()

	s.Shutdown(shutdownCtx)
```

### GRPC server

```go
    // grpc middlewares
    var GRPCMiddlewaresInterceptor = grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	    grpc_recovery.UnaryServerInterceptor(),
	    grpc_ctxtags.UnaryServerInterceptor(),
	    grpc_prometheus.UnaryServerInterceptor,
	    grpc_zap.UnaryServerInterceptor(zap.NewExample()),
    ))

    // context to be Done when SIGINT or SIGTERM is received
	ctx, cancelFunc := context.WithCancel(server.HandleSignals(context.Background()))
	defer cancelFunc()

    // create server.GRPCServer with middlewares
	s, err := server.NewGRPC(server.Config{
		Port: grpcPort,
	}, server.WithGRPCServerOptions(GRPCMiddlewaresInterceptor))
	if err != nil {
		panic(err)
	}

    // register a grpc service desc and server implementation instance
	s.RegisterService(&commonv1.CommonService_ServiceDesc,
		common.New(Server),
	)

    // start serving
	go s.Serve()
    // wait for ctx context to be set done, by SIGINT or SIGTERM
	<-ctx.Done()

    // set a timer for graceful shutdown, if expired will force kill the server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()

	s.Shutdown(shutdownCtx)
```

### GRPC gateway

`GRPCGateway` can be used to add GRPC Gateway generated proxy handlers to `HTTPServer`

```go
    gw, err := server.NewGateway("", grpcClientPort)
	if err != nil {
		panic(err)
	}
    // Use the grpc gateway generated function to register http handlers
	gw.RegisterHandler(ctx, commonv1.RegisterCommonServiceHandlerFromEndpoint)

    // set gateway on HTTPServer with /api prefix
	s.SetGateway("/api", gw)
```

### Mux Server

`MuxServer` can be used to run GRPC and HTTP servers on same port on a host. Internally it uses [cmux](https://github.com/soheilhy/cmux) to route requests to specific server

```go
    // grpc middlewares
    var GRPCMiddlewaresInterceptor = grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	    grpc_recovery.UnaryServerInterceptor(),
	    grpc_ctxtags.UnaryServerInterceptor(),
	    grpc_prometheus.UnaryServerInterceptor,
	    grpc_zap.UnaryServerInterceptor(zap.NewExample()),
    ))

    // context to be Done when SIGINT or SIGTERM is received
	ctx, cancelFunc := context.WithCancel(server.HandleSignals(context.Background()))
	defer cancelFunc()

    // create server.MuxServer with grpc middlewares
	s, err := server.NewMux(server.Config{
		Port: muxPort,
	}, server.WithMuxGRPCServerOptions(GRPCMiddlewaresInterceptor))
	if err != nil {
		panic(err)
	}

    // use same port for grpc-gateway grpc client to proxy requests to
	grpcClientPort := muxPort
    // create, set handlers and set gateway on MuxServer
	gw, err := server.NewGateway("", grpcClientPort)
	if err != nil {
		panic(err)
	}
	gw.RegisterHandler(ctx, commonv1.RegisterCommonServiceHandlerFromEndpoint)
	s.SetGateway("/api", gw)

    // add additional http handlers on MuxServer
	s.RegisterHandler("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}))

    // register grpc service on MuxServer
	s.RegisterService(&commonv1.CommonService_ServiceDesc,
		common.New(Server),
	)

    // start serving
	go s.Serve()
    // wait for ctx context to be set done, by SIGINT or SIGTERM
	<-ctx.Done()
	// clean anything that needs to be closed etc like common server implementation etc
    // set a timer for graceful shutdown, if expired will force kill the server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()

	s.Shutdown(shutdownCtx)
```

### Example

For usage example have a look at this - [example](example/main.go).
