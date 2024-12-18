package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/mcuadros/go-defaults"
	"github.com/raystack/salt/config"
	"github.com/spf13/pflag"
)

type Config struct {
	Server struct {
		Port int    `mapstructure:"port" default:"8000" validate:"required,min=1" cmdx:"port"`
		Host string `mapstructure:"host" cmdx:"host"`
	} `mapstructure:"server" cmdx:"server"`
	LogLevel string `mapstructure:"log_level" cmdx:"log_level"`
}

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed to set environment variable %s: %v", key, err)
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Failed to unset environment variable %s: %v", key, err)
	}
}

func TestDefaultsAreApplied(t *testing.T) {
	cfg := &Config{}
	loader := config.NewLoader()

	loader.Load(cfg)
	if cfg.Server.Port != 8000 || cfg.Server.Host != "" {
		t.Errorf("Defaults were not applied: %+v", cfg)
	}
}

func TestEnvironmentVariableBinding(t *testing.T) {
	cfg := &Config{}
	loader := config.NewLoader()

	setEnv(t, "SERVER_PORT", "9090")
	setEnv(t, "SERVER_HOST", "localhost")
	defer unsetEnv(t, "SERVER_PORT")
	defer unsetEnv(t, "SERVER_HOST")

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected SERVER_PORT to be 9090, got %d", cfg.Server.Port)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected SERVER_HOST to be 'localhost', got %s", cfg.Server.Host)
	}
}

func TestConfigFileLoading(t *testing.T) {
	configFileContent := `
server:
  port: 8080
  host: example.com
log_level: debug
`
	configFilePath := "./test_config.yaml"
	os.WriteFile(configFilePath, []byte(configFileContent), 0644)
	defer os.Remove(configFilePath)

	cfg := &Config{}
	loader := config.NewLoader(config.WithFile(configFilePath))

	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected server.port to be 8080, got %d", cfg.Server.Port)
	}

	if cfg.Server.Host != "example.com" {
		t.Errorf("Expected server.host to be 'example.com', got %s", cfg.Server.Host)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("Expected log_level to be 'debug', got %s", cfg.LogLevel)
	}
}

func TestMissingConfigFile(t *testing.T) {
	cfg := &Config{}
	loader := config.NewLoader(config.WithFile("./nonexistent_config.yaml"))

	if err := loader.Load(cfg); err != nil {
		t.Errorf("Unexpected error for missing config file: %v", err)
	}
}

func TestInvalidConfigurationValidation(t *testing.T) {
	cfg := &Config{}

	setEnv(t, "SERVER_PORT", "0")
	loader := config.NewLoader()
	err := loader.Load(cfg)

	if err == nil {
		t.Fatalf("Expected validation error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid configuration") {
		t.Errorf("Expected validation error message, got: %v", err)
	}
}

func TestEnvOverrideConfig(t *testing.T) {
	// Create a temporary config file with values
	configFileContent := `
server:
  port: 8080
  host: "file-host.com"
log_level: "info"
`
	configFilePath := "./test_config.yaml"
	if err := os.WriteFile(configFilePath, []byte(configFileContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(configFilePath)

	// Set environment variables that should override file values
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_HOST", "env-host.com")
	defer os.Unsetenv("SERVER_PORT")
	defer os.Unsetenv("SERVER_HOST")

	// Define the config struct and loader
	cfg := &Config{}
	loader := config.NewLoader(config.WithFile(configFilePath))

	// Apply defaults
	defaults.SetDefaults(cfg)
	cfg.Server.Port = 3000 // Default value
	cfg.Server.Host = "default-host.com"
	cfg.LogLevel = "debug"

	// Load the configuration
	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate override order
	if cfg.Server.Port != 9090 {
		t.Errorf("Expected SERVER_PORT (env) to override file and defaults, got %d", cfg.Server.Port)
	}

	if cfg.Server.Host != "env-host.com" {
		t.Errorf("Expected SERVER_HOST (env) to override file and defaults, got %s", cfg.Server.Host)
	}

	if cfg.LogLevel != "info" {
		t.Errorf("Expected log_level from file to override defaults, got %s", cfg.LogLevel)
	}
}

func TestFlagsOverrideFileAndEnvVars(t *testing.T) {
	// Create a temporary config file
	configFileContent := `
server:
  port: 8080
  host: "file-host.com"
log_level: "info"
`
	configFilePath := "./test_config.yaml"
	if err := os.WriteFile(configFilePath, []byte(configFileContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(configFilePath)

	// Set environment variables
	setEnv(t, "SERVER_PORT", "9090")
	setEnv(t, "SERVER_HOST", "env-host.com")
	defer unsetEnv(t, "SERVER_PORT")
	defer unsetEnv(t, "SERVER_HOST")

	// Define flags
	flags := pflag.NewFlagSet("test", pflag.ExitOnError)
	flags.Int("server.port", 1000, "Server port")
	flags.String("server.host", "flag-host.com", "Server host")
	flags.String("log_level", "debug", "Log level")

	// Parse command-line flags (simulate CLI args)
	flags.Parse([]string{"--server.port=1234", "--server.host=flag-host.com", "--log_level=trace"})

	// Initialize Loader with flag set
	loader := config.NewLoader(
		config.WithFile(configFilePath),
		config.WithFlags(flags),
	)

	// Load configuration into the struct
	cfg := &Config{}
	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Assert final values
	if cfg.Server.Port != 1234 {
		t.Errorf("Expected Server.Port to be 1234 from flags, got %d", cfg.Server.Port)
	}
	if cfg.Server.Host != "flag-host.com" {
		t.Errorf("Expected Server.Host to be 'flag-host.com' from flags, got %s", cfg.Server.Host)
	}
	if cfg.LogLevel != "trace" {
		t.Errorf("Expected LogLevel to be 'trace' from flags, got %s", cfg.LogLevel)
	}
}

func TestMissingFlags(t *testing.T) {
	// Define flags
	flags := pflag.NewFlagSet("test", pflag.ExitOnError)
	flags.Int("server.port", 8080, "Server port")

	// Initialize Loader with the incomplete flag set
	loader := config.NewLoader(config.WithFlags(flags))

	// Load configuration into the struct
	cfg := &Config{}
	err := loader.Load(cfg)

	// Expect an error because `server.host` and `log.level` flags are missing
	if err == nil {
		t.Fatal("Expected an error due to missing flags, but got nil")
	}
	if !strings.Contains(err.Error(), "missing flag for tag") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestNestedStructFlags(t *testing.T) {
	// Define flags
	flags := pflag.NewFlagSet("test", pflag.ExitOnError)
	flags.Int("server.port", 8080, "Server port")
	flags.String("server.host", "localhost", "Server host")
	flags.String("log_level", "debug", "Log level")

	// Initialize Loader with the flag set
	loader := config.NewLoader(config.WithFlags(flags))

	// Parse flags
	flags.Parse([]string{"--server.port=1234", "--server.host=nested-host.com"})

	// Load configuration into the struct
	cfg := &Config{}
	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Assert nested struct values
	if cfg.Server.Port != 1234 {
		t.Errorf("Expected Server.Port to be 1234, got %d", cfg.Server.Port)
	}
	if cfg.Server.Host != "nested-host.com" {
		t.Errorf("Expected Server.Host to be 'nested-host.com', got %s", cfg.Server.Host)
	}
}

func TestFlagsOnly(t *testing.T) {
	// Define flags
	flags := pflag.NewFlagSet("test", pflag.ExitOnError)
	flags.Int("server.port", 8080, "Server port")
	flags.String("server.host", "localhost", "Server host")
	flags.String("log_level", "info", "Log level")

	// Parse flags
	flags.Parse([]string{"--server.port=9000", "--server.host=flag-only-host", "--log_level=warn"})

	// Initialize Loader with the flag set
	loader := config.NewLoader(config.WithFlags(flags))

	// Load configuration into the struct
	cfg := &Config{}
	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Assert values from flags
	if cfg.Server.Port != 9000 {
		t.Errorf("Expected Server.Port to be 9000, got %d", cfg.Server.Port)
	}
	if cfg.Server.Host != "flag-only-host" {
		t.Errorf("Expected Server.Host to be 'flag-only-host', got %s", cfg.Server.Host)
	}
	if cfg.LogLevel != "warn" {
		t.Errorf("Expected LogLevel to be 'warn', got %s", cfg.LogLevel)
	}
}

func TestLoader_SetAndGetSingleKey(t *testing.T) {
	loader := config.NewLoader()

	// Set a value
	loader.Set("test.key", "testValue")
	loader.Set("test.nested.key", 42)

	value1 := loader.Get("test.key")
	if value1 != "testValue" {
		t.Errorf("Expected value for 'test.key' to be 'testValue', got %v", value1)
	}

	value2 := loader.Get("test.nested.key")
	// Assert the value
	if value2 != 42 {
		t.Errorf("Expected value for 'test.nested.key' to be 42, got %v", value2)
	}
}

func TestLoader_GetKeyFromConfigFile(t *testing.T) {
	// Create a temporary configuration file
	configFileContent := `
server:
  port: 8080
  host: "example.com"
log_level: "info"
`
	configFilePath := "./test_config.yaml"
	if err := os.WriteFile(configFilePath, []byte(configFileContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(configFilePath)

	// Initialize Loader with the config file
	loader := config.NewLoader(config.WithFile(configFilePath))

	// Load the configuration
	cfg := &Config{}
	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Test retrieving keys using Get
	if loader.Get("server.port") != 8080 {
		t.Errorf("Expected 'server.port' to be 8080, got %v", loader.Get("server.port"))
	}

	if loader.Get("server.host") != "example.com" {
		t.Errorf("Expected 'server.host' to be 'example.com', got %v", loader.Get("server.host"))
	}

	if loader.Get("log_level") != "info" {
		t.Errorf("Expected 'log_level' to be 'info', got %v", loader.Get("log_level"))
	}

	// Test retrieving a key that doesn't exist
	if loader.Get("nonexistent.key") != nil {
		t.Errorf("Expected 'nonexistent.key' to be nil, got %v", loader.Get("nonexistent.key"))
	}
}

func TestLoader_SetDoesNotPersistToFile(t *testing.T) {
	// Create a temporary configuration file
	configFileContent := `
server:
  port: 8080
  host: "example.com"
log_level: "info"
`
	configFilePath := "./test_config.yaml"
	if err := os.WriteFile(configFilePath, []byte(configFileContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(configFilePath)

	// Initialize Loader with the config file
	loader := config.NewLoader(config.WithFile(configFilePath))

	// Load the configuration
	cfg := &Config{}
	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Set a new value for a key
	loader.Set("server.port", 9090)

	// Verify the in-memory value has been updated
	if loader.Get("server.port") != 9090 {
		t.Errorf("Expected 'server.port' to be 9090 in memory, got %v", loader.Get("server.port"))
	}

	// Reload the configuration from the file
	newCfg := &Config{}
	newLoader := config.NewLoader(config.WithFile(configFilePath))
	if err := newLoader.Load(newCfg); err != nil {
		t.Fatalf("Failed to reload configuration: %v", err)
	}

	// Verify the original file values are unchanged
	if newLoader.Get("server.port") != 8080 {
		t.Errorf("Expected 'server.port' to remain 8080 in file, got %v", newLoader.Get("server.port"))
	}
	if newLoader.Get("server.host") != "example.com" {
		t.Errorf("Expected 'server.host' to remain 'example.com', got %v", newLoader.Get("server.host"))
	}
	if newLoader.Get("log_level") != "info" {
		t.Errorf("Expected 'log_level' to remain 'info', got %v", newLoader.Get("log_level"))
	}
}

func TestLoader_SetAndSavePersistToFile(t *testing.T) {
	// Create a temporary configuration file
	initialConfig := `
server:
  port: 8080
  host: "example.com"
log_level: "info"
`
	configFilePath := "./test_config.yaml"
	if err := os.WriteFile(configFilePath, []byte(initialConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(configFilePath)

	// Initialize Loader with the config file
	loader := config.NewLoader(config.WithFile(configFilePath))

	// Load the configuration
	cfg := &Config{}
	if err := loader.Load(cfg); err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Set new values for configuration keys
	loader.Set("server.port", 9090)
	loader.Set("server.host", "new-host.com")
	loader.Set("log_level", "debug")

	// Save the updated configuration to the file
	if err := loader.Save(); err != nil {
		t.Fatalf("Failed to save configuration: %v", err)
	}

	// Reload the configuration from the file into a new Loader
	newLoader := config.NewLoader(config.WithFile(configFilePath))
	newCfg := &Config{}
	if err := newLoader.Load(newCfg); err != nil {
		t.Fatalf("Failed to reload configuration: %v", err)
	}

	// Verify that the changes were persisted
	if newLoader.Get("server.port") != 9090 {
		t.Errorf("Expected 'server.port' to be 9090 in file, got %v", newLoader.Get("server.port"))
	}

	if newLoader.Get("server.host") != "new-host.com" {
		t.Errorf("Expected 'server.host' to be 'new-host.com' in file, got %v", newLoader.Get("server.host"))
	}

	if newLoader.Get("log_level") != "debug" {
		t.Errorf("Expected 'log_level' to be 'debug' in file, got %v", newLoader.Get("log_level"))
	}
}
