# Migration Guide

This guide covers migrating from the previous salt version to the new structure.

## Go version

Update `go.mod` to require Go 1.24:

```
go 1.24
```

## Packages removed

| Removed | Replacement |
|---------|-------------|
| `observability/logger` | Use `*slog.Logger` from `log/slog` directly |
| `observability/otelgrpc` | Use `connectrpc.com/otelconnect` |
| `observability/otelhttpclient` | Use `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` |
| `server/mux` | Use `github.com/raystack/salt/server` |
| `db` | Use your preferred DB library (sqlx, pgx, gorm) directly |
| `auth/oidc` | Planned as complete CLI auth solution (#86) |
| `auth/audit` | Planned with standardized schema (#87) |
| `testing/dockertestx` | Use `ory/dockertest/v3` directly |

## Packages moved

| Old | New |
|-----|-----|
| `observability` | `telemetry` |
| `cli/terminator` | `cli/terminal` |
| `cli/prompter` | `cli/prompt` |
| `cli/releaser` | `cli/version` |

## Logger

The custom `logger.Logger` interface and all backends (Zap, Logrus, Slog, Noop) are removed. Use `*slog.Logger` from the Go standard library directly.

```go
// Before
import "github.com/raystack/salt/observability/logger"
l := logger.NewZap()
l := logger.NewLogrus()
l := logger.NewNoop()

// After
import "log/slog"
l := slog.Default()
l := slog.New(slog.NewJSONHandler(os.Stderr, nil))
l := slog.New(slog.DiscardHandler) // noop
```

All salt packages that previously accepted `logger.Logger` now accept `*slog.Logger`.

## Server

The dual-port `server/mux` package is replaced by a single-port `server` package with h2c support.

```go
// Before
import "github.com/raystack/salt/server/mux"
mux.Serve(ctx,
    mux.WithHTTPTarget(":8080", httpServer),
    mux.WithGRPCTarget(":8081", grpcServer),
)

// After
import "github.com/raystack/salt/server"
srv := server.New(
    server.WithAddr(":8080"),
    server.WithHandler("/api/", connectHandler),
)
srv.Start(ctx)
```

H2C and health check (`/ping`) are enabled by default. Use `server.WithoutH2C()` or `server.WithHealthCheck("")` to disable.

## App bootstrap

New `app.Run()` for service bootstrap:

```go
import "github.com/raystack/salt/app"

app.Run(
    app.WithConfig(&cfg, config.WithFile("config.yaml")),
    app.WithLogger(slog.Default()),
    app.WithHTTPMiddleware(middleware.DefaultHTTP(slog.Default())),
    app.WithHandler("/api/", handler),
    app.WithAddr(cfg.Addr),
)
```

HTTP middleware is explicit — use `middleware.DefaultHTTP(logger)` for the standard chain or compose your own. Database connections are managed via `app.WithOnStart` / `app.WithOnStop` hooks.

## CLI bootstrap

New `cli.Execute()` for CLI applications:

```go
// Before
rootCmd := &cobra.Command{Use: "frontier", Short: "identity management"}
mgr := commander.New(rootCmd, commander.WithTopics(topics))
mgr.Init()
rootCmd.AddCommand(serverCmd, configCmd)
rootCmd.Execute()

// After
import "github.com/raystack/salt/cli"

cli.Execute(
    cli.Name("frontier"),
    cli.Description("identity management"),
    cli.Version("0.1.0", "raystack/frontier"),
    cli.Commands(serverCmd, configCmd),
    cli.Topics(topics...),
)
```

Access shared output and prompting in commands:

```go
func newListCmd() *cobra.Command {
    return &cobra.Command{
        Use: "list",
        RunE: func(cmd *cobra.Command, args []string) error {
            out := cli.Output(cmd)
            out.Table(rows)
            return nil
        },
    }
}
```

## Printer

Package-level functions replaced by `Output` type:

```go
// Before
printer.Success("done")
printer.Table(os.Stdout, rows)
printer.JSON(data)
spinner := printer.Spin("loading")

// After
out := printer.NewOutput(os.Stdout)
// or inside a command: out := cli.Output(cmd)

out.Success("done")
out.Table(rows)
out.JSON(data)
spinner := out.Spin("loading")
```

Color formatting functions remain as package-level helpers returning styled strings:

```go
printer.Green("text")
printer.Greenf("count: %d", n)
printer.Icon("success") // ✔
printer.Italic("note")
```

## Telemetry

```go
// Before
import "github.com/raystack/salt/observability"
observability.Init(ctx, cfg, logger)

// After
import "github.com/raystack/salt/telemetry"
telemetry.Init(ctx, cfg, slogLogger)
```

## Middleware

New package for ConnectRPC and HTTP middleware:

```go
import "github.com/raystack/salt/middleware"

// Connect interceptors for your handler
interceptors := middleware.Default(slog.Default())
handler := myv1connect.NewServiceHandler(svc, connect.WithInterceptors(interceptors...))

// HTTP middleware
httpMW := middleware.DefaultHTTP(slog.Default())
```

## Config

```go
// Import path for validator changed
// Before: "github.com/go-playground/validator"
// After:  "github.com/go-playground/validator/v10"

// If you imported go-defaults directly:
// Before: "github.com/mcuadros/go-defaults"
// After:  "github.com/creasty/defaults"
// API change: defaults.SetDefaults(cfg) → defaults.Set(cfg)
```

The config package no longer prints warnings to stdout when a config file is missing.

## Dependency changes

| Removed (direct) | Replacement |
|-------------------|-------------|
| `go.uber.org/zap` | `log/slog` (stdlib) |
| `sirupsen/logrus` | `log/slog` (stdlib) |
| `AlecAivazis/survey/v2` | `charmbracelet/huh` |
| `olekukonko/tablewriter` | `text/tabwriter` (stdlib) |
| `oklog/run` | Removed with `server/mux` |
| `cli/safeexec` | `exec.LookPath` (stdlib) |
| `pkg/errors` | `fmt.Errorf` with `%w` (stdlib) |
| `mcuadros/go-defaults` | `creasty/defaults` |
| `go-playground/validator` v9 | `go-playground/validator/v10` |
| `jmoiron/sqlx` | Use directly if needed |
| `golang-migrate` | Use directly if needed |
| `ory/dockertest` | Use directly if needed |

| Added | Purpose |
|-------|---------|
| `connectrpc.com/connect` | Middleware interceptors |
| `charmbracelet/huh` | Interactive prompts |
| `creasty/defaults` | Struct default values |

| Upgraded | From → To |
|----------|-----------|
| `go-playground/validator` | v9 → v10 |
| `charmbracelet/glamour` | v0.3 → v1.0.0 |
| `muesli/termenv` | v0.11 → v0.16.0 |
| `briandowns/spinner` | v1.18 → v1.23.2 |
| `schollz/progressbar` | v3.8 → v3.19.0 |
| `mattn/go-isatty` | v0.0.19 → v0.0.21 |
