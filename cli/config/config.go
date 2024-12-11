package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mcuadros/go-defaults"
	"github.com/raystack/salt/config"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// Config represents the configuration structure.
type Config struct {
	path  string
	flags *pflag.FlagSet
}

// New creates a new Config instance for the given application.
func New(app string, opts ...Opts) (*Config, error) {
	filePath, err := getConfigFilePath(app)
	if err != nil {
		return nil, fmt.Errorf("failed to determine config file path: %w", err)
	}

	cfg := &Config{path: filePath}
	for _, opt := range opts {
		opt(cfg)
	}

	return cfg, nil
}

// Opts defines a functional option for configuring the Config object.
type Opts func(c *Config)

// WithFlags binds command-line flags to configuration values.
func WithFlags(pfs *pflag.FlagSet) Opts {
	return func(c *Config) {
		c.flags = pfs
	}
}

// Load reads the configuration file into the Config's Data map.
func (c *Config) Load(cfg interface{}) error {
	loaderOpts := []config.Option{config.WithFile(c.path)}

	if c.flags != nil {
		loaderOpts = append(loaderOpts, config.WithFlags(c.flags))
	}

	loader := config.NewLoader(loaderOpts...)
	return loader.Load(cfg)
}

// Init initializes the configuration file with default values.
func (c *Config) Init(cfg interface{}) error {
	defaults.SetDefaults(cfg)

	if fileExists(c.path) {
		return errors.New("configuration file already exists")
	}

	return c.Write(cfg)
}

// Read reads the content of the configuration file as a string.
func (c *Config) Read() (string, error) {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return "", fmt.Errorf("failed to read configuration file: %w", err)
	}
	return string(data), nil
}

// Write writes the given struct to the configuration file in YAML format.
func (c *Config) Write(cfg interface{}) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := ensureDir(filepath.Dir(c.path)); err != nil {
		return err
	}

	if err := os.WriteFile(c.path, data, 0655); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}
	return nil
}

// getConfigFile determines the full path to the configuration file for the application.
func getConfigFilePath(app string) (string, error) {
	dirPath := getConfigDir("raystack")
	if err := ensureDir(dirPath); err != nil {
		return "", err
	}
	return filepath.Join(dirPath, app+".yml"), nil
}

// getConfigDir determines the directory for storing configurations.
func getConfigDir(root string) string {
	switch {
	case envSet("RAYSTACK_CONFIG_DIR"):
		return filepath.Join(os.Getenv("RAYSTACK_CONFIG_DIR"), root)
	case envSet("XDG_CONFIG_HOME"):
		return filepath.Join(os.Getenv("XDG_CONFIG_HOME"), root)
	case runtime.GOOS == "windows" && envSet("APPDATA"):
		return filepath.Join(os.Getenv("APPDATA"), root)
	default:
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", root)
	}
}

// ensureDir ensures that the given directory exists.
func ensureDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, err)
	}
	return nil
}

// envSet checks if an environment variable is set and non-empty.
func envSet(key string) bool {
	return os.Getenv(key) != ""
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
