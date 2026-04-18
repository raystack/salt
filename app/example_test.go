package app_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/raystack/salt/app"
	"github.com/raystack/salt/middleware"
)

func ExampleRun() {
	app.Run(
		app.WithLogger(slog.Default()),
		app.WithHTTPMiddleware(middleware.DefaultHTTP(slog.Default())),
		app.WithHandler("/hello", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, "world")
		})),
		app.WithAddr(":8080"),
	)
}

func ExampleNew() {
	a, err := app.New(
		app.WithLogger(slog.Default()),
		app.WithAddr(":8080"),
	)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.Start(ctx)
}
