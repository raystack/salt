# salt

[![GoDoc reference](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/raystack/salt)
![test workflow](https://github.com/raystack/salt/actions/workflows/test.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/raystack/salt)](https://goreportcard.com/report/github.com/raystack/salt)

Salt is a Golang utility library offering a variety of packages to simplify and enhance application development. It provides modular and reusable components for common tasks, including configuration management, CLI utilities, authentication, logging, and more.

## Installation

To use, run the following command:

```
go get github.com/raystack/salt
```

## Packages

### Configuration
- **`config`**
  Utilities for managing application configurations using environment variables, files, or defaults.

### CLI Utilities
- **`cli/commander`**
  Command execution, completion, help topics, and management tools.

- **`cli/printer`**
  Utilities for formatting and printing output to the terminal.

- **`cli/prompter`**
  Interactive CLI prompts for user input.

- **`cli/terminator`**
  Terminal utilities for browser, pager, and brew helpers.

- **`cli/releaser`**
  Utilities for displaying and managing CLI tool versions.

### Authentication and Security
- **`auth/oidc`**
  Helpers for integrating OpenID Connect authentication flows.

- **`auth/audit`**
  Auditing tools for tracking security events and compliance.

### Server and Infrastructure
- **`server/mux`**
  gRPC-gateway multiplexer for serving gRPC and HTTP on a single port.

- **`server/spa`**
  Single-page application static file handler.

- **`db`**
  Helpers for database connections, migrations, and query execution.

### Observability
- **`observability`**
  OpenTelemetry initialization, metrics, and tracing setup.

- **`observability/logger`**
  Structured logging with Zap and Logrus adapters.

- **`observability/otelgrpc`**
  OpenTelemetry gRPC client interceptors for metrics.

- **`observability/otelhttpclient`**
  OpenTelemetry HTTP client transport for metrics.

### Data Utilities
- **`data/rql`**
  REST query language parser for filters, pagination, sorting, and search.

- **`data/jsondiff`**
  JSON document diffing and reconstruction.

### Development and Testing
- **`testing/dockertestx`**
  Docker-based test environment helpers for Postgres, Minio, SpiceDB, and more.
