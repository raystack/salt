# Salt

[![Go Reference](https://pkg.go.dev/badge/github.com/raystack/salt.svg)](https://pkg.go.dev/github.com/raystack/salt)
![test](https://github.com/raystack/salt/actions/workflows/test.yaml/badge.svg)
![lint](https://github.com/raystack/salt/actions/workflows/lint.yaml/badge.svg)

The standard way to build raystack services and CLIs.

Salt provides `app.Run()` for services and `cli.Init()` / `cli.Execute()` for command-line tools, along with the building blocks they use: configuration, middleware, terminal output, and more.

## Quick start

### Service

```go
package main

import (
    "log/slog"

    "github.com/raystack/salt/app"
    "github.com/raystack/salt/config"
    "github.com/raystack/salt/middleware"
)

func main() {
    var cfg Config

    app.Run(
        app.WithConfig(&cfg, config.WithFile("config.yaml")),
        app.WithLogger(slog.Default()),
        app.WithHTTPMiddleware(middleware.DefaultHTTP(slog.Default())),
        app.WithHandler("/api/", apiHandler),
        app.WithAddr(cfg.Addr),
    )
}
```

H2C and health check at `/ping` enabled by default. HTTP middleware is explicit — you choose what runs.

### CLI

```go
package main

import (
    "github.com/raystack/salt/cli"
    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{Use: "frontier", Short: "identity management"}
    rootCmd.PersistentFlags().String("host", "", "API server host")
    rootCmd.AddCommand(serverCmd, userCmd)

    cli.Init(rootCmd,
        cli.Version("0.1.0", "raystack/frontier"),
    )

    cli.Execute(rootCmd)
}
```

`Init` adds help, shell completion, reference docs, and silences Cobra's default error output. `Execute` runs the command and handles all errors with proper exit codes. Commands access shared I/O via `cli.IO(cmd)`, or the convenience helpers `cli.Output(cmd)` and `cli.Prompter(cmd)`. Use `cli.Test()` in tests for captured, deterministic output.

## Installation

```
go get github.com/raystack/salt
```

Requires Go 1.24+.

## Packages

### Bootstrap

| Package | Description |
|---------|-------------|
| [`app`](app/) | Service lifecycle — config, logger, telemetry, server, graceful shutdown |
| [`cli`](cli/) | CLI lifecycle — init, execute, error handling, help, completion, version check |

### Server & Middleware

| Package | Description |
|---------|-------------|
| [`server`](server/) | HTTP server with h2c, health checks, graceful shutdown |
| [`server/spa`](server/spa/) | Single-page application static file handler |
| [`middleware`](middleware/) | Connect interceptors and HTTP middleware |
| [`middleware/recovery`](middleware/recovery/) | Panic recovery |
| [`middleware/requestid`](middleware/requestid/) | X-Request-ID propagation |
| [`middleware/requestlog`](middleware/requestlog/) | Request logging with duration |
| [`middleware/errorz`](middleware/errorz/) | Error sanitization for clients |
| [`middleware/cors`](middleware/cors/) | CORS with Connect defaults |

### CLI

| Package | Description |
|---------|-------------|
| [`cli/commander`](cli/commander/) | Cobra enhancements — help layout, completion, reference docs, hooks |
| [`cli/printer`](cli/printer/) | Terminal output — styled text, tables, JSON/YAML, spinners, progress bars, markdown |
| [`cli/prompt`](cli/prompt/) | Interactive prompts — select, multi-select, input, confirm |
| [`cli/terminal`](cli/terminal/) | Terminal utilities — TTY detection, browser, pager |
| [`cli/version`](cli/version/) | Version checking against GitHub releases |

### Infrastructure

| Package | Description |
|---------|-------------|
| [`config`](config/) | Configuration from files, env vars, flags, and struct defaults |
| [`telemetry`](telemetry/) | OpenTelemetry initialization — traces and metrics via OTLP |

### Data

| Package | Description |
|---------|-------------|
| [`data/rql`](data/rql/) | REST query language — filters, pagination, sorting, search |
| [`data/jsondiff`](data/jsondiff/) | JSON document diffing and reconstruction |

## Logging

Salt uses `*slog.Logger` from the Go standard library. No custom logger interface — pass `slog.Default()` or any `*slog.Logger` to packages that need it.

```go
// Production
logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

// Tests
logger := slog.New(slog.DiscardHandler)
```

## Migration

See [MIGRATION.md](MIGRATION.md) for upgrading from previous versions.

## License

Apache License 2.0
