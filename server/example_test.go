package server_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/raystack/salt/server"
)

func ExampleNew() {
	srv := server.New(
		server.WithAddr(":8080"),
		server.WithHandler("/hello", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, "world")
		})),
		server.WithLogger(slog.Default()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv.Start(ctx)
}

func ExampleNew_withTimeouts() {
	srv := server.New(
		server.WithAddr(":8080"),
		server.WithReadTimeout(60*time.Second),
		server.WithWriteTimeout(60*time.Second),
		server.WithIdleTimeout(120*time.Second),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv.Start(ctx)
}
