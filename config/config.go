package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/jeremywohl/flatten"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strings"
)

// Loader is responsible for managing configuration
type Loader struct {
	v     *viper.Viper
	flags *pflag.FlagSet
}

// Option defines a functional option for configuring the Loader.
type Option func(c *Loader)

// NewLoader creates a new Loader instance with the provided options.
// It initializes Viper with defaults for YAML configuration files and environment variable handling.
//
// Example:
//
//	loader := config.NewLoader(
//	    config.WithFile("./config.yaml"),
//	    config.WithEnvPrefix("MYAPP"),
//	)
func NewLoader(options ...Option) *Loader {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	loader := &Loader{v: v}
	for _, opt := range options {
		opt(loader)
	}
	return loader
}

// WithFile specifies the configuration file to use.
func WithFile(configFilePath string) Option {
	return func(l *Loader) {
		l.v.SetConfigFile(configFilePath)
	}
}

// WithEnvPrefix specifies a prefix for ENV variables.
func WithEnvPrefix(prefix string) Option {
	return func(l *Loader) {
		l.v.SetEnvPrefix(prefix)
	}
}

// WithFlags specifies a command-line flag set to bind dynamically based on `cmdx` tags.
func WithFlags(flags *pflag.FlagSet) Option {
	return func(l *Loader) {
		l.flags = flags
	}
}

// Load reads the configuration from the file, environment variables, and command-line flags,
// and merges them into the provided configuration struct. It validates the configuration
// using struct tags.
//
// The priority order is:
//  1. Command-line flags
//  2. Environment variables
//  3. Configuration file
//  4. Default values
func (l *Loader) Load(config interface{}) error {
	if err := validateStructPtr(config); err != nil {
		return err
	}

	// Apply default values before reading configuration
	defaults.SetDefaults(config)

	// Bind flags dynamically using reflection on `cmdx` tags if a flag set is provided
	if l.flags != nil {
		if err := bindFlags(l.v, l.flags, reflect.TypeOf(config).Elem(), ""); err != nil {
			return fmt.Errorf("failed to bind flags: %w", err)
		}
	}

	// Bind environment variables for all keys in the config
	keys, err := extractFlattenedKeys(config)
	if err != nil {
		return fmt.Errorf("failed to extract config keys: %w", err)
	}
	for _, key := range keys {
		if err := l.v.BindEnv(key); err != nil {
			return fmt.Errorf("failed to bind environment variable for key %q: %w", key, err)
		}
	}

	// Attempt to read the configuration file
	if err := l.v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Println("Warning: Config file not found. Falling back to defaults and environment variables.")
		}
	}

	// Unmarshal the merged configuration into the provided struct
	if err := l.v.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate the resulting configuration
	if err := validator.New().Struct(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return nil
}

// Get retrieves a configuration value by key.
func (l *Loader) Get(key string) interface{} {
	return l.v.Get(key)
}

// Set updates a configuration value in memory (not persisted to file).
func (l *Loader) Set(key string, value interface{}) {
	l.v.Set(key, value)
}

// Save writes the current configuration to the file specified during initialization.
func (l *Loader) Save() error {
	configFile := l.v.ConfigFileUsed()
	if configFile == "" {
		return errors.New("no configuration file specified for saving")
	}

	settings := l.v.AllSettings()
	content, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile(configFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write configuration to file: %w", err)
	}
	return nil
}

// View returns the current configuration as a formatted JSON string.
func (l *Loader) View() (string, error) {
	settings := l.v.AllSettings()
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format configuration as JSON: %w", err)
	}
	return string(data), nil
}

// validateStructPtr ensures the provided value is a pointer to a struct.
func validateStructPtr(value interface{}) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("load requires a pointer to a struct")
	}
	return nil
}

// extractFlattenedKeys retrieves all keys from the struct in a flattened format.
func extractFlattenedKeys(config interface{}) ([]string, error) {
	var structMap map[string]interface{}
	if err := mapstructure.Decode(config, &structMap); err != nil {
		return nil, err
	}
	flatMap, err := flatten.Flatten(structMap, "", flatten.DotStyle)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(flatMap))
	for k := range flatMap {
		keys = append(keys, k)
	}
	return keys, nil
}

// bindFlags dynamically binds flags to configuration fields based on `cmdx` tags.
func bindFlags(v *viper.Viper, flagSet *pflag.FlagSet, structType reflect.Type, parentKey string) error {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("cmdx")
		if tag == "" {
			continue
		}

		if parentKey != "" {
			tag = parentKey + "." + tag
		}

		if field.Type.Kind() == reflect.Struct {
			// Recurse into nested structs
			if err := bindFlags(v, flagSet, field.Type, tag); err != nil {
				return err
			}
		} else {
			flag := flagSet.Lookup(tag)
			if flag == nil {
				return fmt.Errorf("missing flag for tag: %s", tag)
			}
			if err := v.BindPFlag(tag, flag); err != nil {
				return fmt.Errorf("failed to bind flag for tag: %s, error: %w", tag, err)
			}
		}
	}
	return nil
}
