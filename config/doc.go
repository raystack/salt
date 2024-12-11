/*
Package config provides a flexible and extensible configuration management solution for Go applications.

It integrates configuration files, environment variables, command-line flags, and default values to populate
and validate user-defined structs.

Configuration Precedence:
The `Loader` merges configuration values from multiple sources in the following order of precedence (highest to lowest):
 1. Command-line flags: Defined using `pflag.FlagSet` and dynamically bound via `cmdx` tags.
 2. Environment variables: Dynamically bound to configuration keys, optionally prefixed using `WithEnvPrefix`.
 3. Configuration file: YAML configuration files specified via `WithFile`.
 4. Default values: Struct fields annotated with `default` tags are populated if no other source provides a value.

Defaults:
Default values are specified using the `default` struct tag. Fields annotated with `default` are populated
before any other source (flags, environment variables, or files).

Example:

	type Config struct {
	    ServerPort int    `mapstructure:"server.port" default:"8080"`
	    LogLevel   string `mapstructure:"log.level" default:"info"`
	}

In the absence of higher-priority sources, `ServerPort` will default to `8080` and `LogLevel` to `info`.

Validation:
Validation is performed using the `go-playground/validator` package. Fields annotated with `validate` tags
are validated after merging all configuration sources.

Example:

	type Config struct {
	    ServerPort int    `mapstructure:"server.port" validate:"required,min=1"`
	    LogLevel   string `mapstructure:"log.level" validate:"required,oneof=debug info warn error"`
	}

If validation fails, the `Load` method returns a detailed error indicating the invalid fields.

Annotations:
Configuration structs use the following struct tags to define behavior:
  - `mapstructure`: Maps YAML or environment variables to struct fields.
  - `default`: Provides fallback values for fields when no source overrides them.
  - `validate`: Ensures the final configuration meets application-specific requirements.

Example:

	type Config struct {
	    Server struct {
	        Port int    `mapstructure:"server.port" default:"8080" validate:"required,min=1"`
	        Host string `mapstructure:"server.host" default:"localhost" validate:"required"`
	    } `mapstructure:"server"`

	    LogLevel string `mapstructure:"log.level" default:"info" validate:"required,oneof=debug info warn error"`
	}

The `Loader` will merge all sources, apply defaults, and validate the result in a single call to `Load`.

Features:
  - Merges configurations from multiple sources: flags, environment variables, files, and defaults.
  - Supports nested structs with dynamic field mapping using `cmdx` tags.
  - Validates fields with constraints defined in `validate` tags.
  - Saves and views the final configuration in YAML or JSON formats.

Example Usage:

	type Config struct {
	    ServerPort int    `mapstructure:"server.port" cmdx:"server.port" default:"8080" validate:"required,min=1"`
	    LogLevel   string `mapstructure:"log.level" cmdx:"log.level" default:"info" validate:"required,oneof=debug info warn error"`
	}

	func main() {
	    flags := pflag.NewFlagSet("example", pflag.ExitOnError)
	    flags.Int("server.port", 8080, "Server port")
	    flags.String("log.level", "info", "Log level")

	    loader := config.NewLoader(
	        config.WithFile("./config.yaml"),
	        config.WithEnvPrefix("MYAPP"),
	        config.WithFlags(flags),
	    )

	    flags.Parse(os.Args[1:])

	    cfg := &Config{}
	    if err := loader.Load(cfg); err != nil {
	        log.Fatalf("Failed to load configuration: %v", err)
	    }

	    fmt.Printf("Configuration: %+v\n", cfg)
	}
*/
package config
