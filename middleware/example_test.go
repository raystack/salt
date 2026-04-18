package middleware_test

import (
	"log/slog"
	"net/http"

	"github.com/raystack/salt/middleware"
	"github.com/raystack/salt/middleware/cors"
	"github.com/raystack/salt/middleware/recovery"
	"github.com/raystack/salt/middleware/requestid"
	"github.com/raystack/salt/middleware/requestlog"
)

func ExampleDefault() {
	logger := slog.Default()

	// Use Default() for the standard Connect interceptor chain.
	// Apply to your ConnectRPC handler:
	//
	//   path, handler := myv1connect.NewServiceHandler(svc,
	//       connect.WithInterceptors(middleware.Default(logger)...),
	//   )
	_ = middleware.Default(logger)
}

func ExampleDefaultHTTP() {
	logger := slog.Default()

	// Use DefaultHTTP() for the standard HTTP middleware chain.
	// Apply to app or server:
	//
	//   app.WithHTTPMiddleware(middleware.DefaultHTTP(logger))
	handler := middleware.DefaultHTTP(logger)(http.NotFoundHandler())
	_ = handler
}

func ExampleChainHTTP() {
	logger := slog.Default()

	// Compose a custom HTTP middleware chain.
	chain := middleware.ChainHTTP(
		recovery.HTTPMiddleware(recovery.WithLogger(logger)),
		requestid.HTTPMiddleware(),
		requestlog.HTTPMiddleware(requestlog.WithLogger(logger)),
		cors.Middleware(cors.WithAllowedOrigins("https://myapp.com")),
	)

	handler := chain(http.NotFoundHandler())
	_ = handler
}
