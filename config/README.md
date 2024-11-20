# Config Package

The `config` package simplifies configuration management in Go projects by integrating multiple sources (files, environment variables, flags) and decoding them into structured Go objects. It provides defaults, overrides, and extensibility for various use cases.

## Features

- **Flexible Sources**: Load configuration from files (YAML, JSON, etc.), environment variables, and command-line flags.
- **Defaults and Overrides**: Apply default values, with support for environment variable and flag-based overrides.
- **Powerful Decoding**: Decode nested structures, custom types, and JSON strings into Go structs.
- **Customizable**: Configure behavior with options like file paths, environment variable prefixes, and key replacers.

## Installation

Install the package using:

```bash
go get github.com/raystack/salt/config
```

## Usage

### 1. Basic Configuration Loading

Define your configuration struct:

```go
type Config struct {
    Host string `yaml:"host" default:"localhost"`
    Port int    `yaml:"port" default:"8080"`
}
```

Load the configuration:

```go
package main

import (
    "fmt"
    "github.com/raystack/salt/config"
)

func main() {
    var cfg Config
    loader := config.NewLoader(
        config.WithFile("config.yaml"),
        config.WithEnvPrefix("MYAPP"),
    )
    if err := loader.Load(&cfg); err != nil {
        fmt.Println("Error loading configuration:", err)
        return
    }
    fmt.Printf("Configuration: %+v
", cfg)
}
```

### 2. Using Command-Line Flags

Define your flags and bind them to the configuration struct:

```go
import (
    "github.com/spf13/pflag"
    "github.com/yourusername/config"
)

type Config struct {
    Host string `yaml:"host" cmdx:"host"`
    Port int    `yaml:"port" cmdx:"port"`
}

func main() {
    var cfg Config

    pflags := pflag.NewFlagSet("example", pflag.ExitOnError)
    pflags.String("host", "localhost", "Server host")
    pflags.Int("port", 8080, "Server port")
    pflags.Parse([]string{"--host", "127.0.0.1", "--port", "9090"})

    loader := config.NewLoader(
        config.WithFile("config.yaml"),
        config.WithBindPFlags(pflags, &cfg),
    )
    if err := loader.Load(&cfg); err != nil {
        fmt.Println("Error loading configuration:", err)
        return
    }
    fmt.Printf("Configuration: %+v
", cfg)
}
```

### 3. Environment Variable Overrides

Override configuration values using environment variables:

```go
loader := config.NewLoader(
    config.WithEnvPrefix("MYAPP"),
)
```

Set environment variables like `MYAPP_HOST` or `MYAPP_PORT` to override `host` and `port` values.

## Advanced Features

- **Custom Decode Hooks**: Parse custom formats like JSON strings into maps.
- **Error Handling**: Handles missing files gracefully and provides detailed error messages.
- **Multiple Config Paths**: Search for configuration files in multiple directories using `WithPath`.

## API Reference

### Loader Options

- `WithFile(file string)`: Set the explicit file path for the configuration file.
- `WithPath(path string)`: Add directories to search for configuration files.
- `WithName(name string)`: Set the name of the configuration file (without extension).
- `WithType(type string)`: Set the file type (e.g., "json", "yaml").
- `WithEnvPrefix(prefix string)`: Set a prefix for environment variables.
- `WithBindPFlags(flagSet *pflag.FlagSet, config interface{})`: Bind CLI flags to configuration fields.

### Custom Hooks

- `StringToJsonFunc()`: Decode JSON strings into maps or other structures.

### Struct Tags

- `yaml`: Maps struct fields to YAML keys.
- `default`: Specifies default values for struct fields.
- `cmdx`: Binds struct fields to command-line flags.
