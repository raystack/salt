package cmdx

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/raystack/salt/config"
)

// Environment variables for configuration paths
const (
	RaystackConfigDirEnv = "RAYSTACK_CONFIG_DIR"
	XDGConfigHomeEnv     = "XDG_CONFIG_HOME"
	AppDataEnv           = "AppData"
)

// ConfigLoaderOpt defines a functional option for configuring the Config object.
type ConfigLoaderOpt func(c *Config)

// WithFlags binds command-line flags to configuration values.
func WithFlags(pfs *pflag.FlagSet) ConfigLoaderOpt {
	return func(c *Config) {
		c.boundFlags = pfs
	}
}

// WithLoaderOptions adds custom loader options for configuration loading.
func WithLoaderOptions(opts ...config.LoaderOption) ConfigLoaderOpt {
	return func(c *Config) {
		c.loaderOpts = append(c.loaderOpts, opts...)
	}
}

// SetConfig initializes a new Config object for the specified application.
func SetConfig(app string) *Config {
	return &Config{
		filename: configFile(app),
	}
}

// Config manages the application's configuration file and related operations.
type Config struct {
	filename   string
	boundFlags *pflag.FlagSet
	loaderOpts []config.LoaderOption
}

// File returns the path to the configuration file.
func (c *Config) File() string {
	return c.filename
}

// Defaults populates the given configuration struct with default values.
func (c *Config) Defaults(cfg interface{}) {
	defaults.SetDefaults(cfg)
}

// Init initializes the configuration file with default values.
func (c *Config) Init(cfg interface{}) error {
	defaults.SetDefaults(cfg)

	if fileExists(c.filename) {
		return errors.New("configuration file already exists")
	}

	return c.Write(cfg)
}

// Read reads the content of the configuration file as a string.
func (c *Config) Read() (string, error) {
	data, err := os.ReadFile(c.filename)
	return string(data), err
}

// Write writes the given configuration struct to the configuration file in YAML format.
func (c *Config) Write(cfg interface{}) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if _, err := os.Stat(c.filename); os.IsNotExist(err) {
		_ = os.MkdirAll(configDir("raystack"), 0700)
	}

	if err := os.WriteFile(c.filename, data, 0655); err != nil {
		return err
	}
	return nil
}

// Load loads the configuration from the file and applies the provided loader options.
func (c *Config) Load(cfg interface{}, opts ...ConfigLoaderOpt) error {
	for _, opt := range opts {
		opt(c)
	}

	loaderOpts := []config.LoaderOption{config.WithFile(c.filename)}

	if c.boundFlags != nil {
		loaderOpts = append(loaderOpts, config.WithBindPFlags(c.boundFlags, cfg))
	}
	loaderOpts = append(loaderOpts, c.loaderOpts...)

	loader := config.NewLoader(loaderOpts...)

	return loader.Load(cfg)
}

// configFile determines the full path to the configuration file for the application.
func configFile(app string) string {
	filename := app + ".yml"
	return filepath.Join(configDir("raystack"), filename)
}

// configDir determines the appropriate directory for storing configuration files.
func configDir(root string) string {
	var path string
	if env := os.Getenv(RaystackConfigDirEnv); env != "" {
		path = env
	} else if env := os.Getenv(XDGConfigHomeEnv); env != "" {
		path = filepath.Join(env, root)
	} else if runtime.GOOS == "windows" {
		if env := os.Getenv(AppDataEnv); env != "" {
			path = filepath.Join(env, root)
		}
	} else {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".config", root)
	}

	if !dirExists(path) {
		_ = os.MkdirAll(filepath.Dir(path), 0755)
	}

	return path
}
