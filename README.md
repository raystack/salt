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

## Pacakages

### Configuration and Environment
- **`config`**  
  Utilities for managing application configurations using environment variables, files, or defaults.

### CLI Utilities
- **`cli/cmdx`**  
  Command execution and management tools.

- **`cli/printer`**  
  Utilities for formatting and printing output to the terminal.

- **`cli/prompt`**  
  Interactive CLI prompts for user input.

- **`cli/terminal`**  
  Terminal utilities for colors, cursor management, and formatting.

- **`cli/version`**  
  Utilities for displaying and managing CLI tool versions.

### Authentication and Security
- **`auth/oidc`**  
  Helpers for integrating OpenID Connect authentication flows.

- **`auth/audit`**  
  Auditing tools for tracking security events and compliance.

### Server and Infrastructure
- **`server`**  
  Utilities for setting up and managing HTTP or RPC servers.

- **`db`**  
  Helpers for database connections, migrations, and query execution.

- **`telemetry`**  
  Observability tools for capturing application metrics and traces.

### Development and Testing
- **`dockertestx`**  
  Tools for creating and managing Docker-based testing environments.

### Utilities
- **`log`**  
  Simplified logging utilities for structured and unstructured log messages.

- **`utils`**  
  General-purpose utility functions for common programming tasks.
